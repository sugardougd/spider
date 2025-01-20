package spider

import (
	"context"
	"fmt"
	"golang.org/x/term"
	"net"
)

func RunTCP(config *Config, commands *Commands, ctx context.Context) error {
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

func acceptTCPConnection(listener net.Listener, config *Config, commands *Commands) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\r\n", err)
			break
		}
		go handleTCPConnection(conn, config, commands)
	}
}

func handleTCPConnection(conn net.Conn, config *Config, commands *Commands) {
	defer conn.Close()
	fmt.Printf("[%s]New TCP connection\r\n", conn.RemoteAddr())
	s := New(config, commands)
	s.RunWithTerminal(term.NewTerminal(conn, config.Prompt))
}
