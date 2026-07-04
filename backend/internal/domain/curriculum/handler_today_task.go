package curriculum

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTodayTask(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	task, err := h.repo.GetTodayTask(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed_to_load_today_task"})
		return
	}

	c.JSON(http.StatusOK, TodayTaskResponse{Task: task})
}
