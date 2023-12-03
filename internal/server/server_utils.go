package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func getDateFromURL(r *http.Request) (*time.Time, error) {
	year, err := strconv.Atoi(chi.URLParam(r, "year"))
	if err != nil {
		return nil, err
	}
	month, err := strconv.Atoi(chi.URLParam(r, "month"))
	if err != nil {
		return nil, err
	}
	day, err := strconv.Atoi(chi.URLParam(r, "day"))
	if err != nil {
		return nil, err
	}
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return &date, nil
}
