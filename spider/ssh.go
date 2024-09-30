package spider

import (
	"fmt"
	"github.com/desertbit/readline"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"spider/grumble"
)

type SSHSpider struct {
	cfg *SSHConfig
}

type SSHConfig struct {
	*grumble.Config
	Address      string                           // 监听地址 IP:PORT
	NoClientAuth bool                             // 是否认证客户端
	PasswordFunc func(user, password string) bool // 校验用户名密码
	Banner       string                           // banner
	PrivateFile  string                           // RSA私钥文件
	Commands     CommandsSet
}

// RunSSH run ssh cli
// example golang.org/x/crypto@v0.26.0/ssh/example_test.go:25
func RunSSH(cfg *SSHConfig) error {
	s := &SSHSpider{cfg: cfg}
	if cfg.FuncIsTerminal == nil {
		cfg.FuncIsTerminal = s.FuncIsTerminal
	}
	// SSH 配置
	sshCfg, err := s.newSSHConfig()
	if err != nil {
		return err
	}
	// 监听 TCP 端口
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	fmt.Printf("Listening SSH on %s\r\n", cfg.Address)

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\r\n", err)
			continue
		}
		go s.handleConnection(tcpConn, sshCfg)
	}
}

func (s *SSHSpider) handleConnection(conn net.Conn, sshCfg *ssh.ServerConfig) {
	defer conn.Close()
	// 进行 SSH 握手
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, sshCfg)
	if err != nil {
		fmt.Printf("Failed to establish SSH connection: %v\r\n", err)
		return
	}
	defer sshConn.Close()
	fmt.Printf("[%s@%s]New SSH connection\r\n", sshConn.User(), sshConn.RemoteAddr())

	// 处理通道和请求
	go s.handleRequests(sshConn, reqs)

	for newChannel := range chans {
		go s.handleChannel(sshConn, newChannel)
	}
}

func (s *SSHSpider) handleRequests(sshConn *ssh.ServerConn, reqs <-chan *ssh.Request) {
	//ssh.DiscardRequests(reqs)
	for req := range reqs {
		fmt.Printf("[%s@%s]Received request: %s %t\r\n", sshConn.User(), sshConn.RemoteAddr(), req.Type, req.WantReply)
		if req.WantReply {
			req.Reply(false, nil)
		}
	}
}

func (s *SSHSpider) handleChannel(sshConn *ssh.ServerConn, newChannel ssh.NewChannel) {
	fmt.Printf("[%s@%s]NewChannel type: %s\r\n", sshConn.User(), sshConn.RemoteAddr(), newChannel.ChannelType())
	if newChannel.ChannelType() != "session" {
		newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
		return
	}
	channel, reqs, err := newChannel.Accept()
	if err != nil {
		fmt.Printf("Could not accept channel: %v\r\n", err)
		return
	}
	defer channel.Close()
	go s.handleChannelRequests(sshConn, channel, reqs)

	// 执行命令
	rlCfg := &readline.Config{
		Prompt:         sshConn.User() + "$ ",
		Stdin:          channel,
		Stdout:         channel,
		Stderr:         channel,
		FuncIsTerminal: s.cfg.FuncIsTerminal,
	}
	rl, err := readline.NewEx(rlCfg)
	if err != nil {
		return
	}
	app := newGrumble(s.cfg.Config, s.cfg.Commands)
	if err := app.RunWithReadline(rl); err != nil {
		fmt.Printf("run %s fail. %v\r\n", s.cfg.Name, err)
	}

	// 执行命令
	/*terminal := term.NewTerminal(channel, sshConn.User()+"$ ")
	for {
		line, err := terminal.ReadLine()
		if err != nil {
			break
		}
		command := strings.TrimSpace(line)
		fmt.Printf("[%s@%s]Receive command: %s\r\n", sshConn.User(), sshConn.RemoteAddr(), command)
		if i, err := terminal.Write([]byte(command)); err == nil {
			fmt.Printf("[%s@%s]Write %d bytes\r\n", sshConn.User(), sshConn.RemoteAddr(), i)
		} else {
			fmt.Printf("[%s@%s]Write fail %v", sshConn.User(), sshConn.RemoteAddr(), err)
		}
	}*/
}

func (s *SSHSpider) handleChannelRequests(sshConn *ssh.ServerConn, channel ssh.Channel, reqs <-chan *ssh.Request) {
	for req := range reqs {
		fmt.Printf("[%s@%s]Received channel request: %s %t\r\n", sshConn.User(), sshConn.RemoteAddr(), req.Type, req.WantReply)
		switch req.Type {
		case "pty-req", "shell":
			req.Reply(true, nil)
			break
		default:
			req.Reply(false, nil)
			break

		}
	}
}

func (s *SSHSpider) newSSHConfig() (*ssh.ServerConfig, error) {
	sshCfg := &ssh.ServerConfig{
		NoClientAuth: s.cfg.NoClientAuth,
	}
	if s.cfg.PasswordFunc != nil {
		sshCfg.PasswordCallback = func(conn ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// 这里可以添加你自己的认证逻辑，例如检查用户名和密码
			if s.cfg.PasswordFunc(conn.User(), string(pass)) {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", conn.User())
		}
	}

	// banner
	sshCfg.BannerCallback = func(conn ssh.ConnMetadata) string {
		return s.cfg.Banner
	}

	// 生成一个 SSH 密钥对
	privateBytes, err := os.ReadFile(s.cfg.PrivateFile)
	if err != nil {
		return nil, err
	}
	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, err
	}
	sshCfg.AddHostKey(privateKey)
	return sshCfg, nil
}

func (s *SSHSpider) FuncIsTerminal() bool {
	return true
}
