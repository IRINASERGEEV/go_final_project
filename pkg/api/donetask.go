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
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        writeJSON(w, map[string]string{"error": "id parameter is required"})
        return
    }

    task, err := db.GetTask(idStr)
    if err != nil {
        writeJSON(w, map[string]string{"error": "task not found"})
        return
    }

    now := time.Now()

    // Если нет правила повторения — удаляем
    if task.Repeat == "" {
        if err := db.DeleteTask(idStr); err != nil {
            writeJSON(w, map[string]string{"error": "failed to delete task"})
        } else {
            writeJSON(w, map[string]interface{}{}) // {}
        }
        return
    }

    // Если есть repeat — вычисляем следующую дату
    nextDate, err := repeat.NextDate(now, task.Date, task.Repeat)
    if err != nil {
        writeJSON(w, map[string]string{"error": err.Error()})
        return
    }

    // Если nextDate пустой — удаляем
    if nextDate == "" {
        db.DeleteTask(idStr)
        writeJSON(w, map[string]interface{}{})
        return
    }

    // Обновляем дату
    if err := db.UpdateDate(idStr, nextDate); err != nil {
        writeJSON(w, map[string]string{"error": "failed to update task date"})
        return
    }

    writeJSON(w, map[string]interface{}{}) // {}
}

//Удаление
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        writeJSON(w, map[string]string{"error": "id parameter is required"})
        return
    }

    if err := db.DeleteTask(idStr); err != nil {
        writeJSON(w, map[string]string{"error": "task not found"})
        return
    }

    writeJSON(w, map[string]interface{}{}) // {}
}
