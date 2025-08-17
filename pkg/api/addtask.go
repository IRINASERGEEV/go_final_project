package api

import (
    "encoding/json"  
    "io"
    "net/http"
    "time"
 
    "go_final_project/pkg/db"
    "go_final_project/pkg/repeat"
)

// addTaskHandler — обрабатывает POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
    // Только POST
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Читаем тело
    body, err := io.ReadAll(r.Body)
    if err != nil {
        writeJSON(w, map[string]string{"error": "cannot read request"})
        return
    }

    var task db.Task
    if err := json.Unmarshal(body, &task); err != nil {
        writeJSON(w, map[string]string{"error": "invalid JSON"})
        return
    }

    // Проверка: title обязательный
    if task.Title == "" {
        writeJSON(w, map[string]string{"error": "title is required"})
        return
    }

    now := time.Now()

    // Если дата не указана — ставим сегодня
    if task.Date == "" {
        task.Date = now.Format(dateFormat)
    }

    // Парсим дату
    t, err := time.Parse(dateFormat, task.Date)
    if err != nil {
        writeJSON(w, map[string]string{"error": "invalid date format, expected YYYYMMDD"})
        return
    }

    // Если правило указано — проверяем его
    if task.Repeat != "" {
        _, err := repeat.NextDate(now, task.Date, task.Repeat)
        if err != nil {
            writeJSON(w, map[string]string{"error": "task not found"})
            return
        }
    }

    // Если дата в прошлом
    if afterNow(now, t) {
        if task.Repeat == "" {
            task.Date = now.Format(dateFormat)
        } else {
            next, _ := repeat.NextDate(now, task.Date, task.Repeat)
            task.Date = next
        }
    }

    // Сохраняем
    id, err := db.AddTask(&task)
    if err != nil {
        writeJSON(w, map[string]string{"error": "failed to save task"})
        return
    }

    // Успешный ответ
    writeJSON(w, map[string]string{"id": id})
}

// getTaskHandler обрабатывает GET /api/task?id=...
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        writeJSON(w, map[string]string{"error": "id parameter is required"})
        return
    }

    task, err := db.GetTask(idStr)
    if err != nil {
        if err.Error() == "task not found" {
            writeJSON(w, map[string]string{"error": "task not found"})
        } else {
            writeJSON(w, map[string]string{"error": "database error"})
        }
        return
    }

    writeJSON(w, task)
}

// updateTaskHandler обрабатывает PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
    // Только PUT
    if r.Method != http.MethodPut {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Читаем тело
    body, err := io.ReadAll(r.Body)
    if err != nil {
        writeJSON(w, map[string]string{"error": "cannot read request"})
        return
    }

    var task db.Task
    if err := json.Unmarshal(body, &task); err != nil {
        writeJSON(w, map[string]string{"error": "invalid JSON"})
        return
    }

    // Проверка: ID обязателен
    if task.ID == "" {
        writeJSON(w, map[string]string{"error": "id is required"})
        return
    }

    // Проверка: title обязателен
    if task.Title == "" {
        writeJSON(w, map[string]string{"error": "title is required"})
        return
    }

    now := time.Now()

    // Если дата не указана — ставим сегодня
    if task.Date == "" {
        task.Date = now.Format(dateFormat)
    }

    // Парсим дату
    t, err := time.Parse(dateFormat, task.Date)
    if err != nil {
        writeJSON(w, map[string]string{"error": "invalid date format, expected YYYYMMDD"})
        return
    }

    // Если правило указано — проверяем его
    if task.Repeat != "" {
        _, err := repeat.NextDate(now, task.Date, task.Repeat)
        if err != nil {
            writeJSON(w, map[string]string{"error": err.Error()})
            return
        }
    }

    // Если дата в прошлом
    if afterNow(now, t) {
        if task.Repeat == "" {
            task.Date = now.Format(dateFormat)
        } else {
            next, _ := repeat.NextDate(now, task.Date, task.Repeat)
            task.Date = next
        }
    }

    // Обновляем в БД
    if err := db.UpdateTask(&task); err != nil {
        if err.Error() == "incorrect id for updating task" {
            writeJSON(w, map[string]string{"error": "task not found"})
        } else {
            writeJSON(w, map[string]string{"error": "failed to update task"})
        }
        return
    }

    // Успешно — пустой JSON
    writeJSON(w, map[string]interface{}{})
}
