package server

import (
	"fmt"
	"strings"

	"github.com/amahdian/ai-assistant-be/svc/auth"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/amahdian/ai-assistant-be/server/router"

	"github.com/amahdian/ai-assistant-be/global/env"
	"github.com/amahdian/ai-assistant-be/pkg/logger"
	"github.com/amahdian/ai-assistant-be/storage"
	"github.com/amahdian/ai-assistant-be/storage/pg"
	"github.com/amahdian/ai-assistant-be/svc"
	"github.com/pkg/errors"
)

type Server struct {
	Envs *env.Envs

	Authenticator auth.Authenticator
	Storage       storage.Storage
	Svc           svc.Svc
	Router        *router.Router
}

func NewServer(envs *env.Envs) (*Server, error) {
	s := &Server{
		Envs: envs,
	}
	if err := s.setupLogger(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize logger")
	}
	if err := s.migrateDb(); err != nil {
		return nil, errors.Wrap(err, "failed to migrate the db")
	}
	if err := s.setupStorage(); err != nil {
		return nil, err
	}
	if err := s.setupAuthenticator(); err != nil {
		return nil, errors.Wrap(err, "failed to setup authenticator")
	}
	s.setupServices()
	s.setupRouter()
	return s, nil
}

func (s *Server) Run() (err error) {
	defer func(s *Server) {
		err := s.Close()
		if err != nil {
			logger.Errorf("failed to gracefully shutdown the server and release resources: %v", err)
		}
	}(s)

	err = s.Router.Run(fmt.Sprintf(":%s", s.Envs.Server.HttpPort))
	return err
}

func (s *Server) Close() error {
	if err := logger.Close(); err != nil {
		logger.Errorf("failed to close/sync the logger: %v", err) // can it actually log itself?
		return err
	}
	return nil
}

func (s *Server) setupLogger() error {
	logger.ConfigureFromEnvs(s.Envs)
	return nil
}

func (s *Server) migrateDb() error {
	err := pg.EnsureDatabaseExists(s.Envs.Db.Dsn)
	if err != nil {
		return errors.Wrap(err, "failed to create database")
	}
	migrationsDir := fmt.Sprintf("file://%s/migrations", s.Envs.Server.AssetsDir)
	migrator, err := migrate.New(migrationsDir, s.Envs.Db.Dsn)
	if err != nil {
		return errors.Wrap(err, "failed to open new migration instance")
	}

	err = migrator.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("all migrations are already applied")
			return nil
		}
		return errors.Wrap(err, "failed to run migrations")
	}

	logger.Info("applied migrations to the db")

	return nil
}

func (s *Server) setupStorage() error {
	logLevelEnv := strings.ToLower(s.Envs.Db.LogLevel)
	logLevel := pg.LogLevel(logLevelEnv)
	db, err := pg.OpenGormDb(s.Envs.Db.Dsn, logLevel)
	if err != nil {
		return errors.Wrap(err, "failed to open gorm connection")
	}
	s.Storage = pg.NewStg(db)
	return nil
}

func (s *Server) setupServices() {
	s.Svc = svc.NewSvc(s.Storage, s.Envs)
}

func (s *Server) setupRouter() {
	s.Router = router.NewRouter(
		s.Storage,
		s.Svc,
		s.Envs,
		s.Authenticator)
}

func (s *Server) setupAuthenticator() error {
	s.Authenticator = auth.NewAuthenticator(s.Envs)
	return nil
}
