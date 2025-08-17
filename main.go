// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

const (
	defaultPort = "7540"
	webDir      = "web"
	defaultDB   = "scheduler.db"
)

func main() {
	// Определяем БД: сначала из переменной окружения, потом defaultDB
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDB
	}

	// Инициализация БД
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Printf("Database initialized: %s", dbFile)

	// Проверяем, существует ли папка web
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatalf("Directory %s not found. Make sure it exists.", webDir)
	}

	// Регистрируем API-обработчики
	api.Init()

	// Настройка веб-сервера
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	// Логируем путь к веб-файлам
	absPath, err := filepath.Abs(webDir)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path for web directory: %v", err)
	}
	log.Printf("Serving files from: %s", absPath)

	// Определяем порт: сначала из переменной окружения, потом defaultPort
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	// Запускаем сервер
	addr := ":" + port
	log.Printf("Server is running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
