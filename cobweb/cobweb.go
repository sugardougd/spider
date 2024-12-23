package main

import (
	"flag"
	"fmt"
	"github.com/desertbit/readline"
	"io"
	"net"
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
	fmt.Printf("spider_v1 %s:%d\r\n", host, port)
	if err := tcpClient(host, port); err != nil {
		fmt.Printf("An error occurred: %v", err)
	}
}

func readlineClient(host string, port int) error {
	return readline.DialRemote("tcp", host+":"+strconv.Itoa(port))
}

func tcpClient(host string, port int) error {
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	go func() {
		defer conn.Close()
		for {
			n, _ := io.Copy(conn, os.Stdin)
			if n == 0 {
				break
			}
		}
	}()
	for {
		n, _ := io.Copy(os.Stdout, conn)
		if n == 0 {
			break
		}
	}
	return nil
}
