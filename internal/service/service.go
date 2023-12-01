package service

import (
	"errors"
	"time"

	"github.com/O-Tempora/SberIT/internal/models"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	Db *sqlx.DB
}

func (s *Service) Create(task models.Task) error {
	if task.Deadline.Before(time.Now()) {
		return errors.New("task deadline must be before current time")
	}
	_, err := s.Db.Exec(`insert into tasks (
		(header, description, deadline, done)
		values ($1, $2, $3, $4)
	)`, task.Header, task.Description, task.Deadline, task.Done)
	return err
}
