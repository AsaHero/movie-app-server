package movies

import (
	"context"

	"github.com/AsaHero/movie-app-server/internal/entity"
	"github.com/AsaHero/movie-app-server/internal/repository"
	"gorm.io/gorm"
)

type repo struct {
	repository.BaseRepository[*entity.Movies]
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repo{
		BaseRepository: repository.NewBaseRepository[*entity.Movies](db),
		db:             db,
	}
}

func (r *repo) ListWithFilters(ctx context.Context, limit, page uint64, orderBy, orderDir string, filters entity.MovieFilters) (int64, []entity.Movies, error) {
	db := repository.FromContext(ctx, r.db)

	var movies []entity.Movies
	var total int64

	query := db.Model(&entity.Movies{})

	// Apply search filter
	if filters.Search != nil && *filters.Search != "" {
		searchTerm := "%" + *filters.Search + "%"
		query = query.Where("title ILIKE ?", searchTerm)
	}

	// Apply genres filter
	if len(filters.Genres) > 0 {
		query = query.Joins("JOIN movie_genres ON movies.id = movie_genres.movie_id").
			Where("movie_genres.genre_id IN ?", filters.Genres).
			Group("movies.id")
	}

	// Preload related data
	query = query.Preload("MovieGenres").Preload("MovieGenres.Genre")

	// Get total count for pagination
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	// Apply ordering
	if orderBy != "" {
		if orderDir == "" {
			orderDir = "asc"
		}
		query = query.Order(orderBy + " " + orderDir)
	} else {
		// Default ordering
		query = query.Order("created_at desc")
	}

	// Apply pagination
	offset := (page - 1) * limit
	query = query.Offset(int(offset)).Limit(int(limit))

	// Execute the query
	if err := query.Find(&movies).Error; err != nil {
		return 0, nil, err
	}

	return total, movies, nil
}

func (r *repo) Update(ctx context.Context, movie *entity.Movies) error {
	db := repository.FromContext(ctx, r.db)

	// Update the movie itself (without associations)
	if err := db.Model(movie).Omit("MovieGenres").Updates(map[string]interface{}{
		"title":            movie.Title,
		"release":          movie.Release,
		"plot":             movie.Plot,
		"duration_minutes": movie.DurationMinutes,
		"poster_url":       movie.PosterURL,
		"trailer_url":      movie.TrailerURL,
	}).Error; err != nil {
		return err
	}

	// Replace associations
	if err := db.Model(movie).Association("MovieGenres").Replace(movie.MovieGenres); err != nil {
		return err
	}

	return nil
}
