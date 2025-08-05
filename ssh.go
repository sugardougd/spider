package spider

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"net"
	"os"
)

type PasswordValidator func(user string, password []byte) bool

func RunSSH(ctx context.Context, config *Config, commands *Commands) error {
	// SSH 配置
	sshConfig, err := config.newSSHConfig()
	if err != nil {
		fmt.Printf("Failed to create  ServerConfig: %v\r\n", err)
		return err
	}
	// 监听 TCP 端口
	address := config.SSHConfig.Address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	fmt.Printf("Listening SSH on %s\r\n", address)

	go acceptSSHConnection(ctx, listener, config, sshConfig, commands)

	select {
	case <-ctx.Done():
		return listener.Close()
	}
}

func acceptSSHConnection(ctx context.Context, listener net.Listener, config *Config, sshConfig *ssh.ServerConfig, commands *Commands) {
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\r\n", err)
			break
		}
		childCtx, _ := context.WithCancel(ctx)
		go handleSSHConnection(childCtx, tcpConn, config, sshConfig, commands)
	}
}

func handleSSHConnection(ctx context.Context, conn net.Conn, config *Config, sshConfig *ssh.ServerConfig, commands *Commands) {
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
		go handleSSHChannel(ctx, sshConn, newChannel, config, commands)
	}
}

func handleSSHChannel(ctx context.Context, sshConn *ssh.ServerConn, newChannel ssh.NewChannel, config *Config, commands *Commands) {
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
	s := New(config, commands)
	defer channel.Close()
	go handleSSHChannelRequests(sshConn, reqs, func(width, height int) {
		s.SetSize(width, height)
	})
	if err = s.RunWithTerminal(ctx, term.NewTerminal(channel, config.Prompt)); err != nil {
		fmt.Printf("[%s@%s]Exit Terminal Spider: %v\r\n", sshConn.User(), sshConn.RemoteAddr(), err)
	}
}

func handleSSHChannelRequests(sshConn *ssh.ServerConn, reqs <-chan *ssh.Request, windowChanged func(width, height int)) {
	for req := range reqs {
		fmt.Printf("[%s@%s]Received channel request: %s %t\r\n", sshConn.User(), sshConn.RemoteAddr(), req.Type, req.WantReply)
		switch req.Type {
		case "pty-req", "shell":
			req.Reply(true, nil)
			break
		case "window-change":
			req.Reply(false, nil)
			payload := bytes.NewReader(req.Payload)
			// 0:width; 1:height; 2:widthPixels; 3:heightPixels
			windows := [4]int32{}
			for index := 0; index < len(windows); index++ {
				binary.Read(payload, binary.BigEndian, &windows[index])
			}
			if windows[0] > 0 && windows[1] > 0 {
				windowChanged(int(windows[0]), int(windows[1]))
			}
			break
		default:
			req.Reply(false, nil)
			break
		}
	}
}

func (c *Config) newSSHConfig() (*ssh.ServerConfig, error) {
	sshConfig := &ssh.ServerConfig{
		NoClientAuth: c.SSHConfig.NoClientAuth,
	}
	if c.SSHConfig.PasswordValidator != nil {
		sshConfig.PasswordCallback = func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			// 这里可以添加你自己的认证逻辑，例如检查用户名和密码
			if c.SSHConfig.PasswordValidator(conn.User(), password) {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", conn.User())
		}
	}

	// banner
	sshConfig.BannerCallback = func(conn ssh.ConnMetadata) string {
		return c.SSHConfig.Banner
	}

	// 生成一个 SSH 密钥对
	privateBytes, err := os.ReadFile(c.SSHConfig.PrivateFile)
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
