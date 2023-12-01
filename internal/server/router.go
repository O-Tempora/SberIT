package server

import (
	"encoding/json"
	"net/http"

	"github.com/O-Tempora/SberIT/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// To implement interface http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}, err error) {
	w.WriteHeader(code)
	if err != nil {
		response := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(response)
		s.Logger.Error().Msgf("Resonse: method  %s, URL  %s, code  %d %s, error  %s",
			r.Method, r.URL, code, http.StatusText(code), err.Error())
		return
	}

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
	s.Logger.Info().Msgf("Response: method  %s, URL  %s, Code  %d %s",
		r.Method, r.URL, code, http.StatusText(code))
}

func (s *Server) InitRouter() {
	cors := cors.AllowAll()
	s.Router.Use(cors.Handler)
	s.Router.Route("/task", func(r chi.Router) {
		r.Get("/{id}", nil)
		r.Get("/", s.handleGetAll)
		r.Post("/", s.handleCreateTask)
		r.Put("/{id}", nil)
		r.Delete("/{id}", nil)
	})
}

func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := models.Task{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}
	if err := s.Service.Create(req); err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	s.respond(w, r, http.StatusCreated, nil, nil)
}

func (s *Server) handleGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("HELP"))
	return
}
