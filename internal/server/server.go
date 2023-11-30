package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/O-Tempora/SberIT/config"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type Server struct {
	Config *config.Config
	Db     *pgx.Conn
	Logger *zerolog.Logger
	Router *chi.Mux
}

func InitServer(cf *config.Config) *Server {
	return &Server{
		Config: cf,
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
	s.Logger = &logger
	return s
}

func openDb(host, name string, port int) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := pgx.Connect(
		ctx,
		fmt.Sprintf("postgres://postgres:postgres@%s:%d/%s", host, port, name),
	)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
