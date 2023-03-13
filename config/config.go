package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Databases Databases `yaml:"databases"`
	Loggers   []Logger  `yaml:"loggers"`
	Services  []Service `yaml:"services"`
}

type Databases struct {
	SQLConfig   SQLDatabase   `yaml:"sql"`
	RedisConfig RedisDatabase `yaml:"redis"`
}

type SQLDatabase struct {
	DSN string `yaml:"dsn"`
}

type RedisDatabase struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Logger struct {
	Type  LoggerType `yaml:"type"`
	Level string     `yaml:"level"`
}

type Service struct {
	Title    string   `yaml:"title"`
	Port     uint16   `yaml:"rest_port,omitempty"`
	ExecArgs []string `yaml:"exec_args,omitempty"`
	Prefixes []string `yaml:"prefixes,omitempty"`
}

type LoggerType string

const (
	Console  LoggerType = "console"
	Database LoggerType = "database"
)

var (
	configPath    = "config.yaml"
	defaultConfig = Config{
		Loggers: []Logger{
			{
				Type:  Console,
				Level: "debug",
			},
			{
				Type:  Database,
				Level: "warn",
			},
		},
	}
)

func init() {
	os.Setenv("TZ", "Europe/Kyiv")
	flag.StringVar(&configPath, "p", configPath, "Path to config file")
}

func NewConfig() (c Config) {
	flag.Parse()
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}
	defaultConfig = c
	return
}

func (c Config) GetPortByTitle(title string) uint16 {
	for _, v := range c.Services {
		if v.Title == title {
			return v.Port
		}
	}
	return 0
}
