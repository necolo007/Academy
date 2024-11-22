package main

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"

	"Academy/DB"         // 引入数据库初始化模块
	pb "Academy/pb/user" // 引入生成的 proto 文件
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	// 初始化数据库
	DB.InitConfig()

	// 启动 HTTP Gateway
	mux := runtime.NewServeMux(runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		md := metadata.MD{}
		if token := req.Header.Get("Authorization"); token != "" {
			md["authorization"] = []string{token}
		}
		return md
	}))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, ":8080", opts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("Starting HTTP gateway on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
