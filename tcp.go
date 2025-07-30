package spider

import (
	"context"
	"fmt"
	"golang.org/x/term"
	"net"
)

func RunTCP(ctx context.Context, config *Config, commands *Commands) error {
	address := config.TCPConfig.Address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	fmt.Printf("Listening TCP on %s\r\n", address)

	go acceptTCPConnection(ctx, listener, config, commands)
	select {
	case <-ctx.Done():
		return listener.Close()
	}
}

func acceptTCPConnection(ctx context.Context, listener net.Listener, config *Config, commands *Commands) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\r\n", err)
			break
		}
		childCtx, _ := context.WithCancel(ctx)
		go handleTCPConnection(childCtx, conn, config, commands)
	}
}

func handleTCPConnection(ctx context.Context, conn net.Conn, config *Config, commands *Commands) {
	defer conn.Close()
	fmt.Printf("[%s]New TCP connection\r\n", conn.RemoteAddr())
	s := New(config)
	s.AddCommands(commands)
	if err := s.RunWithTerminal(ctx, term.NewTerminal(conn, config.Prompt)); err != nil {
		fmt.Printf("[%s]Exit Terminal Spider: %v\r\n", conn.RemoteAddr(), err)
	}
}
