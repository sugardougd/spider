package spider

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
)

func spiderGps(c *Context) error {
	c.Spider.Println(os.Args[0])
	c.Spider.Println("PID:", os.Getpid())
	c.Spider.Println("PPID:", os.Getppid())
	return nil
}

func spiderMemory(c *Context) error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	c.Spider.Println("PID:", os.Getpid())
	c.Spider.Println("Sys:", BytesTo(memStats.Sys))
	c.Spider.Println("HeapSys:", BytesTo(memStats.HeapSys))
	c.Spider.Println("HeapAlloc:", BytesTo(memStats.HeapAlloc), fmt.Sprintf("%.2f%%", Percentage(memStats.HeapAlloc, memStats.HeapSys)))
	return nil
}

func spiderStack(c *Context) error {
	stack := string(debug.Stack())
	c.Spider.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
	c.Spider.Println("NumGoroutine:", runtime.NumGoroutine())
	c.Spider.Println(stack)
	return nil
}
