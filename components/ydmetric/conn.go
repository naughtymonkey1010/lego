package ydmetric

import (
	"io"
	"net"
	"time"
)

//udp conn 连接
func UdpConn(host, port string, duration time.Duration) (io.WriteCloser, error) {
	conn, err := net.DialTimeout("udp", host+":"+port, duration)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
