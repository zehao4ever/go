package server

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	ServerPort string `tomal:"port"`
	DBPass     string `toml:"dbpwd"`
	DBName     string `toml:"dbname"`
}

func NewConfig() (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile("server/config.toml", &conf)
	return &conf, err
}
