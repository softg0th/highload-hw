package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockRepo struct{}

func (m *mockRepo) GetLastMessages(ctx context.Context) ([]string, error) {
	return []string{"hello", "world"}, nil
}

type mockHandler struct {
	Repo struct {
		DB interface {
			GetLastMessages(ctx context.Context) ([]string, error)
		}
	}
}

func (h *mockHandler) GetLastMessages(c *gin.Context) {
	ctx := context.Background()
	output, err := h.Repo.DB.GetLastMessages(ctx)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"output": "Unknown error"})
		return
	}
	outputJSON, _ := json.Marshal(output)
	c.JSON(http.StatusOK, gin.H{"output": outputJSON})
}

func TestGetLastMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := &mockHandler{}
	handler.Repo.DB = &mockRepo{}

	handler.GetLastMessages(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}