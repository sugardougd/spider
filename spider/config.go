package spider

// Config specifies the spider options.
type Config struct {
	Name        string // specifies the application name. This field is required.
	Description string // specifies the application description.
	Prompt      string // defines the shell prompt.
}
