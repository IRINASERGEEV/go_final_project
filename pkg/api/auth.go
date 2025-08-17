package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
)

// signInHandler — POST /api/signin
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Читаем JSON
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, map[string]string{"error": "invalid JSON"})
		return
	}

	password := body["password"]
	if password == "" {
		writeJSON(w, map[string]string{"error": "password is required"})
		return
	}

	// Получаем пароль из переменной окружения
	expected := os.Getenv("TODO_PASSWORD")
	if expected == "" {
		// Если пароль не задан — аутентификация не нужна
		writeJSON(w, map[string]string{"error": "authentication not required"})
		return
	}

	// Проверяем пароль
	if password != expected {
		writeJSON(w, map[string]string{"error": "wrong password"})
		return
	}

	// Генерируем "токен" как SHA-256 от пароля
	hash := sha256.Sum256([]byte(expected))
	token := hex.EncodeToString(hash[:])

	// Устанавливаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   8 * 3600, // 8 часов
		HttpOnly: true,
		Secure:   false, // в учебном проекте можно false
	})

	writeJSON(w, map[string]string{"token": token})
}

// auth — middleware для проверки аутентификации
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expectedPass := os.Getenv("TODO_PASSWORD")
		if expectedPass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Пересчитываем ожидаемый токен
		hash := sha256.Sum256([]byte(expectedPass))
		expectedToken := hex.EncodeToString(hash[:])

		if cookie.Value != expectedToken {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
