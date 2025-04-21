package api

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetLastMessages(c *gin.Context) {
	ctx := context.Background()
	output, err := h.Repo.DB.GetLastMessages(ctx)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"output": "Unknown error"})
		return
	}
	outputJSON, err := json.Marshal(output)
	c.JSON(http.StatusOK, gin.H{"output": outputJSON})
	return
}
