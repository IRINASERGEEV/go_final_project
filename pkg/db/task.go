// pkg/db/task.go
package db

import (
    "database/sql"
    "fmt"
    "time"
    "strconv"
)

// Task — структура задачи
type Task struct {
    ID      string    `json:"id"`
    Date    string `json:"date"`
    Title   string `json:"title"`
    Comment string `json:"comment"`
    Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в БД и возвращает ID
func AddTask(task *Task) (string, error) {
    res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
        task.Date, task.Title, task.Comment, task.Repeat)
    if err != nil {
        return "", err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return "", err
    }
    return strconv.FormatInt(id, 10), nil
}

// Tasks возвращает список задач, отсортированных по дате (от ближайших и id, если даты совпадают)
// limit — максимальное количество задач (50)
// search — строка поиска: может быть текст или дата в формате 02.01.2006
func Tasks(limit int, search string) ([]*Task, error) {
    var rows *sql.Rows
    var err error

    if search == "" {
        // Все задачи
        query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date, id LIMIT ?`
        rows, err = db.Query(query, limit)
    } else {
        // Попробуем как дату: DD.MM.YYYY
        if len(search) == 10 && search[2] == '.' && search[5] == '.' {
            if t, errParse := time.Parse("02.01.2006", search); errParse == nil {
                // Дополнительная проверка: день, месяц, год в диапазоне
                day, _ := strconv.Atoi(search[0:2])
                month, _ := strconv.Atoi(search[3:5])
                year, _ := strconv.Atoi(search[6:10])
                if day >= 1 && day <= 31 && month >= 1 && month <= 12 && year >= 1 && year <= 9999 {
                    dateStr := t.Format("20060102")
                    query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date, id LIMIT ?`
                    rows, err = db.Query(query, dateStr, limit)
                    if err != nil {
                        return nil, err
                    }
                    defer rows.Close()
                    return scanRows(rows)
                }
            }
        }

        // Поиск по тексту: title или comment
        likePattern := "%" + search + "%"
        query := `
            SELECT id, date, title, comment, repeat 
            FROM scheduler 
            WHERE title LIKE ? OR comment LIKE ? 
            ORDER BY date, id 
            LIMIT ?`
        rows, err = db.Query(query, likePattern, likePattern, limit)
        if err != nil {
            return nil, err
        }
        defer rows.Close()
        return scanRows(rows)
    }

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    return scanRows(rows)
}


// scanRows сканирует строки из *sql.Rows и возвращает список задач
func scanRows(rows *sql.Rows) ([]*Task, error) {
    var tasks []*Task
    for rows.Next() {
        var t Task
        var id int
        if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
            return nil, err
        }
        t.ID = strconv.Itoa(id)
        tasks = append(tasks, &t)
    }

    if tasks == nil {
        tasks = []*Task{}
    }

    return tasks, rows.Err()
}

// GetTask возвращает задачу по ID
func GetTask(id string) (*Task, error) {
    var task Task
    query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
    err := DB().QueryRow(query, id).Scan(
        &task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("task not found")
    }
    if err != nil {
        return nil, err
    }
    return &task, nil
}

// UpdateTask обновляет задачу в БД
func UpdateTask(task *Task) error {
    query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
    res, err := DB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
    if err != nil {
        return err
    }

    // Проверяем, была ли затронута хотя бы одна строка
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("incorrect id for updating task")
    }

    return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
    res, err := DB().Exec("DELETE FROM scheduler WHERE id = ?", id)
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("task not found")
    }
    return nil
}

// UpdateDate обновляет только дату задачи
func UpdateDate(id string, nextDate string) error {
    res, err := DB().Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("task not found")
    }
    return nil
}
