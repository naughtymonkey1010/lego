package ydmetric

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func initUdpServer(t *testing.T, ip, port string) {
	u, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		t.Errorf("udp server error:%s", err.Error())
	}
	udpConn, err := net.ListenUDP("udp", u)
	if err != nil {
		t.Fatalf("udp server error:%s", err.Error())
	}
	for {
		var data [1024]byte
		n, addr, err := udpConn.ReadFromUDP(data[:])
		if err != nil {
			t.Fatalf("udp server recive error:%s", err.Error())
		}
		go func() {
			// 返回数据
			fmt.Printf("Addr:%s,data:%v count:%d \n", addr, string(data[:n]), n)
			_, err := udpConn.WriteToUDP(data[:n], addr)
			if err != nil {
				fmt.Println("write to udp server failed,err:", err)
			}
		}()
	}

}

func TestUdpConn(t *testing.T) {
	ip := "127.0.0.1"
	port := "15688"

	go func() {
		initUdpServer(t, ip, port)
	}()

	conn, err := UdpConn(ip, port, 5*time.Millisecond)
	if err != nil {
		t.Fatalf("connect udp server error:%s", err.Error())
	}
	_, _ = conn.Write([]byte("hello world"))

	time.Sleep(10 * time.Second)
}
