package handlers

import (
	"strconv"

	"codeinstyle.io/captain/models"
)

// parseUint converts a string to uint, returns 0 if conversion fails
func parseUint(pageID string) uint {
	pid, err := strconv.ParseUint(pageID, 10, 32)
	if err != nil {
		return 0
	}
	return uint(pid)
}

// getNextMenuPosition gets the next available menu position
func (h *AdminHandlers) getNextMenuPosition() int {
	var maxPosition struct{ Max int }
	h.db.Model(&models.MenuItem{}).Select("COALESCE(MAX(position), -1) + 1 as max").Scan(&maxPosition)
	return maxPosition.Max
}
