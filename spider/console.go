package spider

import (
	"fmt"
	"github.com/desertbit/readline"
	"os"
	"spider/grumble"
)

type ConsoleSpider struct {
	cfg *ConsoleConfig
}
type ConsoleConfig struct {
	*grumble.Config
	Commands CommandsSet
}

// RunConsole run console cli
func RunConsole(cfg *ConsoleConfig) error {
	fmt.Println("New Console cli")
	s := &ConsoleSpider{cfg: cfg}
	if cfg.FuncIsTerminal == nil {
		cfg.FuncIsTerminal = s.FuncIsTerminal
	}
	// 执行命令
	rlCfg := &readline.Config{
		Prompt:         "$ ",
		Stdin:          os.Stdin,
		Stdout:         os.Stdout,
		Stderr:         os.Stderr,
		FuncIsTerminal: s.cfg.FuncIsTerminal,
	}
	rl, err := readline.NewEx(rlCfg)
	if err != nil {
		return err
	}
	app := NewGrumble(cfg.Config, cfg.Commands)
	if err := app.RunWithReadline(rl); err != nil {
		fmt.Printf("run %s fail. %v\r\n", cfg.Name, err)
		return err
	}
	return nil
}

func (s *ConsoleSpider) FuncIsTerminal() bool {
	return true
}
