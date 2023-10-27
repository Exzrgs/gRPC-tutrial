package main

import (
	"errors"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func myStreamServerInterceptor1(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[pre stream] my stream server interceptor 1: ", info)
	err := handler(srv, &myServerStreamWrapper{ss, 1})
	log.Println("[post stream] my stream server interceptor 1: ")
	return err
}

type myServerStreamWrapper struct {
	grpc.ServerStream
	number int
}

func (s *myServerStreamWrapper) RecvMsg(m any) error {
	err := s.ServerStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		logMsg := fmt.Sprintf("[pre message] my stream server interceptor %d:", s.number)
		log.Println(logMsg, m)
	}

	return err
}

func (s *myServerStreamWrapper) SendMsg(m any) error {
	logMsg := fmt.Sprintf("[post message] my stream server interceptor %d:", s.number)
	log.Println(logMsg, m)

	err := s.ServerStream.SendMsg(m)
	return err
}

/////////////////////////////////////////////////////////////////////////////////////

func myStreamServerInterceptor2(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[pre stream] my stream server interceptor 2: ", info)
	err := handler(srv, &myServerStreamWrapper{ss, 2})
	log.Println("[post stream] my stream server interceptor 2: ")
	return err
}
