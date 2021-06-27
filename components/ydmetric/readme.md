@see  http://ydwiki.yidian-inc.com/pages/viewpage.action?pageId=7311921

test:
```
8c16GB
go test -test.v -test.bench ^BenchmarkNewMetric1$ -test.run ^$  -test.benchtime=10s  -benchmem 
goos: darwin
goarch: amd64
pkg: gitlab.yidian-inc.com/image/lego/components/ydmetric
BenchmarkNewMetric1
BenchmarkNewMetric1-8           23725195              2730 ns/op             321 B/op          5 allocs/op
PASS
ok      gitlab.yidian-inc.com/image/lego/components/ydmetric    70.563s

```

```
usage:
    ip := "127.0.0.1"
    port := "15688" 
    //只需初始化一次即可
    metric, err := NewMetric(ip, port, "Ydbot.statsd.contech", "image", "stat", 20)
    //如果需要全局tag
    metric.AddTag("host", "127.0.0.1")

    metric.Gauges("upload.b").TimeUnit(TimeUnitMin).Value("3").Done()
    //生成并发送 Ydbot.statsd.contech.image.stat.host-127_0_0_1.upload.b_gauges.1min:3|g
    metric.Timing("upload.b").TimeUnit(TimeUnitMin).Value("3").Done()
    //生成并发送 Ydbot.statsd.contech.image.stat.host-127_0_0_1.upload.b_timing.1min:3|ms
    

    //停止程序
    m.Stop()

```