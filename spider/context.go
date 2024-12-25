package spider

import (
	"fmt"
	"strings"
)

type Context struct {
	Spider     *Spider // Reference to the app.
	Command    *Command
	CommandStr string
	flagValues FlagValues
}

// Stop signalizes the app to exit.
func (context *Context) Stop() error {
	return context.Spider.Stop()
}

func (context *Context) String() string {
	var builder strings.Builder
	builder.WriteString("'" + context.CommandStr + "'\n")
	builder.WriteString(fmt.Sprintf("Command:\n\t%s", context.Command.Name))
	if len(context.flagValues) > 0 {
		builder.WriteString("\nFlags:")
		for name, flag := range context.flagValues {
			builder.WriteString(fmt.Sprintf("\n\t%s: %v", name, flag.value))
			if flag.isDefault {
				builder.WriteString("[*]")
			}
		}
	}
	return builder.String()
}
