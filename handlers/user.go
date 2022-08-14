package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/PabloRosalesJ/go/rest-ws/models"
	"github.com/PabloRosalesJ/go/rest-ws/repository"
	"github.com/PabloRosalesJ/go/rest-ws/server"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type SingUpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func SingUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SingUpLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		found, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			log.Fatal("SingUpHandler.found ", err)
		}

		if found.Id != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(HomeResponse{
				Message: "User alrredy exist",
				Status:  false,
			})
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var user = models.User{
			Email:    request.Email,
			Password: string(hashed),
			Id:       id.String(),
		}

		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			log.Println("Error in DB: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(SingUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SingUpLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		found, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			log.Fatal("SingUpHandler.found ", err)
		}

		if found.Id == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(HomeResponse{
				Message: "Invalid credentials",
				Status:  false,
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(request.Password)); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(HomeResponse{
				Message: "Invalid credentials",
				Status:  false,
			})
			return
		}

		claim := models.AppClaims{
			UserId: found.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenSigned, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "Application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenSigned,
		})
	}
}
