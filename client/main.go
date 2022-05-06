package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/harrypunk/v2lite/proxy/socks/socks5"
)

const (
	localAddr = "127.0.0.1:10051"
)

func main() {
	log.Println("v2lite start")

	// tcp listener
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Printf("failed to listen on %s", localAddr)
		return
	}
	log.Printf("listening on %s", localAddr)
	for {
		localConnection, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err.Error())
			continue
		}
		go func() {
			defer localConnection.Close()
			if _, err = handleConnection(localConnection); err != nil {
				log.Printf("local connection error, %v\n", err)
			}
		}()
	}
}

func handleConnection(conn net.Conn) (io.ReadWriter, error) {
	var buf [128]byte
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			return nil, err
		} else if n == 0 {
			return nil, fmt.Errorf("read zero")
		}
		version := buf[0]
		if version != socks5.Version5 {
			return nil, fmt.Errorf(("not socks5"))
		}

		if _, err = conn.Write([]byte("hello123")); err != nil {
			log.Printf("write err %v", err)
			break
		}
	}
	return conn, nil
}
