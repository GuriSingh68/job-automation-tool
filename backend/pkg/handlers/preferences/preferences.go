package preferences

import (
	"encoding/json"

	"github.com/automation/backend/db"
	errors "github.com/automation/backend/pkg/error"
	"github.com/automation/backend/pkg/types"
	"github.com/gin-gonic/gin"
)

// Store preferences like preferred job titles, locations, etc.
func CreatePreferencesHandler(c *gin.Context) {

	var prefs types.Preference
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	conn := db.GetDB()
	if conn == nil {
		errors.HandleError(c, errors.InternalError("database unavailable", nil))
		return
	}

	// Convert slices to JSON strings before inserting
	jobTitlesJSON, _ := json.Marshal(prefs.JobTitles)
	locationsJSON, _ := json.Marshal(prefs.Locations)
	keywordsJSON, _ := json.Marshal(prefs.Keywords)
	jobTypesJSON, _ := json.Marshal(prefs.JobTypes)

	query := `
		INSERT INTO preferences (user_id, preferred_role, location, keywords, job_type)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := conn.Exec(query, 1, string(jobTitlesJSON), string(locationsJSON), string(keywordsJSON), string(jobTypesJSON))
	if err != nil {
		errors.HandleError(c, errors.InternalError("failed to create preferences", err))
		return
	}

	c.JSON(201, gin.H{"message": "Preferences created successfully"})
}
