package proxy

import (
	"io"
	"log"
	"net"
	"sync"

	"github.com/damoye/moproxy/backend"
	"github.com/damoye/moproxy/config"
)

// Proxy ...
type Proxy struct {
	config  *config.Config
	manager *backend.Manager
}

// New ...
func New(config *config.Config) *Proxy {
	return &Proxy{
		config:  config,
		manager: backend.NewManager(config.Backends),
	}
}

// Run ...
func (proxy *Proxy) Run() {
	go proxy.serveHTTP()
	temp, err := net.Listen("tcp", proxy.config.Address)
	if err != nil {
		log.Fatalln("FATAL: listen:", err)
	}
	ln := temp.(*net.TCPListener)
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Fatalln("FATAL: acceptTCP:", err)
		}
		go proxy.handleConnection(conn)
	}
}

func (proxy *Proxy) handleConnection(frontConn *net.TCPConn) {
	defer func() {
		if err := frontConn.Close(); err != nil {
			log.Println("ERROR: frontConn close:", err)
		}
	}()
	backend := proxy.manager.Get()
	if backend == nil {
		log.Print("ERROR: no backends")
		return
	}
	defer backend.DcreCount()
	temp, err := net.Dial("tcp", backend.Address)
	if err != nil {
		log.Println("ERROR: dial:", err)
		return
	}
	backConn := temp.(*net.TCPConn)
	defer func() {
		if err := backConn.Close(); err != nil {
			log.Println("ERROR: backConn close:", err)
		}
	}()
	log.Printf("INFO: frontend: %s, backend: %s", frontConn.RemoteAddr(), backConn.RemoteAddr())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go pipe(backConn, frontConn, &wg)
	go pipe(frontConn, backConn, &wg)
	wg.Wait()
	log.Printf("INFO: frontend: %s ended", frontConn.RemoteAddr())
}

func pipe(dst *net.TCPConn, src *net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Println("ERROR: copy:", err)
	}
	err = dst.CloseWrite()
	if err != nil {
		log.Println("ERROR: closeWrite:", err)
	}
}
