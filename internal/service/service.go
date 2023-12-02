package service

import (
	"time"

	"github.com/O-Tempora/SberIT/internal/models"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	Db *sqlx.DB
}

func (s *Service) Create(task models.Task) error {
	if task.Deadline.Before(time.Now()) {
		task.Deadline = time.Now().Add(24 * time.Hour)
	}
	_, err := s.Db.Exec(`insert into tasks
		(header, description, deadline, done)
		values ($1, $2, $3, $4)`,
		task.Header, task.Description, task.Deadline, task.Done)
	return err
}

func (s *Service) GetAll() ([]models.Task, error) {
	var tasks []models.Task
	if err := s.Db.Select(&tasks, `select * from tasks`); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) Get(id int) (*models.Task, error) {
	var task *models.Task
	if err := s.Db.Select(task, `select * from tasks where id = $1`, id); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) Delete(id int) error {
	if _, err := s.Db.Exec(`delete from tasks where id = $1`, id); err != nil {
		return err
	}
	return nil
}
