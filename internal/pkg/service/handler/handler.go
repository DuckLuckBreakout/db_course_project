package handler

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/service"
	"net/http"
)

type Handler struct {
	UseCase service.UseCase
}

func (h Handler) Clear(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h Handler) Status(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func NewHandler(useCase service.UseCase) service.Handler {
	return &Handler{
		UseCase: useCase,
	}
}
