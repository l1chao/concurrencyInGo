package main_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	// var ch chan int // 一种管道只能够装一种类型的数据。未手动初始化的channel零值为nil。

	ch1 := make(chan int)    // 不带缓冲区的管道
	ch2 := make(chan string) // 带缓冲区的管道

	// 管道操作符 - 非参数时
	ch1 <- 10 // 写入数据
	go func() {
		data := <-ch1 // 读出数据
		fmt.Println(data)
	}()

	func(chan<- int) { // 表示只能写。读会报错。注意<-是紧跟chan的
		//...
	}(ch1) // 传参的时候不用加管道操作符
	func(ch <-chan int) { // 表示只能读
		//...
	}(ch1)

	// 关闭性与可读/可写。
	// 对关闭了的、无缓存的、无剩余元素的channel可读。只是读空值。
	close(ch2)
	value, ok := <-ch2 // value为零值，ok为false
	fmt.Println(value, ok)

	// 对关闭了的、有缓存的、有剩余元素的channel可读。感觉就像正常读一样。
	var ch3 = make(chan string, 5)
	ch3 <- "good"
	ch3 <- "world"
	close(ch3)
	value, ok = <-ch3
	fmt.Println(value, ok)

	value, ok = <-ch3
	fmt.Println(value, ok)

	// 不能对关闭了的channel写！
	// ch3<-"addtional" // 报错

	// 对nil channel进行读写将会永久阻塞
	var ch4 chan int
	data := <-ch4
	fmt.Println(data)
}

func Test2(t *testing.T) {
	// channel是一个队列，读取是消耗读。
	// 用cap查看缓存区大小，用len查看当前channel元素。
	var ch3 = make(chan string, 5)
	ch3 <- "good"
	ch3 <- "world"
	fmt.Println(len(ch3), cap(ch3))

	value, ok := <-ch3
	fmt.Println(value, ok)
	fmt.Println(len(ch3), cap(ch3))

	value, ok = <-ch3
	fmt.Println(value, ok)
	fmt.Println(len(ch3), cap(ch3))
}

type Person struct{} // 结构体类型

func Test3(t *testing.T) {
	// var s struct{}      // 结构体类型
	// var s1 = struct{}{} // 结构体实例

	// var f func(int) (int, error)       // 函数类型
	// var f = func(a int) (int, error) { // 函数实例
	// 	return 1, nil
	// }
}

// 如果要等好几个事件结束了之后，当前主协程才能够结束，那就需要用waitgroup来实现
func TestFoo1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(5)

	for i := range 5 {
		go foo1_helper(&wg, i+1)
	}
	wg.Wait()
	fmt.Println("main goroutine is over.")
}

func foo1_helper(wg *sync.WaitGroup, id int) {
	defer wg.Done() // 如果要主协程等foo1_helper协程完毕了之后才能结束，那么helper的Done就应该是函数要结束的时候才调用。
	fmt.Println("Welcome! No.", id)
}

func Test4(t *testing.T) {

}

// kubernetes控制器
func waitForStopOrTimeout(stopCh <-chan struct{}, timeout time.Duration) <-chan struct{} {
	stopChWithTimeout := make(chan struct{})
	go func() {
		select {
		case <-stopCh: // 自然结束
		case <-time.After(timeout): //最长等待时间
		}
		close(stopChWithTimeout)
	}()
	return stopChWithTimeout
}
