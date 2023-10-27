package main

import (
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

func myStreamClientInterceptor1(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Println("[pre] my stream client interceptor 1", method)
	stream, err := streamer(ctx, desc, cc, method, opts...)
	return &myClientStreamWrapper{stream}, err
}

type myClientStreamWrapper struct {
	grpc.ClientStream
}

func (s *myClientStreamWrapper) SendMsg(m any) error {
	log.Println("[pre message] my stream client interceptor 1", m)
	err := s.ClientStream.SendMsg(m)
	return err
}

func (s *myClientStreamWrapper) RecvMsg(m any) error {
	err := s.ClientStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		log.Println("[post message] my stream client interceptor 1", m)
	}

	return err
}

func (s *myClientStreamWrapper) CloseSend() error {
	err := s.ClientStream.CloseSend()

	log.Println("[post] my stream client interceptor 1")
	return err
}
