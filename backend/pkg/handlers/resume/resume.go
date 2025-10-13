package resume

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/automation/backend/db"
	errors "github.com/automation/backend/pkg/error"
	"github.com/gin-gonic/gin"
)

type ParsedResume struct {
	Raw      map[string]string      `json:"raw"`
	Sections map[string]interface{} `json:"sections"`
	Metadata map[string]string      `json:"metadata"`
}

// UploadResumeHandler handles resume file upload and triggers background parsing
func UploadResumeHandler(c *gin.Context) {
	userID := 1 // TODO: Get from authenticated user session

	file, err := c.FormFile("resume")
	if err != nil {
		errors.HandleError(c, errors.BadRequest("missing resume file"))
		return
	}

	if err := os.MkdirAll("uploads/resumes", os.ModePerm); err != nil {
		errors.HandleError(c, errors.InternalError("failed to create upload directory", err))
		return
	}

	savePath := filepath.Join("uploads/resumes", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		errors.HandleError(c, errors.InternalError("failed to save file", err))
		return
	}

	conn := db.GetDB()
	if conn == nil {
		errors.HandleError(c, errors.InternalError("database unavailable", nil))
		return
	}

	res, err := conn.Exec(
		`INSERT INTO resumes (user_id, file_path, original_json) VALUES (?, ?, ?)`,
		userID, savePath, "{}",
	)
	if err != nil {
		errors.HandleError(c, errors.InternalError("failed to insert resume", err))
		return
	}

	resumeID, err := res.LastInsertId()
	if err != nil {
		errors.HandleError(c, errors.InternalError("failed to get resume id", err))
		return
	}

	go parseResumeAsync(resumeID, savePath)

	c.JSON(http.StatusAccepted, gin.H{
		"message":   "resume uploaded successfully",
		"resume_id": resumeID,
	})
}

// parseResumeAsync runs the Python parser and updates the database
func parseResumeAsync(resumeID int64, filePath string) {
	cmd := exec.Command("python", "scripts/resume_parser.py", filePath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("parser failed for resume %d: %v, output: %s", resumeID, err, string(output))
		return
	}

	var parsed ParsedResume
	if err := json.Unmarshal(output, &parsed); err != nil {
		log.Printf("invalid json from parser for resume %d: %v, output: %s", resumeID, err, string(output))
		return
	}

	jsonData, err := json.Marshal(parsed)
	if err != nil {
		log.Printf("failed to marshal parsed data for resume %d: %v", resumeID, err)
		return
	}

	conn := db.GetDB()
	if conn == nil {
		log.Printf("database unavailable for resume %d", resumeID)
		return
	}

	result, err := conn.Exec(
		`UPDATE resumes SET original_json = ? WHERE id = ?`,
		string(jsonData), resumeID,
	)
	if err != nil {
		log.Printf("failed to update resume %d: %v", resumeID, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get rows affected for resume %d: %v", resumeID, err)
		return
	}

	if rowsAffected == 0 {
		log.Printf("no rows updated for resume %d", resumeID)
		return
	}

	log.Printf("successfully parsed resume %d", resumeID)
}

// ListResumesHandler returns all resumes for the current user
func ListResumesHandler(c *gin.Context) {
	conn := db.GetDB()
	if conn == nil {
		errors.HandleError(c, errors.InternalError("database unavailable", nil))
		return
	}

	rows, err := conn.Query(`SELECT id, file_path, original_json FROM resumes`)
	if err != nil {
		errors.HandleError(c, errors.InternalError("failed to query resumes", err))
		return
	}
	defer rows.Close()

	var resumes []map[string]interface{}
	for rows.Next() {
		var id int64
		var filePath, jsonData string

		if err := rows.Scan(&id, &filePath, &jsonData); err != nil {
			log.Printf("failed to scan resume row: %v", err)
			continue
		}

		resumes = append(resumes, map[string]interface{}{
			"id":        id,
			"file_path": filepath.Base(filePath),
			"parsed":    jsonData != "{}",
		})
	}

	if err := rows.Err(); err != nil {
		errors.HandleError(c, errors.InternalError("error iterating resumes", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"resumes": resumes})
}
