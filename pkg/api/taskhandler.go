package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task_scheduler/pkg/db"
	"time"
)

// getTaskHandler GET /api/task?id
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idData := r.FormValue("id")
	if idData == "" {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "no id"})
		return
	}

	id, err := strconv.ParseInt(idData, 10, 64)
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "fail id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "issue not found"}) // Issue not found
		return
	}

	JSON(w, http.StatusOK, task)
}

// updateTaskHandler PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	// Проверка title
	if task.Title == "" {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "title is empty"})
		return
	}

	// Проверка даты
	now := time.Now()
	if err := checkDate(&task, now); err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "date missmatch"})
		return
	}

	// Обновление задачи
	err := db.UpdateTask(&task)
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "fail update"})
		return
	}

	//
	JSON(w, http.StatusOK, map[string]interface{}{})
}

// doneTaskHandler обрабатывает POST /api/task/done?id=...
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// проверка метода
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}

	idData := r.FormValue("id")
	if idData == "" {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "no id"})
		return
	}

	id, err := strconv.ParseInt(idData, 10, 64)
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "fail id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		JSON(w, http.StatusNotFound, map[string]string{"error": "issue not found"})
		return
	}

	// Если repeat пустой — удаляем
	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			JSON(w, http.StatusBadRequest, map[string]string{"error": "failed to delete"})
			return
		}

		JSON(w, http.StatusOK, map[string]interface{}{})
		return
	}

	// следующая дата
	now := time.Now()
	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "fail next date"})
		return
	}

	// Обновляем дату
	err = db.UpdateDate(nextDate, id)
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "fail to update"})
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{})
}

// deleteTaskHandler
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idData := r.FormValue("id")
	if idData == "" {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "no id"})
		return
	}

	id, err := strconv.ParseInt(idData, 10, 64)
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "fail id"})
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		JSON(w, http.StatusNotFound, map[string]string{"error": "fail to delete"})
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{})
}
