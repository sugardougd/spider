package spider

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"net"
	"os"
)

type PasswordValidator func(user, password string) bool

func RunSSH(config *Config, commands *Commands, ctx context.Context) error {
	// SSH 配置
	sshConfig, err := config.newSSHConfig()
	if err != nil {
		fmt.Printf("Failed to create  ServerConfig: %v\r\n", err)
		return err
	}
	// 监听 TCP 端口
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return err
	}
	fmt.Printf("Listening SSH on %s\r\n", config.Address)

	go acceptSSHConnection(listener, config, sshConfig, commands)

	select {
	case <-ctx.Done():
		return listener.Close()
	}
}

func acceptSSHConnection(listener net.Listener, config *Config, sshConfig *ssh.ServerConfig, commands *Commands) {
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\r\n", err)
			break
		}
		go handleSSHConnection(tcpConn, config, sshConfig, commands)
	}
}

func handleSSHConnection(conn net.Conn, config *Config, sshConfig *ssh.ServerConfig, commands *Commands) {
	defer conn.Close()
	// 进行 SSH 握手
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, sshConfig)
	if err != nil {
		fmt.Printf("Failed to establish SSH connection: %v\r\n", err)
		return
	}
	defer sshConn.Close()
	fmt.Printf("[%s@%s]New SSH connection\r\n", sshConn.User(), sshConn.RemoteAddr())

	// 处理通道和请求
	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		go handleSSHChannel(sshConn, newChannel, config, commands)
	}
}

func handleSSHChannel(sshConn *ssh.ServerConn, newChannel ssh.NewChannel, config *Config, commands *Commands) {
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
	go handleSSHChannelRequests(sshConn, reqs)
	s := New(config, commands)
	s.RunWithTerminal(term.NewTerminal(channel, config.Prompt))
}

func handleSSHChannelRequests(sshConn *ssh.ServerConn, reqs <-chan *ssh.Request) {
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

func (s *Config) newSSHConfig() (*ssh.ServerConfig, error) {
	sshConfig := &ssh.ServerConfig{
		NoClientAuth: s.NoClientAuth,
	}
	if s.passwordValidator != nil {
		sshConfig.PasswordCallback = func(conn ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// 这里可以添加你自己的认证逻辑，例如检查用户名和密码
			if s.passwordValidator(conn.User(), string(pass)) {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", conn.User())
		}
	}

	// banner
	sshConfig.BannerCallback = func(conn ssh.ConnMetadata) string {
		return s.Banner
	}

	// 生成一个 SSH 密钥对
	privateBytes, err := os.ReadFile(s.PrivateFile)
	if err != nil {
		return nil, err
	}
	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, err
	}
	sshConfig.AddHostKey(privateKey)
	return sshConfig, nil
}
