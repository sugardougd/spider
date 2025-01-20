package spider

type ConfigOption struct {
	Name              string            // specifies the application name. This field is required.
	Description       string            // specifies the application description.
	Prompt            string            // defines the shell prompt.
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
