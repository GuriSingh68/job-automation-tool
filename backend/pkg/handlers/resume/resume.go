package resume

import (
	"github.com/gin-gonic/gin"
)

func ListResumesHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"resumes": []string{"resume1.pdf", "resume2.pdf"},
	})
}
