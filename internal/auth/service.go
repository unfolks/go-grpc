package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidToken = errors.New("invalid token")
)

type service struct {
	secretKey []byte
	policies  []Policy
	repo      UserRepository
}

func NewService(secretKey string, repo UserRepository) Service {
	s := &service{
		secretKey: []byte(secretKey),
		repo:      repo,
	}
	s.initPolicies()
	return s
}

func (s *service) initPolicies() {
	s.policies = []Policy{
		{
			SubjectRole:  "admin",
			Action:       ActionCreate,
			ResourceType: "*",
			Condition: func(sub Subject, res Resource) bool {
				return true
			},
		},
		{
			SubjectRole:  "admin",
			Action:       ActionRead,
			ResourceType: "*",
			Condition: func(sub Subject, res Resource) bool {
				return true
			},
		},
		{
			SubjectRole:  "admin",
			Action:       ActionUpdate,
			ResourceType: "*",
			Condition: func(sub Subject, res Resource) bool {
				return true
			},
		},
		{
			SubjectRole:  "admin",
			Action:       ActionDelete,
			ResourceType: "*",
			Condition: func(sub Subject, res Resource) bool {
				return true
			},
		},
		// User management policies
		{
			SubjectRole:  "user",
			Action:       ActionRead,
			ResourceType: "user",
			Condition: func(sub Subject, res Resource) bool {
				return true
			},
		},
		// Example ABAC policy: Owner can update their own resource
		{
			SubjectRole:  "user",
			Action:       ActionUpdate,
			ResourceType: "customer",
			Condition: func(sub Subject, res Resource) bool {
				ownerID, ok := res.Attributes["owner_id"].(string)
				return ok && ownerID == sub.ID
			},
		},
		// Common policy: Users can read everything
		{
			SubjectRole:  "user",
			Action:       ActionRead,
			ResourceType: "*",
			Condition: func(sub Subject, res Resource) bool {
				return true
			},
		},
	}
}

func (s *service) Authorize(ctx context.Context, sub Subject, act Action, res Resource) (bool, error) {
	for _, p := range s.policies {
		if p.SubjectRole == sub.Role && (p.ResourceType == "*" || p.ResourceType == res.Type) && (p.Action == act) {
			if p.Condition(sub, res) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (s *service) GenerateToken(ctx context.Context, sub Subject) (string, error) {
	claims := jwt.MapClaims{
		"sub":  sub.ID,
		"role": sub.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"attr": sub.Attributes,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *service) ValidateToken(ctx context.Context, tokenStr string) (Subject, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return Subject{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		subID, _ := claims["sub"].(string)
		role, _ := claims["role"].(string)
		attrs, _ := claims["attr"].(map[string]interface{})

		return Subject{
			ID:         subID,
			Role:       role,
			Attributes: attrs,
		}, nil
	}

	return Subject{}, ErrInvalidToken
}

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", ErrUnauthorized
	}

	// In a real app, use bcrypt.CompareHashAndPassword here.
	// For this demonstration, we'll assume the password matches if it's the same as the hash (as seeded).
	if u.PasswordHash != password {
		return "", ErrUnauthorized
	}

	sub := Subject{
		ID:         u.ID,
		Username:   u.Username,
		Role:       u.Role,
		Attributes: u.Attributes,
	}

	return s.GenerateToken(ctx, sub)
}

func (s *service) CreateUser(ctx context.Context, sub Subject, user *User) error {
	authorized, err := s.Authorize(ctx, sub, ActionCreate, Resource{Type: "user"})
	if err != nil {
		return err
	}
	if !authorized {
		return ErrUnauthorized
	}

	return s.repo.Create(ctx, user)
}

func (s *service) UpdateUser(ctx context.Context, sub Subject, user *User) error {
	authorized, err := s.Authorize(ctx, sub, ActionUpdate, Resource{Type: "user"})
	if err != nil {
		return err
	}
	if !authorized {
		return ErrUnauthorized
	}

	return s.repo.Update(ctx, user)
}

func (s *service) GetUser(ctx context.Context, sub Subject, id string) (*User, error) {
	authorized, err := s.Authorize(ctx, sub, ActionRead, Resource{Type: "user"})
	if err != nil {
		return nil, err
	}
	if !authorized {
		return nil, ErrUnauthorized
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) ListUsers(ctx context.Context, sub Subject) ([]*User, error) {
	authorized, err := s.Authorize(ctx, sub, ActionRead, Resource{Type: "user"})
	if err != nil {
		return nil, err
	}
	if !authorized {
		return nil, ErrUnauthorized
	}

	return s.repo.List(ctx)
}
