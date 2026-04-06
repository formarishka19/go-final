package api

import (
	"errors"
	"go-final/pkg/db"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	var err error
	if repeat == "" {
		return "", nil
	}
	date, err := time.Parse(DATEFORMAT, dstart)
	if err != nil {
		return "", errors.New("error of parsing dstart to time value")
	}
	params := strings.Split(repeat, " ")
	years := 0
	days := 0
	switch params[0] {
	case "d":

		if len(params) < 2 {
			return "", errors.New("invalid format of repeat rule")
		}
		days, err = strconv.Atoi(params[1])
		if err != nil {
			return "", errors.New("invalid quanity of days in rule")
		}
		if days > 400 {
			return "", errors.New("invalid quanity of days in rule, nust be less than 400")
		}
	case "y":
		years = 1
	default:
		return "", errors.New("invalid format of rule")
	}
	for {
		date = date.AddDate(years, 0, days)
		if AfterNow(date, now) {
			break
		}
	}
	return date.Format(DATEFORMAT), nil
}

func AfterNow(date, now time.Time) bool {
	nowStr := now.Format(DATEFORMAT)
	dateStr := date.Format(DATEFORMAT)
	if dateStr > nowStr {
		return true
	}
	return false
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if len(task.Date) == 0 {
		task.Date = now.Format(DATEFORMAT)
	}
	t, err := time.Parse(DATEFORMAT, task.Date)
	if err != nil {
		return err
	}
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}
	if AfterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(DATEFORMAT)
		} else {
			task.Date = next
		}
	}
	return nil
}

func CheckID(id string) bool {
	_, err := strconv.Atoi(id)
	if id == "" {
		return false
	}
	if err != nil {
		return false
	}
	return true
}
