package users

import (
	"context"

	"github.com/AsaHero/movie-app-server/internal/entity"
	"github.com/AsaHero/movie-app-server/internal/repository"
	"github.com/AsaHero/movie-app-server/pkg/database/postgres"
)

type Repository interface {
	repository.BaseRepository[*entity.Users]
	FindByLogin(ctx context.Context, login string) (*entity.Users, error)
}

func (r *repo) FindByLogin(ctx context.Context, login string) (*entity.Users, error) {
	db := repository.FromContext(ctx, r.db)
	var user *entity.Users

	if err := db.Where("username = ?", login).Or("email = ?", login).First(&user).Error; err != nil {
		return nil, postgres.Error(err, "FindByLogin", &entity.Users{})
	}

	return user, nil
}
