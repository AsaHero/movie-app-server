package entity

type Genres struct {
	ID   int64 `gorm:"primary_key"`
	Name string

	MovieGenres []MovieGenres `gorm:"foreignKey:GenreID"`
}
