package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/harrypunk/v2lite/proxy/socks"
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
		if version != socks.Version5 {
			return nil, fmt.Errorf(("not socks5"))
		}

		// Write socks version
		if _, err = conn.Write([]byte{socks.Version5, socks.AuthNone}); err != nil {
			return nil, fmt.Errorf("failed to write socks5")
		}

		// Read cmd
		n, err = conn.Read(buf[:])
		log.Printf("read command n: %v", n)
		if err != nil {
			return nil, fmt.Errorf("failed to read cmd %w", err)
		}
		if n < 7 {
			break
		}
	}
	return conn, nil
}
