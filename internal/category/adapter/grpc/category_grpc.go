package grpc

import (
	"context"
	"hex-postgres-grpc/internal/auth"
	"hex-postgres-grpc/internal/category/domain"
	categorypb "hex-postgres-grpc/proto/category"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	categorypb.UnimplementedCategoryServiceServer
	service domain.Service
}

func NewCategoryGRPCServer(service domain.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) CreateCategory(ctx context.Context, req *categorypb.CreateCategoryRequest) (*categorypb.CreateCategoryResponse, error) {
	// For gRPC, we might need a generic user for now if interceptor isn't setting subject
	sub, _ := auth.SubjectFromContext(ctx)
	userID := sub.ID
	if userID == "" {
		userID = "grpc-system" // fallback
	}

	cat, err := s.service.CreateCategory(ctx, req.Name, userID)
	if err != nil {
		return nil, err
	}
	return &categorypb.CreateCategoryResponse{
		Category: &categorypb.CategoryMessage{
			Id:        cat.ID,
			Name:      cat.Name,
			CreatedAt: timestamppb.New(cat.CreatedAt),
			UpdatedAt: func() *timestamppb.Timestamp {
				if cat.UpdatedAt != nil {
					return timestamppb.New(*cat.UpdatedAt)
				}
				return nil
			}(),
		},
	}, nil
}

func (s *Server) GetCategory(ctx context.Context, req *categorypb.GetCategoryRequest) (*categorypb.GetCategoryResponse, error) {
	cat, err := s.service.GetCategory(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &categorypb.GetCategoryResponse{
		Category: &categorypb.CategoryMessage{
			Id:        cat.ID,
			Name:      cat.Name,
			CreatedAt: timestamppb.New(cat.CreatedAt),
			UpdatedAt: func() *timestamppb.Timestamp {
				if cat.UpdatedAt != nil {
					return timestamppb.New(*cat.UpdatedAt)
				}
				return nil
			}(),
		},
	}, nil
}

func (s *Server) UpdateCategory(ctx context.Context, req *categorypb.UpdateCategoryRequest) (*categorypb.UpdateCategoryResponse, error) {
	sub, _ := auth.SubjectFromContext(ctx)
	userID := sub.ID
	if userID == "" {
		userID = "grpc-system"
	}

	cat, err := s.service.UpdateCategory(ctx, req.Id, req.Name, userID)
	if err != nil {
		return nil, err
	}
	return &categorypb.UpdateCategoryResponse{
		Category: &categorypb.CategoryMessage{
			Id:        cat.ID,
			Name:      cat.Name,
			CreatedAt: timestamppb.New(cat.CreatedAt),
			UpdatedAt: func() *timestamppb.Timestamp {
				if cat.UpdatedAt != nil {
					return timestamppb.New(*cat.UpdatedAt)
				}
				return nil
			}(),
		},
	}, nil
}

func (s *Server) DeleteCategory(ctx context.Context, req *categorypb.DeleteCategoryRequest) (*categorypb.DeleteCategoryResponse, error) {
	sub, _ := auth.SubjectFromContext(ctx)
	userID := sub.ID
	if userID == "" {
		userID = "grpc-system"
	}

	if err := s.service.DeleteCategory(ctx, req.Id, userID); err != nil {
		return nil, err
	}
	return &categorypb.DeleteCategoryResponse{Success: true}, nil
}

func (s *Server) ListCategories(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error) {
	cats, err := s.service.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	var pbCats []*categorypb.CategoryMessage
	for _, cat := range cats {
		pbCats = append(pbCats, &categorypb.CategoryMessage{
			Id:        cat.ID,
			Name:      cat.Name,
			CreatedAt: timestamppb.New(cat.CreatedAt),
			UpdatedAt: func() *timestamppb.Timestamp {
				if cat.UpdatedAt != nil {
					return timestamppb.New(*cat.UpdatedAt)
				}
				return nil
			}(),
		})
	}

	return &categorypb.ListCategoriesResponse{Categories: pbCats}, nil
}
