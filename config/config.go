package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

// Config ...
type Config struct {
	Address     string   `json:"address"`
	Backends    []string `json:"backends"`
	HTTPAddress string   `json:"http_address"`
}

var defaultConfig = Config{
	Address: ":8080",
	Backends: []string{
		"127.0.0.1:6379",
		"127.0.0.1:9221",
	},
	HTTPAddress: ":8081",
}

// GenerateConfig ...
func GenerateConfig(configPath string) (*Config, error) {
	flag.Parse()
	if configPath == "" {
		return &defaultConfig, nil
	}
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = json.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
