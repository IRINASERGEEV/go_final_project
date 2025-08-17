package api

import "net/http"

// Init — регистрирует все обработчики
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler) // nextdate.go
	http.HandleFunc("/api/signin", SignInHandler)    //
	//http.HandleFunc("/api/task", taskHandler)     // POST /api/task
	http.HandleFunc("/api/tasks", auth(tasksHandler))        // GET /api/tasks
	http.HandleFunc("/api/task", auth(taskHandler))          // маршрут с поддержкой GET/PUT/POST
	http.HandleFunc("/api/task/done", auth(doneTaskHandler)) // donetask.go
}

// taskHandler — маршрутизация по HTTP-методу
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
