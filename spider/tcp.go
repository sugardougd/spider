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

	go acceptTCPConnection(listener, config, commands, ctx)
	select {
	case <-ctx.Done():
		return listener.Close()
	}
}

func acceptTCPConnection(listener net.Listener, config *Config, commands *Commands, ctx context.Context) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\r\n", err)
			break
		}
		childCtx, _ := context.WithCancel(ctx)
		go handleTCPConnection(conn, config, commands, childCtx)
	}
}

func handleTCPConnection(conn net.Conn, config *Config, commands *Commands, ctx context.Context) {
	defer conn.Close()
	fmt.Printf("[%s]New TCP connection\r\n", conn.RemoteAddr())
	s := New(config)
	s.AddCommands(commands)
	if err := s.RunWithTerminal(term.NewTerminal(conn, config.Prompt), ctx); err != nil {
		fmt.Printf("[%s]Exit Terminal Spider: %v\r\n", conn.RemoteAddr(), err)
	}
}
