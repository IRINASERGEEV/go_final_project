package db

import (
    "database/sql"
    "os"

    _ "modernc.org/sqlite"
)

// Глобальная переменная для соединения с БД
var db *sql.DB

// Схема: создание таблицы и индекса
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
`

// Init инициализирует базу данных: открывает соединение и создаёт таблицу, если нужно
func Init(dbFile string) error {
    // Проверяем, существует ли файл БД
    _, err := os.Stat(dbFile)
    dbExists := err == nil
    if err != nil && !os.IsNotExist(err) {
        return err
    }

    // Открываем соединение с SQLite
    var errOpen error
    db, errOpen = sql.Open("sqlite", dbFile)
    if errOpen != nil {
        return errOpen
    }

    // Если БД не существовала — создаём таблицу и индекс
    if !dbExists {
        _, err = db.Exec(schema)
        if err != nil {
            return err
        }
    }

    return nil
}

// DB возвращает текущее соединение с базой данных
func DB() *sql.DB {
    return db
}
