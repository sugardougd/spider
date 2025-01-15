package spider

import (
	"fmt"
	"strconv"
	"strings"
)

type FlagParser func(flag, value string, found bool, args []string) (*FlagValue, []string, error)

type Flags struct {
	list    []*Flag
	parsers map[string]FlagParser
}

type Flag struct {
	Short   string
	Long    string
	Help    string // help message for the flag.
	Usage   string // define how to use the flag.
	Require bool
	Default interface{}
}

type FlagValue struct {
	value     interface{}
	isDefault bool
}

type FlagValues map[string]*FlagValue

func (flags *Flags) parse(command *Command, args []string, required bool, flagValues FlagValues) (remaining []string, err error) {
	for len(args) > 0 {
		name := args[0]
		if !strings.HasPrefix(name, "-") {
			break
		}
		// find flag
		pos := strings.Index(name, "=")
		value := ""
		found := false
		if pos > 0 {
			value = name[pos+1:]
			found = true
			name = name[:pos]
		}
		flag := flags.Find(name)
		if flag == nil {
			break
		}
		args = args[1:]
		// find FlagParser
		fullName := flag.fullName(command.FullName())
		parser, ok := flags.parsers[flag.Long]
		if ok {
			var flagValue *FlagValue
			flagValue, args, err = parser(name, value, found, args)
			if err != nil {
				break
			}
			flagValues[fullName] = flagValue
		}
	}
	remaining = args

	// check require argument and set the default values
	for _, flag := range flags.list {
		fullName := flag.fullName(command.FullName())
		if _, ok := flagValues[fullName]; ok {
			continue
		}
		if flag.Require && required {
			err = fmt.Errorf("missing flag '-%s'", flag.Short)
			break
		}
		flagValues[fullName] = &FlagValue{
			value:     flag.Default,
			isDefault: true,
		}
	}
	return
}

func (flags *Flags) Find(f string) *Flag {
	for _, flag := range flags.list {
		if flag.match(f) {
			return flag
		}
	}
	return nil
}

func (flags *Flags) MatchPrefix(f string) []*Flag {
	var match []*Flag
	for _, flag := range flags.list {
		if flag.matchPrefix(f) {
			match = append(match, flag)
		}
	}
	return match
}

func (flags *Flags) register(flag *Flag, fp FlagParser) error {
	err := flag.validate()
	if err != nil {
		return err
	}
	for _, f := range flags.list {
		if len(flag.Short) > 0 && len(f.Short) > 0 && flag.Short == f.Short {
			return fmt.Errorf("flag '%s' registered twice", flag.Short)
		}
		if flag.Long == f.Long {
			return fmt.Errorf("flag '%s' registered twice", flag.Long)
		}
	}
	flags.list = append(flags.list, flag)
	if flags.parsers == nil {
		flags.parsers = make(map[string]FlagParser)
	}
	flags.parsers[flag.Long] = fp
	return nil
}

// Bool register a bool flag
func (flags *Flags) Bool(flag *Flag) error {
	return flags.register(flag, parseBoolFlag)
}

// Int register an int flag
func (flags *Flags) Int(flag *Flag) error {
	return flags.register(flag, parseIntFlag)
}

// Float64 register an float64 flag
func (flags *Flags) Float64(flag *Flag) error {
	return flags.register(flag, parseFloat64Flag)
}

// String register a string flag
func (flags *Flags) String(flag *Flag) error {
	return flags.register(flag, parseStringFlag)
}

func (flag *Flag) match(f string) bool {
	return (len(flag.Short) > 0 && f == "-"+flag.Short) ||
		(len(flag.Long) > 0 && f == "--"+flag.Long)
}

func (flag *Flag) matchPrefix(f string) bool {
	return (len(flag.Short) > 0 && strings.HasPrefix("-"+flag.Short, f)) ||
		(len(flag.Long) > 0 && strings.HasPrefix("--"+flag.Long, f))
}

func (flag *Flag) fullName(parent string) string {
	return parent + "." + flag.Long
}

func (flag *Flag) validate() error {
	if len(flag.Short) > 1 {
		return fmt.Errorf("invalid short flag: '%s': must be a single character", flag.Short)
	} else if strings.HasPrefix(flag.Short, "-") {
		return fmt.Errorf("invalid short flag: '%s': must not start with a '-'", flag.Short)
	} else if len(flag.Long) == 0 {
		return fmt.Errorf("empty long flag: long='%s'", flag.Long)
	} else if strings.HasPrefix(flag.Long, "-") {
		return fmt.Errorf("invalid long flag: '%s': must not start with a '-'", flag.Long)
	} else if len(flag.Help) == 0 {
		return fmt.Errorf("empty flag help message")
	}
	return nil
}

func parseBoolFlag(flag, value string, found bool, args []string) (*FlagValue, []string, error) {
	if found {
		v, err := strconv.ParseBool(value)
		if err != nil {
			return nil, args, fmt.Errorf("invalid boolean value for flag: %s", flag)
		}
		return &FlagValue{
			value:     v,
			isDefault: false,
		}, args, nil
	}
	return &FlagValue{
		value:     true,
		isDefault: false,
	}, args, nil
}

func parseIntFlag(flag, value string, found bool, args []string) (*FlagValue, []string, error) {
	vStr := value
	if !found {
		if len(args) > 0 {
			vStr = args[0]
			args = args[1:]
		} else {
			return nil, args, fmt.Errorf("missing int value for flag: %s", flag)
		}
	}
	b, err := strconv.Atoi(vStr)
	if err != nil {
		return nil, args, fmt.Errorf("invalid int value for flag: %s", flag)
	}
	return &FlagValue{
		value:     b,
		isDefault: false,
	}, args, nil
}

func parseFloat64Flag(flag, value string, found bool, args []string) (*FlagValue, []string, error) {
	vStr := value
	if !found {
		if len(args) > 0 {
			vStr = args[0]
			args = args[1:]
		} else {
			return nil, args, fmt.Errorf("missing float value for flag: %s", flag)
		}
	}
	b, err := strconv.ParseFloat(vStr, 64)
	if err != nil {
		return nil, args, fmt.Errorf("invalid float value for flag: %s", flag)
	}
	return &FlagValue{
		value:     b,
		isDefault: false,
	}, args, nil
}

func parseStringFlag(flag, value string, found bool, args []string) (*FlagValue, []string, error) {
	vStr := value
	if !found {
		if len(args) > 0 {
			vStr = args[0]
			args = args[1:]
		} else {
			return nil, args, fmt.Errorf("missing string value for flag: %s", flag)
		}
	}
	return &FlagValue{
		value:     vStr,
		isDefault: false,
	}, args, nil
}

// Bool parent is Command FullName
func (flags FlagValues) Bool(parent, long string) (bool, error) {
	fullName := parent + "." + long
	flagVal := flags[fullName]
	if flagVal == nil {
		return false, fmt.Errorf("missing flag value: flag '%s' not registered", long)
	}
	v, ok := flagVal.value.(bool)
	if !ok {
		return false, fmt.Errorf("failed to assert flag '%s' to bool", long)
	}
	return v, nil
}

func (flags FlagValues) Int(parent, long string) (int, error) {
	fullName := parent + "." + long
	flagVal := flags[fullName]
	if flagVal == nil {
		return 0, fmt.Errorf("missing flag value: flag '%s' not registered", long)
	}
	v, ok := flagVal.value.(int)
	if !ok {
		return 0, fmt.Errorf("failed to assert flag '%s' to int", long)
	}
	return v, nil
}

func (flags FlagValues) Float64(parent, long string) (float64, error) {
	fullName := parent + "." + long
	flagVal := flags[fullName]
	if flagVal == nil {
		return 0, fmt.Errorf("missing flag value: flag '%s' not registered", long)
	}
	v, ok := flagVal.value.(float64)
	if !ok {
		return 0, fmt.Errorf("failed to assert flag '%s' to float64", long)
	}
	return v, nil
}

func (flags FlagValues) String(parent, long string) (string, error) {
	fullName := parent + "." + long
	flagVal := flags[fullName]
	if flagVal == nil {
		return "", fmt.Errorf("missing flag value: flag '%s' not registered", long)
	}
	v, ok := flagVal.value.(string)
	if !ok {
		return "", fmt.Errorf("failed to assert flag '%s' to string", long)
	}
	return v, nil
}
