package Interceptor

import (
	"Academy/utils"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

// AuthInterceptor 是一个用于身份验证的 gRPC 服务器拦截器
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	var token string
	if val, ready := md["authorization"]; ready {
		token = val[0]
	}
	log.Println("token:", token)
	userClaims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	ctx = context.WithValue(ctx, "userClaims", userClaims)

	return handler(ctx, req)
}
