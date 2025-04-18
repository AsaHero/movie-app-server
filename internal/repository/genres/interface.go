package genres

import (
	"github.com/AsaHero/movie-app-server/internal/entity"
	"github.com/AsaHero/movie-app-server/internal/repository"
)

type Repository interface {
	repository.BaseRepository[*entity.Genres]
}
