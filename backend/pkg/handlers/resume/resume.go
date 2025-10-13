package resume

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/automation/backend/db"
	"github.com/gin-gonic/gin"
)

type ParsedResume struct {
	Raw      map[string]string      `json:"raw"`
	Sections map[string]interface{} `json:"sections"`
	Metadata map[string]string      `json:"metadata"`
}

func ListResumesHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"resumes": []string{"resume1.pdf", "resume2.pdf"},
	})
}

// Upload handler (async parsing)
func UploadResumeHandler(c *gin.Context) {
	//userID := c.PostForm("user_id")
	userID := 1 // Placeholder for now
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing resume file"})
		return
	}

	os.MkdirAll("uploads/resumes", os.ModePerm)
	savePath := filepath.Join("uploads/resumes", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Insert placeholder in DB
	conn := db.GetDB()
	if conn == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database not initialized"})
		return
	}
	res, err := conn.Exec(`
		INSERT INTO resumes (user_id, file_path, original_json)
		VALUES (?, ?, ?)`,
		userID, savePath, "{}")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert into db"})
		return
	}

	resumeID, _ := res.LastInsertId()

	// Run parser in background goroutine
	go func(resumeID int64, savePath string) {
		fmt.Println("🧠 [DEBUG] Starting parser for:", savePath)

		cmd := exec.Command("python", "scripts/resume_parser.py", savePath)
		out, err := cmd.CombinedOutput() // capture both stdout and stderr

		if err != nil {
			fmt.Println("❌ [ERROR] Parser failed:", err)
			fmt.Println("🔍 [STDERR]", string(out))
			conn.Exec("INSERT INTO logs (message, level) VALUES (?, ?)", fmt.Sprintf("Parser error: %v", err), "error")
			return
		}

		var parsed ParsedResume
		if err := json.Unmarshal(out, &parsed); err != nil {
			fmt.Println("❌ [ERROR] Failed to parse JSON:", err)
			fmt.Println("🔍 [OUTPUT]", string(out))
			conn.Exec("INSERT INTO logs (message, level) VALUES (?, ?)", fmt.Sprintf("JSON parse error: %v", err), "error")
			return
		}

		fmt.Println("✅ [DEBUG] Parsed resume successfully, updating DB...")

		_, err = conn.Exec(`
        UPDATE resumes
        SET original_json = ?
        WHERE id = ?`,
			toJSON(parsed), resumeID)
		if err != nil {
			fmt.Println("❌ [ERROR] DB update failed:", err)
		}

		conn.Exec("INSERT INTO logs (message, level) VALUES (?, ?)",
			fmt.Sprintf("Resume %d parsed successfully", resumeID), "info")
	}(resumeID, savePath)

	c.JSON(http.StatusAccepted, gin.H{
		"message":   "Resume uploaded successfully — parsing in background",
		"resume_id": resumeID,
	})
}

func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
