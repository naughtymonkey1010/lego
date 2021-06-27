package ydmetric

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {
	ip := "127.0.0.1"
	port := "15688"

	go func() {
		initUdpServer(t, ip, port)
	}()

	m, err := NewMetric(ip, port, "Ydbot.statsd.contech", "image", "stat", 20)

	if err != nil {
		t.Fatalf("new metric err:%s", err.Error())
	}

	t.Logf("cap:%d len:%d poollen:%d", m.Cap(), m.Len(), m.SendPoolRunning())

	t.Logf("%s", m.Counter("upload.a").Debug())

	m.Counter("upload.a").Done()
	m.Gauges("upload.b").TimeUnit(TimeUnitMin).Value("3").Done()
	m.Sets("upload.c").TimeUnit(TimeUnitHour).Value("4").Done()
	m.Timing("upload.d").TimeUnit(TimeUnitDay).Value("100").Done()
	m.Rate("upload.r").Value("200").Done()

	m.Stop()
	time.Sleep(5 * time.Second)
	t.Logf("cap:%d len:%d poollen:%d", m.Cap(), m.Len(), m.SendPoolRunning())

	time.Sleep(5 * time.Second)

}

func TestTag_String(t *testing.T) {
	ip := "127.0.0.1"
	port := "15688"

	go func() {
		initUdpServer(t, ip, port)
	}()

	metric, err := NewMetric(ip, port, "Ydbot.statsd.contech", "image", "stat", 20)

	if err != nil {
		t.Fatalf("new metric err:%s", err.Error())
	}

	t.Logf("cap:%d len:%d poollen:%d", metric.Cap(), metric.Len(), metric.SendPoolRunning())

	metric.AddTag("host", "127.0.0.1")

	t.Log(metric.Gauges("upload.b").Value("1").Debug())

	assert.Equal(t, metric.Gauges("upload.b").Value("1").Debug(), "Ydbot.statsd.contech.image.stat.host-127_0_0_1.upload.b_gauges.1sec:1|g\n")

	metric.ResetTag()
	assert.Equal(t, metric.Gauges("upload.b").Value("1").Debug(), "Ydbot.statsd.contech.image.stat.upload.b_gauges.1sec:1|g\n")
	t.Log(metric.Gauges("upload.b").Value("1").Debug())

	t.Log(metric.Gauges("upload.b").TimeUnit(TimeUnitMin).Value("3").Debug())

	t.Log(metric.Gauges("upload.b").TimeUnit(TimeUnitMin).Value("3").Debug())
	t.Log(metric.Timing("upload.b").TimeUnit(TimeUnitMin).Value("3").Debug())

}

func genrateItem(m *Metric, n int) {
	for i := 0; i < n; i++ {
		m.Gauges("upload.b").TimeUnit(TimeUnitMin).Value("3").Done()
	}
}

func BenchmarkNewMetric1(b *testing.B) {
	ip := "127.0.0.1"
	port := "15688"

	m, err := NewMetric(ip, port, "Ydbot.statsd.contech", "image", "stat", 20)

	if err != nil {
		b.Fatalf("new metric err:%s", err.Error())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		genrateItem(m, 1) // run fib(30) b.N times
	}
}

func BenchmarkNewMetric10(b *testing.B) {
	ip := "127.0.0.1"
	port := "15688"

	m, err := NewMetric(ip, port, "Ydbot.statsd.contech", "image", "stat", 20)

	if err != nil {
		b.Fatalf("new metric err:%s", err.Error())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		genrateItem(m, 10) // run fib(30) b.N times
	}
}

func BenchmarkNewMetric100(b *testing.B) {
	ip := "127.0.0.1"
	port := "15688"

	m, err := NewMetric(ip, port, "Ydbot.statsd.contech", "image", "stat", 20)

	if err != nil {
		b.Fatalf("new metric err:%s", err.Error())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		genrateItem(m, 100) // run fib(30) b.N times
	}
}

func BenchmarkNewMetric1000(b *testing.B) {
	ip := "127.0.0.1"
	port := "15688"

	m, err := NewMetric(ip, port, "Ydbot.statsd.contech", "image", "stat", 20)

	if err != nil {
		b.Fatalf("new metric err:%s", err.Error())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		genrateItem(m, 1000) // run fib(30) b.N times
	}
}

func BenchmarkNewMetric10000(b *testing.B) {
	ip := "127.0.0.1"
	port := "15688"

	m, err := NewMetric(ip, port, "Ydbot.statsd.contech", "image", "stat", 20)

	if err != nil {
		b.Fatalf("new metric err:%s", err.Error())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		genrateItem(m, 10000) // run fib(30) b.N times
	}
}
