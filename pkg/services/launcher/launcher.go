package launcher

import (
	"lab_2/pkg/models/settings"
	ls "lab_2/pkg/services/launcher/service"
)

type Launcher struct {
	*ls.Service
}

func NewLauncher(e *settings.Essential) *Launcher {
	return &Launcher{
		Service: ls.NewService(e),
	}
}

func (l *Launcher) Start() error {
	if err := l.init(); err != nil {
		return err
	}
	//set http services

	return l.Service.Start()
}

func (l *Launcher) Stop() error {
	return nil
}

func (l *Launcher) init() (err error) {
	return nil
}
