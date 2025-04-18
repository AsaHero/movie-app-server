package entity

type MovieGenres struct {
	MovieID int64 `gorm:"column:movie_id;primary_key"`
	GenreID int64 `gorm:"column:genre_id;primary_key"`

	Movie *Movies `gorm:"foreignKey:ID;references:MovieID"`
	Genre *Genres `gorm:"foreignKey:ID;references:GenreID"`
}
