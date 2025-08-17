// pkg/repeat/repeat.go
package repeat

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// afterNow возвращает true, если date > now (по дате)
func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return dateOnly.After(nowOnly) // date > now
}

// NextDate возвращает следующую дату > now по правилу repeat
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil // задача не повторяется
	}

	// Парсим начальную дату
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", errors.New("invalid date format, expected YYYYMMDD")
	}

	// Разбиваем правило
	parts := strings.Split(strings.TrimSpace(repeat), " ")
	if len(parts) == 0 {
		return "", errors.New("empty repeat rule")
	}

	switch parts[0] {
	case "y":
		// Ежегодно: +1 год, пока > now
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format("20060102"), nil

	case "d":
		// Каждые N дней
		if len(parts) != 2 {
			return "", errors.New("invalid daily repeat format: use 'd <days>'")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("invalid number of days")
		}
		if interval <= 0 || interval > 400 {
			return "", errors.New("days must be 1-400")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format("20060102"), nil

	case "w":
		// По дням недели: w 1,3,5
		if len(parts) != 2 {
			return "", errors.New("invalid weekly rule: use 'w <1-7>'")
		}
		validDays := [8]bool{} // 1..7
		for _, s := range strings.Split(parts[1], ",") {
			day, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("invalid weekday: must be 1-7")
			}
			validDays[day] = true
		}

		// Идём по дням, пока не найдём подходящий
		for {
			date = date.AddDate(0, 0, 1)
			if !afterNow(date, now) {
				continue
			}
			// В Go: Sunday=0, Monday=1, ..., Saturday=6
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7 // воскресенье
			}
			if validDays[weekday] {
				return date.Format("20060102"), nil
			}
		}

	case "m":
		// Правило "m" требует как минимум один список дней
		if len(parts) < 2 {
			return "", errors.New("invalid monthly rule: days required")
		}

		// Флаги для специальных дней: -1 = последний день месяца, -2 = предпоследний
		hasLast := false
		hasSecondLast := false

		// Множество обычных дней месяца (например, 1, 15, 25)
		specifiedDays := make(map[int]bool)

		// Разбираем строку с днями (например, "1,15,-1")
		for _, s := range strings.Split(parts[1], ",") {
			s = strings.TrimSpace(s) // Убираем пробелы

			if s == "-1" {
				// Отмечаем, что нужно учитывать последний день месяца
				hasLast = true
			} else if s == "-2" {
				// Отмечаем, что нужно учитывать предпоследний день месяца
				hasSecondLast = true
			} else {
				// Обычный день месяца (например, "15")
				day, err := strconv.Atoi(s)
				if err != nil || day < 1 || day > 31 {
					return "", errors.New("invalid day: must be 1-31 or -1, -2")
				}
				specifiedDays[day] = true // Добавляем в множество
			}
		}

		// Множество указанных месяцев (опционально, например, "1,3,6")
		specifiedMonths := make(map[int]bool)
		if len(parts) >= 3 {
			for _, s := range strings.Split(parts[2], ",") {
				month, err := strconv.Atoi(strings.TrimSpace(s))
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("invalid month: must be 1-12")
				}
				specifiedMonths[month] = true
			}
		}

		// Начинаем с начальной даты задачи
		date, _ := time.Parse("20060102", dstart)

		// Идём по дням вперёд, пока не найдём подходящую дату
		for {
			date = date.AddDate(0, 0, 1) // +1 день

			// Пропускаем даты, которые не больше текущего времени (now)
			if !afterNow(date, now) {
				continue
			}

			month := int(date.Month())

			// Если указаны месяцы, проверяем, входит ли текущий месяц
			if len(specifiedMonths) > 0 && !specifiedMonths[month] {
				continue
			}

			day := date.Day()

			// Проверяем специальные дни
			if hasLast && isLastDayOfMonth(date) {
				// Это последний день месяца — подходит
				return date.Format("20060102"), nil
			}
			if hasSecondLast && isSecondLastDayOfMonth(date) {
				// Это предпоследний день месяца — подходит
				return date.Format("20060102"), nil
			}
			if specifiedDays[day] {
				// Это указанный обычный день — подходит
				return date.Format("20060102"), nil
			}
		}
	}
	return "", errors.New("internal error: unreachable")
}

// isLastDayOfMonth проверяет, является ли дата последним днём месяца
func isLastDayOfMonth(date time.Time) bool {
	next := date.AddDate(0, 0, 1)
	return next.Month() != date.Month()
}

// isSecondLastDayOfMonth проверяет, является ли дата предпоследним днём месяца
func isSecondLastDayOfMonth(date time.Time) bool {
	next := date.AddDate(0, 0, 2)
	return next.Month() != date.Month()
}
