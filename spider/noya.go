package spider

import (
	"fmt"
	"math/rand"
	"spider/grumble"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func NoyaCommands() *grumble.Commands {
	commands := grumble.Commands{}
	nc := &grumble.Command{
		Name: "noya",
		Help: "noya tools",
		Run: func(c *grumble.Context) error {
			c.App.PrintCommandHelp(c.Command, true)
			return nil
		},
	}
	// add
	add := &grumble.Command{
		Name: "math-add",
		Help: "displays math-add problem",
		Flags: func(f *grumble.Flags) {
			f.Int("t", "top", 10, "the sum top")
			f.Int("c", "count", 10, "the count of problem")
		},
		Run: func(c *grumble.Context) error {
			top := c.Flags.Int("top")
			count := c.Flags.Int("count")
			out := make(chan string)
			go calculate(count, []mathOperator{&mathAddOperator{top}}, out)
			for o := range out {
				c.App.Println(o)
			}
			return nil
		},
	}
	// aub
	sub := &grumble.Command{
		Name: "math-sub",
		Help: "displays math-sub problem",
		Flags: func(f *grumble.Flags) {
			f.Int("t", "top", 10, "the sum top")
			f.Int("c", "count", 10, "the count of problem")
		},
		Run: func(c *grumble.Context) error {
			top := c.Flags.Int("top")
			count := c.Flags.Int("count")
			out := make(chan string)
			go calculate(count, []mathOperator{&mathSubOperator{top}}, out)
			for o := range out {
				c.App.Println(o)
			}
			return nil
		},
	}
	// add & sub
	addSub := &grumble.Command{
		Name: "math-add-sub",
		Help: "displays math-add & math-sub problem",
		Flags: func(f *grumble.Flags) {
			f.Int("t", "top", 10, "the sum top")
			f.Int("c", "count", 10, "the count of problem")
		},
		Run: func(c *grumble.Context) error {
			top := c.Flags.Int("top")
			count := c.Flags.Int("count")
			out := make(chan string)
			go calculate(count, []mathOperator{&mathAddOperator{top}, &mathSubOperator{top}}, out)
			for o := range out {
				c.App.Println(o)
			}
			return nil
		},
	}
	nc.AddCommand(add)
	nc.AddCommand(sub)
	nc.AddCommand(addSub)
	commands.Add(nc)
	return &commands
}

func calculate(c int, operators []mathOperator, out chan string) {
	for i := 0; i < c; i++ {
		r.Intn(len(operators))
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
	sum int
}

func (add *mathAddOperator) calculate() (string, error) {
	num1 := r.Intn(add.sum)
	num2 := add.sum - num1
	//return fmt.Sprintf("%2d + %2d = ", num1, num2), nil
	return fmt.Sprintf("%d + %d = ", num1, num2), nil
}

type mathSubOperator struct {
	sum int
}

func (sub *mathSubOperator) calculate() (string, error) {
	num1 := r.Intn(sub.sum)
	//return fmt.Sprintf("%2d - %2d = ", sub.sum, num1), nil
	return fmt.Sprintf("%d - %d = ", sub.sum, num1), nil
}
