package main

import (
    "io"
    "fmt"
	"net"
    "time"
    "sync"
	"google.golang.org/grpc"
	pb "stream_grpc/proto"
)

const PORT = ":50055"

type server struct {
    pb.UnimplementedGreeterServer
}

// Server-side streaming RPC
// 服务端流模式
func (s *server) GetStream(req *pb.StreamRequestData, stream pb.Greeter_GetStreamServer) error {
    for i := 0; i < 5; i ++ {
        data := &pb.StreamResponseData{
            Data: fmt.Sprintf("server send %v", i),
        }
        if err := stream.Send(data); err != nil {
            fmt.Println("server send stream failed: %v", err)
            return err
        }
        time.Sleep(time.Second)
    }
	return nil
}

// Client-side Steaming RPC
// 客户端流模式
func (s *server) PutStream(stream pb.Greeter_PutStreamServer) error {
    for {
        data, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.StreamResponseData{
                Data: "server will be close",
            })
        }
        if err != nil {
            fmt.Println(err)
            return err
        }
        fmt.Println("server recv info: ", data.Data)
    }
}

// Bidirectional streaming RPC
// 双向流模式
func (s *server) AllStream(stream pb.Greeter_AllStreamServer) error {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            in, err := stream.Recv()
            if err == io.EOF{
                fmt.Println("client has been end", time.Now())
                return
            }
            if err != nil {
                fmt.Println("server recv has err: %v", err)
                return
            }

            fmt.Println("server recv info ", in)
        }
    }()

    for i := 0; i < 5; i ++ {
        data := &pb.StreamResponseData{
            Data: fmt.Sprintf("server send %v", i),
        }
        if err := stream.Send(data); err != nil {
            return err
        }
    }

    wg.Wait()
    return nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
    err = s.Serve(lis)
    if err != nil {
        panic(err)
    }
}
