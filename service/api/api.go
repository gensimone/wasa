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

	RootMedia string
	Media     string

	DefaultUserPhoto  string
	DefaultGroupPhoto string
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
	if cfg.RootMedia == "" {
		return nil, errors.New("root directory is required")
	}
	if cfg.Media == "" {
		return nil, errors.New("media directory is required")
	}
	if cfg.DefaultUserPhoto == "" {
		return nil, errors.New("default user photo is required")
	}
	if cfg.DefaultGroupPhoto == "" {
		return nil, errors.New("default group photo is required")
	}

	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router:            router,
		baseLogger:        cfg.Logger,
		db:                cfg.Database,
		rootMedia:         cfg.RootMedia,
		media:             cfg.Media,
		defaultUserPhoto:  cfg.DefaultUserPhoto,
		defaultGroupPhoto: cfg.DefaultGroupPhoto,
	}, nil
}

type _router struct {
	router     *httprouter.Router
	baseLogger logrus.FieldLogger
	db         database.AppDatabase

	// Various path used internally by some API handlers.
	rootMedia string
	media     string

	// Default user and group photos.
	defaultUserPhoto  string
	defaultGroupPhoto string
}
