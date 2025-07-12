package main_test

import "fmt"

// 下面是一种错误的写法。往ch1里面写数据将会卡住。data读数据必须在另一个协程里面完成，而不能在同一协程完成。
// 如果直接运行，将会报错：fatal error: all goroutines are asleep - deadlock!
func foo() {
	ch1 := make(chan int) // 不带缓冲区的管道
	ch1 <- 10
	data := <-ch1
	fmt.Println(data)
}

func foo_solution() {
	ch1 := make(chan int) // 不带缓冲区的管道
	ch1 <- 10

	go func() {
		data := <-ch1
		fmt.Println(data)
	}()
}
