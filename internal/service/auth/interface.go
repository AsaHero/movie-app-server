package auth

import (
	"context"

	"github.com/AsaHero/movie-app-server/internal/entity"
)

type Service interface {
	LoginByUsername(ctx context.Context, username, password string) (*entity.Users, error)
	Login(ctx context.Context, login, password string) (*entity.Users, error)
	Register(ctx context.Context, name, email, password string) (*entity.Users, error)
}
