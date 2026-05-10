package routes

import (
	"notes-app/internal/handlers"
	"notes-app/internal/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Setup registers all application routes on the given Gin engine.
func Setup(r *gin.Engine) {
	r.Use(metrics.PrometheusMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Home — list all notes.
	r.GET("/", handlers.IndexHandler)

	// Note CRUD routes.
	notes := r.Group("/notes")
	{
		notes.GET("/new", handlers.NewNoteHandler)
		notes.POST("", handlers.CreateNoteHandler)
		notes.GET("/:id", handlers.ViewNoteHandler)
		notes.GET("/:id/edit", handlers.EditNoteHandler)
		notes.POST("/:id/update", handlers.UpdateNoteHandler)
		notes.POST("/:id/delete", handlers.DeleteNoteHandler)
	}
}
