package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

// Config ...
type Config struct {
	Address  string   `json:"address"`
	Backends []string `json:"backends"`
}

var defaultConfig = Config{
	Address: ":8080",
	Backends: []string{
		"127.0.0.1:6379",
		"127.0.0.1:9221",
	},
}

// GenerateConfig ...
func GenerateConfig(configPath string) (*Config, error) {
	flag.Parse()
	if configPath == "" {
		log.Print("config path is empty, use default config")
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
