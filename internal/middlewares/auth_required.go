package middlewares

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexlwc11/simple_auth_go/internal/models"
	"github.com/gin-gonic/gin"
)

type TokenProvider func(tokenValue string) (*models.SessionToken, error)

func AuthRequired(tokenProvider TokenProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		tokenValue, err := extractBearerTokenFromHeader(authorizationHeader)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		token, err := tokenProvider(tokenValue)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if err := isTokenExpired(*token); err != nil {
			log.Println(errors.New("invalid token"))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("user_id", token.UserID)
		c.Next()
	}
}

func isTokenExpired(token models.SessionToken) error {
	if token.ExpiredAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	return nil
}

func extractBearerTokenFromHeader(header string) (string, error) {
	if header == "" {
		return "", errors.New("invalid header value")
	}

	token := strings.Split(header, " ")
	if len(token) != 2 {
		return "", errors.New("unable to extract token from header")
	}

	return token[1], nil
}
