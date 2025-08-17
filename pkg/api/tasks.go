// pkg/api/tasks.go
package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

// TasksResp — структура ответа
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Только GET
	if r.Method != http.MethodGet {
		writeJSON(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	//задачи по умолчанию или поиск
	limit := 50
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(limit, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": "failed to fetch tasks"}, http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON: {"tasks": [...]}
	writeJSON(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
