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

func (s *Service) Create(task models.Task) (int, error) {
	if task.Deadline.Before(time.Now()) {
		task.Deadline = time.Now().Add(24 * time.Hour)
	}

	var id int
	err := s.Db.Get(&id, `insert into tasks
		(header, description, deadline, done)
		values ($1, $2, $3, $4)
		returning id`,
		task.Header, task.Description, task.Deadline, task.Done)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *Service) GetList(page, take *int, done *bool) ([]models.Task, error) {
	var tasks []models.Task
	var err error

	// if pagination was not set
	if page == nil || take == nil {
		err = s.Db.Select(
			&tasks,
			`select * from tasks`,
		)
	} else {
		// if pagination was set with status
		if done != nil {
			err = s.Db.Select(
				&tasks,
				`select * from tasks where done = $1 limit $2 offset $3`,
				*done, *take, *take*(*page-1),
			)
		} else { // if pagination was set without status
			err = s.Db.Select(
				&tasks,
				`select * from tasks limit $1 offset $2`,
				*take, *take*(*page-1),
			)
		}
	}
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) Get(id int) (*models.Task, error) {
	var task models.Task
	if err := s.Db.Get(&task, `select * from tasks where id = $1`, id); err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *Service) Delete(id int) error {
	if _, err := s.Db.Exec(`delete from tasks where id = $1`, id); err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(id int, task models.Task) error {
	if task.Deadline.Before(time.Now()) {
		return errors.New("task deadline can not be earlier than today")
	}
	if _, err := s.Db.Exec(`update tasks set header=$1, description=$2, deadline=$3, done=$4 where id = $5`,
		task.Header, task.Description, task.Deadline, task.Done, id); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetByDateAndStatus(date time.Time, done, statusWasSet bool) ([]models.Task, error) {
	var err error
	var tasks []models.Task

	if statusWasSet {
		err = s.Db.Select(&tasks, `select * from tasks where deadline = $1 and done = $2`, date, done)
	} else {
		err = s.Db.Select(&tasks, `select * from tasks where deadline = $1`, date)
	}

	if err != nil {
		return nil, err
	}
	return tasks, nil
}
