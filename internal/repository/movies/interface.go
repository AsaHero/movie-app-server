package movies

import (
	"context"

	"github.com/AsaHero/movie-app-server/internal/entity"
	"github.com/AsaHero/movie-app-server/internal/repository"
)

type Repository interface {
	repository.BaseRepository[*entity.Movies]
	ListWithFilters(ctx context.Context, limit, page uint64, orderBy, orderDir string, filters entity.MovieFilters) (int64, []entity.Movies, error)
}
