package main

import (
	"lab_2/internal/base"
	"lab_2/pkg/services/api"

	"github.com/rs/zerolog/log"
)

func main() {
	app, err := base.NewApplication(api.NewAPI, base.NameAPI, false)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	if err := app.Start(); err != nil {
		log.Error().Msg(err.Error())
	}
	if err := app.Stop(); err != nil {
		log.Error().Msg(err.Error())
	}
}