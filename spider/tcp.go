package spider

import (
	"fmt"
	"github.com/desertbit/readline"
	"net"
	"spider/grumble"
)

type TCPSpider struct {
	cfg *TCPConfig
}

type TCPConfig struct {
	*grumble.Config
	Address  string // 监听地址 IP:PORT
	Commands CommandsSet
}

// RunTCP run tcp cli
func RunTCP(cfg *TCPConfig) error {
	s := &TCPSpider{cfg: cfg}
	if cfg.FuncIsTerminal == nil {
		cfg.FuncIsTerminal = s.FuncIsTerminal
	}
	// 监听 TCP 端口
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	fmt.Printf("Listening on %s\r\n", cfg.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\r\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *TCPSpider) handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("[%s]New TCP connection\r\n", conn.RemoteAddr())
	// 执行命令
	rlCfg := &readline.Config{
		Prompt:         "$ ",
		Stdin:          conn,
		Stdout:         conn,
		Stderr:         conn,
		FuncIsTerminal: s.cfg.FuncIsTerminal,
	}
	rl, err := readline.NewEx(rlCfg)
	if err != nil {
		return
	}
	app := NewGrumble(s.cfg.Config, s.cfg.Commands)
	if err := app.RunWithReadline(rl); err != nil {
		fmt.Printf("run %s fail. %v\r\n", s.cfg.Name, err)
	}
}

func (s *TCPSpider) FuncIsTerminal() bool {
	return false
}
