package spider

import (
	"strings"
	"unicode"
)

const (
	CharTab = 9
)

func (s *Spider) autoComplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	if key != CharTab || !s.Config.Interactive {
		return
	}
	prefix := line[:pos]
	args := strings.Fields(prefix)
	keyWord := ""
	if len(args) > 0 && !unicode.IsSpace(rune(line[pos-1])) {
		// the previous character of pressed TAB is not space
		keyWord = args[len(args)-1]
		args = args[:len(args)-1]
	}

	var commands Commands
	var flags Flags
	var suggestions []string

	if len(args) == 0 {
		commands = s.Commands
	} else {
		command, remaining, err := s.Commands.parse(args, false, make(FlagValues))
		if err != nil || command == nil || len(remaining) > 0 {
			// not match,not suggestions
			return
		}
		commands = command.Children
		flags = command.flags
	}

	// suggestion by key word(keyWord)
	if strings.HasPrefix(keyWord, "-") {
		// suggestion flags
		for _, flag := range flags.list {
			if len(flag.Short) > 0 {
				short := "-" + flag.Short
				if len(short) > len(keyWord) && strings.HasPrefix(short, keyWord) {
					suggestions = append(suggestions, short)
				}
			}
			long := "--" + flag.Long
			if len(long) > len(keyWord) && strings.HasPrefix(long, keyWord) {
				suggestions = append(suggestions, long)
			}
		}
	} else if len(keyWord) > 0 {
		// suggestion command
		for _, cmd := range commands.list {
			if cmd.Name == keyWord {
				suggestions = commandSuggestions(suggestions, cmd)
			} else if strings.HasPrefix(cmd.Name, keyWord) {
				suggestions = append(suggestions, cmd.Name)
			}
			for _, alias := range cmd.Aliases {
				if alias == keyWord {
					suggestions = commandSuggestions(suggestions, cmd)
				} else if strings.HasPrefix(alias, keyWord) {
					suggestions = append(suggestions, alias)
				}
			}
		}
	} else {
		// suggestion by space
		// suggestion sub-command
		for _, cmd := range commands.list {
			suggestions = append(suggestions, cmd.Name)
			for _, alias := range cmd.Aliases {
				suggestions = append(suggestions, alias, keyWord)
			}
		}
		// suggestion flags
		for _, flag := range flags.list {
			if len(flag.Short) > 0 {
				suggestions = append(suggestions, "-"+flag.Short)
			}
			suggestions = append(suggestions, "--"+flag.Long)
		}
	}

	if ok = len(suggestions) > 0; !ok {
		return
	}
	if len(suggestions) == 1 {
		// only one suggestion, append to the end
		newLine = prefix[:len(prefix)-len(keyWord)]
		if len(newLine) == 0 || unicode.IsSpace(rune(newLine[len(newLine)-1])) {
			newPos = len(newLine) + len(suggestions[0])
			newLine += suggestions[0] + line[pos:]
		} else {
			newPos = len(newLine) + len(suggestions[0]) + 1
			newLine += BLANK + suggestions[0] + line[pos:]
		}
	} else {
		newPos = pos
		newLine = line
		s.Print(s.Config.Prompt)
		s.Println(newLine)
		s.Println(strings.Join(suggestions, TAB2))
	}

	return
}

func commandSuggestions(suggestions []string, command *Command) []string {
	for _, child := range command.Children.list {
		suggestions = append(suggestions, child.Name)
	}
	for _, flag := range command.flags.list {
		if len(flag.Short) > 0 {
			suggestions = append(suggestions, "-"+flag.Short)
		}
		suggestions = append(suggestions, "--"+flag.Long)
	}
	return suggestions
}
