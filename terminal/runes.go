package terminal

import (
	"bytes"
	"unicode"
)

var TabWidth = 4
var zeroWidth = []*unicode.RangeTable{
	unicode.Mn,
	unicode.Me,
	unicode.Cc,
	unicode.Cf,
}

var doubleWidth = []*unicode.RangeTable{
	unicode.Han,
	unicode.Hangul,
	unicode.Hiragana,
	unicode.Katakana,
}

var runes = Runes{}

type Runes struct{}

func (Runes) Width(r rune) int {
	if r == '\t' {
		return TabWidth
	}
	if unicode.IsOneOf(zeroWidth, r) {
		return 0
	}
	if unicode.IsOneOf(doubleWidth, r) {
		return 2
	}
	return 1
}

func (Runes) WidthAll(r []rune) (length int) {
	for i := 0; i < len(r); i++ {
		length += runes.Width(r[i])
	}
	return
}

func (Runes) Backspace(r []rune) []byte {
	return bytes.Repeat([]byte{'\b'}, runes.WidthAll(r))
}

func (Runes) Copy(r []rune) []rune {
	n := make([]rune, len(r))
	copy(n, r)
	return n
}
