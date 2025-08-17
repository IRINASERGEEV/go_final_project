// pkg/api/donetask.go
package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/repeat"
)

// doneTaskHandler обрабатывает POST /api/task/done?id=...
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, map[string]string{"error": "id parameter is required"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(idStr)
	if err != nil {
		if err.Error() == "task not found" {
			writeJSON(w, map[string]string{"error": "task not found"}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": "database error"}, http.StatusInternalServerError)
		}
		return
	}

	now := time.Now()

	// Если нет правила повторения — удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(idStr); err != nil {
			writeJSON(w, map[string]string{"error": "failed to delete task"}, http.StatusInternalServerError)
		} else {
			writeJSON(w, map[string]interface{}{}, http.StatusOK) // {}
		}
		return
	}

	// Если есть repeat — вычисляем следующую дату
	nextDate, err := repeat.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Если nextDate пустой — удаляем
	if nextDate == "" {
		if err := db.DeleteTask(idStr); err != nil {
			writeJSON(w, map[string]interface{}{}, http.StatusInternalServerError)
		}
		writeJSON(w, map[string]interface{}{}, http.StatusOK)
		return
	}

	// Обновляем дату
	if err := db.UpdateDate(idStr, nextDate); err != nil {
		writeJSON(w, map[string]string{"error": "failed to update task date"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK) // {}
}

// Удаление
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, map[string]string{"error": "id parameter is required"}, http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(idStr); err != nil {
		if err.Error() == "task not found" {
			writeJSON(w, map[string]string{"error": "task not found"}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": "failed to delete task"}, http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK) // {}
}
