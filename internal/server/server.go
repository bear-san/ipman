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

func (s *IPManServer) ListAddresses(_ context.Context, _ *emptypb.Empty) (*grpc.ListAddressesResponse, error) {
	addresses, err := s.IPRepo.GetAddresses()
	if err != nil {
		return nil, err
	}

	response := grpc.ListAddressesResponse{
		Addresses: make([]*grpc.IPAddress, 0),
	}
	for _, address := range addresses {
		var addressType grpc.AddressType
		switch address.AddressType {
		case ip_repo.IP_ADDRESS_TYPE_LOCAL:
			addressType = grpc.AddressType_LOCAL
		case ip_repo.IP_ADDRESS_TYPE_GLOBAL:
			addressType = grpc.AddressType_GLOBAL
		default:
			return nil, status.Errorf(codes.Unavailable, "invalid Address Type: %s", address.AddressType)
		}

		response.Addresses = append(response.Addresses, &grpc.IPAddress{
			Address:           address.Address,
			Subnet:            address.Subnet,
			GatewayAddress:    address.GatewayAddress,
			AddressType:       addressType,
			Using:             address.Using,
			AutoAssignEnabled: address.AutoAssignEnabled,
			Description:       address.Description,
		})
	}

	return &response, nil
}

func (s *IPManServer) AssignAddress(_ context.Context, in *grpc.AssignAddressRequest) (*grpc.AssignAddressResponse, error) {
	var addressType string
	switch in.AddressType {
	case grpc.AddressType_LOCAL:
		addressType = ip_repo.IP_ADDRESS_TYPE_LOCAL
	case grpc.AddressType_GLOBAL:
		addressType = ip_repo.IP_ADDRESS_TYPE_GLOBAL
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid Address Type: %s", in.GetAddressType().String())
	}
	assignedIPAddress, err := s.IPRepo.AssignIPAddress(addressType, in.Description)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to assign IP Address: %v", err.Error())
	}
	return &grpc.AssignAddressResponse{
		Address: &grpc.IPAddress{
			Address:           assignedIPAddress.Address,
			Subnet:            assignedIPAddress.Subnet,
			GatewayAddress:    assignedIPAddress.GatewayAddress,
			AddressType:       in.AddressType,
			Using:             assignedIPAddress.Using,
			AutoAssignEnabled: assignedIPAddress.AutoAssignEnabled,
			Description:       assignedIPAddress.Description,
		},
	}, nil
}

func (s *IPManServer) ReleaseAddress(_ context.Context, in *grpc.ReleaseAddressRequest) (*emptypb.Empty, error) {
	var addressType string
	switch in.Address.AddressType {
	case grpc.AddressType_LOCAL:
		addressType = ip_repo.IP_ADDRESS_TYPE_LOCAL
	case grpc.AddressType_GLOBAL:
		addressType = ip_repo.IP_ADDRESS_TYPE_GLOBAL
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid Address Type: %s", in.Address.GetAddressType().String())
	}
	err := s.IPRepo.ReleaseIPAddress(ip_repo.IPAddress{
		Address:           in.Address.Address,
		Subnet:            in.Address.Subnet,
		GatewayAddress:    in.Address.GatewayAddress,
		AddressType:       addressType,
		Using:             in.Address.Using,
		AutoAssignEnabled: in.Address.AutoAssignEnabled,
		Description:       in.Address.Description,
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to release IP Address: %v", err.Error())
	}
	return &emptypb.Empty{}, nil
}
