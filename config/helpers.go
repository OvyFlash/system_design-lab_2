package config

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

func getLogLevel() zerolog.Level {
	level := "info"
	for _, v := range defaultConfig.Loggers {
		if v.Type == Console {
			level = v.Level
			break
		}
	}
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "-":
		return zerolog.NoLevel
	case "disabled":
		return zerolog.Disabled
	}
	panic("log level does not exist")
}

func GetLogger(prefix string) zerolog.Logger {
	pr := fmt.Sprintf("[%s]", prefix)
	return zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		FormatCaller: func(i interface{}) string {
			return pr
		},
	}).Level(getLogLevel()).With().Timestamp().Logger()
}

func GetGormLogLevel() logger.LogLevel {
	level := "warn"
	for _, v := range defaultConfig.Loggers {
		if v.Type == Database {
			level = v.Level
			break
		}
	}
	switch level {
	case "info":
		return logger.Info
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "silent":
		return logger.Silent
	}
	panic("log level does not exist")
}

func GetLocation() *time.Location {
	location, err := time.LoadLocation("Europe/Kyiv")
	if err != nil {
		return time.Local
	}
	return location
}
