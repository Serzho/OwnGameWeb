package handlers

import "OwnGameWeb/internal/services"

type PlayHandler struct {
	service *services.PlayService
}

func NewPlayHandler(s *services.PlayService) *PlayHandler {
	return &PlayHandler{service: s}
}
