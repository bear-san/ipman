package server

import (
	"context"
	"github.com/bear-san/ipman/pkg/grpc"
	"github.com/bear-san/ipman/pkg/ip_repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IPManServer struct {
	grpc.UnimplementedIPAddressServiceServer
	IPRepo *ip_repo.IPRepo
}

func (s *IPManServer) AssignAddress(context.Context, *grpc.AssignAddressRequest) (*grpc.AssignAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignAddress not implemented")
}

func (s *IPManServer) ReleaseAddress(context.Context, *grpc.ReleaseAddressRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseAddress not implemented")
}
