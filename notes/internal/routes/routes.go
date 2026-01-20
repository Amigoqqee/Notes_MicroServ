package routes

import (
	"notes/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(noteHandler *handler.Handler) *gin.Engine {
	router := gin.Default()

	noteAPI := router.Group("/notes")
	noteAPI.Use(noteHandler.GetJWTMiddleware())
	{
		noteAPI.POST("/note", noteHandler.CreateNote)
		noteAPI.GET("/note/:id", noteHandler.GetNoteByID)
		noteAPI.PUT("/note/:id", noteHandler.UpdateNote)
		noteAPI.DELETE("/note/:id", noteHandler.DeleteNote)
		noteAPI.GET("/notes", noteHandler.GetAllNotes)
	}
	return router
}
