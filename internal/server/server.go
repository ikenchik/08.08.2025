package server

import (
	"downloader/internal/config"
	"downloader/internal/handlers"
	"downloader/internal/models"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func NewServer(newConfig config.Config) {
	app := models.NewAppState(newConfig)
	r := mux.NewRouter()

	r.HandleFunc("/tasks", handlers.CreateTaskHandler(app)).Methods("POST")

	r.HandleFunc("/tasks/{id}/urls", handlers.AddURLHandler(app)).Methods("POST")

	r.HandleFunc("/tasks/{id}", handlers.GetTaskHandler(app)).Methods("GET")

	r.HandleFunc("/download/{filename}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := vars["filename"]
		filePath := filepath.Join(".", filename)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/zip")
		http.ServeFile(w, r, filePath)

	}).Methods("GET")

	log.Printf("Server starting on port %s...", newConfig.Port)
	log.Fatal(http.ListenAndServe(":"+newConfig.Port, r))
}
