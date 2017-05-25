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

func pipe(dst *net.TCPConn, src *net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Println("copy front to back:", err)
	}
	err = dst.CloseWrite()
	if err != nil {
		log.Println("backConn closeWrite:", err)
	}
}

func handleConnection(frontConn *net.TCPConn) {
	defer func() {
		if err := frontConn.Close(); err != nil {
			log.Println("frontConn close:", err)
		}
	}()
	backend := backendManager.Get()
	if backend == nil {
		log.Print("no backends")
		return
	}
	defer backend.Decr()
	temp, err := net.Dial("tcp", backend.Address)
	if err != nil {
		log.Println("dial:", err)
		return
	}
	backConn := temp.(*net.TCPConn)
	defer func() {
		if err := backConn.Close(); err != nil {
			log.Println("backConn close:", err)
		}
	}()
	log.Printf("FRONTEND: %s, BACKEND: %s", frontConn.RemoteAddr(), backConn.RemoteAddr())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go pipe(backConn, frontConn, &wg)
	go pipe(frontConn, backConn, &wg)
	wg.Wait()
	log.Printf("FRONTEND: %s ENDED", frontConn.RemoteAddr())
}

func main() {
	config, err := config.GenerateConfig(*configPath)
	if err != nil {
		panic(err)
	}
	backendManager = backend.NewManager(config.Backends)
	temp, err := net.Listen("tcp", config.Address)
	ln := temp.(*net.TCPListener)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Println("accept:", err)
		}
		go handleConnection(conn)
	}
}
