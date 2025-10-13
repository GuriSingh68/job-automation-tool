package routes

import (
	perferences "github.com/automation/backend/pkg/handlers/preferences"
	resumes "github.com/automation/backend/pkg/handlers/resume"
	"github.com/gin-gonic/gin"
)

func RegisterResumeRoutes(rg *gin.RouterGroup) {
	resume := rg.Group("/resume")
	{
		resume.POST("/upload", resumes.UploadResumeHandler)
		resume.GET("/list", resumes.ListResumesHandler)
		resume.DELETE(":id", resumes.DeleteResume)
	}
}

func RegisterPreferncesRoutes(rg *gin.RouterGroup) {
	preferences := rg.Group("/preferences")
	{
		//preferences.GET("/", preferences.GetPreferencesHandler)
		preferences.POST("/", perferences.CreatePreferencesHandler)
	}
}
