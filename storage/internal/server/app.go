package server

import (
	"github.com/gin-gonic/gin"
	"storage/internal/api"
	"storage/internal/repository"
)

type App struct {
	Repo   *repository.Repository
	Router *gin.Engine
}

func NewApp(repo *repository.Repository) *App {
	return &App{
		Repo:   repo,
		Router: gin.Default(),
	}
}

func (a *App) setupRouter(handler *api.Handler) {

	api := a.Router.Group("/api")
	{
		api.GET("/", handler.GetLastMessages)
	}
}

func (a *App) SetupApp(port string) error {
	handler := api.NewHandler(a.Repo)
	a.setupRouter(handler)
	return a.Router.Run(":" + port)
}
