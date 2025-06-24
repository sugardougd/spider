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

// BoolList registers a bool list argument.
func (a *Args) BoolList(arg *Arg) {
	a.register(arg, parseBoolListArg)
}

// Int registers an int argument.
func (a *Args) Int(arg *Arg) {
	a.register(arg, parseIntArg)
}

// IntList registers an int list argument.
func (a *Args) IntList(arg *Arg) {
	a.register(arg, parseIntListArg)
}

// Float64 registers a float argument.
func (a *Args) Float64(arg *Arg) {
	a.register(arg, parseFloat64Arg)
}

// Float64List registers a float list argument.
func (a *Args) Float64List(arg *Arg) {
	a.register(arg, parseFloat64ListArg)
}

// String registers a string argument.
func (a *Args) String(arg *Arg) {
	a.register(arg, parseStringArg)
}

// StringList registers a string list argument.
func (a *Args) StringList(arg *Arg) {
	a.register(arg, parseStringListArg)
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

func parseBoolListArg(arg string, args []string) (*ArgValue, []string, error) {
	if len(args) == 0 {
		return nil, args, fmt.Errorf("missing argument '%s'", arg)
	}
	var list []bool
	for _, a := range args {
		v, err := strconv.ParseBool(a)
		if err != nil {
			return nil, args, fmt.Errorf("invalid boolean value for argument: %s", arg)
		}
		list = append(list, v)
	}
	args = args[:0]
	return &ArgValue{
		value:     list,
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

func parseIntListArg(arg string, args []string) (*ArgValue, []string, error) {
	if len(args) == 0 {
		return nil, args, fmt.Errorf("missing argument '%s'", arg)
	}
	var list []int
	for _, a := range args {
		v, err := strconv.Atoi(a)
		if err != nil {
			return nil, args, fmt.Errorf("invalid int value for argument: %s", arg)
		}
		list = append(list, v)
	}
	args = args[:0]
	return &ArgValue{
		value:     list,
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

func parseFloat64ListArg(arg string, args []string) (*ArgValue, []string, error) {
	if len(args) == 0 {
		return nil, args, fmt.Errorf("missing argument '%s'", arg)
	}
	var list []float64
	for _, a := range args {
		v, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return nil, args, fmt.Errorf("invalid float value for argument: %s", arg)
		}
		list = append(list, v)
	}
	args = args[:0]
	return &ArgValue{
		value:     list,
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

func parseStringListArg(arg string, args []string) (*ArgValue, []string, error) {
	if len(args) == 0 {
		return nil, args, fmt.Errorf("missing argument '%s'", arg)
	}
	vStr := args
	args = args[:0]
	return &ArgValue{
		value:     vStr,
		isDefault: false,
	}, args, nil
}

func (args ArgValues) Bool(name string) (bool, error) {
	argVal := args[name]
	if argVal == nil {
		return false, fmt.Errorf("missing arg value: arg '%s' not registered", name)
	}
	if argVal.value == nil {
		return false, nil
	}
	if v, ok := argVal.value.(bool); ok {
		return v, nil
	}
	return false, fmt.Errorf("failed to assert arg '%s' to bool", name)
}

func (args ArgValues) BoolList(name string) ([]bool, error) {
	argVal := args[name]
	if argVal == nil {
		return nil, fmt.Errorf("missing arg value: arg '%s' not registered", name)
	}
	if argVal.value == nil {
		return nil, nil
	}
	if v, ok := argVal.value.([]bool); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to assert arg '%s' to bool list", name)
}

func (args ArgValues) Int(name string) (int, error) {
	argVal := args[name]
	if argVal == nil {
		return 0, fmt.Errorf("missing arg value: arg '%s' not registered", name)
	}
	if argVal.value == nil {
		return 0, nil
	}
	if v, ok := argVal.value.(int); ok {
		return v, nil
	}
	return 0, fmt.Errorf("failed to assert arg '%s' to int", name)
}

func (args ArgValues) IntList(name string) ([]int, error) {
	argVal := args[name]
	if argVal == nil {
		return nil, fmt.Errorf("missing arg value: arg '%s' not registered", name)
	}
	if argVal.value == nil {
		return nil, nil
	}
	if v, ok := argVal.value.([]int); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to assert arg '%s' to int list", name)
}

func (args ArgValues) Float64(name string) (float64, error) {
	argVal := args[name]
	if argVal == nil {
		return 0, fmt.Errorf("missing arg value: arg '%s' not registered", name)
	}
	if argVal.value == nil {
		return 0, nil
	}
	if v, ok := argVal.value.(float64); ok {
		return v, nil
	}
	return 0, fmt.Errorf("failed to assert arg '%s' to float64", name)
}

func (args ArgValues) Float64List(name string) ([]float64, error) {
	argVal := args[name]
	if argVal == nil {
		return nil, fmt.Errorf("missing arg value: arg '%s' not registered", name)
	}
	if argVal.value == nil {
		return nil, nil
	}
	if v, ok := argVal.value.([]float64); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to assert arg '%s' to float64 list", name)
}

func (args ArgValues) String(name string) (string, error) {
	argVal := args[name]
	if argVal == nil {
		return "", fmt.Errorf("missing arg value: arg '%s' not registered", name)
	}
	if argVal.value == nil {
		return "", nil
	}
	if v, ok := argVal.value.(string); ok {
		return v, nil
	}
	return "", fmt.Errorf("failed to assert arg '%s' to string", name)
}

func (args ArgValues) StringList(name string) ([]string, error) {
	argVal := args[name]
	if argVal == nil {
		return nil, fmt.Errorf("missing arg value: arg '%s' not registered", name)
	}
	if argVal.value == nil {
		return nil, nil
	}
	if v, ok := argVal.value.([]string); ok {
		return v, nil
	}
	return nil, fmt.Errorf("failed to assert arg '%s' to string list", name)
}
