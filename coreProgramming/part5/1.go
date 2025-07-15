package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
)

func GenerateIntA(done chan struct{}) chan int {
	ch := make(chan int) // 无缓存
	go func() {
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

func main() {
	done := make(chan struct{})
	ch := GenerateIntA(done)

	fmt.Println(<-ch)
	fmt.Println(<-ch)

	close(done)

	fmt.Println(<-ch)
	fmt.Println(<-ch)

	println("NumGoroutine=", runtime.NumGoroutine())
}
