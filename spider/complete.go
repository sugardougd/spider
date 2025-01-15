package spider

import (
	"strings"
)

const (
	CharTab = 9
)

func (s *Spider) autoComplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	if key != CharTab {
		return
	}
	prefix := line[:pos]
	if len(prefix) == 0 {
		return
	}

	args := strings.Fields(prefix)
	flagValues := make(FlagValues)
	command, args, err := s.Commands.parse(args, true, flagValues)
	if err != nil {
		return
	}
	if len(args) != 1 {
		// still multiple parameters that have not been matched
		return
	}

	arg := args[0]
	var matched string
	if command == nil {
		// match command
		if match := s.Commands.MatchPrefix(arg); len(match) == 1 {
			ok = true
			matched = match[0].Name
		}
	} else {
		if strings.HasPrefix(arg, "-") {
			// match flag
			if match := command.flags.MatchPrefix(arg); len(match) == 1 {
				ok = true
				if strings.HasPrefix(arg, "--") {
					matched = "--" + match[0].Long
				} else {
					matched = "-" + match[0].Short
				}
			}
		} else {
			// match sub command
			if match := command.Children.MatchPrefix(arg); len(match) == 1 {
				ok = true
				matched = match[0].Name
			}
		}
	}

	if ok {
		newLine = prefix
		if after, found := strings.CutPrefix(matched, arg); found {
			newLine += after
		}
		newPos = len(newLine)
		newLine += line[pos:]
	}
	return
}
