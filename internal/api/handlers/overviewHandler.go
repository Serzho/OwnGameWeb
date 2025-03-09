package handlers

import "github.com/gin-gonic/gin"

type OverviewHandler struct{}

func NewOverviewHandler() *OverviewHandler {
	return &OverviewHandler{}
}

func (h *OverviewHandler) IndexPage(_ *gin.Context) {}
