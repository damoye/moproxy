package main

import (
	"flag"
	"io"
	"log"
	"net"
	"sync"

	"github.com/damoye/moproxy/backend"
	"github.com/damoye/moproxy/config"
)

var (
	configPath     = flag.String("config", "", "config file path")
	backendManager *backend.Manager
)

func handleConnection(frontConn net.Conn) {
	defer frontConn.Close()
	backend := backendManager.Get()
	if backend == nil {
		log.Print("no backends")
		return
	}
	defer backend.Decr()
	backConn, err := net.Dial("tcp", backend.Address)
	if err != nil {
		log.Println("dial:", err)
		return
	}
	defer backConn.Close()
	log.Printf("FRONTEND: %s, BACKEND: %s", frontConn.RemoteAddr(), backConn.RemoteAddr())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := io.Copy(backConn, frontConn)
		if err != nil {
			log.Println("copy front to back:", err)
		}
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(frontConn, backConn)
		if err != nil {
			log.Println("copy back to front:", err)
		}
	}()
	wg.Wait()
}

func main() {
	config, err := config.GenerateConfig(*configPath)
	if err != nil {
		panic(err)
	}
	backendManager = backend.NewManager(config.Backends)
	ln, err := net.Listen("tcp", config.Address)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}
