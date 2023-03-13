package api

import (
	"lab_2/pkg/models/settings"
	"lab_2/pkg/services/api/handler"
	"lab_2/pkg/services/api/service"
)

type API struct {
	*settings.Essential
	handler *handler.Handler
	service *service.Service
}

func NewAPI(e *settings.Essential) *API {
	service := service.NewService(e)
	return &API{
		handler: handler.NewHandler(e, service),
		service: service,
	}
}

func (l *API) Start() error {
	if err := l.service.Start(l.handler); err != nil {
		return err
	}
	//set http services
	return nil
}

func (l *API) Stop() error {
	return l.service.Stop()
}
