package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"codeinstyle.io/captain/types"
	"github.com/gin-gonic/gin"
)

type SkillsHandlers struct{}

func NewSkillsHandlers() *SkillsHandlers {
	return &SkillsHandlers{}
}

func (h *SkillsHandlers) GetSkills(c *gin.Context) {
	data, err := os.ReadFile("data/skills.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read skills data"})
		return
	}

	var skills []types.SkillSection
	if err := json.Unmarshal(data, &skills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse skills data"})
		return
	}

	c.JSON(http.StatusOK, skills)
}
