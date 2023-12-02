package models

import "time"

type Task struct {
	Id          int       `json:"id"`
	Header      string    `json:"header"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Done        bool      `json:"done"`
}
