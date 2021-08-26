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
// 服务端流模式
func serverSideRPC() {
    conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    defer fmt.Println("conn close")
    c := pb.NewGreeterClient(conn)

    data := &pb.StreamRequestData{
        Data: "client send one msg",
    }

    stream, _ := c.GetStream(context.Background(), data)
    for {
        data, err := stream.Recv() // socket send recv
        if err == io.EOF {
            fmt.Println("client recv EOF")
            break
        } else if err != nil {
            fmt.Println("client recv has failed, %v", err)
            break
        }
        fmt.Println(data.Data)
    }
}

// Client-side streaming RPC
// 客户端流模式
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
        fmt.Println("client put has error: %v", err)
    }

    for i := 0; i <= 5; i ++ {
        data := &pb.StreamRequestData{
            Data: fmt.Sprintf("client send stream %d", i),
        }
        if err := stream.Send(data); err != nil{
            fmt.Println("client stream send has error %v", err)
        }
        time.Sleep(time.Second)
    }

    reply, err := stream.CloseAndRecv()
    if err != nil {
        fmt.Println("client recv has error %v", err)
    }
    fmt.Println("client recv info ", reply)
}

// 双向流模式
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
                fmt.Println("server has been close", time.Now())
                return
            }
            if err != nil {
                fmt.Println("failed to receive err: %v", err)
            }
            fmt.Println("client recv data: ", in)
        }
    }()

    for i := 0; i <= 5; i ++ {
        data := &pb.StreamRequestData{
            Data: fmt.Sprintf("client send %v", i),
        }
        if err := stream.Send(data); err != nil {
            fmt.Println("Failed to send a note %v", err)
        }
        time.Sleep(time.Second)
    }

    stream.CloseSend()
    wg.Wait()
}

func main(){
    //serverSideRPC()
    //clientSideRPC()
    bidirectionRPC()
}
