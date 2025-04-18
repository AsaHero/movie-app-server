package handlers

import (
	"github.com/AsaHero/movie-app-server/delivery/api/validation"
	"github.com/AsaHero/movie-app-server/internal/service/auth"
	"github.com/AsaHero/movie-app-server/internal/service/genres"
	"github.com/AsaHero/movie-app-server/internal/service/movies"
	"github.com/AsaHero/movie-app-server/internal/service/users"
	"github.com/AsaHero/movie-app-server/pkg/config"
)

type HandlerOptions struct {
	Config        *config.Config
	Validator     *validation.Validator
	AuthService   auth.Service
	UsersService  users.Service
	MoviesSerive  movies.Service
	GenresService genres.Service
}
