package commands

import (
	"fmt"
	"math/rand"
	"spider/spider"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func NoyaCommand() *spider.Command {
	noya := &spider.Command{
		Name:        "noya",
		Description: "noya tools",
		Usage:       "noya...",
		Run: func(c *spider.Context) error {
			c.Spider.PrintCommandHelp(c.Command)
			return nil
		},
	}
	// add
	add := &spider.Command{
		Name:        "math-add",
		Description: "displays math-add problem",
		Usage:       "noya math-add",
		Flags: func(f *spider.Flags) {
			f.Int(&spider.Flag{Short: "l", Long: "low", Help: "the low value", Default: 10})
			f.Int(&spider.Flag{Short: "t", Long: "top", Help: "the top value", Default: 10})
			f.Int(&spider.Flag{Short: "c", Long: "count", Help: "the count of problem", Default: 10})
		},
		Run: func(c *spider.Context) error {
			low, err := c.FlagValues.Int(c.Command, "low")
			if err != nil {
				return err
			}
			top, err := c.FlagValues.Int(c.Command, "top")
			if err != nil {
				return err
			}
			count, err := c.FlagValues.Int(c.Command, "count")
			if err != nil {
				return err
			}
			out := make(chan string)
			go calculate(count, []mathOperator{
				&mathAddOperator{mathParams{low: low, top: top}},
			}, out)
			for o := range out {
				c.Spider.Println(o)
			}
			return nil
		},
	}

	// keep-add
	addKeep := &spider.Command{
		Name:        "math-add-keep",
		Description: "displays math-add-keep problem",
		Usage:       "noya math-add-keep",
		Flags: func(f *spider.Flags) {
			f.Int(&spider.Flag{Short: "l", Long: "low", Help: "the low value", Default: 10})
			f.Int(&spider.Flag{Short: "t", Long: "top", Help: "the top value", Default: 10})
			f.Int(&spider.Flag{Short: "c", Long: "count", Help: "the count of problem", Default: 10})
		},
		Run: func(c *spider.Context) error {
			low, err := c.FlagValues.Int(c.Command, "low")
			if err != nil {
				return err
			}
			top, err := c.FlagValues.Int(c.Command, "top")
			if err != nil {
				return err
			}
			count, err := c.FlagValues.Int(c.Command, "count")
			if err != nil {
				return err
			}
			out := make(chan string)
			go calculate(count, []mathOperator{
				&mathAddKeepOperator{mathParams{low: low, top: top}},
			}, out)
			for o := range out {
				c.Spider.Println(o)
			}
			return nil
		},
	}

	// sub
	sub := &spider.Command{
		Name:        "math-sub",
		Description: "displays math-sub problem",
		Usage:       "noya math-sub",
		Flags: func(f *spider.Flags) {
			f.Int(&spider.Flag{Short: "l", Long: "low", Help: "the low value", Default: 10})
			f.Int(&spider.Flag{Short: "t", Long: "top", Help: "the top value", Default: 10})
			f.Int(&spider.Flag{Short: "c", Long: "count", Help: "the count of problem", Default: 10})
		},
		Run: func(c *spider.Context) error {
			low, err := c.FlagValues.Int(c.Command, "low")
			if err != nil {
				return err
			}
			top, err := c.FlagValues.Int(c.Command, "top")
			if err != nil {
				return err
			}
			count, err := c.FlagValues.Int(c.Command, "count")
			if err != nil {
				return err
			}
			out := make(chan string)
			go calculate(count, []mathOperator{
				&mathSubOperator{mathParams{low: low, top: top}},
			}, out)
			for o := range out {
				c.Spider.Println(o)
			}
			return nil
		},
	}

	// sub
	subKeep := &spider.Command{
		Name:        "math-sub-keep",
		Description: "displays math-sub-keep problem",
		Usage:       "noya math-sub-keep",
		Flags: func(f *spider.Flags) {
			f.Int(&spider.Flag{Short: "l", Long: "low", Help: "the low value", Default: 10})
			f.Int(&spider.Flag{Short: "t", Long: "top", Help: "the top value", Default: 10})
			f.Int(&spider.Flag{Short: "c", Long: "count", Help: "the count of problem", Default: 10})
		},
		Run: func(c *spider.Context) error {
			low, err := c.FlagValues.Int(c.Command, "low")
			if err != nil {
				return err
			}
			top, err := c.FlagValues.Int(c.Command, "top")
			if err != nil {
				return err
			}
			count, err := c.FlagValues.Int(c.Command, "count")
			if err != nil {
				return err
			}
			out := make(chan string)
			go calculate(count, []mathOperator{
				&mathSubKeepOperator{mathParams{low: low, top: top}},
			}, out)
			for o := range out {
				c.Spider.Println(o)
			}
			return nil
		},
	}

	// add & sub
	addSub := &spider.Command{
		Name:        "math-add-sub",
		Description: "displays math-add & math-sub problem",
		Usage:       "noya math-add-sub",
		Flags: func(f *spider.Flags) {
			f.Int(&spider.Flag{Short: "l", Long: "low", Help: "the low value", Default: 10})
			f.Int(&spider.Flag{Short: "t", Long: "top", Help: "the top value", Default: 10})
			f.Int(&spider.Flag{Short: "c", Long: "count", Help: "the count of problem", Default: 10})
		},
		Run: func(c *spider.Context) error {
			low, err := c.FlagValues.Int(c.Command, "low")
			if err != nil {
				return err
			}
			top, err := c.FlagValues.Int(c.Command, "top")
			if err != nil {
				return err
			}
			count, err := c.FlagValues.Int(c.Command, "count")
			if err != nil {
				return err
			}
			out := make(chan string)
			go calculate(count, []mathOperator{
				&mathAddOperator{mathParams{low: low, top: top}},
				&mathSubOperator{mathParams{low: low, top: top}},
			},
				out)
			for o := range out {
				c.Spider.Println(o)
			}
			return nil
		},
	}

	// add & sub
	addSubKeep := &spider.Command{
		Name:        "math-add-sub-keep",
		Description: "displays math-addkeep & math-subkeep problem",
		Usage:       "noya math-add-sub-keep",
		Flags: func(f *spider.Flags) {
			f.Int(&spider.Flag{Short: "l", Long: "low", Help: "the low value", Default: 10})
			f.Int(&spider.Flag{Short: "t", Long: "top", Help: "the top value", Default: 10})
			f.Int(&spider.Flag{Short: "c", Long: "count", Help: "the count of problem", Default: 10})
		},
		Run: func(c *spider.Context) error {
			low, err := c.FlagValues.Int(c.Command, "low")
			if err != nil {
				return err
			}
			top, err := c.FlagValues.Int(c.Command, "top")
			if err != nil {
				return err
			}
			count, err := c.FlagValues.Int(c.Command, "count")
			if err != nil {
				return err
			}
			out := make(chan string)
			go calculate(count, []mathOperator{
				&mathAddKeepOperator{mathParams{low: low, top: top}},
				&mathSubKeepOperator{mathParams{low: low, top: top}},
				&mathAddSubKeepOperator{mathParams{low: low, top: top}},
			},
				out)
			for o := range out {
				c.Spider.Println(o)
			}
			return nil
		},
	}

	noya.AddCommand(add)
	noya.AddCommand(addKeep)
	noya.AddCommand(sub)
	noya.AddCommand(subKeep)
	noya.AddCommand(addSub)
	noya.AddCommand(addSubKeep)
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

type mathParams struct {
	low int
	top int
}

type mathAddOperator struct {
	mathParams
}

func (add *mathAddOperator) calculate() (string, error) {
	num1 := r.Intn(add.top+1-add.low) + add.low // [low,top]
	num2 := r.Intn(num1-1) + 1                  // [1, num1-1]
	num3 := num1 - num2
	return fmt.Sprintf("%d + %d = ", num2, num3), nil
}

type mathAddKeepOperator struct {
	mathParams
}

func (add *mathAddKeepOperator) calculate() (string, error) {
	num1 := r.Intn(add.top+1-add.low) + add.low // [low,top]
	num2 := r.Intn(num1-3) + 2                  // [2, num1-2]
	num3 := r.Intn(num1-num2-1) + 1
	num4 := num1 - num2 - num3
	return fmt.Sprintf("%d + %d + %d = ", num2, num3, num4), nil
}

type mathSubOperator struct {
	mathParams
}

func (sub *mathSubOperator) calculate() (string, error) {
	num1 := r.Intn(sub.top+1-sub.low) + sub.low // [low,top]
	num2 := r.Intn(num1-1) + 1                  // [1, num1-1]
	return fmt.Sprintf("%d - %d = ", num1, num2), nil
}

type mathSubKeepOperator struct {
	mathParams
}

func (sub *mathSubKeepOperator) calculate() (string, error) {
	num1 := r.Intn(sub.top+1-sub.low) + sub.low // [low,top]
	num2 := r.Intn(num1-3) + 2                  // [2, num1-2]
	num3 := r.Intn(num1-num2-1) + 1
	return fmt.Sprintf("%d - %d - %d = ", num1, num2, num3), nil
}

type mathAddSubKeepOperator struct {
	mathParams
}

func (o *mathAddSubKeepOperator) calculate() (string, error) {
	num1 := r.Intn(o.top+1-o.low) + o.low // [low,top]
	num2 := r.Intn(num1-3) + 2            // [2, num1-2]
	num3 := r.Intn(num1-num2-1) + 1
	num4 := num1 - num2 - num3
	t := r.Intn(4)
	if t == 0 {
		return fmt.Sprintf("%d + %d + %d = ", num2, num3, num4), nil
	}
	if t == 1 {
		return fmt.Sprintf("%d + %d - %d = ", num2, num3+num4, num4), nil
	}
	if t == 2 {
		return fmt.Sprintf("%d - %d + %d = ", num2+num3, num3, num4), nil
	}
	return fmt.Sprintf("%d - %d - %d = ", num1, num2, num3), nil
}
