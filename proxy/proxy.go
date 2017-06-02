package proxy

import (
	"io"
	"net"
	"sync"

	"github.com/damoye/llog"
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
		panic(err)
	}
	ln := temp.(*net.TCPListener)
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			panic(err)
		}
		go proxy.handleConnection(conn)
	}
}

func (proxy *Proxy) handleConnection(clientConn *net.TCPConn) {
	defer func() {
		if err := clientConn.Close(); err != nil {
			llog.Error("close client connection: ", err)
		}
	}()
	backend := proxy.manager.Get()
	if backend == nil {
		llog.Error("no backends")
		return
	}
	defer backend.DcreCount()
	temp, err := net.Dial("tcp", backend.Address)
	if err != nil {
		llog.Error("dial backend: ", err)
		return
	}
	serverConn := temp.(*net.TCPConn)
	defer func() {
		if err := serverConn.Close(); err != nil {
			llog.Error("close server connection: ", err)
		}
	}()
	llog.Info("connected ", clientConn.RemoteAddr(), " to ", serverConn.RemoteAddr())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go pipe(serverConn, clientConn, &wg)
	go pipe(clientConn, serverConn, &wg)
	wg.Wait()
	llog.Info("disconnected ", clientConn.RemoteAddr(), " to ", serverConn.RemoteAddr())
}

func pipe(dst *net.TCPConn, src *net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := io.Copy(dst, src)
	if err != nil {
		llog.Error("copy: ", err)
	}
	err = dst.CloseWrite()
	if err != nil {
		llog.Error("close write", err)
	}
}
