package movies

import (
	"context"

	"github.com/AsaHero/movie-app-server/internal/entity"
)

type Service interface {
	Create(ctx context.Context, movie *entity.Movies, movieGenres []*entity.MovieGenres) error
	Update(ctx context.Context, movie *entity.Movies) error
	List(ctx context.Context, limit, page uint64, orderBy, orderDir string, filters entity.MovieFilters) (int64, []entity.Movies, error)
	GetByID(ctx context.Context, id int64) (*entity.Movies, error)
	Delete(ctx context.Context, id int64) error
}
