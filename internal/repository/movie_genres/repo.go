package movie_genres

import (
	"github.com/AsaHero/movie-app-server/internal/entity"
	"github.com/AsaHero/movie-app-server/internal/repository"
	"gorm.io/gorm"
)

type repo struct {
	repository.BaseRepository[*entity.MovieGenres]
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repo{
		BaseRepository: repository.NewBaseRepository[*entity.MovieGenres](db),
		db:             db,
	}
}
