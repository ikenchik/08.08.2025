package handlers

import (
	"downloader/internal/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type response struct {
	ID      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateTaskHandler(app *models.AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := app.CreateTask()

		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response{ID: task.ID})
	}
}

func AddURLHandler(app *models.AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskID := vars["id"]
		url := r.FormValue("url")

		if url == "" {
			http.Error(w, "URL required", http.StatusBadRequest)
			return
		}

		if err := app.AddURL(taskID, url); err != nil {
			http.Error(w, err.Error(), getStatusCode(err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetTaskHandler(app *models.AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskID := vars["id"]
		task, err := app.GetTask(taskID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		respTask := *task

		if respTask.Archive != nil {
			scheme := "http"

			if r.TLS != nil {
				scheme = "https"
			}

			baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)
			absoluteURL := baseURL + *respTask.Archive
			respTask.Archive = &absoluteURL
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respTask)
	}
}

func getStatusCode(err error) int {
	switch err.Error() {
	case "task not found":
		return http.StatusNotFound
	case "maximum files reached":
		return http.StatusForbidden
	case "invalid file type":
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
