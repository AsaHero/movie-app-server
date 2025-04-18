package movies

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AsaHero/movie-app-server/delivery/api/handlers"
	"github.com/AsaHero/movie-app-server/delivery/api/middlewares"
	"github.com/AsaHero/movie-app-server/delivery/api/models"
	"github.com/AsaHero/movie-app-server/delivery/api/outerr"
	"github.com/AsaHero/movie-app-server/delivery/api/validation"
	"github.com/AsaHero/movie-app-server/internal/entity"
	"github.com/AsaHero/movie-app-server/internal/service/genres"
	"github.com/AsaHero/movie-app-server/internal/service/movies"
	"github.com/AsaHero/movie-app-server/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/shogo82148/pointer"
)

type handler struct {
	config        *config.Config
	validator     *validation.Validator
	moviesService movies.Service
	genresService genres.Service
}

func New(router *gin.RouterGroup, opt *handlers.HandlerOptions) {
	handler := handler{
		config:        opt.Config,
		validator:     opt.Validator,
		moviesService: opt.MoviesSerive,
		genresService: opt.GenresService,
	}

	router.Use(middlewares.BearerAuth(opt.Config.Token.Secret))

	router.POST("/", handler.CreateMovie)
	router.GET("/", handler.GetAllMovies)
	router.GET("/:id", handler.GetMovie)
	router.PUT("/:id", handler.UpdateMovie)
	router.DELETE("/:id", handler.DeleteMovie)

	router.GET("/genres", handler.GetAllGenres)
}

// @Security ApiKeyAuth
// @Summary Create movie
// @Description Create movie
// @Tags Movies
// @Accept json
// @Produce json
// @Param request body models.CreateMovieRequest true "Create movie request"
// @Success 201 {object} models.Movie
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /movies [post]
func (h *handler) CreateMovie(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	releaseDate, err := time.Parse(time.DateOnly, req.Release)
	if err != nil {
		outerr.BadRequest(c, "Invalid release date, format should be YYYY-MM-DD")
	}

	movie := &entity.Movies{
		Title:           req.Title,
		Release:         releaseDate,
		Plot:            req.Plot,
		DurationMinutes: req.DurationMinutes,
		PosterURL:       req.PosterURL,
		TrailerURL:      req.TrailerURL,
	}

	genres := make([]*entity.MovieGenres, 0, len(req.Genres))

	for _, genreID := range req.Genres {
		genres = append(genres, &entity.MovieGenres{
			GenreID: int64(genreID),
		})
	}

	err = h.moviesService.Create(ctx, movie, genres)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, models.Empty{})
}

// @Security ApiKeyAuth
// @Summary Get all movies
// @Description Get all movies
// @Tags Movies
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param genres query []string false "Filter by genres" collectionFormat(csv)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param order_by query string false "Order by field" Enums(title,release,created_at)
// @Param order_dir query string false "Order direction" Enums(asc,desc)
// @Success 200 {object} models.GetAllMoviesResponse
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /movies [get]
func (h *handler) GetAllMovies(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.GetAllMoviesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	var genreIDs []int
	if req.Genres != "" {
		genresList := strings.Split(req.Genres, ",")
		for _, g := range genresList {
			id, err := strconv.Atoi(g)
			if err != nil {
				outerr.BadRequest(c, "Invalid genre ID format")
				return
			}
			genreIDs = append(genreIDs, id)
		}
	}

	total, movies, err := h.moviesService.List(ctx,
		uint64(*req.Limit), uint64(*req.Page),
		pointer.StringValue(req.OrderBy), pointer.StringValue(req.OrderDir),
		entity.MovieFilters{
			Search: req.Search,
			Genres: genreIDs,
		},
	)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	response := models.GetAllMoviesResponse{
		Total:  total,
		Movies: make([]models.Movie, 0, len(movies)),
	}

	for _, movie := range movies {
		mov := models.Movie{
			ID:              movie.ID,
			Title:           movie.Title,
			Release:         movie.Release.Format(time.RFC3339),
			Plot:            movie.Plot,
			DurationMinutes: movie.DurationMinutes,
			PosterURL:       movie.PosterURL,
			TrailerURL:      movie.TrailerURL,
			Genres:          make([]string, 0, len(movie.MovieGenres)),
		}

		for _, genre := range movie.MovieGenres {
			mov.Genres = append(mov.Genres, genre.Genre.Name)
		}

		response.Movies = append(response.Movies, mov)
	}

	c.JSON(http.StatusOK, response)
}

// @Security ApiKeyAuth
// @Summary Get movie by id
// @Description Get movie by id
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie id"
// @Success 200 {object} models.Movie
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /movies/{id} [get]
func (h *handler) GetMovie(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString("user_id")
	if userID == "" {
		outerr.Unauthorized(c, "user_id is required")
		return
	}

	value := c.Param("id")
	if value == "" {
		outerr.BadRequest(c, "id is required")
		return
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		outerr.BadRequest(c, "Invalid id")
		return
	}

	movie, err := h.moviesService.GetByID(ctx, id)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	response := models.Movie{
		ID:              movie.ID,
		Title:           movie.Title,
		Release:         movie.Release.Format(time.RFC3339),
		Plot:            movie.Plot,
		DurationMinutes: movie.DurationMinutes,
		PosterURL:       movie.PosterURL,
		TrailerURL:      movie.TrailerURL,
		Genres:          make([]string, 0, len(movie.MovieGenres)),
	}

	for _, genre := range movie.MovieGenres {
		response.Genres = append(response.Genres, genre.Genre.Name)
	}

	c.JSON(http.StatusOK, response)
}

// @Security ApiKeyAuth
// @Summary Update movie
// @Description Update movie
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie id"
// @Param request body models.UpdateMovieRequest true "Update movie request"
// @Success 200 {object} models.Empty
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /movies/{id} [put]
func (h *handler) UpdateMovie(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString("user_id")
	if userID == "" {
		outerr.Unauthorized(c, "user_id is required")
		return
	}

	value := c.Param("id")
	if value == "" {
		outerr.BadRequest(c, "id is required")
		return
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		outerr.BadRequest(c, "Invalid id")
		return
	}

	var req models.UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	releaseDate, err := time.Parse(time.RFC3339, req.Release)
	if err != nil {
		outerr.BadRequest(c, "Invalid release date")
		return
	}

	movie := &entity.Movies{
		ID:              id,
		Title:           req.Title,
		Release:         releaseDate,
		Plot:            req.Plot,
		DurationMinutes: req.DurationMinutes,
		PosterURL:       req.PosterURL,
		TrailerURL:      req.TrailerURL,
	}

	for _, genreID := range req.Genres {
		movie.MovieGenres = append(movie.MovieGenres, entity.MovieGenres{
			MovieID: movie.ID,
			GenreID: int64(genreID),
		})
	}

	if err := h.moviesService.Update(ctx, movie); err != nil {
		outerr.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Empty{})
}

// @Security ApiKeyAuth
// @Summary Delete movie
// @Description Delete movie
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie id"
// @Success 200 {object} models.Empty
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /movies/{id} [delete]
func (h *handler) DeleteMovie(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString("user_id")
	if userID == "" {
		outerr.Unauthorized(c, "user_id is required")
		return
	}

	value := c.Param("id")
	if value == "" {
		outerr.BadRequest(c, "id is required")
		return
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		outerr.BadRequest(c, "Invalid id")
		return
	}

	if err := h.moviesService.Delete(ctx, id); err != nil {
		outerr.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Empty{})
}

// @Security ApiKeyAuth
// @Summary Get all genres
// @Description Get all genres
// @Tags Genres
// @Accept json
// @Produce json
// @Success 200 {object} models.GetAllGenresResponse
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /movies/genres [get]
func (h *handler) GetAllGenres(r *gin.Context) {
	ctx := r.Request.Context()

	genres, err := h.genresService.GetAll(ctx)
	if err != nil {
		outerr.HandleError(r, err)
		return
	}

	response := models.GetAllGenresResponse{
		Genres: make([]models.Gener, 0, len(genres)),
	}

	for _, genre := range genres {
		response.Genres = append(response.Genres, models.Gener{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	r.JSON(http.StatusOK, response)
}
