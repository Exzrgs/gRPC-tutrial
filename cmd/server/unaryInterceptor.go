package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func myUnaryServerInterceptor1(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Println("[pre] my unary server interceptor 1: ", req)
	res, err := handler(ctx, req)
	log.Println("[post] my unary server interceptor 1: ", res)
	return res, err
}
