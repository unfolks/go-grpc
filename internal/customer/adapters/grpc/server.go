package grpc

import (
	"context"
	"hex-postgres-grpc/internal/customer/usecase"
	customerpb "hex-postgres-grpc/proto/customer"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	customerpb.UnimplementedCustomerServiceServer
	createCustomer *usecase.CreateCustomer
}

func NewServer(createCustomer *usecase.CreateCustomer) *Server {
	return &Server{
		createCustomer: createCustomer,
	}
}

func (s *Server) CreateCustomer(ctx context.Context, req *customerpb.CreateCustomerRequest) (*customerpb.CreateCustomerResponse, error) {
	c, err := s.createCustomer.Execute(ctx, req.Name, req.Email, req.Address)
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

//func (s *server) GetCustomer(ctx context.Context)
