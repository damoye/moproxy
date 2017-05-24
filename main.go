package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

var (
	backends = []string{"127.0.0.1:6379", "127.0.0.1:9221"}
	r        = rand.New(rand.NewSource(time.Now().Unix()))
)

func getBackend() string {
	return backends[r.Intn(len(backends))]
}

func handleConnection(frontConn net.Conn) {
	defer frontConn.Close()
	backend := getBackend()
	backConn, err := net.Dial("tcp", backend)
	if err != nil {
		log.Println("dial:", err)
		return
	}
	defer backConn.Close()
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
	ln, err := net.Listen("tcp", ":8080")
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
