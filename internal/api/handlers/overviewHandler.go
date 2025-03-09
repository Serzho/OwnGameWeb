package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type OverviewHandler struct{}

func NewOverviewHandler() *OverviewHandler {
	return &OverviewHandler{}
}

func (h *OverviewHandler) IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
