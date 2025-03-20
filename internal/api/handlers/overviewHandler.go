package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OverviewHandler struct{}

func NewOverviewHandler() *OverviewHandler {
	return &OverviewHandler{}
}

func (h *OverviewHandler) IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
