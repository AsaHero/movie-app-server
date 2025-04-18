package users

import (
	"context"
	"time"

	"github.com/AsaHero/movie-app-server/internal/entity"
	"github.com/AsaHero/movie-app-server/internal/repository/users"
	"github.com/google/uuid"
)

type service struct {
	contextTimeout time.Duration
	userRepo       users.Repository
}

func New(contextTimeout time.Duration, userRepo users.Repository) Service {
	return &service{
		contextTimeout: contextTimeout,
		userRepo:       userRepo,
	}
}

func (s *service) Create(ctx context.Context, user *entity.Users) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	s.beforeCreate(user)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id string) (*entity.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	user, err := s.userRepo.FindOne(ctx, map[string]any{"id": id})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service) beforeCreate(user *entity.Users) {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}
}

func (service) beforeUpdate(user *entity.Users) {
	user.UpdatedAt = time.Now()
}
