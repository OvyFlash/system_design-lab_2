package main

import (
	"lab_2/internal/base"
	"lab_2/pkg/services/launcher"

	"github.com/rs/zerolog/log"
)

func main() {
	app, err := base.NewApplication(launcher.NewLauncher, base.NameLauncher, false)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	if err := app.Start(); err != nil {
		log.Error().Msg(err.Error())
	}
}
