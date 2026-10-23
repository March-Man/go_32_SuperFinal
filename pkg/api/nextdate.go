package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

// NextDate calculates the next date of a task based on the repeat rule.
// now is the current time, dstart is the start date, and repeat is the rule.
// It returns the date in 20060102 format, which is strictly greater than now.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неверная дата: %w", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат правила")
	}

	switch parts[0] {
	case "d":
		return nextDateByDays(now, date, parts)
	case "y":
		return nextDateByYear(now, date)
	case "w":
		return nextDateByWeekdays(now, date, parts)
	case "m":
		return nextDateByMonths(now, date, parts)
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", parts[0])
	}
}

// afterNow returns true if date is strictly greater than now (ignoring time).
func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	dateOnly := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}

// nextDateByDays handles the "d N" rule.
func nextDateByDays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат правила d")
	}

	days, err := strconv.Atoi(parts[1])
	if err != nil || days <= 0 || days > 400 {
		return "", fmt.Errorf("неверный интервал дней: %s", parts[1])
	}

	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

// nextDateByYear handles the "y" rule.
func nextDateByYear(now, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

// nextDateByWeekdays handles the "w N,N,..." rule.
func nextDateByWeekdays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат правила w")
	}

	var weekdays [8]bool
	for _, s := range strings.Split(parts[1], ",") {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 7 {
			return "", fmt.Errorf("неверный день недели: %s", s)
		}
		weekdays[n] = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		if !afterNow(date, now) {
			continue
		}
		wd := int(date.Weekday())
		if wd == 0 {
			wd = 7
		}
		if weekdays[wd] {
			return date.Format(dateFormat), nil
		}
	}
}

// nextDateByMonths handles the "m D,D,... [M,M,...]" rule.
func nextDateByMonths(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("неверный формат правила m")
	}

	daySet := make(map[int]bool)
	for _, s := range strings.Split(parts[1], ",") {
		n, err := strconv.Atoi(s)
		if err != nil || n == 0 || n < -2 || n > 31 {
			return "", fmt.Errorf("неверный день месяца: %s", s)
		}
		daySet[n] = true
	}

	var monthSet [13]bool
	if len(parts) == 3 {
		for _, s := range strings.Split(parts[2], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 12 {
				return "", fmt.Errorf("неверный месяц: %s", s)
			}
			monthSet[n] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			monthSet[i] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)
		if !afterNow(date, now) {
			continue
		}

		y, m, d := date.Date()
		month := int(m)

		if !monthSet[month] {
			continue
		}

		if daySet[d] {
			return date.Format(dateFormat), nil
		}

		lastDay := daysInMonth(y, month)
		if daySet[-1] && d == lastDay {
			return date.Format(dateFormat), nil
		}
		if daySet[-2] && d == lastDay-1 {
			return date.Format(dateFormat), nil
		}
	}
}

// daysInMonth returns the number of days in the given month.
func daysInMonth(year, month int) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
