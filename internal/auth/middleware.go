package auth

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	SubjectContextKey contextKey = "subject"
)

func (s *service) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		sub, err := s.ValidateToken(r.Context(), parts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), SubjectContextKey, sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *service) GRPCUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return handler(ctx, req)
	}

	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, status.Error(codes.Unauthenticated, "invalid auth header")
	}

	sub, err := s.ValidateToken(ctx, parts[1])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	newCtx := context.WithValue(ctx, SubjectContextKey, sub)
	return handler(newCtx, req)
}

// Helper to get subject from context
func SubjectFromContext(ctx context.Context) (Subject, bool) {
	sub, ok := ctx.Value(SubjectContextKey).(Subject)
	return sub, ok
}
