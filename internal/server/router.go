package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/O-Tempora/SberIT/internal/models"
	"github.com/go-chi/chi/v5"

	_ "github.com/O-Tempora/SberIT/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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
	s.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", s.Config.Port)), //The url pointing to API definition
	))

	s.Router.Route("/tasks", func(r chi.Router) {
		r.Get("/{id}", s.handleGet)
		r.Get("/", s.handleGetList)
		r.Get("/byDate/{year}-{month}-{day}", s.handleGetByDate)
		r.Post("/", s.handleCreateTask)
		r.Put("/{id}", s.handleUpdate)
		r.Delete("/{id}", s.handleDelete)
	})
}

// CreateTask godoc
//
//	@Summary		Create task
//	@Description	Creates task with fields in body param and returns inserted id if successfull
//	@Tags			Create
//	@Accept			json
//	@Produce		json
//	@Param			task	body	models.Task	true	"Task data"
//	@Router			/tasks [post]
//	@Success		200	{integer}		Id
//	@Failure		400	{string}	error
//	@Failure		500	{string}	error
func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := models.Task{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}
	id, err := s.Service.Create(req)
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	s.respond(w, r, http.StatusCreated, id, nil)
}

// GetList godoc
//
//	@Summary		Get task list
//	@Description	Returns list of tasks with optional pagination (page + take) and optional filter by status (done)
//	@Tags			GetList
//	@Accept			json
//	@Produce		json
//	@Param			done	query	bool	false	"Task status"
//	@Param			page	query	int		false	"Page number"
//	@Param			take	query	int		false	"Page size"
//	@Router			/tasks [get]
//	@Success		200	{array}		models.Task
//	@Failure		400	{string}	error
//	@Failure		500	{string}	error
func (s *Server) handleGetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var err error
	var tasks []models.Task
	var done *bool

	// See if "done" parameter was set (optional)
	if r.URL.Query().Get("done") != "" {
		buf_done, err := strconv.ParseBool(r.URL.Query().Get("done"))
		if err != nil {
			s.respond(w, r, http.StatusBadRequest, nil, err)
			return
		}
		done = &buf_done
	}

	// Pagination parameters were not set, so selecting tasks only by status
	if r.URL.Query().Get("page") == "" && r.URL.Query().Get("take") == "" {
		tasks, err = s.Service.GetList(done)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, nil, err)
			return
		}
		s.respond(w, r, http.StatusOK, tasks, nil)
		return
	}

	// Any of pagination parameters were set
	page, pageErr := strconv.Atoi(r.URL.Query().Get("page"))
	take, takeErr := strconv.Atoi(r.URL.Query().Get("take"))
	if err = errors.Join(takeErr, pageErr); err != nil {
		s.respond(w, r, http.StatusBadRequest, nil, err)
		return
	}

	tasks, err = s.Service.GetListWithPagination(page, take, done)
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, nil, err)
		return
	}
	s.respond(w, r, http.StatusOK, tasks, nil)
	return
}

// GetTask godoc
//
//	@Summary		Get task by id
//	@Description	Returns task with id from id path vparam. Returns error if no task with such id exists
//	@Tags			Get
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Task id"
//	@Router			/tasks/{id} [get]
//	@Success		200	{object}	models.Task
//	@Failure		400	{string}	error
//	@Failure		500	{string}	error
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

// DeleteTask godoc
//
//	@Summary		Delete task by id
//	@Description	Deletes task by id from id path param
//	@Tags			Delete
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Task id"
//	@Router			/tasks/{id} [delete]
//	@Success		200
//	@Failure		400	{string}	error
//	@Failure		500	{string}	error
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

// UpdateTask godoc
//
//	@Summary		Update task
//	@Description	Updates task using data from body and with id from path param. If some fields of body struct are omitted, they will be overwritten by default values
//	@Tags			Update
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int			true	"Task id"
//	@Param			task	body	models.Task	true	"Task data"
//	@Router			/tasks/{id} [put]
//	@Success		200
//	@Failure		400	{string}	error
//	@Failure		500	{string}	error
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

// GetByDate godoc
//
//	@Summary		Get tasks by date
//	@Description	Returns tasks by date from path params and optional filter by status (done)
//	@Tags			GetList
//	@Accept			json
//	@Produce		json
//	@Param			year	path	int		true	"Year"
//	@Param			month	path	int		true	"Month"
//	@Param			day		path	int		true	"Day"
//	@Param			done	query	bool	false	"Task status"
//	@Router			/tasks/byDate/{year}-{month}-{day} [get]
//	@Success		200	{array}		models.Task
//	@Failure		400	{string}	error
//	@Failure		500	{string}	error
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
