package helpers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PabloRosalesJ/go/rest-ws/models"
	"github.com/PabloRosalesJ/go/rest-ws/repository"
	"github.com/PabloRosalesJ/go/rest-ws/server"
	"github.com/golang-jwt/jwt"
)

const UNAUTHORIZED string = "Unauthorized"

func AuthUser(s server.Server, w http.ResponseWriter, r *http.Request) (models.User, error) {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.Config().JWTSecret), nil
	})

	var user = models.User{}

	if err != nil {
		return user, err
	}

	if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
		user, err := repository.GetUserById(r.Context(), claims.UserId)
		if err != nil {
			return models.User{}, err
		}

		if user.Id == "" {
			return *user, errors.New("Unauthorized")
		}

		return *user, nil
	} else {
		return user, err
	}
}

func AuthUserId(s server.Server, r *http.Request) (*models.AppClaims, error) {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))

	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.Config().JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New(UNAUTHORIZED)
}
