package main

import (
	"Academy/DB"
	I "Academy/Interceptor"
	"Academy/Model"
	"Academy/global"
	pb "Academy/pb/product"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type ProductServiceServer struct {
	pb.UnimplementedProductCatalogServiceServer
}

func ConvertToGrpcProduct(p *Model.Product) *pb.Product {
	return &pb.Product{
		Id:          uint32(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Picture:     p.Picture,
		Price:       float32(p.Price),
		Sort:        p.Sort,
	}
}

func (P *ProductServiceServer) ListProducts(ctx context.Context, req *pb.ListProductsReq) (*pb.ListProductsResp, error) {
	DB.InitConfig()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(400, "failed to extract metadata from context")
	}
	log.Println("Metadata:", md)

	var Product []Model.Product
	offset := int(req.PageSize) * int(req.Page-1)
	//按照分类过滤
	if req.Sort != "" {
		global.Db.WithContext(ctx).Where("Sort LIKE ?", fmt.Sprintf("%%%s%%", req.Sort)).Find(&Product)
	}
	//分页查询
	if err := global.Db.WithContext(ctx).Offset(offset).Limit(int(req.PageSize)).Find(&Product).Error; err != nil {
		return nil, err
	}

	//回复resp
	var grpcProduct []*pb.Product
	for _, p := range Product {
		grpcProduct = append(grpcProduct, ConvertToGrpcProduct(&p))
	}
	return &pb.ListProductsResp{Products: grpcProduct}, nil
}

func (P *ProductServiceServer) GetProduct(ctx context.Context, req *pb.GetProductReq) (*pb.GetProductResp, error) {
	DB.InitConfig()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(400, "failed to extract metadata from context")
	}
	log.Println("Metadata:", md)

	var product Model.Product
	if err := global.Db.WithContext(ctx).Where("ID = ? OR Name = ?", req.Id, req.Name).Find(&product); err.Error != nil {
		return nil, err.Error
	}
	return &pb.GetProductResp{Product: ConvertToGrpcProduct(&product)}, nil
}

func (P *ProductServiceServer) SearchProducts(ctx context.Context, req *pb.SearchProductsReq) (*pb.SearchProductsResp, error) {
	DB.InitConfig()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(400, "failed to extract metadata from context")
	}
	log.Println("Metadata:", md)

	var product []Model.Product
	err := global.Db.WithContext(ctx).Where("Name LIKE ? OR Description LIKE ?", fmt.Sprintf("%%%s%%", req.Query), fmt.Sprintf("%%%s%%", req.Query)).Find(&product)
	if err.Error != nil {
		return nil, err.Error
	}
	var GrpcProduct []*pb.Product
	for _, p := range product {
		GrpcProduct = append(GrpcProduct, ConvertToGrpcProduct(&p))
	}
	return &pb.SearchProductsResp{Results: GrpcProduct}, nil
}
func (P *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.CreateProductReq) (*pb.CreateProductResp, error) {
	DB.InitConfig()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "failed to extract metadata from context")
	}
	log.Println("Metadata:", md)
	//检查用户身份
	userClaims, ok := ctx.Value("userClaims").(*Model.UserClaims)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "failed to extract user claims from context")
	}
	if userClaims.Role != "shop" {
		return nil, status.Errorf(codes.PermissionDenied, "only shop can create product")
	}
	//创建商品
	product := &Model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       float64(req.Price),
		Sort:        req.Sort,
	}
	if res := global.Db.Create(&product); res.Error != nil {
		log.Println("Error:", res.Error)
		return nil, status.Errorf(codes.InvalidArgument, "failed to create product")
	}
	return &pb.CreateProductResp{Id: uint32(product.ID), Success: true}, nil
}

func main() {
	l, _ := net.Listen("tcp", ":5002")
	s := grpc.NewServer(grpc.UnaryInterceptor(I.AuthInterceptor))
	pb.RegisterProductCatalogServiceServer(s, &ProductServiceServer{})
	log.Println("Server is running on port 5002")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
