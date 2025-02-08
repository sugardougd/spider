package main

import (
	"fmt"
	"os"
	"spider/spider"
	"spider/spider/example/commands"
)

func main() {
	config := spider.NewConfig(
		spider.ConfigName("spider"),
		spider.ConfigDescription("spider is a tool to list and diagnose Go processes"),
		spider.ConfigPrompt("spider > "))
	s := spider.New(config, spider.NewCommands(commands.NoyaCommand(), commands.TestCommand()))
	cmd := []string{
		"help",
		//"help test",
		//"test",
		//"test -a -b 1 -c 1.1 -d 1 false 2 3 4",
		//"test --aaa --bbb 1 --ccc 1.1 --ddd 1 true 2 3 4",
		//"test -a=false -b=1 -c=1.1 -d=1 true 2",
		//"test --aaa=false --bbb=1 --ccc=1.1 --ddd=1 true",
		//"test subtest false 2 3 4",
		//"test subtest -a -b 1 -c 1.1 -d 1 false 2 3 4",
		//"test subtest --aaa --bbb 1 --ccc 1.1 --ddd 1 true 2 3 4",
		//"test subtest -a=false -b=1 -c=1.1 -d=1 true 2",
		//"test subtest --aaa=false --bbb=1 --ccc=1.1 --ddd=1 true",
		//"test --aaa=false --bbb=1 --ccc=1.1 --ddd=1 subtest --aaa=false --bbb=1 --ccc=1.1 --ddd=1",
		//"err",
		//"exit",
	}
	for _, c := range cmd {
		if err := s.RunCommand(c); err != nil {
			fmt.Println("Run fail: ", err)
			os.Exit(0)
		}
	}
}
