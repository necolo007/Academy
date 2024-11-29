package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"sync"

	"Academy/DB" // 引入数据库初始化模块
	"Academy/pb/product"
	"Academy/pb/user" // 引入生成的 proto 文件
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	// 初始化数据库
	DB.InitConfig()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go RegisterGateway(&wg)
	wg.Wait()
}

func RegisterGateway(wg *sync.WaitGroup) {
	defer wg.Done()
	mux := runtime.NewServeMux()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err1 := user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:5001", opts)
	err2 := product.RegisterProductCatalogServiceHandlerFromEndpoint(ctx, mux, "localhost:5002", opts)
	if err1 != nil || err2 != nil {
		log.Fatalln("failed to register gateway for user service: ")
	}

	log.Println("Starting HTTP gateway")
	if m := http.ListenAndServe(":8080", mux); m != nil {
		log.Fatalf("failed to serve HTTP: %v", m)
	}
}
