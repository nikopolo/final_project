package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102"

// для определения date позже now
func isAfterNow(date, now time.Time) bool {
	return date.Format(layout) > now.Format(layout)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// проверки на пустой repeat
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	repeat = strings.TrimSpace(repeat)
	if len(repeat) == 0 {
		return "", errors.New("repeat is empty")
	}

	paramsRepeat := strings.Split(repeat, " ")
	if len(paramsRepeat) == 0 {
		return "", errors.New("repeat is empty")
	}

	// парсинг даты
	startDate, err := time.Parse(layout, dstart)
	if err != nil {
		return "", fmt.Errorf("error dstart: %v", err)
	}

	// параметр repeat "y"
	if paramsRepeat[0] == "y" {
		if len(paramsRepeat) != 1 {
			return "", errors.New("error format")
		}
		date := startDate
		for {
			date = date.AddDate(1, 0, 0)
			if isAfterNow(date, now) {
				return date.Format(layout), nil
			}
		}
	}

	// параметр repeat "d"
	if paramsRepeat[0] == "d" {
		if len(paramsRepeat) != 2 {
			return "", errors.New("error format")
		}
		interval, err := strconv.Atoi(string(paramsRepeat[1]))
		if err != nil {
			return "", fmt.Errorf("error day interval: %v", err)
		}
		if interval <= 0 || interval > 400 {
			return "", fmt.Errorf("days must be between 1 and 400")
		}

		date := startDate
		for {
			date = date.AddDate(0, 0, interval)
			if isAfterNow(date, now) {
				return date.Format(layout), nil
			}
		}
	}

	return "", fmt.Errorf("unsupported operation")
}

// запрос к /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {

	// параметры
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	if dateParam == "" {
		http.Error(w, "empty param date", http.StatusBadRequest)
		return
	}
	if repeatParam == "" {
		http.Error(w, "empty param repeat", http.StatusBadRequest)
		return
	}

	// если не задана текущая дата
	var now time.Time
	if nowParam == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(layout, nowParam)
		if err != nil {
			http.Error(w, "error 'now'", http.StatusBadRequest)
			return
		}
	}

	// вызов NextDate
	next, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}

// запрос к /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method error", http.StatusMethodNotAllowed)
	}
}
