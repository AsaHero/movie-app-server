package genres

import (
	"context"
	"time"

	"github.com/AsaHero/movie-app-server/internal/entity"
	"github.com/AsaHero/movie-app-server/internal/inerr"
	"github.com/AsaHero/movie-app-server/internal/repository/genres"
)

type service struct {
	contextTimeout time.Duration
	genresRepo     genres.Repository
}

func New(contextTimeout time.Duration, genresRepo genres.Repository) Service {
	return &service{
		contextTimeout: contextTimeout,
		genresRepo:     genresRepo,
	}
}

func (s *service) GetAll(ctx context.Context) ([]*entity.Genres, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	_, genres, err := s.genresRepo.FindAll(ctx, 0, 0, "", map[string]any{})
	if err != nil {
		return nil, inerr.Err(err)
	}

	return genres, err
}
