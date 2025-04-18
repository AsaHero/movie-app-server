package movies

import (
	"context"
	"time"

	"github.com/AsaHero/movie-app-server/internal/entity"
	"github.com/AsaHero/movie-app-server/internal/inerr"
	"github.com/AsaHero/movie-app-server/internal/repository/movie_genres"
	"github.com/AsaHero/movie-app-server/internal/repository/movies"
)

type service struct {
	contextTimeout time.Duration
	movieRepo      movies.Repository
	movieGeresRepo movie_genres.Repository
}

func New(contextTimeout time.Duration, movieRepo movies.Repository, movieGenresRepo movie_genres.Repository) Service {
	return &service{
		contextTimeout: contextTimeout,
		movieRepo:      movieRepo,
		movieGeresRepo: movieGenresRepo,
	}
}

func (s *service) Create(ctx context.Context, movie *entity.Movies, movieGenres []*entity.MovieGenres) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	s.beforeCreate(movie)

	err := s.movieRepo.WithTransaction(ctx, func(ctx context.Context) error {
		// First create the movie
		if err := s.movieRepo.Create(ctx, movie); err != nil {
			return err
		}

		for i := range movieGenres {
			movieGenres[i].MovieID = movie.ID
		}

		if len(movieGenres) > 0 {

			if err := s.movieGeresRepo.BatchCreate(ctx, movieGenres); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return inerr.Err(err) // print error and pass
	}

	return nil
}
func (s *service) Update(ctx context.Context, movie *entity.Movies) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if err := s.movieRepo.Update(ctx, movie); err != nil {
		return inerr.Err(err)
	}

	return nil
}

func (s *service) List(ctx context.Context, limit, page uint64, orderBy, orderDir string, filters entity.MovieFilters) (int64, []entity.Movies, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if limit > 100 {
		limit = 100
	}

	if page < 1 {
		page = 1
	}

	total, movies, err := s.movieRepo.ListWithFilters(ctx, limit, page, orderBy, orderDir, filters)
	if err != nil {
		return 0, nil, inerr.Err(err)
	}

	return total, movies, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*entity.Movies, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	movie, err := s.movieRepo.FindOne(ctx, map[string]any{"id": id}, "MovieGenres", "MovieGenres.Genre")
	if err != nil {
		return nil, inerr.Err(err)
	}

	return movie, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err := s.movieRepo.Delete(ctx, map[string]any{"id": id})
	if err != nil {
		return inerr.Err(err)
	}

	return nil
}

func (s *service) beforeCreate(m *entity.Movies) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = time.Now()
	}
}

func (s *service) beforeUpdate(m *entity.Movies) {
	m.UpdatedAt = time.Now()
}
