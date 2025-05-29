package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gabrielnov/tcc-api/internal/dto"
	"github.com/gabrielnov/tcc-api/internal/service"
)

const DIRECTORY = "./scanned_code"

type ResponsePayload struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req dto.RequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Print(err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	codeAnalysisService := service.NewCodeAnalysisService(service.NewLlmService(), service.NewBanditService(), service.NewFileManager())

	response, err := codeAnalysisService.Run(req)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Failed to run Bandit analysis", http.StatusInternalServerError)
		return
	}

	respJson, err := json.Marshal(response)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Failed to marshal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func main() {
	http.HandleFunc("/scan", scanHandler)

	log.Print("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
