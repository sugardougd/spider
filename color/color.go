package color

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const escape = "\x1b"

const (
	Reset        Attribute = iota // 重置
	Bold                          // 加粗或增加强度
	Faint                         // 减弱（灰色/暗色）
	Italic                        // 斜体
	Underline                     // 下划线
	BlinkSlow                     // 慢速闪烁（不广泛支持）
	BlinkRapid                    // 快速闪烁（不广泛支持）
	ReverseVideo                  // 反显（前景色和背景色交换）
	Concealed                     // 隐藏（不广泛支持）
	CrossedOut                    // 删除线（不广泛支持）
)

const (
	ResetBold       Attribute = iota + 22 // 重置
	ResetItalic                           // 重置
	ResetUnderline                        // 重置
	ResetBlinking                         // 重置
	_                                     // 重置
	ResetReversed                         // 重置
	ResetConcealed                        // 重置
	ResetCrossedOut                       // 重置
)

const (
	FgBlack   Attribute = iota + 30 // 黑色
	FgRed                           // 红色
	FgGreen                         // 绿色
	FgYellow                        // 黄色
	FgBlue                          // 蓝色
	FgMagenta                       // 品红（紫色）
	FgCyan                          // 青色
	FgWhite                         // 白色
)

const (
	FgHiBlack   Attribute = iota + 90 // 亮黑
	FgHiRed                           // 亮红
	FgHiGreen                         // 亮绿
	FgHiYellow                        // 亮黄
	FgHiBlue                          // 亮蓝
	FgHiMagenta                       // 亮品红
	FgHiCyan                          // 亮青
	FgHiWhite                         // 亮白
)

const (
	BgBlack   Attribute = iota + 40 // 黑色背景
	BgRed                           // 红色背景
	BgGreen                         // 绿色背景
	BgYellow                        // 黄色背景
	BgBlue                          // 蓝色背景
	BgMagenta                       // 品红背景
	BgCyan                          // 青色背景
	BgWhite                         // 白色背景
)

const (
	BgHiBlack   Attribute = iota + 100 // 亮黑背景
	BgHiRed                            // 亮红背景
	BgHiGreen                          // 亮绿背景
	BgHiYellow                         // 亮黄背景
	BgHiBlue                           // 亮蓝背景
	BgHiMagenta                        // 亮品红背景
	BgHiCyan                           // 亮青背景
	BgHiWhite                          // 亮白背景
)

var resetAttributes = map[Attribute]Attribute{
	Bold:         ResetBold,
	Faint:        ResetBold,
	Italic:       ResetItalic,
	Underline:    ResetUnderline,
	BlinkSlow:    ResetBlinking,
	BlinkRapid:   ResetBlinking,
	ReverseVideo: ResetReversed,
	Concealed:    ResetConcealed,
	CrossedOut:   ResetCrossedOut,
}

// Color defines a custom color object which is defined by SGR parameters.
type Color struct {
	params  []Attribute
	noColor bool
}

type Attribute int

func New(value ...Attribute) *Color {
	c := &Color{
		params: make([]Attribute, 0),
	}
	c.Add(value...)
	return c
}

// Add is used to chain SGR parameters. Use as many as parameters to combine
// and create custom color objects. Example: Add(color.FgRed, color.Underline).
func (c *Color) Add(value ...Attribute) *Color {
	c.params = append(c.params, value...)
	c.noColor = len(c.params) == 0
	return c
}

func (c *Color) Accept(other *Color) *Color {
	c.Add(other.params...)
	return c
}

// DisableColor disables the color output
func (c *Color) DisableColor() {
	c.noColor = true
}

// EnableColor enables the color output
func (c *Color) EnableColor() {
	c.noColor = false
}

func (c *Color) Printf(w io.Writer, format string, a ...interface{}) (int, error) {
	return fmt.Fprint(w, c.Sprintf(format, a...))
}

func (c *Color) Printfln(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprint(w, c.Sprintln(a...))
}

// Sprintf return string with the colors attributes.
func (c *Color) Sprintf(format string, a ...interface{}) string {
	return c.wrap(fmt.Sprintf(format, a...))
}

// Sprintln return string with the colors attributes.
func (c *Color) Sprintln(a ...interface{}) string {
	return c.wrap(fmt.Sprintln(a...))
}

// wrap wraps the s string with the colors attributes.
func (c *Color) wrap(s string) string {
	if c.noColor {
		return s
	}
	return c.format() + s + c.unformat()
}

func (c *Color) format() string {
	return fmt.Sprintf("%s[%sm", escape, c.sequence())
}

func (c *Color) unformat() string {
	//for each element in sequence let's use the specific reset escape, or the generic one if not found
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(Reset))
		ra, ok := resetAttributes[v]
		if ok {
			format[i] = strconv.Itoa(int(ra))
		} else {
			return fmt.Sprintf("%s[%dm", escape, Reset)
		}
	}
	return fmt.Sprintf("%s[%sm", escape, strings.Join(format, ";"))
}

// sequence returns a formatted SGR sequence to be plugged into a "\x1b[...m"
func (c *Color) sequence() string {
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(v))
	}

	return strings.Join(format, ";")
}
