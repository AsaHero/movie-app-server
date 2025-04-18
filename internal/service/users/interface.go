package users

import (
	"context"

	"github.com/AsaHero/movie-app-server/internal/entity"
)

type Service interface {
	Create(ctx context.Context, user *entity.Users) error
	GetByID(ctx context.Context, id string) (*entity.Users, error)
}
