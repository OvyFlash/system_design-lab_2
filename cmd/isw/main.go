package main

import (
	"lab_2/internal/base"
	"lab_2/pkg/services/isw_reader"

	"github.com/rs/zerolog/log"
)

func main() {
	app, err := base.NewApplication(isw_reader.NewISWReader, base.NameISW, true)
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
