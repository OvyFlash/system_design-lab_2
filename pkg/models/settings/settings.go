package settings

import (
	"lab_2/config"
	"lab_2/internal/databases/postgres"
	"lab_2/internal/databases/redis"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

type Essential struct {
	Title         string
	Storage       *postgres.Storage
	Redis         *redis.Redis
	Config        config.Config
	ctx           context.Context
	ServerEnabled bool
	Server        *http.Server
	Router        *httprouter.Router
	zerolog.Logger
}

func NewEssential(c config.Config, title string, httpServer bool) (*Essential, error) {
	e := &Essential{
		Title:         title,
		Config:        c,
		ctx:           context.Background(),
		ServerEnabled: httpServer,
		Server:        nil,
		Router:        httprouter.New(),
		Logger:        config.GetLogger(title),
	}

	err := e.initDatabases()
	return e, err
}

func (e *Essential) initDatabases() (err error) {
	e.Storage, err = postgres.NewStorage(e.Config.Databases.SQLConfig)
	if err != nil {
		return err
	}
	e.Redis, err = redis.NewRedis(e.Config.Databases.RedisConfig)
	if err != nil {
		return err
	}
	return nil
}

func (e *Essential) GetContext() context.Context {
	if e.ctx == nil {
		e.ctx = context.Background()
	}
	return e.ctx
}
