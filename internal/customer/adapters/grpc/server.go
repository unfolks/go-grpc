package grpc

import (
	"context"
	"hex-postgres-grpc/internal/auth"
	"hex-postgres-grpc/internal/customer/domain"
	customerpb "hex-postgres-grpc/proto/customer"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	customerpb.UnimplementedCustomerServiceServer
	service domain.Service
	auth    auth.Service
}

func NewServer(service domain.Service, auth auth.Service) *Server {
	return &Server{
		service: service,
		auth:    auth,
	}
}

func (s *Server) CreateCustomer(ctx context.Context, req *customerpb.CreateCustomerRequest) (*customerpb.CreateCustomerResponse, error) {
	sub, ok := auth.SubjectFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	authorized, err := s.auth.Authorize(ctx, sub, auth.ActionCreate, auth.Resource{Type: "customer"})
	if err != nil || !authorized {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

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

func (s *Server) GetCustomer(ctx context.Context, req *customerpb.GetCustomerRequest) (*customerpb.GetCustomerResponse, error) {
	sub, ok := auth.SubjectFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	c, err := s.service.GetCustomer(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get customer: %v", err)
	}

	if c == nil {
		return nil, status.Error(codes.NotFound, "customer not found")
	}

	// ABAC check: Owner or Admin
	authorized, err := s.auth.Authorize(ctx, sub, auth.ActionRead, auth.Resource{
		Type: "customer",
		ID:   c.ID,
		Attributes: map[string]interface{}{
			"owner_id": c.ID, // Assuming customer ID is owner ID
		},
	})
	if err != nil || !authorized {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	return &customerpb.GetCustomerResponse{
		Customer: &customerpb.CustomerMessage{
			Id:        c.ID,
			Name:      c.Name,
			Email:     c.Email,
			Address:   c.Address,
			CreatedAt: timestamppb.New(c.CreatedAt),
		},
	}, nil
}

func (s *Server) ListCustomers(ctx context.Context, req *customerpb.ListCustomersRequest) (*customerpb.ListCustomersResponse, error) {
	sub, ok := auth.SubjectFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	// Basic check: Admin or authenticated user can list
	authorized, err := s.auth.Authorize(ctx, sub, auth.ActionRead, auth.Resource{Type: "customer"})
	if err != nil || !authorized {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	customers, err := s.service.ListCustomers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list customers: %v", err)
	}

	var msgs []*customerpb.CustomerMessage
	for _, c := range customers {
		msgs = append(msgs, &customerpb.CustomerMessage{
			Id:        c.ID,
			Name:      c.Name,
			Email:     c.Email,
			Address:   c.Address,
			CreatedAt: timestamppb.New(c.CreatedAt),
		})
	}

	return &customerpb.ListCustomersResponse{
		Customers: msgs,
	}, nil
}
