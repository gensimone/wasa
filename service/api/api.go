package api

import (
	"errors"
	"net/http"

	"github.com/gensimone/WASA-project/service/database"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger   logrus.FieldLogger
	Database database.AppDatabase
	Uploads  string
}

type Router interface {
	Handler() http.Handler
	Close() error
}

func New(cfg Config) (Router, error) {
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}
	if cfg.Uploads == "" {
		return nil, errors.New("uploads is required")
	}

	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router:     router,
		baseLogger: cfg.Logger,
		db:         cfg.Database,
		uploads:    cfg.Uploads,
	}, nil
}

type _router struct {
	router     *httprouter.Router
	baseLogger logrus.FieldLogger
	db         database.AppDatabase
	uploads    string
}
