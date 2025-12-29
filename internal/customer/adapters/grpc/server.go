package grpc

import (
	"context"
	"hex-postgres-grpc/internal/customer/domain"
	customerpb "hex-postgres-grpc/proto/customer"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	customerpb.UnimplementedCustomerServiceServer
	service domain.Service
}

func NewServer(service domain.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) CreateCustomer(ctx context.Context, req *customerpb.CreateCustomerRequest) (*customerpb.CreateCustomerResponse, error) {
	c, err := s.service.CreateCustomer(ctx, req.Name, req.Email, req.Address)
	if err != nil {
		return nil, err
	}

	return &customerpb.CreateCustomerResponse{
		Customer: &customerpb.CustomerMessage{
			Id:        c.ID,
			Name:      c.Name,
			Email:     c.Email,
			Address:   c.Address,
			CreatedAt: timestamppb.New(c.CreatedAt),
		},
	}, nil
}
