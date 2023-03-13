package service

import (
	"bufio"
	"fmt"
	"io"
	"lab_2/pkg/models/settings"
	"lab_2/pkg/services/launcher/models"
	"os/exec"
)

type Service struct {
	*settings.Essential
	serviceControllers []*models.ServiceController
}

func NewService(e *settings.Essential) *Service {
	return &Service{
		Essential: e,
	}
}

func (s *Service) Start() error {
	if err := s.init(); err != nil {
		return err
	}
	return nil
}

func (s *Service) init() error {
	s.initServices()
	return nil
}

func (s *Service) initServices() {
	for _, v := range s.Config.Services {
		if v.Title == s.Title {
			continue
		}
		s.serviceControllers = append(s.serviceControllers, &models.ServiceController{
			Title:    v.Title,
			ExecArgs: v.ExecArgs,
		})
		s.startService(s.serviceControllers[len(s.serviceControllers)-1])
	}
}

func (s *Service) startService(service *models.ServiceController) (err error) {
	s.Info().Msgf("starting service '%s'", service.Title)
	args := service.GetExecArgs()
	service.Exec = exec.Command(args[0], args[1:]...)
	if service.Exec == nil {
		return models.ErrCouldNotExec
	}
	service.StdoutReader, err = service.Exec.StdoutPipe()
	if err != nil {
		return err
	}
	service.StderrReader, err = service.Exec.StderrPipe()
	if err != nil {
		return err
	}
	service.Scanner = bufio.NewScanner(io.MultiReader(service.StdoutReader, service.StderrReader))
	if err = service.Exec.Start(); err != nil {
		s.Error().Msgf("cannot start %s", err.Error())
		return
	}
	go s.serviceLogs(service)
	go s.serviceWaiter(service)
	return nil
}

func (s *Service) serviceLogs(service *models.ServiceController) {
	for service.Scanner != nil {
		if !service.Scanner.Scan() {
			return
		}
		fmt.Println(service.Scanner.Text())
	}
}

func (s *Service) serviceWaiter(service *models.ServiceController) {
	state, err := service.Exec.Process.Wait()
	if err != nil {
		s.Error().Msgf("cannot read '%s' exit state: %s", service.Title, err.Error())
	} else {
		s.Error().Msgf("service '%s' exited with code %d", service.Title, state.ExitCode())
	}
}

// func (s *Service) clearProcess(service *models.ServiceController) {
// 	// Ignoring reader close errors because readers may be closed before this call
// 	if service.StdoutReader != nil {
// 		_ = service.StdoutReader.Close()
// 		service.StdoutReader = nil
// 	}
// 	if service.StderrReader != nil {
// 		_ = service.StderrReader.Close()
// 		service.StderrReader = nil
// 	}
// 	service.Exec = nil
// 	service.Scanner = nil
// 	// service.ShouldStop = false
// }
