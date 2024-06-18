package cmd

import (
	"github.com/alexlwc11/simple_auth_go/internal/apis"
	"github.com/alexlwc11/simple_auth_go/internal/middlewares"
	"github.com/alexlwc11/simple_auth_go/internal/models"
	"github.com/alexlwc11/simple_auth_go/internal/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProceedSchemaMigration(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.SessionToken{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.RefreshToken{}); err != nil {
		return err
	}

	return nil
}

func CreateAuthHandler(db *gorm.DB) apis.AuthHandler {
	return apis.NewAuthHandlerImpl(
		repositories.NewUserRepositoryImpl(db),
		repositories.NewSessionTokenRepositoryImpl(db),
		repositories.NewRefreshTokenRepositoryImpl(db),
	)
}

func AuthRequired(db *gorm.DB) gin.HandlerFunc {
	return middlewares.AuthRequired(
		repositories.NewSessionTokenRepositoryImpl(db).FindByValue,
	)
}
