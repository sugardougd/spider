package spider

type ExecutedHook func(*Context, error)

type ConfigOption struct {
	Name         string       // specifies the application name. This field is required.
	Description  string       // specifies the application description.
	Prompt       string       // defines the shell prompt.
	Interactive  bool         // cannot auto complete if not
	ExecutedHook ExecutedHook // 执行命令后回调函数
	Welcome      string       // welcome message
	TCPConfig    TCPConfig
	SSHConfig    SSHConfig
}

type TCPConfig struct {
	Address string // 监听地址 IP:PORT
}

type SSHConfig struct {
	Address           string            // 监听地址 IP:PORT
	NoClientAuth      bool              // 是否认证客户端
	PasswordValidator PasswordValidator // 校验用户名密码
	Banner            string            // banner
	PrivateFile       string            // RSA私钥文件
}

// Config specifies the spider options.
type Config struct {
	ConfigOption
}

type ConfigOptional func(option *ConfigOption)

func WithName(name string) ConfigOptional {
	return func(option *ConfigOption) {
		option.Name = name
	}
}

func WithDescription(description string) ConfigOptional {
	return func(option *ConfigOption) {
		option.Description = description
	}
}

func WithPrompt(prompt string) ConfigOptional {
	return func(option *ConfigOption) {
		option.Prompt = prompt
	}
}

func WithInteractive(interactive bool) ConfigOptional {
	return func(option *ConfigOption) {
		option.Interactive = interactive
	}
}

func WithExecutedHook(hook ExecutedHook) ConfigOptional {
	return func(option *ConfigOption) {
		option.ExecutedHook = hook
	}
}

func WithWelcome(welcome string) ConfigOptional {
	return func(option *ConfigOption) {
		option.Welcome = welcome
	}
}

func WithTCPConfig(config TCPConfig) ConfigOptional {
	return func(option *ConfigOption) {
		option.TCPConfig = config
	}
}

func WithSSHConfig(config SSHConfig) ConfigOptional {
	return func(option *ConfigOption) {
		option.SSHConfig = config
	}
}

func NewConsoleConfig(name, description, prompt, welcome string, hook ExecutedHook) *Config {
	return NewConfig(WithName(name),
		WithDescription(description),
		WithInteractive(true),
		WithPrompt(prompt),
		WithWelcome(welcome),
		WithExecutedHook(hook))
}

func NewTCPConfig(name, description, prompt, welcome string, hook ExecutedHook, config TCPConfig) *Config {
	return NewConfig(WithName(name),
		WithDescription(description),
		WithInteractive(false),
		WithPrompt(prompt),
		WithWelcome(welcome),
		WithExecutedHook(hook),
		WithTCPConfig(config))
}

func NewSSHConfig(name, description, prompt, welcome string, hook ExecutedHook, config SSHConfig) *Config {
	return NewConfig(WithName(name),
		WithDescription(description),
		WithInteractive(true),
		WithPrompt(prompt),
		WithWelcome(welcome),
		WithExecutedHook(hook),
		WithSSHConfig(config))
}

func NewConfig(optional ...ConfigOptional) *Config {
	option := &ConfigOption{}
	for _, o := range optional {
		o(option)
	}
	return &Config{*option}
}

func (c *Config) Option(optional ...ConfigOptional) {
	for _, o := range optional {
		o(&c.ConfigOption)
	}
}
