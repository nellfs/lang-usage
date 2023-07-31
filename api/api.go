package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nellfs/lang-usage/storage"
	"github.com/nellfs/lang-usage/types"
)

type APIServer struct {
	user       string
	listenAddr string
	storage    storage.Storage
}

func NewAPIServer(user string, listenAddr string, storage storage.Storage) *APIServer {
	return &APIServer{
		user:       user,
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *APIServer) Run() {
	http.HandleFunc("/", apiHandler)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type Response struct {
	Message string `json:"message"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Create a response message
	response := Response{Message: "Hello, this is a simple API!"}

	// Encode the response as JSON and send it
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *APIServer) handleGetCodeReport(w http.ResponseWriter, r *http.Request) (*types.CodeReport, error) {
	return nil, nil
}
