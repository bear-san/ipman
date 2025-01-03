package main

import (
	"github.com/bear-san/ipman/internal/server"
	"google.golang.org/grpc"
	"net"

	ipmanpb "github.com/bear-san/ipman/pkg/grpc"
)

func main() {
	ipmanServer := server.IPManServer{}
	grpcSrv := grpc.NewServer()

	ipmanpb.RegisterIPAddressServiceServer(grpcSrv, &ipmanServer)

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	if err := grpcSrv.Serve(listener); err != nil {
		panic(err)
	}
}
