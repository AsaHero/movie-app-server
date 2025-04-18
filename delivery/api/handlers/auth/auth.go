package auth

import (
	"net/http"

	"github.com/AsaHero/movie-app-server/delivery/api/handlers"
	"github.com/AsaHero/movie-app-server/delivery/api/models"
	"github.com/AsaHero/movie-app-server/delivery/api/outerr"
	"github.com/AsaHero/movie-app-server/delivery/api/validation"
	"github.com/AsaHero/movie-app-server/internal/service/auth"
	"github.com/AsaHero/movie-app-server/internal/service/users"
	"github.com/AsaHero/movie-app-server/pkg/config"
	"github.com/AsaHero/movie-app-server/pkg/security"
	"github.com/gin-gonic/gin"
)

type handler struct {
	config       *config.Config
	validator    *validation.Validator
	authService  auth.Service
	usersService users.Service
}

func New(router *gin.RouterGroup, opt *handlers.HandlerOptions) {
	handler := handler{
		config:       opt.Config,
		validator:    opt.Validator,
		authService:  opt.AuthService,
		usersService: opt.UsersService,
	}

	router.POST("/login", handler.Login)
	router.POST("/register", handler.Register)
	router.POST("/refresh", handler.Refresh)
}

// Login godoc

// @Summary      Login
// @Description  Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /auth/login [post]
func (h *handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	user, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	accessToken, refreshToken, err := security.GenerateTokenPair(user.ID, h.config.Token.Secret)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshToken godoc

// @Summary      Refresh token
// @Description  Refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /auth/refresh [post]
func (h *handler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	tokenClaims, err := security.ParseRefreshToken(req.RefreshToken, h.config.Token.Secret)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	// get the user from the token claims
	user, err := h.usersService.GetByID(ctx, tokenClaims.UserID)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	accessToken, refreshToken, err := security.GenerateTokenPair(user.ID, h.config.Token.Secret)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Register godoc

// @Summary      Register
// @Description  Register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param request body models.RegisterRequest true "Register request"
// @Success 201 {object} models.LoginResponse
// @Failure 400 {object} outerr.ErrorResponse
// @Failure 500 {object} outerr.ErrorResponse
// @Router /auth/register [post]
func (h *handler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		outerr.BadRequest(c, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		outerr.HandleError(c, err)
		return
	}

	user, err := h.authService.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	accessToken, refreshToken, err := security.GenerateTokenPair(user.ID, h.config.Token.Secret)
	if err != nil {
		outerr.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
