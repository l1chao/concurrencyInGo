package main_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"testing"
	"time"
)

func GenerateIntA(done chan struct{}) chan int {
	defer fmt.Println("G func exited.")
	ch := make(chan int) // 无缓存
	go func() {
		defer fmt.Println("goroutine of GenerateIntA exited.")

		for {
			select {
			case ch <- func() int { // 返回值二合一：用一个函数包裹即可
				n, _ := rand.Int(rand.Reader, big.NewInt(100)) // 添加上限值
				return int(n.Int64())                          // 转换为int
			}():
			case <-done:
				close(ch)
				return
			}
		}

	}()
	return ch
}

func Test1(t *testing.T) {
	done := make(chan struct{})
	ch := GenerateIntA(done)

	fmt.Println(<-ch)
	fmt.Println(<-ch)

	close(done)
	time.Sleep(1 * time.Second)

	fmt.Println(<-ch)
	fmt.Println(<-ch)

	time.Sleep(2 * time.Second)
	println("NumGoroutine=", runtime.NumGoroutine())
}

// 模拟后台工作协程
func backgroundWorker(done chan struct{}) chan string {
	results := make(chan string)

	go func() {
		defer close(results) // 确保关闭结果通道
		defer fmt.Println("Worker: Cleaning up resources...")

		counter := 0
		for {
			select {
			case <-time.After(500 * time.Millisecond): // 每500ms执行一次任务
				counter++
				result := fmt.Sprintf("Result #%d", counter)
				results <- result
				fmt.Println("Worker: Produced", result)

			case <-done: // 收到退出信号
				fmt.Println("Worker: Received shutdown signal")
				return // 退出协程
			}
		}
	}()

	return results
}

func Test2(t *testing.T) {
	done := make(chan struct{}) // 退出通知通道
	results := backgroundWorker(done)

	// 主程序从结果通道读取3次
	for i := 0; i < 3; i++ {
		result := <-results
		fmt.Println("Main: Received", result)
	}

	// 发送退出通知
	fmt.Println("Main: Sending shutdown signal...")
	close(done)

	// 尝试读取剩余结果
	fmt.Println("Main: Checking for remaining results...")
	for result := range results {
		fmt.Println("Main: Received leftover", result)
	}

	// 等待以确保协程已退出
	time.Sleep(300 * time.Millisecond)
	fmt.Println("Main: Program completed")
}
