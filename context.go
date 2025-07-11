package spider

import (
	"context"
	"fmt"
	"strings"
)

type Context struct {
	Spider     *Spider // Reference to the app.
	Command    *Command
	CommandStr string
	FlagValues FlagValues
	ArgValues  ArgValues
	Ctx        context.Context
}

// Stop signalizes the app to exit.
func (context *Context) Stop() error {
	return context.Spider.stop()
}

func (context *Context) String() string {
	var builder strings.Builder
	builder.WriteString("'" + context.CommandStr + "'\n")
	builder.WriteString(fmt.Sprintf("Command:\n\t%s", context.Command.Name))
	if len(context.FlagValues) > 0 {
		builder.WriteString("\nFlags:")
		for name, flag := range context.FlagValues {
			builder.WriteString(fmt.Sprintf("\n\t%s: %v", name.long, flag.value))
			if flag.isDefault {
				builder.WriteString("[*]")
			}
		}
	}
	if len(context.ArgValues) > 0 {
		builder.WriteString("\nArgument:")
		for name, arg := range context.ArgValues {
			builder.WriteString(fmt.Sprintf("\n\t%s: %v", name, arg.value))
			if arg.isDefault {
				builder.WriteString("[*]")
			}
		}
	}
	return builder.String()
}
