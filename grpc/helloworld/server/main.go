package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "grpc/pb"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer
// 继承了父类UnimplementedGreeterServer
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context,
	in *pb.HelloRequest) (*pb.HelloReply, error) {

	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// 监听所有网卡的50051端口
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建grpc服务
	s := grpc.NewServer()

	/** 注册接口服务
	 *  以定义proto时的service为单位进行注册，服务中可以有多个方法
	 *  proto编译时会为每个service生成Register***Server方法
	 */
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	// 将监听交给rpc服务处理
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
