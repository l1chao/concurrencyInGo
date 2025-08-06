package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// 1.等待1个goroutine结束
func worker(args ...interface{}) {
	if len(args) == 0 {
		return
	}

	interval, ok := args[0].(int)
	if !ok {
		return
	}

	time.Sleep(time.Duration(interval) * time.Second)
}

func Spawn1(f func(args ...interface{}), args ...interface{}) chan struct{} {
	c := make(chan struct{})
	go func() {
		f(args)
		c <- struct{}{}
	}()

	return c
}

func main() {
	done := Spawn1(worker, 5)
	<-done
	fmt.Println("worker done.")
}

// 2. 获取goroutine退出状态

func Worker1(args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("invalid args.")
	}

	interval, ok := args[0].(int)
	if !ok {
		return errors.New("invalid interval args.")
	}

	time.Sleep(time.Second * time.Duration(interval))
	return errors.New("OK")
}

func Spawn2(f func(args ...interface{}) error, args ...interface{}) chan error {
	c := make(chan error)
	go func() {
		c <- f(args...)
	}()
	return c
}

func main() {
	done := Spawn2(Worker1, 5)
	println("Spawn worker1")
	err := <-done
	println("worker1 done:", err)

	done = Spawn2(Worker1)
	println("spawn worker2")
	err = <-done
	println("worker2 done:", err)
}

// 3. 等待多个goroutine退出

// func worker()

func spawnGroup(n int, f func(args ...interface{}), args ...interface{}) chan struct{} {
	c := make(chan struct{})
	var wg sync.WaitGroup

	for i := range n {
		wg.Add(1)
		go func() {
			f(args, i+1)
			println(fmt.Sprintf("worker-%d", i), "done")

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		c <- struct{}{}
	}()

	return c
}

func main() {
	done := spawnGroup(5, worker, 3)
	println("spawn a group of workers.")
	<-done
	println("group workers done")
}

// 4.超时取消等待
func main() {
	done := spawnGroup(5, worker, 30)
	println("spawn a group of workers.")

	timer := time.NewTimer(time.Second * 5)
	defer timer.Stop()

	select {
	case <-done:
		println("workers done.")
	case <-timer.C:
		println("workers timeout exit!")
	}
}
