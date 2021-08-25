package main

import (
	"google.golang.org/grpc"
	"net"
	"stream_grpc/proto"
)

const PORT = ":50052"

type Server struct {
}

// Server-side streaming RPC
func (s *Server) GetStream(*StreamRequestData, Greeter_GetStreamServer) error {
	return nil
}

// Client-side Steaming RPC
func (s *Server) PutStream(Greeter_PutStreamServer) error {
	return nil
}

// Bidirectional streaming RPC
func (s *Server) AllStream(Greeter_AllStreamServer) error {
	return nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})
}
