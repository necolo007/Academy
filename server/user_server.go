package main

import (
	"Academy/Model"
	"Academy/global"
	pb "Academy/pb/user"
	"Academy/utils"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"time"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterResp, error) {
	var UserNow, ExistingUser Model.User
	//查找用户
	if err := global.Db.WithContext(ctx).Where("email=? OR username = ?", req.Username).Find(&ExistingUser); err == nil {
		return &pb.RegisterResp{Success: false, UserId: 0}, nil
	}
	//存储用户信息
	UserNow = Model.User{Username: req.Username, Password: utils.Hash(req.Password), Email: req.Email}
	if err := global.Db.WithContext(ctx).Create(&UserNow); err != nil {
		return nil, err.Error
	}
	return &pb.RegisterResp{Success: true, UserId: int32(UserNow.ID)}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var UserNow Model.User
	//查找用户
	if err := global.Db.WithContext(ctx).Where("email =?", req.Email).First(&UserNow); err != nil {
		return &pb.LoginResp{Success: false, UserId: 0}, err.Error
	}

	//验证密码
	if utils.Hash(req.Password) != UserNow.Password {
		return &pb.LoginResp{Success: false, UserId: 0}, nil
	}
	token, _ := utils.GenerateToken(int32(UserNow.ID))
	return &pb.LoginResp{Success: true, UserId: int32(UserNow.ID), Token: token}, nil
}

func callService(client *UserServiceServer) {
	// 设置元数据
	md := metadata.New(map[string]string{
		"authorization": "Bearer <your_token>",
		"trace-id":      "12345",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 调用服务
	resp, err := client.Register(ctx, &pb.RegisterReq{Username: "Necolo007", Password: "123456", Email: "123@qq.com", ConfirmPassword: "123456"})
	if err != nil {
		log.Fatalf("Failed to call service: %v", err)
	}
	log.Printf("Response: %v", resp)
}

func main() {
	l, _ := net.Listen("tcp", ":8080")
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServiceServer{})
	err := s.Serve(l)
	if err != nil {
		log.Println(err)
	}
}
