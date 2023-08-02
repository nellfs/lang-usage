package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nellfs/lang-usage/storage"
)

type APIServer struct {
	user       string
	listenAddr string
	storage    storage.Storage
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSONResponse(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(user string, listenAddr string, storage storage.Storage) *APIServer {
	return &APIServer{
		user:       user,
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *APIServer) Run() {
	http.HandleFunc("/", makeHTTPHandleFunc(s.handleData)) // convert to makehttphandlerfunc

	fmt.Printf("Server is running on http://localhost%s", s.listenAddr)
	err := http.ListenAndServe(s.listenAddr, nil)
	if err != nil {
		panic(err)
	}
}

type Response struct {
	Message string `json:"message"`
}

func (s *APIServer) handleData(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
    return s.handleGetLanguage(w, r)
	}

	return fmt.Errorf("Method not allowed: %s", r.Method)
}

//return all languages data by default
func (s *APIServer) handleGetLanguage(w http.ResponseWriter, r *http.Request) error {
  language := "PotrickLanguage"
  languages, err := s.storage.GetLanguages(&language)
  if err != nil {
    return err
  }

  return WriteJSONResponse(w, http.StatusOK, languages)
}

func (s *APIServer) handleGetCodeReport(w http.ResponseWriter, r *http.Request)  error {
  return nil  
}
