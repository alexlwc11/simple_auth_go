package apis

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/alexlwc11/simple_auth_go/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler interface {
	Register(*gin.Context)
	SignIn(*gin.Context)
	Refresh(*gin.Context)
}

type AuthHandlerImpl struct {
	UserRepository         repositories.UserRepository
	SessionTokenRepository repositories.SessionTokenRepository
	RefreshTokenRepository repositories.RefreshTokenRepository
}

func NewAuthHandlerImpl(
	userRepo repositories.UserRepository,
	sessionTokenRepo repositories.SessionTokenRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
) AuthHandler {
	return &AuthHandlerImpl{
		UserRepository:         userRepo,
		SessionTokenRepository: sessionTokenRepo,
		RefreshTokenRepository: refreshTokenRepo,
	}
}

// InDto for [Register] and [SignIn]
// TODO support other sign up methods e.g. email & password
type userInfoInDto struct {
	DeviceUUID string `json:"device_uuid"`
}

// OutDto for [Register] and [SignIn]
type sessionTokenOutDto struct {
	SessionToken     string    `json:"session_token"`
	SessionExpiredAt time.Time `json:"session_expired_at"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiredAt time.Time `json:"refresh_expired_at"`
}

// Register godoc
//
//	@Summary		Register
//	@Description	New user registration
//	@Tags			Auth
//	@Accept			json
//	@Param			user_info	body	userInfoInDto	true	"User info for registration"
//	@Produce		json
//	@Success		200	{object}	sessionTokenOutDto
//	@Failure		500
//	@Router			/register [post]
func (ahi *AuthHandlerImpl) Register(c *gin.Context) {
	// Create new user with device UUID
	var userCred userInfoInDto
	if err := c.BindJSON(&userCred); err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := uuid.Validate(userCred.DeviceUUID); err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Create a new user with the provided user info
	user, err := ahi.UserRepository.CreateWithDeviceUUID(userCred.DeviceUUID)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	outDto, err := ahi.generateTokenOutDtoWithUserId(user.ID)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, outDto)
}

// SignIn godoc
//
//	@Summary		Sign in
//	@Description	Existing user sign in
//	@Tags			Auth
//	@Accept			json
//	@Param			user_info	body	userInfoInDto	true	"User info for signing in"
//	@Produce		json
//	@Success		200	{object}	sessionTokenOutDto
//	@Failure		500
//	@Router			/sign_in [post]
func (ahi *AuthHandlerImpl) SignIn(c *gin.Context) {
	var userCred userInfoInDto
	if err := c.BindJSON(&userCred); err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := uuid.Validate(userCred.DeviceUUID); err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Find the existing user with the provided user info
	user, err := ahi.UserRepository.FindByDeviceUUID(userCred.DeviceUUID)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	outDto, err := ahi.generateTokenOutDtoWithUserId(user.ID)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, outDto)
}

type refreshInDto struct {
	RefreshToken string `json:"refresh_token"`
}

// Refresh godoc
//
//	@Summary		Refresh
//	@Description	Get new set of tokens with refresh token
//	@Tags			Auth
//	@Accept			json
//	@Param			refresh_token	body	refreshInDto	true	"Refresh token"
//	@Produce		json
//	@Success		200	{object}	sessionTokenOutDto
//	@Failure		500
//	@Router			/refresh [post]
func (ahi *AuthHandlerImpl) Refresh(c *gin.Context) {
	var refreshToken refreshInDto
	if err := c.BindJSON(&refreshToken); err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := ahi.RefreshTokenRepository.FindByValue(refreshToken.RefreshToken)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if token.ExpiredAt.Before(time.Now()) {
		log.Println(errors.New("token expired"))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	outDto, err := ahi.generateTokenOutDtoWithUserId(token.UserID)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, outDto)
}

func (ahi *AuthHandlerImpl) generateTokenOutDtoWithUserId(userID uint) (sessionTokenOutDto, error) {
	// Create session token with user ID
	sessionToken, err := ahi.SessionTokenRepository.CreateWithUserId(userID)
	if err != nil {
		return sessionTokenOutDto{}, err
	}

	// Create refresh token with user ID
	refreshToken, err := ahi.RefreshTokenRepository.CreateWithUserId(userID)
	if err != nil {
		return sessionTokenOutDto{}, err
	}

	return sessionTokenOutDto{
		SessionToken:     sessionToken.Value,
		SessionExpiredAt: sessionToken.ExpiredAt,
		RefreshToken:     refreshToken.Value,
		RefreshExpiredAt: refreshToken.ExpiredAt,
	}, nil
}
