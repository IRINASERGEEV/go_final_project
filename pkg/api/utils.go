package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// writeJSON — универсальный JSON-ответ
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// afterNow возвращает true, если t1 > t2 (по дате, без времени)
func afterNow(t1, t2 time.Time) bool {
	d1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	d2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())
	return !d1.Before(d2) // t1 > t2
}
