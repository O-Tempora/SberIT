package service

import (
	"log"
	"testing"
	"time"

	"github.com/O-Tempora/SberIT/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var service Service

func TestMain(m *testing.M) {
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:6969/sber?sslmode=disable")
	if err != nil {
		log.Fatal(err.Error())
	}
	service.Db = db

	_, err = db.Exec(
		`create table tasks(
			id serial4 PRIMARY KEY NOT NULL,
			header text,
			description text,
			deadline date,
			done bool
		);
		insert into tasks
		(header, description, deadline, done)
		values
		('Header1', 'Description1', '2023-12-04', false),
		('Header2', 'Description2', '2023-12-04', true),
		('Header3', 'Description3', '2023-12-05', true);
	`)
	if err != nil {
		log.Fatal(err.Error())
	}
	m.Run()
}

func TestGetOne(t *testing.T) {
	var test_cases = []struct {
		id       int
		err      error
		expected models.Task
		actual   models.Task
	}{
		{
			id: 1,
			expected: models.Task{
				Id:          1,
				Header:      "Header1",
				Description: "Description1",
				Deadline:    time.Date(2023, time.Month(12), 4, 0, 0, 0, 0, time.Local),
				Done:        false,
			},
		},
		{
			id: 2,
			expected: models.Task{
				Id:          2,
				Header:      "Header2",
				Description: "Description2",
				Deadline:    time.Date(2023, time.Month(12), 4, 0, 0, 0, 0, time.Local),
				Done:        true,
			},
		},
		{
			id: 3,
			expected: models.Task{
				Id:          3,
				Header:      "Header3",
				Description: "Description3",
				Deadline:    time.Date(2023, time.Month(12), 5, 0, 0, 0, 0, time.Local),
				Done:        true,
			},
		},
	}

	for _, tc := range test_cases {
		task, err := service.Get(tc.id)
		if assert.Nil(t, err) {
			tc.actual = *task
			assert.Equal(t, tc.expected.Id, tc.actual.Id)
			assert.Equal(t, tc.expected.Header, tc.actual.Header)
			assert.Equal(t, tc.expected.Description, tc.actual.Description)
			assert.Equal(t, tc.expected.Done, tc.actual.Done)
		}
	}
}

func TestGetByDate(t *testing.T) {
	var test_cases = []struct {
		date     time.Time
		done     bool
		wasSet   bool
		err      error
		expected []models.Task
	}{
		{
			done:   true,
			wasSet: true,
			date:   time.Date(2023, time.Month(12), 4, 0, 0, 0, 0, time.Local),
			expected: []models.Task{
				{
					Id:          2,
					Header:      "Header2",
					Description: "Description2",
					Deadline:    time.Date(2023, time.Month(12), 4, 0, 0, 0, 0, time.Local),
					Done:        true,
				},
			},
		},
		{
			done:     false,
			wasSet:   true,
			date:     time.Date(2023, time.Month(12), 5, 0, 0, 0, 0, time.Local),
			expected: nil,
		},
		{
			done:   false,
			wasSet: false,
			date:   time.Date(2023, time.Month(12), 4, 0, 0, 0, 0, time.Local),
			expected: []models.Task{
				{
					Id:          1,
					Header:      "Header1",
					Description: "Description1",
					Deadline:    time.Date(2023, time.Month(12), 4, 0, 0, 0, 0, time.Local),
					Done:        false,
				},
				{
					Id:          2,
					Header:      "Header2",
					Description: "Description2",
					Deadline:    time.Date(2023, time.Month(12), 4, 0, 0, 0, 0, time.Local),
					Done:        true,
				},
			},
		},
	}

	for _, tc := range test_cases {
		actual, err := service.GetByDateAndStatus(tc.date, tc.done, tc.wasSet)
		if assert.Nil(t, err) {
			assert.Equal(t, len(tc.expected), len(actual))
		}
	}
}

func TestGetList(t *testing.T) {
	var fl bool = false
	var tr bool = true
	var test_cases = []struct {
		done            *bool
		expected_length int
	}{
		{
			done:            nil,
			expected_length: 3,
		},
		{
			done:            &fl,
			expected_length: 1,
		},
		{
			done:            &tr,
			expected_length: 2,
		},
	}

	for _, tc := range test_cases {
		tasks, err := service.GetList(tc.done)
		if assert.Nil(t, err) {
			assert.Equal(t, tc.expected_length, len(tasks))
		}
	}
}

func TestGetListWithPagination(t *testing.T) {
	var fl bool = false
	var tr bool = true
	var test_cases = []struct {
		done            *bool
		page            int
		take            int
		expected_length int
	}{
		{
			done:            nil,
			page:            1,
			take:            2,
			expected_length: 2,
		},
		{
			done:            &fl,
			page:            2,
			take:            1,
			expected_length: 0,
		},
		{
			done:            &tr,
			page:            1,
			take:            3,
			expected_length: 2,
		},
	}

	for _, tc := range test_cases {
		tasks, err := service.GetListWithPagination(tc.page, tc.take, tc.done)
		if assert.Nil(t, err) {
			assert.Equal(t, tc.expected_length, len(tasks))
		}
	}
}

func TestUpdateOwerwrites(t *testing.T) {
	var tc = models.Task{
		Id:       1,
		Deadline: time.Now().Add(24 * time.Hour),
	}
	err := service.Update(tc.Id, tc)
	task, _ := service.Get(tc.Id)
	if assert.Nil(t, err) {
		assert.Equal(t, "", task.Header)
		assert.Equal(t, "", task.Description)
		assert.False(t, task.Done)
	}
}

func TestUpdateInvalidDeadline(t *testing.T) {
	var tc = models.Task{
		Id:       1,
		Deadline: time.Now().Add(-24 * time.Hour),
	}
	err := service.Update(tc.Id, tc)
	assert.ErrorIs(t, err, errInvalidDeadline)
}

func TestDelete(t *testing.T) {
	var test_cases = []struct {
		id        int
		remaining int
	}{
		{
			id:        1,
			remaining: 2,
		},
		{
			id:        20,
			remaining: 2,
		},
	}
	for _, tc := range test_cases {
		err := service.Delete(tc.id)
		if assert.Nil(t, err) {
			res, err := service.GetList(nil)
			if assert.Nil(t, err) {
				assert.Equal(t, tc.remaining, len(res))
			}
		}
	}
}

func TestCreate(t *testing.T) {
	var test_cases = []struct {
		task models.Task
		id   int
	}{
		{
			id: 4,
		},
		{
			id: 5,
		},
		{
			id: 6,
		},
	}

	for _, tc := range test_cases {
		id, err := service.Create(tc.task)
		if assert.Nil(t, err) {
			assert.Equal(t, id, tc.id)
		}
	}
}
