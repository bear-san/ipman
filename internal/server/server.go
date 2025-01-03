package server

import (
	"context"
	"github.com/bear-san/ipman/pkg/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IPManServer struct {
	grpc.UnimplementedIPAddressServiceServer
}

func (s *IPManServer) AssignAddress(context.Context, *grpc.AssignAddressRequest) (*grpc.AssignAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignAddress not implemented")
}
