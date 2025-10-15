package api

import (
	"net/http"

	"task_scheduler/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler GET /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	const MAX_RECORDS = 50

	if r.Method != http.MethodGet {
		http.Error(w, "Error method", http.StatusBadRequest)
		return
	}

	tasks, err := db.Tasks(MAX_RECORDS) // максимальное количество записей
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "tasks"})
		return
	}

	JSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}
