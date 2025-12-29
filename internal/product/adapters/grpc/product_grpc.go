package grpc

import (
	"context"
	product "hex-postgres-grpc/internal/product/domain"
	productpb "hex-postgres-grpc/proto/product"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	productpb.UnimplementedProductServiceServer
	service product.Service
}

func NewProductGRPCServer(service product.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {
	p, err := s.service.CreateProduct(ctx, req.Name, req.Price)
	if err != nil {
		return nil, err
	}
	return &productpb.CreateProductResponse{
		Product: &productpb.ProductMessage{
			Id:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			CreatedAt: timestamppb.New(p.CreatedAt),
		},
	}, nil
}

func (s *Server) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &productpb.GetProductResponse{
		Product: &productpb.ProductMessage{
			Id:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			CreatedAt: timestamppb.New(p.CreatedAt),
		},
	}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error) {
	p, err := s.service.UpdateProduct(ctx, req.Id, req.Name, req.Price)
	if err != nil {
		return nil, err
	}
	return &productpb.UpdateProductResponse{
		Product: &productpb.ProductMessage{
			Id:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			CreatedAt: timestamppb.New(p.CreatedAt),
		},
	}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error) {
	if err := s.service.DeleteProduct(ctx, req.Id); err != nil {
		return nil, err
	}
	return &productpb.DeleteProductResponse{Success: true}, nil
}

func (s *Server) ListProducts(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
	products, err := s.service.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	var pbProducts []*productpb.ProductMessage
	for _, p := range products {
		pbProducts = append(pbProducts, &productpb.ProductMessage{
			Id:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			CreatedAt: timestamppb.New(p.CreatedAt),
		})
	}

	return &productpb.ListProductsResponse{Products: pbProducts}, nil
}
