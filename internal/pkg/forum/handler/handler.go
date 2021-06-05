package handler

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum"
)

type Handler struct {
	UseCase forum.UseCase
}

func NewHandler(useCase forum.UseCase) forum.Handler {
	return &Handler{
		UseCase: useCase,
	}
}
