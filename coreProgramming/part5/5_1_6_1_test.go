package main_test

import (
	"sync"
	"time"
)

// 什么是fan-out和fan-in

// 扇出
// Fan-out: 一个输入源 -> 多个工作者协程
func fanOut(in <-chan int, workers int) []<-chan int {
	outs := make([]<-chan int, workers) // 将一个输入：in channel变为workers个channel，然后这些channel能够被多个协程并发处理。

	for i := range workers { // 通过workers个协程里面的channel来处理in里面的内容。
		out := make(chan int)
		go func(id int) {
			defer close(out) // 保证了所有通道在结束使用后被关闭。
			for value := range in {
				time.Sleep(time.Millisecond * 100) // 模拟处理任务
				result := value * 2                // 可以进行初级加工
				out <- result
			}
		}(i)
		outs[i] = out
	}
	return outs
}

func fanOut1(in <-chan int, workers int) []<-chan int {
	outs := make([]<-chan int, workers)
	for i := range workers {
		out := make(chan int)
		go func() {
			defer close(out)
			for value := range in {
				time.Sleep(100 * time.Millisecond)
				out <- value * 2
			}
		}()
		outs[i] = out
	}
	return outs
}

// 扇入：扇入模式是将多个通道的数据汇总到一个通道中。
// Fan-in: 多个通道 -> 一个输出通道
func fanIn(ins ...<-chan int) <-chan int {
	out := make(chan int) // 总输出，无缓冲
	var wg sync.WaitGroup

	collector := func(in <-chan int) {
		defer wg.Done()
		for value := range in { // 当前in关闭了之后，就会执行defer。
			out <- value
		}
	}

	wg.Add(len(ins))
	for _, in := range ins {
		go collector(in)
	}

	// 关闭out
	go func() {
		wg.Wait()  // 等ins里面的所有in channel都被关闭的时候，就能通过这里的阻塞。
		close(out) // 为了让扇入的out在完成使命（接收所有的ins里面的数据）之后会正常关闭，这里必须设置同步（接收完所有的ins在前，关闭out在后）
	}()

	return out
}

func fanIn1(ins ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	collector := func(in <-chan int) {

	}
}
