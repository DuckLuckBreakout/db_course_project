package handler

import (
	"encoding/json"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/post"
	"github.com/DuckLuckBreakout/db_course_project/internal/tools/http_utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	UseCase post.UseCase
}

func (h Handler) UpdateDetails(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}
	defer r.Body.Close()

	var updatePost models.Post
	err = json.Unmarshal(body, &updatePost)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusConflict)
		return
	}
	updatePost.Id = int64(id)

	_, err = h.UseCase.UpdateDetails(&updatePost)
	if err == errors.ErrUserNotFound {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusConflict)
		return
	}

	resultData := make(map[string]interface{})
	resultData["author"] = updatePost.Author
	resultData["created"] = updatePost.Created
	resultData["forum"] = updatePost.Forum
	resultData["id"] = updatePost.Id
	resultData["message"] = updatePost.Message
	resultData["thread"] = updatePost.Thread
	resultData["thread"] = updatePost.Thread

	http_utils.SetJSONResponse(w, updatePost, http.StatusOK)
}

func (h Handler) Details(w http.ResponseWriter, r *http.Request) {

	postId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusConflict)
		return
	}

	result := make(map[string]interface{})

	related := r.URL.Query().Get("related")
	resultData, err := h.UseCase.Details(postId)

	if err == errors.ErrUserNotFound {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
		return
	}
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusConflict)
		return
	}
	postMap := make(map[string]interface{})
	postMap["author"] = resultData.Author
	postMap["created"] = resultData.Created
	postMap["forum"] = resultData.Forum
	postMap["id"] = resultData.Id
	if resultData.IsEdited {
		postMap["isEdited"] = resultData.IsEdited
	}
	postMap["message"] = resultData.Message
	postMap["thread"] = resultData.Thread
	result["post"] = postMap
	if strings.Contains(related, "user") {
		resultData, err := h.UseCase.DetailsUser(postId)
		if err == errors.ErrUserNotFound {
			http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
			return
		}
		result["author"] = resultData
		if err != nil {
			http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusConflict)
			return
		}

	}

	if strings.Contains(related, "thread") {
		resultData, err := h.UseCase.DetailsThread(postId)
		if err == errors.ErrUserNotFound {
			http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
			return
		}
		result["thread"] = resultData
		if err != nil {
			http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusConflict)
			return
		}

	}

	if strings.Contains(related, "forum") {
		resultData, err := h.UseCase.DetailsForum(postId)
		if err == errors.ErrUserNotFound {
			http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusNotFound)
			return
		}
		result["forum"] = resultData
		if err != nil {
			http_utils.SetJSONResponse(w, errors.ErrUserNotFound, http.StatusConflict)
			return
		}

	}
	http_utils.SetJSONResponse(w, result, http.StatusOK)

}

func NewHandler(userUCase post.UseCase) post.Handler {
	return &Handler{
		UseCase: userUCase,
	}
}
