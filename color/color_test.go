package color

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestColor_Sprintf(t *testing.T) {
	fmt.Println(New().Sprintf("line"))
	fmt.Println(New(Faint).Sprintf("line"))
	fmt.Println(New(FgRed).Sprintf("line"))
	fmt.Println(New(FgRed, BgWhite).Sprintf("line"))
	fmt.Println(New(FgRed, BgWhite, ReverseVideo).Sprintf("line %v", time.Now()))
}

func TestColor_Sprintln(t *testing.T) {
	fmt.Println(New(FgRed).Sprintln("Hello"))
}

func TestColor_Printf(t *testing.T) {
	New(FgRed).Printf(os.Stdout, "Hello %v", time.Now())
}

func TestColor_Printfln(t *testing.T) {
	New(FgRed).Printfln(os.Stdout, "Hello")
}

func TestColor_Accept(t *testing.T) {
	underline := New(Underline)
	hiGreen := New(FgHiGreen)
	underline.Accept(hiGreen).Printfln(os.Stdout, "Hello")
}

func TestColor_SecondFormat(t *testing.T) {
	underline := New(Underline)
	hiGreen := New(FgHiGreen)
	line := underline.Sprintf("line1")
	fmt.Println(line)
	line = hiGreen.Sprintf(line + "line2")
	fmt.Println(line)
}
