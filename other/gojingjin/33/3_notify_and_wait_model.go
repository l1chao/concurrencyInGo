package main

import (
	"fmt"
	"sync"
	"time"
)

func Worker3(j int) {
	time.Sleep(time.Duration(j) * time.Second)
}

func Spawn3(f func(int)) chan string {
	quit := make(chan string)

	go func() {
		var jobs chan int
		for {
			select {
			case j := <-jobs:
				f(j)
			case <-quit:
				quit <- "OK"
				return
			}
		}
	}()

	return quit
}

func main() {
	quit := Spawn3(Worker3)
	println("spawn a worker goroutine.")

	time.Sleep(5 * time.Second)

	println("notify the worker to exit...")
	quit <- "exit"
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	select {
	case state := <-quit:
		println("worker done with ", state)
	case <-timer.C:
		println("wait worker exit timeout")
	}
}

// 通知并等待多个goroutine

// func worker()

func spawnGroup(n int, f func(int)) chan struct{} {
	quit := make(chan struct{})
	job := make(chan int)
	var wg sync.WaitGroup

	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("worker-%d", i)
			for {
				j, ok := <-job // 持续监听时遇到关闭，则收到零值+false；对一个已经关闭了的读取，则收到零值+false。
				if !ok {
					println(name, "done")
					return
				}
				worker(j)
			}
		}(i)
	}

	go func() {
		<-quit
		close(job)
		wg.Wait()
		quit <- struct{}{}
	}()

	return quit
}

func main() {
	quit := spawnGroup(5, worker)
	println("spawn a group of workers")

	time.Sleep(5 * time.Second)

	println("notify the workers to exit...")
	quit <- struct{}{}

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	select {
	case <-timer.C:
		println("wait timeout.")
	case state := <-quit:
		println("worker done with", state)
	}
}
