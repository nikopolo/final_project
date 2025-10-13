package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"task_scheduler/pkg/db"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// json.NewEncoder(w).Encode(data) с ним не проходят тесты
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "error json marshal", http.StatusBadRequest)
		return
	}

	w.Write(jsonData)
}

// addTaskHandler POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Десериализация JSON
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

	// Добавление задачи
	id, err := db.AddTask(&task)
	if err != nil {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "add task"})
		return
	}

	// Возврат ID
	JSON(w, http.StatusOK, map[string]int64{"id": id})
}

// проверка даты
func checkDate(task *db.Task, now time.Time) error {
	dateLayout := "20060102"

	// дата не указана — сегодня
	if task.Date == "" {
		task.Date = now.Format(dateLayout)
		return nil
	}

	// Проверка формат
	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	if isAfterNow(t, now) || t.Format(dateLayout) == now.Format(dateLayout) {
		return nil
	}

	if task.Repeat == "" {
		task.Date = now.Format(dateLayout)
	} else {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = next
	}

	return nil
}

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
