package main

import (
	"flag"
	"log"

	"github.com/damoye/moproxy/config"
	"github.com/damoye/moproxy/proxy"
)

func main() {
	configPath := flag.String("config", "", "config file path")
	config, err := config.GenerateConfig(*configPath)
	if err != nil {
		log.Fatalln("FATAL: generate config:", err)
	}
	proxy := proxy.New(config)
	proxy.Run()
}
