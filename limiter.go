package main

import (
	"fmt"
	"sync"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

const (
	Duration     = 3
	Period       = 300 * time.Millisecond
	DefaultLimit = 0.8
)

type Limiter struct {
	st       time.Time
	proc     *process.Process
	cpuUsage []float64
	lim      float64
	wg       *sync.WaitGroup

	done <-chan bool
}

func NewLimiter(pid int32, lim float64, wg *sync.WaitGroup) *Limiter {
	proc, err := process.NewProcess(pid)
	if err != nil {
		panic(err)
	}

	return &Limiter{proc: proc, cpuUsage: make([]float64, Duration), lim: lim, wg: wg}
}

func (l *Limiter) limit() {
	fmt.Println("Started Limiter")
	tick := time.NewTicker(Period)

	busy2, all2 := l.getBusy()
	counter := 0
	go func() {
		defer l.wg.Done()
		defer tick.Stop()

		for {
			select {
			case <-l.done:
				return

			case <-tick.C:
				fmt.Println("going to calculate cpu usage")
				busy1, all1 := busy2, all2
				busy2, all2 = l.getBusy()
				cpuUsage := getCPUUsage(busy1, all1, busy2, all2)
				l.cpuUsage[counter] = cpuUsage
				avgUsage := average(l.cpuUsage)
				fmt.Println("Avg Usage::", avgUsage)

				if avgUsage > l.lim {
					// send a stop signal
					fmt.Println("Stopping::", avgUsage)
					l.sendSignal(syscall.SIGSTOP)
				} else {
					// send a continue signal
					fmt.Println("Resuming::", avgUsage)
					l.sendSignal(syscall.SIGCONT)
				}
				counter += 1
				if counter > Duration-1 {
					counter = 0
				}
			default:
				// fmt.Println("I am doing nothing")
				// do nothing
			}
		}
	}()
}

func (l *Limiter) sendSignal(sig syscall.Signal) {
	err := l.proc.SendSignal(sig)
	if err != nil {
		panic(err)
	}
}

func average(m []float64) (avg float64) {
	for _, elem := range m {
		avg += elem
	}
	avg /= float64(len(m))
	return
}

func (l *Limiter) getBusy() (busy, all float64) {
	busy, _ = busyFromTimes(getProcessCpuTimes(l.proc))
	_, all = busyFromTimes(getGlobalCpuTimes())
	return
}

func getGlobalCpuTimes() cpu.TimesStat {
	ts, err := cpu.Times(false)
	if err != nil {
		panic(err)
	}
	return ts[0]
}

func getProcessCpuTimes(proc *process.Process) cpu.TimesStat {
	t, err := proc.Times()
	if err != nil {
		panic(err)
	}
	return *t
}

func busyFromTimes(t cpu.TimesStat) (busy, all float64) {
	busy = t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Guest + t.GuestNice
	all = busy + t.Idle
	return
}

func getCPUUsage(busy1, all1, busy2, all2 float64) float64 {
	if all1 == all2 {
		return 0.0
	}
	usage := ((busy2 - busy1) / (all2 - all1)) * 100.0
	return usage
}
