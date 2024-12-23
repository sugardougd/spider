package main

import (
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"strings"
	"terminal"
	"time"
)

func main() {
	TestTerminal()
	//TestMockTerminal()
	//TestOsTerminal()
}

func TestTerminal() {
	config := terminal.Config{
		Prompt: "$ ",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	ter := terminal.New(&config)

	for {
		line, err := ter.Readline()
		if err != nil {
			fmt.Println("read err: " + err.Error())
			break
		}
		command := strings.TrimSpace(line)
		if _, err := ter.WriteLine([]byte("type: " + command)); err != nil {
			fmt.Printf("Write fail %v", err)
		}
		ter.WritePrompt()
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
			line, err := ter.Readline()
			if err != nil {
				fmt.Println("read err: " + err.Error())
				break
			}
			command := strings.TrimSpace(line)
			if i, err := ter.WriteLine([]byte("type: " + command)); err == nil {
				fmt.Printf("Write %d bytes\r\n", i)
			} else {
				fmt.Printf("Write fail %v", err)
			}
			ter.WritePrompt()
		}
	}()
	time.Sleep(time.Second * 10)
	ter.Close()
	time.Sleep(time.Second * 5)
}

func TestOsTerminal() {
	terminal := term.NewTerminal(&StdReadWriter{os.Stdin, os.Stdout}, "$ ")
	for {
		line, err := terminal.ReadLine()
		if err != nil {
			break
		}
		command := strings.TrimSpace(line)
		if i, err := terminal.Write([]byte("type: " + command)); err == nil {
			fmt.Printf("Write %d bytes\r\n", i)
		} else {
			fmt.Printf("Write fail %v", err)
		}
	}
}

type StdReadWriter struct {
	in  io.Reader
	out io.Writer
}

func (std *StdReadWriter) Read(p []byte) (n int, err error) {
	n, err = std.in.Read(p)
	return n, err
}

func (std *StdReadWriter) Write(p []byte) (n int, err error) {
	n, err = std.out.Write(p)
	return n, err
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
