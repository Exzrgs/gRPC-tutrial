package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	hellopb "mygrpc/pkg/grpc"
	"os"

	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	scanner *bufio.Scanner
	client  hellopb.GreetingServiceClient
)

func main() {
	fmt.Println("start gRPC client")

	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,

		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("connection fatal")
		return
	}
	defer conn.Close()

	client = hellopb.NewGreetingServiceClient(conn)

	for {
		fmt.Println("1: Normal")
		fmt.Println("2: ServerStream")
		fmt.Println("3: ClientStream")
		fmt.Println("4: BiStream")
		fmt.Println("5: exit")
		fmt.Println("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello()
		case "2":
			HelloServerStream()
		case "3":
			HelloClientStream()
		case "4":
			HelloBiStreams()
		case "5":
			fmt.Println("bye")
			goto M
		}
	}
M:
}

func getName() string {
	fmt.Println("please enter your name")
	scanner.Scan()
	name := scanner.Text()
	return name
}

func Hello() {
	name := getName()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			fmt.Printf("code: %d message: %s\n", stat.Code(), stat.Message())
			fmt.Printf("details: %s\n", stat.Details())
		} else {
			fmt.Println(err)
		}
		return
	}

	fmt.Println(res.GetMessage())
	return
}

func HelloServerStream() {
	name := getName()

	req := &hellopb.HelloRequest{
		Name: name,
	}

	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			println("all response are recieved")
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res)
	}
}

/*
リクエストを複数送る
受け取る
表示
*/
func HelloClientStream() {
	count := 4

	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < count; i++ {
		name := getName()
		req := &hellopb.HelloRequest{
			Name: name,
		}
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.GetMessage())
	return
}

/*
何個も送る。おなじ
*/
func HelloBiStreams() {
	stream, err := client.HelloBiStreams(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	var sendNum = 3
	fmt.Printf("enter %d names\n", sendNum)

	var sendEnd bool
	var sendCount = 0
	for {
		if sendEnd != true {
			name := getName()
			if err := stream.Send(&hellopb.HelloRequest{
				Name: name,
			}); err != nil {
				fmt.Println(err)
				return
			}
			sendCount++
			if sendCount == sendNum {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			if sendEnd == true {
				fmt.Println("got all response")
				return
			}
			continue
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.GetMessage())
	}
}
