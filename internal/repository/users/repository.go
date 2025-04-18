package users

import (
	"github.com/AsaHero/movie-app-server/internal/entity"
	"github.com/AsaHero/movie-app-server/internal/repository"
	"gorm.io/gorm"
)

// repository implements the Repository interface
type repo struct {
	repository.BaseRepository[*entity.Users]
	db *gorm.DB
}

// New creates a new user repository
func New(db *gorm.DB) Repository {
	return &repo{
		BaseRepository: repository.NewBaseRepository[*entity.Users](db),
		db:             db,
	}
}
