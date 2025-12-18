package grpc

import (
	"context"
	order "hex-postgres-grpc/internal/order/domain"
	orderpb "hex-postgres-grpc/proto/order"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	orderpb.UnimplementedORderServiceServer
	svc order.Service
}

func NewOrderGRPCServer(svc order.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	o, err := s.svc.CreateOrder(ctx, req.Amount)
	if err != nil {

		return nil, err
	}
	return &orderpb.CreateOrderResponse{
		Order: &orderpb.OrderMessage{
			Id:        o.ID,
			Amount:    o.Amount,
			CreatedAt: timestamppb.New(o.CreatedAt),
		},
	}, nil
}

func (s *Server) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	o, err := s.svc.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &orderpb.GetOrderResponse{
		Order: &orderpb.OrderMessage{
			Id:        o.ID,
			Amount:    o.Amount,
			CreatedAt: timestamppb.New(o.CreatedAt),
		},
	}, nil
}

func (s *Server) UpdateOrder(ctx context.Context, req *orderpb.UpdateOrderRequest) (*orderpb.UpdateOrderResponse, error) {
	o, err := s.svc.UpdateOrder(ctx, req.Id, req.Amount)
	if err != nil {
		return nil, err
	}
	return &orderpb.UpdateOrderResponse{
		Order: &orderpb.OrderMessage{
			Id:        o.ID,
			Amount:    o.Amount,
			CreatedAt: timestamppb.New(o.CreatedAt),
		},
	}, nil
}

func (s *Server) DeleteOrder(ctx context.Context, req *orderpb.DeleteOrderRequest) (*orderpb.DeleteOrderResponse, error) {
	if err := s.svc.DeleteOrder(ctx, req.Id); err != nil {
		return nil, err
	}
	return &orderpb.DeleteOrderResponse{Success: true}, nil
}

func (s *Server) ListOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	orders, err := s.svc.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	var pbOrders []*orderpb.OrderMessage
	for _, o := range orders {
		pbOrders = append(pbOrders, &orderpb.OrderMessage{
			Id:        o.ID,
			Amount:    o.Amount,
			CreatedAt: timestamppb.New(o.CreatedAt),
		})
	}

	return &orderpb.ListOrdersResponse{Orders: pbOrders}, nil
}
