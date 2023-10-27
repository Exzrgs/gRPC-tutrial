package main

import (
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

func myStreamServerInterceptor1(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[pre stream] my stream server interceptor 1: ", info)
	err := handler(srv, &myServerStreamWrapper1{ss})
	log.Println("[post stream] my stream server interceptor 1: ")
	return err
}

type myServerStreamWrapper1 struct {
	grpc.ServerStream
}

func (s *myServerStreamWrapper1) RecvMsg(m any) error {
	err := s.ServerStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		log.Println("[pre message] my stream server interceptor 1: ", m)
	}

	return err
}

func (s *myServerStreamWrapper1) SendMsg(m any) error {
	log.Println("[post message] my stream server interceptor 1: ", m)

	err := s.ServerStream.SendMsg(m)
	return err
}
