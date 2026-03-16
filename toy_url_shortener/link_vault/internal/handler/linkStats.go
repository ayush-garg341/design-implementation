package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linkvault/internal/analytics"
)

type LinkStatsHandler struct {
	analytics analytics.Client
}

func NewLinkStatsHandler(a analytics.Client) *LinkStatsHandler {
	return &LinkStatsHandler{analytics: a}
}

// GET /api/stats/link/:code — returns link + its analytics

func (h *LinkStatsHandler) GetLinkStats(c *gin.Context) {
	shortCode := c.Param("code")

	ctx := c.Request.Context()

	stats, err := h.analytics.GetStats(ctx, shortCode)
	if err != nil {
		if errors.Is(err, analytics.ErrUnavailable) {
			// analytics service is down — still return the link, just no stats
			c.JSON(http.StatusOK, gin.H{
				"link":  shortCode,
				"stats": nil,
				"note":  "analytics temporarily unavailable",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"link": shortCode,
		"stats": gin.H{
			"total_clicks": stats.ClickCount,
			"long_url":     stats.LongUrl,
			"user_id":      stats.UserId,
		},
	})
}
