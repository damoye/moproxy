package main

import (
	"encoding/json"
	"flag"

	"github.com/damoye/llog"
	"github.com/damoye/moproxy/config"
	"github.com/damoye/moproxy/proxy"
)

func main() {
	configPath := flag.String("config", "", "config file path")
	config, err := config.GenerateConfig(*configPath)
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	llog.Info("config: ", string(b))
	proxy := proxy.New(config)
	proxy.Run()
}
