// pkg/api/nextdate.go
package api

import (
    "fmt"
    "net/http"
    "time"

    "go_final_project/pkg/repeat" 
)

const dateFormat = "20060102"

// nextDayHandler обрабатывает GET /api/nextdate?now &date &repeat
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
    // Получаем параметры
    nowStr := r.FormValue("now")
    dateStr := r.FormValue("date")
    repeatRule := r.FormValue("repeat")

    // Проверяем обязательные параметры
    if dateStr == "" {
        http.Error(w, "missing 'date' parameter", http.StatusBadRequest)
        return
    }
    if repeatRule == "" {
        http.Error(w, "missing 'repeat' parameter", http.StatusBadRequest)
        return
    }

    // Парсим now: если не указан — используем текущую дату
    var now time.Time
    if nowStr != "" {
        parsed, err := time.Parse(dateFormat, nowStr)
        if err != nil {
            http.Error(w, "invalid 'now' format, expected YYYYMMDD", http.StatusBadRequest)
            return
        }
        now = parsed
    } else {
        now = time.Now()
    }

    // Вызываем функцию из pkg/repeat
    nextDate, err := repeat.NextDate(now, dateStr, repeatRule)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Если nextDate пустой — задача не повторяется
    if nextDate == "" {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "") // пустой ответ
        return
    }

    // Возвращаем дату в формате 20060102
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, nextDate)
}
