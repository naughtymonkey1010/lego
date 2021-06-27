package grpc_client

import (
	"google.golang.org/grpc"
)

//grpc 客户端初始化
//var kacp = keepalive.ClientParameters{
//	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
//	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
//	PermitWithoutStream: true,             // send pings even without active streams
//}

type Client struct {
	Setting *Setting
	Conn    *grpc.ClientConn
}

type Setting struct {
	address string
}

func NewGrpcClient(addr string) (*Client, error) {
	setting := &Setting{
		address: addr,
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{Setting: setting, Conn: conn}, nil
}

func (c *Client) GetConn() *grpc.ClientConn {
	return c.Conn
}

func (c *Client) Close() {
	_ = c.Conn.Close()
}
