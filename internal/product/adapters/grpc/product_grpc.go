package grpc

import (
	"context"
	"hex-postgres-grpc/internal/product/usecase"
	productpb "hex-postgres-grpc/proto/product"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	productpb.UnimplementedProductServiceServer
	createProduct *usecase.CreateProduct
	getProduct    *usecase.GetProduct
	listProducts  *usecase.ListProducts
	updateProduct *usecase.UpdateProduct
	deleteProduct *usecase.DeleteProduct
}

func NewProductGRPCServer(
	createProduct *usecase.CreateProduct,
	getProduct *usecase.GetProduct,
	listProducts *usecase.ListProducts,
	updateProduct *usecase.UpdateProduct,
	deleteProduct *usecase.DeleteProduct,
) *Server {
	return &Server{
		createProduct: createProduct,
		getProduct:    getProduct,
		listProducts:  listProducts,
		updateProduct: updateProduct,
		deleteProduct: deleteProduct,
	}
}

func (s *Server) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {
	p, err := s.createProduct.Execute(ctx, req.Name, req.Price)
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
	p, err := s.getProduct.Execute(ctx, req.Id)
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
	p, err := s.updateProduct.Execute(ctx, req.Id, req.Name, req.Price)
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
	if err := s.deleteProduct.Execute(ctx, req.Id); err != nil {
		return nil, err
	}
	return &productpb.DeleteProductResponse{Success: true}, nil
}

func (s *Server) ListProducts(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
	products, err := s.listProducts.Execute(ctx)
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
