package service

import "errors"

var (
	errInvalidDeadline = errors.New("task deadline can not be earlier than today")
)
