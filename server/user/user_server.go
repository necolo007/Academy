package main

import (
	"Academy/DB"
	"Academy/Model"
	"Academy/global"
	pb "Academy/pb/user"
	"Academy/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"sync"
	"time"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterResp, error) {
	DB.InitConfig()
	// 检查具体字段是否为空
	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		return nil, fmt.Errorf("email, password, and confirm_password are required")
	}
	// 提取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract metadata from context")
	}
	fmt.Println("Metadata:", md)

	var UserNow Model.User
	//查找用户
	if res := global.Db.Where("email=? OR username = ?", req.Email, req.Username).Find(&UserNow); res.Error != nil {
		return &pb.RegisterResp{Success: false, UserId: 0}, nil
	}
	//存储用户信息
	UserNow = Model.User{Username: req.Username, Password: utils.Hash(req.Password), Email: req.Email}
	if err := global.Db.Create(&UserNow); err != nil {
		return nil, err.Error
	}
	return &pb.RegisterResp{Success: true, UserId: int32(UserNow.ID)}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	DB.InitConfig()
	if req.Email == "" || req.Password == "" || req.Username == "" {
		return nil, fmt.Errorf("email, password, and username are required")
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	// 提取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract metadata from context")
	}
	fmt.Println("Metadata:", md)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var UserNow Model.User
	//查找用户
	if res := global.Db.Where("email =? OR username = ?", req.Email, req.Username).First(&UserNow); res.Error != nil {
		return &pb.LoginResp{Success: false, UserId: 0}, res.Error
	}
	//验证密码
	if utils.Hash(req.Password) != UserNow.Password {
		return &pb.LoginResp{Success: false, UserId: 0}, nil
	}
	token, _ := utils.GenerateToken(UserNow)
	return &pb.LoginResp{Success: true, UserId: int32(UserNow.ID), Token: token}, nil
}

func main() {
	l, _ := net.Listen("tcp", ":5001")
	log.Println("Server is running on port 5001")
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServiceServer{})
	err := s.Serve(l)
	if err != nil {
		log.Println(err)
	}
}

func RegisterGrpcServer(wg *sync.WaitGroup) {
	defer wg.Done()
	s := grpc.NewServer()
	UserServer := &UserServiceServer{}
	pb.RegisterUserServiceServer(s, UserServer)
}
