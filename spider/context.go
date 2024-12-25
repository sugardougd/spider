package spider

type Context struct {
	Spider     *Spider // Reference to the app.
	Command    *Command
	flagValues FlagValues
}

// Stop signalizes the app to exit.
func (c *Context) Stop() error {
	return c.Spider.Stop()
}
