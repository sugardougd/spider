package spider

import (
	"testing"
)

func TestSpider_autoComplete(t *testing.T) {
	config := NewConfig(
		ConfigName("spider"),
		ConfigDescription("spider is a tool to list and diagnose Go processes"),
		ConfigPrompt("spider > "))
	s := New(config, testCommand())

	tests := []struct {
		line    string
		pos     int
		key     rune
		newLine string
		newPos  int
		ok      bool
	}{
		{
			line:    "",
			pos:     0,
			key:     CharTab,
			newLine: "",
			newPos:  0,
			ok:      true,
		},
		{
			line:    "hello",
			pos:     5,
			key:     CharTab,
			newLine: "",
			newPos:  0,
			ok:      false,
		},
		{
			line:    "command",
			pos:     7,
			key:     CharTab,
			newLine: "command",
			newPos:  7,
			ok:      true,
		},
		{
			line:    "command ",
			pos:     8,
			key:     CharTab,
			newLine: "command ",
			newPos:  8,
			ok:      true,
		},
		{
			line:    "command -",
			pos:     9,
			key:     CharTab,
			newLine: "command -",
			newPos:  9,
			ok:      true,
		},
		{
			line:    "command -o",
			pos:     10,
			key:     CharTab,
			newLine: "command -o",
			newPos:  10,
			ok:      true,
		},
		{
			line:    "command --",
			pos:     10,
			key:     CharTab,
			newLine: "command --",
			newPos:  10,
			ok:      true,
		},
		{
			line:    "command --f",
			pos:     11,
			key:     CharTab,
			newLine: "command --f",
			newPos:  11,
			ok:      true,
		},
		{
			line:    "command --flag1",
			pos:     15,
			key:     CharTab,
			newLine: "command --flag1",
			newPos:  15,
			ok:      true,
		},
		{
			line:    "command sub",
			pos:     11,
			key:     CharTab,
			newLine: "command sub",
			newPos:  11,
			ok:      true,
		},
		{
			line:    "command subcommand",
			pos:     11,
			key:     CharTab,
			newLine: "command subcommand",
			newPos:  11,
			ok:      true,
		},
		{
			line:    "command sub-command1 -",
			pos:     22,
			key:     CharTab,
			newLine: "command sub-command1 -",
			newPos:  22,
			ok:      true,
		},
		{
			line:    "command sub-command ",
			pos:     20,
			key:     CharTab,
			newLine: "command sub-command ",
			newPos:  20,
			ok:      true,
		},
		{
			line:    "command1",
			pos:     7,
			key:     CharTab,
			newLine: "command1",
			newPos:  7,
			ok:      true,
		},
		{
			line:    "help ",
			pos:     5,
			key:     CharTab,
			newLine: "help ",
			newPos:  5,
			ok:      true,
		},
		{
			line:    "exit ",
			pos:     5,
			key:     CharTab,
			newLine: "exit ",
			newPos:  5,
			ok:      true,
		},
	}

	for _, test := range tests {
		newLine, newPos, ok := s.autoComplete(test.line, test.pos, test.key)
		if newLine != test.newLine || newPos != test.newPos || ok != test.ok {
			t.Fatalf("'%s' auto-complete fail. newLine=%s;newPos=%d;ok=%t", test.line, newLine, newPos, ok)
		}
	}
}
