package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/O-Tempora/SberIT/internal/models"
	"github.com/go-chi/chi/v5"
)

// To implement http.Handler
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
	s.Router.Route("/tasks", func(r chi.Router) {
		r.Get("/{id}", s.handleGet)
		r.Get("/", s.handleGetList)
		r.Get("/byDate/{year}-{month}-{day}", s.handleGetByDate)
		r.Post("/", s.handleCreateTask)
		r.Put("/{id}", s.handleUpdate)
		r.Delete("/{id}", s.handleDelete)
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

func (s *Server) handleGetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var err error
	var tasks []models.Task
	var page, take *int
	var done *bool

	buf_done, err := strconv.ParseBool(r.URL.Query().Get("done"))
	if err == nil {
		done = &buf_done
	}
	buf_page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err == nil {
		buf_take, err := strconv.Atoi(r.URL.Query().Get("take"))
		if err == nil {
			page = &buf_page
			take = &buf_take
		}
	}
	tasks, err = s.Service.GetList(page, take, done)
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, tasks, nil)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}
	task, err := s.Service.Get(id)
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, task, nil)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}
	if err := s.Service.Delete(id); err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, nil, nil)
}

func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}
	req := models.Task{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}
	if err := s.Service.Update(id, req); err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, nil, nil)
}

func (s *Server) handleGetByDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error
	var done bool
	var tasks []models.Task

	date, err := getDateFromURL(r)
	if err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}

	// if status was not set - get all by date
	if r.URL.Query().Get("done") == "" {
		tasks, err = s.Service.GetByDateAndStatus(*date, false, false)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, tasks, nil)
		return
	}
	// if status was set - get all by date and status
	done, err = strconv.ParseBool(r.URL.Query().Get("done"))
	if err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}
	tasks, err = s.Service.GetByDateAndStatus(*date, done, true)
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, tasks, nil)
}
