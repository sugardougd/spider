package spider

type ExecutedHook func(*Context, error)

type ConfigOption struct {
	Name              string            // specifies the application name. This field is required.
	Description       string            // specifies the application description.
	Prompt            string            // defines the shell prompt.
	Interactive       bool              // cannot auto complete if not
	ExecutedHook      ExecutedHook      // 执行命令后回调函数
	Address           string            // 监听地址 IP:PORT. for tcp & ssh
	NoClientAuth      bool              // 是否认证客户端. for ssh
	passwordValidator PasswordValidator // 校验用户名密码. for ssh
	Banner            string            // banner. for ssh
	PrivateFile       string            // RSA私钥文件. for ssh
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

func ConfigAddress(address string) ConfigOptional {
	return func(option *ConfigOption) {
		option.Address = address
	}
}

func ConfigNoClientAuth(noClientAuth bool) ConfigOptional {
	return func(option *ConfigOption) {
		option.NoClientAuth = noClientAuth
	}
}

func ConfigPasswordValidator(passwordValidator PasswordValidator) ConfigOptional {
	return func(option *ConfigOption) {
		option.passwordValidator = passwordValidator
	}
}

func ConfigBanner(banner string) ConfigOptional {
	return func(option *ConfigOption) {
		option.Banner = banner
	}
}

func ConfigPrivateFile(privateFile string) ConfigOptional {
	return func(option *ConfigOption) {
		option.PrivateFile = privateFile
	}
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
