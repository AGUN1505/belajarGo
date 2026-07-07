package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type RunRequest struct {
	Code string `json:"code"`
}

type RunResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	})

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RunRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create temp file
		tmpDir, err := os.MkdirTemp("", "go-run-*")
		if err != nil {
			json.NewEncoder(w).Encode(RunResponse{Error: "Failed to create temp dir: " + err.Error()})
			return
		}
		defer os.RemoveAll(tmpDir)

		tmpFile := filepath.Join(tmpDir, "main.go")
		if err := os.WriteFile(tmpFile, []byte(req.Code), 0644); err != nil {
			json.NewEncoder(w).Encode(RunResponse{Error: "Failed to write temp file: " + err.Error()})
			return
		}

		// Run go run main.go
		cmd := exec.Command("go", "run", tmpFile)
		out, err := cmd.CombinedOutput()
		
		var res RunResponse
		if err != nil {
			res.Error = string(out)
		} else {
			res.Output = string(out)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})

	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
