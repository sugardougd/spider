package spider

import (
	"strings"
	"unicode"
)

const (
	CharTab = 9
)

func (s *Spider) autoComplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	if key != CharTab {
		return
	}
	prefix := line[:pos]
	args := strings.Fields(prefix)
	tip := ""
	if len(args) > 0 && !unicode.IsSpace(rune(line[pos-1])) {
		tip = args[len(args)-1]
		args = args[:len(args)-1]
	}

	// auto-completion for help
	if len(args) > 0 && args[0] == "help" {
		args = args[1:]
	}

	var commands Commands
	var flags Flags
	var suggestions []string

	if len(args) == 0 {
		commands = s.Commands
	} else {
		command, args, err := s.Commands.parse(args, false, make(FlagValues))
		if err != nil || command == nil || len(args) > 0 {
			// len(args) > 0 : still multiple parameters that have not been suggestions
			return
		}
		commands = command.Children
		flags = command.flags
	}

	if len(tip) > 0 {
		for _, cmd := range commands.list {
			if strings.HasPrefix(cmd.Name, tip) {
				suggestions = append(suggestions, strings.TrimPrefix(cmd.Name, tip))
			}
			for _, alias := range cmd.Aliases {
				if strings.HasPrefix(alias, tip) {
					suggestions = append(suggestions, strings.TrimPrefix(alias, tip))
				}
			}
		}
		for _, flag := range flags.list {
			if len(flag.Short) > 0 {
				short := "-" + flag.Short
				if strings.HasPrefix(short, tip) {
					suggestions = append(suggestions, strings.TrimPrefix(short, tip))
				}
			}
			long := "--" + flag.Long
			if strings.HasPrefix(long, tip) {
				suggestions = append(suggestions, strings.TrimPrefix(long, tip))
			}
		}
	} else {
		for _, cmd := range commands.list {
			if strings.HasPrefix(cmd.Name, tip) {
				suggestions = append(suggestions, strings.TrimPrefix(cmd.Name, tip))
			}
			for _, alias := range cmd.Aliases {
				if strings.HasPrefix(alias, tip) {
					suggestions = append(suggestions, strings.TrimPrefix(alias, tip))
				}
			}
		}
		for _, flag := range flags.list {
			if len(flag.Short) > 0 {
				suggestions = append(suggestions, "-"+flag.Short)
			}
			suggestions = append(suggestions, "--"+flag.Long)
		}
	}

	if len(suggestions) > 0 {
		ok = true
		if len(suggestions) == 1 {
			newPos = pos + len(suggestions[0])
			newLine = prefix + suggestions[0] + line[pos:]
		} else {
			for i, s := range suggestions {
				suggestions[i] = tip + s
			}
			var buffer strings.Builder
			buffer.WriteString(s.Config.Prompt)
			buffer.WriteString(line)
			buffer.WriteRune('\n')
			buffer.WriteString(strings.Join(suggestions, TAB2))
			s.Println(buffer.String())
			newPos = pos
			newLine = line
		}
	}
	return
}
