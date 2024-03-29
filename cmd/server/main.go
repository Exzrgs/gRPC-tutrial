package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	hellopb "mygrpc/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func (*myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	// わざとエラーを発生させるコード
	// stat := status.New(codes.Unknown, "Unknown error occurred")
	// stat, _ = stat.WithDetails(&errdetails.DebugInfo{
	// 	Detail: "Detail reason of error",
	// })
	// err := stat.Err()[

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println(md)
	}

	headerMD := metadata.New(map[string]string{
		"type": "unary",
		"from": "server",
		"in":   "header",
	})
	if err := grpc.SetHeader(ctx, headerMD); err != nil {
		return nil, err
	}

	trailerMD := metadata.New(map[string]string{
		"type": "unary",
		"from": "server",
		"in":   "trailer",
	})
	if err := grpc.SetTrailer(ctx, trailerMD); err != nil {
		return nil, err
	}

	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (*myServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.GreetingService_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hellopb.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello %s!", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (*myServer) HelloClientStream(stream hellopb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello %v!", nameList)
			return stream.SendAndClose(&hellopb.HelloResponse{
				Message: message,
			})
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}

/*
毎回送信するように変えてる
*/
func (*myServer) HelloBiStreams(stream hellopb.GreetingService_HelloBiStreamsServer) error {
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Println(md)
	}

	headerMD := metadata.New(map[string]string{
		"type": "stream",
		"from": "server",
		"in":   "header",
	})
	if err := stream.SendHeader(headerMD); err != nil {
		return err
	}

	trailerMD := metadata.New(map[string]string{
		"type": "stream",
		"from": "server",
		"in":   "trailer",
	})
	stream.SetTrailer(trailerMD)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all request are received")
			return nil
		}
		if err != nil {
			return err
		}

		message := fmt.Sprintf("Hello, %v!", req.GetName())
		if err := stream.Send(&hellopb.HelloResponse{
			Message: message,
		}); err != nil {
			return err
		}
	}
}

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			myStreamServerInterceptor1,
			myStreamServerInterceptor2,
		),
	)
	hellopb.RegisterGreetingServiceServer(s, NewMyServer())

	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
