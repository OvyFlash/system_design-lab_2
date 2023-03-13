package service

import (
	"fmt"
	"lab_2/config"
	"lab_2/pkg/models/settings"
	"lab_2/pkg/services/api/models"
	"net/http"
)

type Service struct {
	gatewayServerHTTP *http.Server
	*settings.Essential
	apiTree ServiceRoutingTree
}

func NewService(e *settings.Essential) *Service {
	return &Service{
		Essential: e,
	}
}

func (s *Service) Start(handler http.Handler) error {
	if err := s.init(handler); err != nil {
		return err
	}
	go func() {
		s.Info().Msgf("starting HTTP server at :%d", s.Config.GetPortByTitle(s.Title))
		if err := s.gatewayServerHTTP.ListenAndServe(); err != http.ErrServerClosed {
			s.Error().Msgf("http server: %s", err.Error())
		}
	}()
	return nil
}

func (s *Service) Stop() error {
	return s.gatewayServerHTTP.Close()
}

func (s *Service) init(handler http.Handler) error {
	s.initializeHTTPServer(handler)
	return s.initializePrefixesList()
}

func (s *Service) initializePrefixesList() (err error) {
	s.apiTree = NewServiceRoutingTree()
	for _, serviceItem := range s.Config.Services {
		if len(serviceItem.Prefixes) != 0 {
			foundService, err := s.apiTree.FindByPort(serviceItem.Port)
			if err == nil {
				err = fmt.Errorf(
					"api collision: services '%s' and '%s' both use the same port: %d",
					serviceItem.Title, foundService.Name, serviceItem.Port)
				return err
			}
			err = s.addServiceAPI(&serviceItem)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (s *Service) initializeHTTPServer(handler http.Handler) {
	s.gatewayServerHTTP = &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", s.Config.GetPortByTitle(s.Title)),
	}
}

func (s *Service) addServiceAPI(serviceItem *config.Service) (err error) {
	for _, prefix := range serviceItem.Prefixes {
		err = s.apiTree.Add(
			prefix,
			models.ProxyService{
				Name: serviceItem.Title,
				Port: serviceItem.Port,
			},
		)
		if err != nil {
			return
		}
	}
	return
}

func (s *Service) GetServiceWithEndpoint(endpoint string) (service models.ProxyService, err error) {
	service, err = s.apiTree.Find(endpoint)
	if err != nil {
		err = fmt.Errorf("service not found for endpoint '%s'", endpoint)
	}
	return
}
