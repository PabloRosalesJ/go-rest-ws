package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/PabloRosalesJ/go/rest-ws/helpers"
	apiResponse "github.com/PabloRosalesJ/go/rest-ws/helpers/api"
	"github.com/PabloRosalesJ/go/rest-ws/models"
	"github.com/PabloRosalesJ/go/rest-ws/repository"
	"github.com/PabloRosalesJ/go/rest-ws/server"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type PostRequest struct {
	PostContent string `json:"postContent"`
}

type PostResponse struct {
	Id          string      `json:"id"`
	PostContent string      `json:"postContent"`
	User        models.User `json:"user"`
}

func CreatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := helpers.AuthUser(s, w, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var request = PostRequest{}
		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var newPost = models.Post{
			Id:          id.String(),
			PostContent: request.PostContent,
			UserId:      user.Id,
		}

		err = repository.InsertPost(r.Context(), &newPost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(PostResponse{
			Id:          newPost.Id,
			PostContent: newPost.PostContent,
			User: models.User{
				Id:    user.Id,
				Email: user.Email,
			},
		})
	}
}

func GetPostById(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repository.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if post == nil {
			apiResponse.ResponseOk(w, 404, "Model not found")
			return
		}

		apiResponse.ResponseOk(w, 200, post)
	}
}

func UpdatePost(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := helpers.AuthUserId(s, r)
		if err != nil {
			apiResponse.ResponseErr(w, http.StatusInternalServerError, err.Error(), []string{})
			return
		}

		userId := claims.UserId
		postId := mux.Vars(r)["id"]
		var request = PostRequest{}

		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			apiResponse.ResponseErr(w, http.StatusBadRequest, "Invalid request", []string{err.Error()})
			return
		}

		err = repository.UpdatePost(r.Context(), postId, userId, request.PostContent)

		if err != nil {
			apiResponse.ResponseErr(w, http.StatusInternalServerError, err.Error(), []string{})
			return
		}

		apiResponse.ResponseOk(w, http.StatusOK, "Update successfull")
	}

}
