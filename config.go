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

func ConfigName(name string) ConfigOptional {
	return func(option *ConfigOption) {
		option.Name = name
	}
}

func ConfigDescription(description string) ConfigOptional {
	return func(option *ConfigOption) {
		option.Description = description
	}
}

func ConfigPrompt(prompt string) ConfigOptional {
	return func(option *ConfigOption) {
		option.Prompt = prompt
	}
}

func ConfigInteractive(interactive bool) ConfigOptional {
	return func(option *ConfigOption) {
		option.Interactive = interactive
	}
}

func ConfigExecutedHook(hook ExecutedHook) ConfigOptional {
	return func(option *ConfigOption) {
		option.ExecutedHook = hook
	}
}

func ConfigWelcome(welcome string) ConfigOptional {
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
	return NewConfig(ConfigName(name),
		ConfigDescription(description),
		ConfigInteractive(true),
		ConfigPrompt(prompt),
		ConfigWelcome(welcome),
		ConfigExecutedHook(hook))
}

func NewTCPConfig(name, description, prompt, welcome string, hook ExecutedHook, config TCPConfig) *Config {
	return NewConfig(ConfigName(name),
		ConfigDescription(description),
		ConfigInteractive(false),
		ConfigPrompt(prompt),
		ConfigWelcome(welcome),
		ConfigExecutedHook(hook),
		WithTCPConfig(config))
}

func NewSSHConfig(name, description, prompt, welcome string, hook ExecutedHook, config SSHConfig) *Config {
	return NewConfig(ConfigName(name),
		ConfigDescription(description),
		ConfigInteractive(true),
		ConfigPrompt(prompt),
		ConfigWelcome(welcome),
		ConfigExecutedHook(hook),
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
