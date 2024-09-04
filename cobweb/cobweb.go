package main

import (
	"flag"
	"fmt"
	"github.com/desertbit/readline"
	"os"
	"strconv"
)

const (
	defaultHost = ""
	defaultPort = 4322
)

// main cobweb <host> <port>
func main() {
	flag.Parse()
	args := flag.Args()
	host, port := defaultHost, defaultPort
	if len(args) > 0 {
		host = args[0]
	}
	if len(args) > 1 {
		if p, err := strconv.Atoi(args[1]); err == nil {
			port = p
		} else {
			fmt.Fprintf(os.Stderr, "illage %s, usage: cobweb <host> <port>\r\n", args[1])
			os.Exit(1)
		}
	}
	fmt.Printf("spider %s:%d\r\n", host, port)
	if err := readline.DialRemote("tcp", host+":"+strconv.Itoa(port)); err != nil {
		fmt.Errorf("An error occurred: %s \n", err.Error())
	}
}
