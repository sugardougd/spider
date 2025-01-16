package spider

import (
	"context"
	"fmt"
	"golang.org/x/term"
	"net"
)

type TCPConfig struct {
	*Config
	Address string // 监听地址 IP:PORT
}

func RunTCP(config *TCPConfig, commands *Commands, ctx context.Context) error {
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return err
	}
	fmt.Printf("Listening TCP on %s\r\n", config.Address)

	go acceptTCPConnection(listener, config, commands)
	select {
	case <-ctx.Done():
		return listener.Close()
	}
}

func acceptTCPConnection(listener net.Listener, config *TCPConfig, commands *Commands) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\r\n", err)
			break
		}
		go handleTCPConnection(conn, config, commands)
	}
}

func handleTCPConnection(conn net.Conn, config *TCPConfig, commands *Commands) {
	defer conn.Close()
	fmt.Printf("[%s]New TCP connection\r\n", conn.RemoteAddr())
	s := New(config.Config, commands)
	s.runWithTerminal(term.NewTerminal(conn, config.Prompt))
}
