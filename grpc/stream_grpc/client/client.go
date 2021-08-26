package main

import (
    "fmt"
    "time"
    "context"
    "sync"
    "io"
    "google.golang.org/grpc"
    pb "stream_grpc/proto"
)

// Server-side streaming RPC
func serverSideRPC() {
    conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    defer fmt.Println("conn close")
    c := pb.NewGreeterClient(conn)

    res, _ := c.GetStream(context.Background(), &pb.StreamRequestData{Data: "client send"})
    for {
        a, err := res.Recv() // socket send recv
        if err == io.EOF {
            fmt.Println("get stream end")
            break
        } else if err != nil {
            fmt.Println(err)
            break
        }
        fmt.Println(a.Data)
    }
}

// Client-side streaming RPC
func clientSideRPC() {
    conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    defer fmt.Println("conn close")
    c := pb.NewGreeterClient(conn)

    stream, err := c.PutStream(context.Background())
    if err != nil {
        fmt.Println("put stream has error: %v", err)
    }
    for i := 0; i <= 5; i ++ {
        if err := stream.Send(&pb.StreamRequestData{
            Data: fmt.Sprintf("client send stream %d", i),
        }); err != nil{
            fmt.Println("client stream send %v", err)
        }
        time.Sleep(time.Second)
    }
    reply, err := stream.CloseAndRecv()
    if err != nil {
        fmt.Println("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
    }
    fmt.Println("recv server reply: %v", reply)
}

func bidirectionRPC() {
    conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    defer fmt.Println("conn close")
    c := pb.NewGreeterClient(conn)

    stream, _ := c.AllStream(context.Background())

    var wg sync.WaitGroup
    wg.Add(1)
    go func(){
        defer wg.Done()
        for{
            in, err := stream.Recv()
            if err == io.EOF {
                fmt.Println("server has been close, %v", time.Now())
                return
            }
            if err != nil {
                fmt.Println("failed to receive err: %v", err)
            }
            fmt.Println("get message from server: %v", in)
        }
    }()

    for i := 0; i <= 5; i ++ {
        if err := stream.Send(&pb.StreamRequestData{
            Data: "client send",
        }); err != nil {
            fmt.Println("Failed to send a note %v", err)
        }
        time.Sleep(time.Second)
    }

    wg.Wait()
    stream.CloseSend()
}

func main(){
    //serverSideRPC()
    //clientSideRPC()
    bidirectionRPC()
}
