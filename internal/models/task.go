package models

import "time"

type Task struct {
	Id          int        `json:"id,omitempty"`
	Header      string     `json:"header,omitempty"`
	Description string     `json:"description,omitempty"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	Done        bool       `json:"done,omitempty"`
}
