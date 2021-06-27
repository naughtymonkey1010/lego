package ydmetric

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

const (
	//队列长度 百万
	QueueLength = 10e6
)

type Metric struct {
	//upd发送配置
	udpConfig UdpConfig
	//com 项目组织
	com Com
	//metricItems Pool
	metricItemsPool   sync.Pool
	stringBuilderPool sync.Pool
	//异步发送使用队列
	queue chan string
	//发送容量
	sendGoroutineCap int
	//发送池
	sendPool *ants.PoolWithFunc

	//关闭标志
	closeChan chan struct{}

	lock sync.Mutex
}

type Com struct {
	//team组
	team    string
	group   string
	project string
	tags    []Tag
	prefix  string
}

//发送配置项
type UdpConfig struct {
	host     string
	port     string
	duration time.Duration
}

type Tag struct {
	k, v string
}

func (t *Tag) String() string {
	return checkTagString(fmt.Sprintf("%s-%s", t.k, t.v))
}

//异步发送初始化
func NewMetric(host, port, team, group, project string, cap int) (*Metric, error) {
	if len(host) < 1 {
		host = os.Getenv("K8S_YIDIAN_LOCAL_IP")
	}
	if len(port) < 0 {
		port = "15688"
	}
	//链接尝试
	conn, err := UdpConn(host, port, 500*time.Millisecond)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("udp host or port error can not connected error:%s", err.Error()))
	}
	_ = conn.Close()
	if len(team) < 0 {
		return nil, errors.New("arg team must not empty")
	}
	if len(group) < 0 {
		return nil, errors.New("arg group must not empty")
	}
	if len(project) < 0 {
		return nil, errors.New("arg project must not empty")
	}
	m := &Metric{
		udpConfig: UdpConfig{
			host:     host,
			port:     port,
			duration: 5 * time.Second,
		},
		com: Com{
			team:    team,
			group:   group,
			project: project,
			tags:    make([]Tag, 0),
			prefix:  fmt.Sprintf("%s.%s.%s.", team, group, project),
		},
		//全局队列
		queue: make(chan string, QueueLength),

		closeChan: make(chan struct{}),

		stringBuilderPool: sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
		lock: sync.Mutex{},
	}

	m.metricItemsPool = sync.Pool{
		New: func() interface{} {
			return &MetricItems{
				metric: m,
			}
		},
	}

	if cap < 1 {
		cap = 20
	}
	m.sendGoroutineCap = cap

	//实例化队列池
	p, err := NewConsumerPool(m.sendGoroutineCap, m.Process)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("init send pool error err:%s", err.Error()))
	}
	m.sendPool = p
	//启动一半的
	m.start()
	return m, nil
}

func (m *Metric) start() {
	go m.monitor()
	m.initSendPool(m.sendGoroutineCap / 4)
}

//停止程序
func (m *Metric) Stop() {
	m.lock.Lock()
	defer m.lock.Unlock()
	//关闭队列
	close(m.closeChan)
	//关闭Queue 防止程序报错
	//close(m.queue)
	m.sendPool.Release()
}

//监控队列 数量不够增加消费数量
func (m *Metric) monitor() {
	for {
		select {
		case <-time.After(60 * time.Second):
			m.lock.Lock()
			//超过一半
			if m.Len() >= m.Cap()/2 {
				m.initSendPool(m.sendPool.Cap() - m.sendPool.Running())
			} else if m.Len() >= m.Cap()/4 {
				//超过1/4
				m.initSendPool((m.sendPool.Cap() - m.sendPool.Running()) / 2)
			} else if m.sendPool.Running() < m.sendGoroutineCap/5 {
				//保持20%运行
				m.initSendPool(m.sendGoroutineCap/5 - m.sendPool.Running())
			}
			m.lock.Unlock()
		case <-m.closeChan:
			goto END
		}
	}
END:
}

func (m *Metric) initSendPool(num int) {
	for i := 0; i < num; i++ {
		_ = m.sendPool.Invoke("")
	}
}

func (m *Metric) Len() int {
	return len(m.queue)
}

func (m *Metric) Cap() int {
	return cap(m.queue)
}

func (m *Metric) SendPoolRunning() int {
	return m.sendPool.Running()
}

//增加tag项
func (m *Metric) AddTag(k, v string) {
	t := buildTag(k, v)
	m.com.tags = append(m.com.tags, t)
	m.tagsFormat()
}

//删除tag项目
func (m *Metric) ResetTag() {
	m.com.tags = m.com.tags[0:0]
	//重新组装
	m.tagsFormat()
}

func (m *Metric) tagsFormat() {
	prefix := fmt.Sprintf("%s.%s.%s.", m.com.team, m.com.group, m.com.project)
	for _, v := range m.com.tags {
		prefix += v.String() + "."
	}
	m.com.prefix = prefix
}

//单位时间累计计数
func (m *Metric) Counter(item string) *MetricItems {
	items := m.metricItemsPool.Get().(*MetricItems)
	items.Item(item)
	items.itemSuffix = "_counter"
	//默认时间
	items.timeUnit = TimeUnitSec
	//默认value
	items.value = "1"
	items.suffix = "c"
	items.buffer = m.stringBuilderPool.Get().(*strings.Builder)

	return items
}

//一直累加 如统计当前总用户数
func (m *Metric) Gauges(item string) *MetricItems {
	items := m.metricItemsPool.Get().(*MetricItems)
	items.Item(item)
	items.itemSuffix = "_gauges"
	//默认时间
	items.timeUnit = TimeUnitSec
	//默认value
	items.value = "1"
	items.suffix = "g"
	items.buffer = m.stringBuilderPool.Get().(*strings.Builder)
	return items
}

//某个metric unique事件的个数
func (m *Metric) Sets(item string) *MetricItems {
	items := m.metricItemsPool.Get().(*MetricItems)
	items.Item(item)
	items.itemSuffix = "_sets"
	//默认时间
	items.timeUnit = TimeUnitSec
	//默认value
	items.value = "1"
	items.suffix = "s"
	items.buffer = m.stringBuilderPool.Get().(*strings.Builder)
	return items
}

//耗时操作
func (m *Metric) Timing(item string) *MetricItems {
	items := m.metricItemsPool.Get().(*MetricItems)
	items.Item(item)
	items.itemSuffix = "_timing"
	//默认时间
	items.timeUnit = TimeUnitSec
	//默认value
	items.value = "1"
	items.suffix = "ms"
	items.buffer = m.stringBuilderPool.Get().(*strings.Builder)
	return items
}

//概率指标
func (m *Metric) Rate(item string) *MetricItems {
	items := m.metricItemsPool.Get().(*MetricItems)
	items.Item(item)
	items.itemSuffix = "_rate"
	//默认时间
	items.timeUnit = TimeUnitSec
	//默认value
	items.value = "1"
	items.suffix = "rate"
	items.buffer = m.stringBuilderPool.Get().(*strings.Builder)
	return items
}

//处理发送程序
func (m *Metric) Process(args interface{}) {
	_ = args
	//进行连接
	conn, err := UdpConn(m.udpConfig.host, m.udpConfig.port, m.udpConfig.duration)
	if err != nil {
		return
	}
	for {
		select {
		//接收进行发送
		case item, ok := <-m.queue:
			if ok {
				_, _ = conn.Write([]byte(item))
			} else {
				goto END
			}
		case <-time.After(1200 * time.Second):
			_ = conn.Close()
			goto END
		case <-m.closeChan:
			_ = conn.Close()
			goto END
		}
	}
END:
}

func buildTag(k, v string) Tag {
	return Tag{
		k: k,
		v: v,
	}
}

//when tags contains "." convert it to "_"
func checkTagString(s string) string {
	if strings.Contains(s, ".") {
		return strings.Replace(s, ".", "_", -1)
	}
	return s
}
