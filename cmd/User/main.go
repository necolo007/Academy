package User

import (
	pb "Academy/pb/user"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	l, _ := net.Listen("tcp", "8080")
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server.UserServiceServer{})
	err := s.Serve(l)
	if err != nil {
		log.Println(err)
	}
}
