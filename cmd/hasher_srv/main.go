package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Serhii1Epam/enhanceHttpServer/pkg/api"
	"github.com/Serhii1Epam/enhanceHttpServer/pkg/hasher"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type hsr_server struct {
	pb.UnimplementedHasherServer
}

func (s hsr_server) GetHash(ctx context.Context, in *pb.GetHashReq) (out *pb.GetHashResp, err error) {
	log.Printf("GetHash() received: name[%v], pass[%v]", in.GetName(), in.GetPass())
	out.Hash, err = hasher.NewHasher(in.GetPass()).HashPassword()
	return out, err
}

func (s hsr_server) CheckHash(ctx context.Context, in *pb.CheckHashReq) (out *pb.CheckHashResp, err error) {
	log.Printf("CheckHash() received: hash[%v], pass[%v]", in.GetHash(), in.GetPass())
	err = nil
	out.Resp = hasher.NewHasher(in.GetPass()).CheckPasswordHash(in.GetHash())
	if !out.Resp {
		err = errors.New("Invalid Hash.")
	}
	return out, err
}

func main() {
	fmt.Println("Start Hasher service server...")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &server{})
	pb.RegisterHasherServer(s, &hsr_server{})
	log.Printf("Server listening at %v...", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
