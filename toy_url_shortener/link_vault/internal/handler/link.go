package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linkvault/internal/service"
)

type CreateShortLinkBody struct {
	LongUrl string `json:"long_url"`
}

type LinkHandler struct {
	service *service.LinkService
}

func NewLinkHandler(s *service.LinkService) *LinkHandler {
	return &LinkHandler{service: s}
}

func (h *LinkHandler) Create(c *gin.Context) {
	var shortlink CreateShortLinkBody
	err := c.BindJSON(&shortlink)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	newLink, err := h.service.CreateShortLink(c.Request.Context(), shortlink.LongUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, newLink)
}

func (h *LinkHandler) List(c *gin.Context) {
	links, err := h.service.AllLinks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, links)
}
