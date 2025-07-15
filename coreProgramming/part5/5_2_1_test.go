package main_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func GenerateIntB() chan int {
	ch := make(chan int, 10)

	go func() {
		for {
			ch <- func() int {
				num, _ := rand.Int(rand.Reader, big.NewInt(100))
				return int(num.Int64())
			}()
		}
	}()
	return ch
}

func Test5_2_1_1(t *testing.T) {
	ch := GenerateIntB()
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

// 多个goroutine增强型生成器

func GenerateC() chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			ch <- func() int {
				num, _ := rand.Int(rand.Reader, big.NewInt(100))
				return int(num.Int64())
			}()
		}
	}()
	return ch
}

func GenerateInt() chan int {
	ch := make(chan int, 20)

	go func() {
		for {
			select {
			case ch <- <-GenerateIntB():
			case ch <- <-GenerateC():
			}
		}
	}()
	return ch
}
func Test5_2_1_2(t *testing.T) {
	ch := GenerateInt()
	for range 100 {
		fmt.Println(<-ch)
	}
}
