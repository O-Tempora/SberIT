package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/O-Tempora/SberIT/config"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Server struct {
	Config *config.Config
	Db     *sqlx.DB
	Logger zerolog.Logger
	Router *chi.Mux
}

func InitServer(cf *config.Config) *Server {
	return &Server{
		Config: cf,
		Logger: zerolog.New(os.Stdout),
	}
}

func (s *Server) WithDb(host, name string, port int) *Server {
	db, err := openDb(host, name, port)
	if err != nil {
		log.Fatal(err.Error())
	}
	s.Db = db
	return s
}

func (s *Server) WithLogger(srcs ...io.Writer) *Server {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        io.MultiWriter(srcs...),
		NoColor:    false,
		TimeFormat: time.ANSIC,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatTimestamp: func(i interface{}) string {
			t, _ := time.Parse(time.RFC3339, fmt.Sprintf("%s", i))
			return t.Format(time.RFC1123)
		},
	}).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	s.Logger = logger
	return s
}

func openDb(host, name string, port int) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("postgres://postgres:postgres@%s:%d/%s?sslmode=disable", host, port, name)
	db, err := sqlx.Open(
		"postgres",
		connStr,
	)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
