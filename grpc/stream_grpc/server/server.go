package main

import (
    "io"
    "fmt"
	"net"
    "time"
	"google.golang.org/grpc"
	pb "stream_grpc/proto"
)

const PORT = ":50055"

type server struct {
    pb.UnimplementedGreeterServer
}

// Server-side streaming RPC
func (s *server) GetStream(req *pb.StreamRequestData, res pb.Greeter_GetStreamServer) error {
    i := 0
    for {
        i ++
        _ = res.Send(&pb.StreamResponseData{
            Data: fmt.Sprintf("server send %v", i),
        })
        time.Sleep(time.Second)
        if i >= 10 {
            break
        }
    }
	return nil
}

// Client-side Steaming RPC
func (s *server) PutStream(stream pb.Greeter_PutStreamServer) error {
    for {
        data, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.StreamResponseData{
                Data: "server stream will be close",
            })
        }
        if err != nil {
            fmt.Println(err)
            return err
        }
        fmt.Println(data.Data)
    }
}

// Bidirectional streaming RPC
func (s *server) AllStream(stream pb.Greeter_AllStreamServer) error {
    for {
        in, err := stream.Recv()
        if err == io.EOF{
            fmt.Println("client has been end, %v", time.Now())
            return nil
        }
        if err != nil {
            fmt.Println("server recv has err: %v", err)
            return err
        }

        fmt.Println("stream recv info %v", in)

        if err := stream.Send(&pb.StreamResponseData{
            Data: fmt.Sprintf("server send"),
        }); err != nil {
            return err
        }
    }
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
