package repositories

import (
	"time"

	"github.com/alexlwc11/simple_auth_go/internal/constants"
	"github.com/alexlwc11/simple_auth_go/internal/models"
	"github.com/alexlwc11/simple_auth_go/internal/utils"

	"gorm.io/gorm"
)

type SessionTokenRepository interface {
	CreateWithUserId(userId uint) (*models.SessionToken, error)
	FindByValue(value string) (*models.SessionToken, error)
}

type SessionTokenRepositoryImpl struct {
	DB *gorm.DB
}

func NewSessionTokenRepositoryImpl(db *gorm.DB) SessionTokenRepository {
	return &SessionTokenRepositoryImpl{DB: db}
}

func (str *SessionTokenRepositoryImpl) CreateWithUserId(userId uint) (*models.SessionToken, error) {
	expiredDuration := time.Now().Add(constants.SessionTokenValidTime)

	tokenValue, err := utils.GenerateToken()
	if err != nil {
		return &models.SessionToken{}, err
	}

	token := models.SessionToken{
		Token: models.Token{
			UserID:    userId,
			Value:     tokenValue,
			ExpiredAt: expiredDuration,
		},
	}

	dbError := str.DB.Create(&token).Error
	if dbError != nil {
		return &models.SessionToken{}, dbError
	}

	return &token, nil
}

func (str *SessionTokenRepositoryImpl) FindByValue(value string) (*models.SessionToken, error) {
	var token models.SessionToken
	err := str.DB.Where("value = ?", value).First(&token).Error
	if err != nil {
		return &models.SessionToken{}, err
	}

	return &token, nil
}
