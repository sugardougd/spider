package spider

import (
	"fmt"
	"strconv"
)

type ArgParser func(arg string, args []string) (*ArgValue, []string, error)

type Args struct {
	list    []*Arg
	parsers map[string]ArgParser
}

type Arg struct {
	Name    string
	Help    string
	Default interface{}
	Require bool
}

type ArgValue struct {
	value     interface{}
	isDefault bool
}

type ArgValues map[string]*ArgValue

func (a *Args) parse(args []string, argValues ArgValues) (remaining []string, err error) {
	for _, arg := range a.list {
		if len(args) == 0 {
			break
		}
		if parser, ok := a.parsers[arg.Name]; ok {
			var argValue *ArgValue
			argValue, args, err = parser(arg.Name, args)
			if err != nil {
				break
			}
			argValues[arg.Name] = argValue
		}
	}
	remaining = args

	// check require argument and set the default values
	for _, arg := range a.list {
		if _, ok := argValues[arg.Name]; ok {
			continue
		}
		if arg.Require {
			err = fmt.Errorf("missing argument '%s'", arg.Name)
			break
		}
		argValues[arg.Name] = &ArgValue{
			value:     arg.Default,
			isDefault: true,
		}
	}
	return
}

func (a *Args) register(arg *Arg, ap ArgParser) error {
	err := arg.validate()
	if err != nil {
		return err
	}
	for _, a := range a.list {
		if a.Name == arg.Name {
			return fmt.Errorf("argument '%s' registered twice", arg.Name)
		}
	}
	a.list = append(a.list, arg)
	if a.parsers == nil {
		a.parsers = make(map[string]ArgParser)
	}
	a.parsers[arg.Name] = ap
	return nil
}

// Bool registers a bool argument.
func (a *Args) Bool(arg *Arg) {
	a.register(arg, parseBoolArg)
}

// Int registers an int argument.
func (a *Args) Int(arg *Arg) {
	a.register(arg, parseIntArg)
}

// Float64 registers a float argument.
func (a *Args) Float64(arg *Arg) {
	a.register(arg, parseFloat64Arg)
}

// String registers a string argument.
func (a *Args) String(arg *Arg) {
	a.register(arg, parseStringArg)
}

func (arg *Arg) validate() error {
	if len(arg.Name) == 0 {
		return fmt.Errorf("empty argument name")
	} else if len(arg.Help) == 0 {
		return fmt.Errorf("empty flag help message")
	}
	return nil
}

func parseBoolArg(arg string, args []string) (*ArgValue, []string, error) {
	if len(args) == 0 {
		return nil, args, fmt.Errorf("missing argument '%s'", arg)
	}
	vStr := args[0]
	args = args[1:]
	v, err := strconv.ParseBool(vStr)
	if err != nil {
		return nil, args, fmt.Errorf("invalid boolean value for argument: %s", arg)
	}
	return &ArgValue{
		value:     v,
		isDefault: false,
	}, args, nil
}

func parseIntArg(arg string, args []string) (*ArgValue, []string, error) {
	if len(args) == 0 {
		return nil, args, fmt.Errorf("missing argument '%s'", arg)
	}
	vStr := args[0]
	args = args[1:]
	v, err := strconv.Atoi(vStr)
	if err != nil {
		return nil, args, fmt.Errorf("invalid int value for argument: %s", arg)
	}
	return &ArgValue{
		value:     v,
		isDefault: false,
	}, args, nil
}

func parseFloat64Arg(arg string, args []string) (*ArgValue, []string, error) {
	if len(args) == 0 {
		return nil, args, fmt.Errorf("missing argument '%s'", arg)
	}
	vStr := args[0]
	args = args[1:]
	v, err := strconv.ParseFloat(vStr, 64)
	if err != nil {
		return nil, args, fmt.Errorf("invalid float value for argument: %s", arg)
	}
	return &ArgValue{
		value:     v,
		isDefault: false,
	}, args, nil
}

func parseStringArg(arg string, args []string) (*ArgValue, []string, error) {
	if len(args) == 0 {
		return nil, args, fmt.Errorf("missing argument '%s'", arg)
	}
	vStr := args[0]
	args = args[1:]
	return &ArgValue{
		value:     vStr,
		isDefault: false,
	}, args, nil
}
