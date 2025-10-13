package routes

import (
	handlers "github.com/automation/backend/pkg/handlers/resume"
	"github.com/gin-gonic/gin"
)

func RegisterResumeRoutes(rg *gin.RouterGroup) {
	resume := rg.Group("/resume")
	{
		resume.POST("/upload", handlers.UploadResumeHandler)
		resume.GET("/list", handlers.ListResumesHandler)
	}
}
