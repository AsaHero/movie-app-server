package models

import "time"

type Movie struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Release         string    `json:"release"`
	Plot            *string   `json:"plot"`
	DurationMinutes int16     `json:"duration_minutes"`
	PosterURL       string    `json:"poster_url"`
	TrailerURL      string    `json:"trailer_url"`
	Genres          []string  `json:"genres"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateMovieRequest struct {
	Title           string  `json:"title" validate:"required,min=2,max=255"`
	Release         string  `json:"release" validate:"required"`
	Plot            *string `json:"plot"`
	DurationMinutes int16   `json:"duration_minutes" validate:"required,min=1,max=500"`
	PosterURL       string  `json:"poster_url" validate:"required"`
	TrailerURL      string  `json:"trailer_url" validate:"required"`
	Genres          []int   `json:"genres" validate:"required"`
}

type UpdateMovieRequest struct {
	Title           string  `json:"title" validate:"required,min=2,max=255"`
	Release         string  `json:"release" validate:"required"`
	Plot            *string `json:"plot"`
	DurationMinutes int16   `json:"duration_minutes" validate:"required,min=1,max=500"`
	PosterURL       string  `json:"poster_url" validate:"required"`
	TrailerURL      string  `json:"trailer_url" validate:"required"`
	Genres          []int   `json:"genres" validate:"required"`
}

type GetAllMoviesRequest struct {
	Page     *int    `form:"page" validate:"min=1"`
	Limit    *int    `form:"limit" validate:"min=1,max=100"`
	OrderBy  *string `form:"order_by"`
	OrderDir *string `form:"order_dir"`
	Search   *string `form:"search"`
	Genres   string  `form:"genres"`
}

type GetAllMoviesResponse struct {
	Movies []Movie `json:"movies"`
	Total  int64   `json:"total"`
}

type Gener struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetAllGenresResponse struct {
	Genres []Gener `json:"genres"`
}
