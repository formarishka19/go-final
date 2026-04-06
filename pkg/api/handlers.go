package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-final/pkg/db"
	"net/http"
	"time"
)

const contentType = "application/json; charset=UTF-8"

type TaskListResp struct {
	Tasks []*db.Task `json:"tasks"`
}

type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func WriteJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", contentType)
	if err, ok := data.(error); ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
	} else {
		if data != "" {
			json.NewEncoder(w).Encode(Response{ID: fmt.Sprint(data)})
		} else {
			json.NewEncoder(w).Encode(Response{})
		}

	}
}

func ListTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests allowed", http.StatusMethodNotAllowed)
		return
	}
	tasks, err := db.Tasks(50)
	if err != nil {
		WriteJson(w, err)
		return
	}
	w.Header().Set("Content-Type", contentType)
	json.NewEncoder(w).Encode(TaskListResp{
		Tasks: tasks,
	})
}

func ViewTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if taskID == "" {
		err := errors.New("no ID in request")
		WriteJson(w, err)
		return
	}
	t := db.Task{ID: taskID}
	err := db.GetTask(&t)
	if err != nil {
		WriteJson(w, err)
		return
	}
	w.Header().Set("Content-Type", contentType)
	json.NewEncoder(w).Encode(t)
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}
	taskID := r.FormValue("id")
	if !CheckID(taskID) {
		err := fmt.Errorf("no task ID in request")
		WriteJson(w, err)
		return
	}
	t := db.Task{ID: taskID}
	err := db.GetTask(&t)
	if err != nil {
		WriteJson(w, err)
		return
	}
	if t.Repeat == "" {
		err := db.DeleteTask(&t)
		if err != nil {
			WriteJson(w, err)
			return
		}
		WriteJson(w, "")
		return
	}
	now := time.Now()
	nextDate, err := NextDate(now, t.Date, t.Repeat)
	if err != nil {
		err = fmt.Errorf("error nextdate calculating %w", err)
		WriteJson(w, err)
		return
	}
	t.Date = nextDate
	err = db.UpdateTask(&t)
	if err != nil {
		err = fmt.Errorf("update date in completed task %w", err)
		WriteJson(w, err)
		return
	}
	WriteJson(w, "")
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("id")
	if !CheckID(taskID) {
		err := fmt.Errorf("no task ID in request")
		WriteJson(w, err)
		return
	}
	t := db.Task{ID: taskID}
	err := db.DeleteTask(&t)
	if err != nil {
		WriteJson(w, err)
		return
	}
	WriteJson(w, "")
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		WriteJson(w, err)
		return
	}
	if !CheckID(t.ID) {
		err := fmt.Errorf("invalid task id")
		WriteJson(w, err)
		return
	}
	if len(t.Title) == 0 {
		err := fmt.Errorf("empty title in request")
		WriteJson(w, err)
		return
	}
	err = checkDate(&t)
	if err != nil {
		WriteJson(w, err)
		return
	}
	err = db.UpdateTask(&t)
	if err != nil {
		WriteJson(w, err)
		return
	}
	WriteJson(w, "")
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		WriteJson(w, err)
		return
	}
	if len(t.Title) == 0 {
		err := fmt.Errorf("empty title in request")
		WriteJson(w, err)
		return
	}
	err = checkDate(&t)
	if err != nil {
		WriteJson(w, err)
		return
	}
	taskId, err := db.AddTask(&t)
	if err != nil {
		WriteJson(w, err)
		return
	}
	WriteJson(w, taskId)
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests allowed", http.StatusMethodNotAllowed)
		return
	}
	now, err := time.Parse(DATEFORMAT, r.FormValue("now"))
	if err != nil {
		now = time.Now()
	}
	nDate, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, "Server error in next date calculating", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nDate))
}
