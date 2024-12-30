package main

import (
	"fmt"
	"spider/terminal"
	"strings"
	"time"
)

func main() {
	TestTerminal()
	//TestMockTerminal()
}

func TestTerminal() {
	ter := terminal.NewConsole("$ ")

	for {
		select {
		case line := <-ter.Readline():
			command := strings.TrimSpace(line)
			if _, err := ter.WriteLine([]byte("type: " + command)); err != nil {
				fmt.Printf("Write fail %v", err)
			}
			ter.WritePrompt()
		case <-ter.Done():
			break
		}
	}
}

func TestMockTerminal() {
	stdio := &MockReadWriter{
		in: make(chan byte),
	}
	config := terminal.Config{
		Prompt: "$ ",
		Stdin:  stdio,
		Stdout: stdio,
		Stderr: stdio,
	}
	ter := terminal.New(&config)
	go func() {
		stdio.in <- byte('1')
		stdio.in <- terminal.KeyEnter
		stdio.in <- byte('2')
		stdio.in <- terminal.KeyEnter
	}()
	go func() {
		for {
			select {
			case line := <-ter.Readline():
				command := strings.TrimSpace(line)
				if i, err := ter.WriteLine([]byte("type: " + command)); err == nil {
					fmt.Printf("Write %d bytes\r\n", i)
				} else {
					fmt.Printf("Write fail %v", err)
				}
				ter.WritePrompt()
			case <-ter.Done():
				break
			}
		}
	}()
	time.Sleep(time.Second * 10)
	ter.Stop()
	time.Sleep(time.Second * 5)
}

type MockReadWriter struct {
	in chan byte
}

func (mock *MockReadWriter) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	p[0] = <-mock.in
	n = 1
	return n, err
}

func (mock *MockReadWriter) Close() error {
	close(mock.in)
	fmt.Println("close mock.in")
	return nil
}

func (mock *MockReadWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	fmt.Print(string(p))
	return n, err
}
