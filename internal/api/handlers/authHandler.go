package handlers

import (
	"OwnGameWeb/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) SignIn(_ *gin.Context) {}

func (h *AuthHandler) SignUp(_ *gin.Context) {}

func (h *AuthHandler) SignUpPage(_ *gin.Context) {}

func (h *AuthHandler) SignInPage(_ *gin.Context) {}

func (h *AuthHandler) RecoverPassword(_ *gin.Context) {}

func (h *AuthHandler) RecoverPasswordPage(_ *gin.Context) {}
