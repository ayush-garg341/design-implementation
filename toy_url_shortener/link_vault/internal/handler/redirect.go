package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linkvault/internal/service"
)

type RedirectRequestBody struct {
	ShortCode string `json:"short_code"`
}

type RedirectHandler struct {
	service *service.RedirectService
}

func NewRedirectHandler(s *service.RedirectService) *RedirectHandler {
	return &RedirectHandler{service: s}
}

func (h *RedirectHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortcode")
	longUrl, err := h.service.GetRedirectUrl(c.Request.Context(), shortCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Redirect(http.StatusFound, *longUrl)
}
