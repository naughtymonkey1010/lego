package grpc_server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//var kaep = keepalive.EnforcementPolicy{
//	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
//	PermitWithoutStream: true,            // Allow pings even when there are no active streams
//}
//
//var kasp = keepalive.ServerParameters{
//	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
//	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
//	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
//	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
//	Timeout:               5 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
//}

type GrpcServer struct {
	Setting *Setting
	Server  *grpc.Server
}

type Setting struct {
	Host string
	Port int
}

func NewGrpcServer(setting *Setting) (*GrpcServer, error) {
	//初始化grpc服务
	//grpcServer := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	grpcServer := grpc.NewServer()
	//注册路由
	//RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	return &GrpcServer{Setting: setting, Server: grpcServer}, nil
}

//注册server
func (s *GrpcServer) RegisterService(f func(s *grpc.Server, srv interface{}), ss interface{}) {
	f(s.Server, ss)
}

func (s *GrpcServer) Run() error {
	//反射注册 用于查询接口功能
	reflection.Register(s.Server)

	addr := fmt.Sprintf("%s:%d", s.Setting.Host, s.Setting.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	err = s.Server.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *GrpcServer) RunAsync() {
	go func() {
		_ = s.Run()
	}()
}

func (s *GrpcServer) GracefulShutdown() {
	s.Server.GracefulStop()
}
