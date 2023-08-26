package handler

import (
	"github.com/perpetua1g0d/watch_organizer/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}
