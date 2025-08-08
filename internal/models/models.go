package models

import (
	"archive/zip"
	"bytes"
	"downloader/internal/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	StatusPending    = "PENDING"
	StatusProcessing = "PROCESSING"
	StatusCompleted  = "COMPLETED"
	StatusFailed     = "FAILED"
)

type Task struct {
	ID        string            `json:"id"`
	Status    string            `json:"status"`
	URLs      []string          `json:"urls"`
	Errors    map[string]string `json:"errors,omitempty"`
	Archive   *string           `json:"archive,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

type AppState struct {
	mu          sync.Mutex
	Tasks       map[string]*Task `json:"tasks"`
	ActiveTasks int              `json:"active_tasks"`
	Config      config.Config    `json:"config"`
}

func NewAppState(cfg config.Config) *AppState {
	return &AppState{
		Tasks:  make(map[string]*Task),
		Config: cfg,
	}
}

func (app *AppState) CreateTask() (*Task, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.ActiveTasks >= app.Config.MaxTasks {
		return nil, fmt.Errorf("server is busy")
	}

	taskID := generateID()
	task := &Task{
		ID:        taskID,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
	app.Tasks[taskID] = task
	app.ActiveTasks++

	return task, nil
}

func (app *AppState) AddURL(taskID, url string) error {
	app.mu.Lock()
	defer app.mu.Unlock()
	task, exists := app.Tasks[taskID]

	if !exists {
		return fmt.Errorf("task not found")
	}

	if len(task.URLs) >= app.Config.MaxFiles {
		return fmt.Errorf("maximum files reached")
	}

	ext := strings.ToLower(filepath.Ext(url))
	if !isAllowed(ext, app.Config.AllowedExts) {
		return fmt.Errorf("invalid file type")
	}

	task.URLs = append(task.URLs, url)

	if len(task.URLs) == app.Config.MaxFiles {
		go app.processTask(task)
	}

	return nil
}

func (app *AppState) GetTask(taskID string) (*Task, error) {
	app.mu.Lock()
	defer app.mu.Unlock()
	task, exists := app.Tasks[taskID]

	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	return task, nil
}

func (app *AppState) processTask(task *Task) {
	app.mu.Lock()
	task.Status = StatusProcessing
	app.mu.Unlock()

	archiveFileName := fmt.Sprintf("%s.zip", task.ID)
	errs := make(map[string]string)
	archiveFile, err := os.Create(archiveFileName)

	if err != nil {
		app.handleTaskError(task, err)
		return
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, url := range task.URLs {
		if err := downloadAndAddToZip(zipWriter, url); err != nil {
			errs[url] = err.Error()
		}
	}

	app.mu.Lock()
	defer app.mu.Unlock()

	task.Status = StatusCompleted
	task.Errors = errs

	archivePath := "/download/" + archiveFileName
	task.Archive = &archivePath
	app.ActiveTasks--

	fmt.Printf("Task %s completed!\nlink: http://localhost:%s%s\n\n", task.ID, app.Config.Port, archivePath)
}

func (app *AppState) handleTaskError(task *Task, err error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	task.Status = StatusFailed
	task.Errors = map[string]string{"system": err.Error()}
	app.ActiveTasks--
}

func downloadAndAddToZip(zipWriter *zip.Writer, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	filename := filepath.Base(url)
	entry, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(entry, bytes.NewReader(body))
	return err
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func isAllowed(ext string, allowed []string) bool {
	for _, a := range allowed {
		if "."+a == ext {
			return true
		}
	}
	return false
}
