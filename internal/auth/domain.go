package auth

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type Subject struct {
	ID         string
	Username   string
	Role       string
	Attributes map[string]interface{}
}

type User struct {
	ID           string
	Username     string
	PasswordHash string
	Role         string
	Attributes   map[string]interface{}
}

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	List(ctx context.Context) ([]*User, error)
}

type Resource struct {
	Type       string
	ID         string
	Attributes map[string]interface{}
}

type Policy struct {
	SubjectRole  string
	Action       Action
	ResourceType string
	// Condition defines a function that must be true for the policy to apply
	// This is where the "Attribute" part of ABAC comes in.
	Condition func(sub Subject, res Resource) bool
}

type Service interface {
	Authorize(ctx context.Context, sub Subject, act Action, res Resource) (bool, error)
	GenerateToken(ctx context.Context, sub Subject) (string, error)
	ValidateToken(ctx context.Context, token string) (Subject, error)
	Login(ctx context.Context, username, password string) (string, error)
	HTTPMiddleware(next http.Handler) http.Handler
	GRPCUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)

	// User management
	CreateUser(ctx context.Context, sub Subject, user *User) error
	UpdateUser(ctx context.Context, sub Subject, user *User) error
	GetUser(ctx context.Context, sub Subject, id string) (*User, error)
	ListUsers(ctx context.Context, sub Subject) ([]*User, error)
}
