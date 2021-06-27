package ydmetric

import (
	"strings"
)

const (
	//时间单位
	TimeUnitSec  = "1sec"
	TimeUnitMin  = "1min"
	TimeUnitHour = "1hour"
	TimeUnitDay  = "1day"
)

//接口类型
type Item interface {
	Item(item string)
	Value(value string)
	TimeUnit(unit string)
	Done()
}

type MetricItems struct {
	item       string
	itemSuffix string
	value      string
	//时间单位
	timeUnit string
	//后缀
	suffix string
	buffer *strings.Builder

	metric *Metric
}

func (m *MetricItems) Item(item string) *MetricItems {
	m.item = checkItem(item)
	return m
}

//设置时间单位 默认值 1sec
func (m *MetricItems) TimeUnit(unit string) *MetricItems {
	m.timeUnit = unit
	return m
}

//可设置值 默认值 1
func (m *MetricItems) Value(value string) *MetricItems {
	m.value = value
	return m
}

//debug 打印值
func (m *MetricItems) Debug() string {
	m.format()
	s := m.buffer.String()
	m.reset()
	return s
}

//设置完成
func (m *MetricItems) Done() {
	m.format()
	m.metric.queue <- m.buffer.String()
	m.reset()
}

func (m *MetricItems) reset() {
	m.restoreBuffer()
	m.item = ""
	m.itemSuffix = ""
	m.value = ""
	m.timeUnit = ""
	m.suffix = ""
	m.buffer = nil
	m.metric.metricItemsPool.Put(m)
}

func (m *MetricItems) restoreBuffer() {
	m.buffer.Reset()
	m.metric.stringBuilderPool.Put(m.buffer)
}

//生成字符串
func (m *MetricItems) format() *MetricItems {
	m.buffer.WriteString(m.metric.com.prefix)
	m.buffer.WriteString(m.item)
	m.buffer.WriteString(m.itemSuffix)
	m.buffer.WriteByte('.')
	m.buffer.WriteString(m.timeUnit)
	m.buffer.WriteByte(':')
	m.buffer.WriteString(m.value)
	m.buffer.WriteByte('|')
	m.buffer.WriteString(m.suffix)
	m.buffer.WriteByte('\n')
	return m
}

//检查item 将- 转换成_
func checkItem(s string) string {
	if strings.Contains(s, "-") {
		return strings.Replace(s, "-", "_", -1)
	}
	return s
}
