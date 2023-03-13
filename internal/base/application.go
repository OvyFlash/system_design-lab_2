package base

import (
	"context"
	"lab_2/config"
	"lab_2/pkg/models/settings"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type Application[T Servicer] struct {
	*settings.Essential
	service T

	ctx    context.Context
	cancel context.CancelFunc
}

func NewApplication[T Servicer](newService func(*settings.Essential) T, title string, httpServer bool) (*Application[T], error) {
	e, err := settings.NewEssential(config.NewConfig(), title, httpServer)
	if err != nil {
		return nil, err
	}
	a := &Application[T]{
		Essential: e,
		service:   newService(e),
	}
	return a, nil
}

func (a *Application[T]) Start() error {
	if err := a.init(); err != nil {
		return err
	}
	if err := a.service.Start(); err != nil {
		return err
	}
	a.startHTTPServer()

	<-a.gracefulExit().Done()

	return a.Stop()
}

func (a *Application[T]) Stop() error {
	if err := a.service.Stop(); err != nil {
		return err
	}
	return a.stopHTTPServer()
}

func (a *Application[T]) init() error {
	a.ctx, a.cancel = context.WithCancel(a.GetContext())
	a.initializeHTTPServer()
	return nil
}

func (a *Application[T]) initializeHTTPServer() {
	if !a.ServerEnabled {
		return
	}
	// Initialize HTTP router
	a.Router = &httprouter.Router{
		RedirectTrailingSlash: true, RedirectFixedPath: true,
		HandleMethodNotAllowed: true, HandleOPTIONS: true,
	}
	a.Router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Headers", "*")
			header.Set("Access-Control-Expose-Headers", "*")
		}
		w.WriteHeader(http.StatusNoContent)
	})
	// Initialize HTTP server
	a.Server = &http.Server{
		Handler: loggingHandler{a.Router, a.Logger},
		Addr:    ":" + strconv.Itoa(int(a.Config.GetPortByTitle(a.Title))),
	}
}

func (a *Application[T]) startHTTPServer() {
	if !a.ServerEnabled {
		return
	}
	a.Info().Msgf("starting HTTP server at :%d", a.Config.GetPortByTitle(a.Title))
	go func() {
		a.Info().Msgf("%v", a.Router)
		if err := a.Server.ListenAndServe(); err != http.ErrServerClosed {
			a.Error().Msgf("http server:", err)
		}
	}()
}

func (a *Application[T]) stopHTTPServer() error {
	if !a.ServerEnabled {
		return nil
	}
	a.Info().Msgf("stopping HTTP server at :%d", a.Config.GetPortByTitle(a.Title))
	return a.Server.Shutdown(a.GetContext())
}

func (a *Application[T]) gracefulExit() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		log.Info().Msgf("system call:%+v", <-c)
		a.cancel()
	}()
	return a.ctx
}
