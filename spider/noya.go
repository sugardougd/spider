package spider

import (
	"fmt"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func NoyaCommand() *Command {
	noya := &Command{
		Name:        "noya",
		Description: "noya tools",
		Run: func(c *Context) error {
			c.Spider.PrintCommandHelp(c.Command)
			return nil
		},
	}
	// add
	add := &Command{
		Name:        "math-add",
		Description: "displays math-add problem",
		Flags: func(f *Flags) {
			f.Int(&Flag{Short: "l", Long: "low", Help: "the low value", Default: 10})
			f.Int(&Flag{Short: "t", Long: "top", Help: "the top value", Default: 10})
			f.Int(&Flag{Short: "c", Long: "count", Help: "the count of problem", Default: 10})
		},
		Run: func(c *Context) error {
			fullName := c.Command.FullName()
			low, err := c.FlagValues.Int(fullName, "low")
			if err != nil {
				return err
			}
			top, err := c.FlagValues.Int(fullName, "top")
			if err != nil {
				return err
			}
			count, err := c.FlagValues.Int(fullName, "count")
			if err != nil {
				return err
			}
			out := make(chan string)
			go calculate(count, []mathOperator{&mathAddOperator{low, top}}, out)
			for o := range out {
				c.Spider.Println(o)
			}
			return nil
		},
	}
	// aub
	sub := &Command{
		Name:        "math-sub",
		Description: "displays math-sub problem",
		Flags: func(f *Flags) {
			f.Int(&Flag{Short: "l", Long: "low", Help: "the low value", Default: 10})
			f.Int(&Flag{Short: "t", Long: "top", Help: "the top value", Default: 10})
			f.Int(&Flag{Short: "c", Long: "count", Help: "the count of problem", Default: 10})
		},
		Run: func(c *Context) error {
			fullName := c.Command.FullName()
			low, err := c.FlagValues.Int(fullName, "low")
			if err != nil {
				return err
			}
			top, err := c.FlagValues.Int(fullName, "top")
			if err != nil {
				return err
			}
			count, err := c.FlagValues.Int(fullName, "count")
			if err != nil {
				return err
			}
			out := make(chan string)
			go calculate(count, []mathOperator{&mathSubOperator{low, top}}, out)
			for o := range out {
				c.Spider.Println(o)
			}
			return nil
		},
	}
	// add & sub
	addSub := &Command{
		Name:        "math-add-sub",
		Description: "displays math-add & math-sub problem",
		Flags: func(f *Flags) {
			f.Int(&Flag{Short: "l", Long: "low", Help: "the low value", Default: 10})
			f.Int(&Flag{Short: "t", Long: "top", Help: "the top value", Default: 10})
			f.Int(&Flag{Short: "c", Long: "count", Help: "the count of problem", Default: 10})
		},
		Run: func(c *Context) error {
			fullName := c.Command.FullName()
			low, err := c.FlagValues.Int(fullName, "low")
			if err != nil {
				return err
			}
			top, err := c.FlagValues.Int(fullName, "top")
			if err != nil {
				return err
			}
			count, err := c.FlagValues.Int(fullName, "count")
			if err != nil {
				return err
			}
			out := make(chan string)
			go calculate(count, []mathOperator{&mathAddOperator{low, top}, &mathSubOperator{low, top}}, out)
			for o := range out {
				c.Spider.Println(o)
			}
			return nil
		},
	}
	noya.AddCommand(add)
	noya.AddCommand(sub)
	noya.AddCommand(addSub)
	return noya
}

func calculate(c int, operators []mathOperator, out chan string) {
	for i := 0; i < c; i++ {
		if s, err := operators[r.Intn(len(operators))].calculate(); err == nil {
			out <- s
		}
	}
	close(out)
}

type mathOperator interface {
	calculate() (string, error)
}

type mathAddOperator struct {
	low int
	top int
}

func (add *mathAddOperator) calculate() (string, error) {
	num1 := r.Intn(add.top+1-add.low) + add.low
	num2 := r.Intn(num1 + 1)
	num3 := num1 - num2
	return fmt.Sprintf("%d + %d = ", num2, num3), nil
}

type mathSubOperator struct {
	low int
	top int
}

func (sub *mathSubOperator) calculate() (string, error) {
	num1 := r.Intn(sub.top+1-sub.low) + sub.low
	num2 := r.Intn(num1 + 1)
	return fmt.Sprintf("%d - %d = ", num1, num2), nil
}
