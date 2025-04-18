package genres

import (
	"context"

	"github.com/AsaHero/movie-app-server/internal/entity"
)

type Service interface {
	GetAll(ctx context.Context) ([]*entity.Genres, error)
}
