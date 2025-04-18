package entity

import "time"

type Movies struct {
	ID              int64 `gorm:"primary_key"`
	Title           string
	Release         time.Time
	Plot            *string
	DurationMinutes int16
	PosterURL       string
	TrailerURL      string
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Relations
	MovieGenres []MovieGenres `gorm:"foreignKey:MovieID"`
	Genres      []Genres      `gorm:"many2many:movie_genres;joinForeignKey:MovieID;joinReferences:GenreID"`
}
