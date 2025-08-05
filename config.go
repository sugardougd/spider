package spider

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

type ExecutedHook func(*Context, error)

// Config specifies the spider options.
type Config struct {
	Name         string       // specifies the application name. This field is required.
	Description  string       // specifies the application description.
	Prompt       string       // defines the shell prompt.
	Interactive  bool         // cannot auto complete if not
	ExecutedHook ExecutedHook // 执行命令后回调函数
	Welcome      string       // welcome message
}

type TCPConfig struct {
	*Config
	Address string // 监听地址 IP:PORT
}

type SSHConfig struct {
	*Config
	Address           string            // 监听地址 IP:PORT
	NoClientAuth      bool              // 是否认证客户端
	PasswordValidator PasswordValidator // 校验用户名密码
	Banner            string            // banner
	PrivateFile       string            // RSA私钥文件
}

func (c *SSHConfig) newSSHConfig() (*ssh.ServerConfig, error) {
	sshConfig := &ssh.ServerConfig{
		NoClientAuth: c.NoClientAuth,
	}
	if c.PasswordValidator != nil {
		sshConfig.PasswordCallback = func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			// 这里可以添加你自己的认证逻辑，例如检查用户名和密码
			if c.PasswordValidator(conn.User(), password) {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", conn.User())
		}
	}

	// banner
	sshConfig.BannerCallback = func(conn ssh.ConnMetadata) string {
		return c.Banner
	}

	// 生成一个 SSH 密钥对
	privateBytes, err := os.ReadFile(c.PrivateFile)
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
