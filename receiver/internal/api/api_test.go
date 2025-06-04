package api

import (
	"receiver/internal/infra"
	"receiver/internal/service"
	"testing"
)

func TestNewHandler(t *testing.T) {
	i := &infra.Infra{}
	s := &service.Service{}

	h := NewHandler(i, s)
	if h == nil {
		t.Fatal("Expected handler to be created")
	}
}