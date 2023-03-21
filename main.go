package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/shirou/gopsutil/process"
)

func main() {

	pid := flag.Int("p", -1, "target process")
	lim := flag.Float64("l", 0, "limit between 0 - 1")

	flag.Parse()

	if *pid == -1 {
		fmt.Println("invalid input, target process is either not provided or invalid pid")
		os.Exit(1)
	}

	if *lim == 0 {
		*lim = DefaultLimit
	}

	fmt.Println("I am here")

	p := process.Process{Pid: int32(*pid)}
	if run, err := p.IsRunning(); err != nil {
		fmt.Printf("Error, %v\n", err)
	} else {
		fmt.Printf("If process running, %v\n", run)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	l := NewLimiter(int32(*pid), *lim*100, &wg)
	l.limit()
	wg.Wait()
}
