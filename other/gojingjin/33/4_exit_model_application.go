package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type GracefullyShutdowner interface {
	Shutdown(waitTimeout time.Duration) error
}

func ConcurrentShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	c := make(chan struct{})

	go func() {
		var wg sync.WaitGroup
		for _, g := range shutdowners {
			wg.Add(1)
			go func(shutdowner GracefullyShutdowner) {
				defer wg.Done()
				shutdowner.Shutdown(waitTimeout)
			}(g)
		}

		wg.Wait()
		c <- struct{}{} // 实际上这个c和上面的wg的作用似乎重合了。但是没有，因为wg是用来确保所有子goroutine都完成了的，而c则是用来实现超时取消机制的！
	}()

	timer := time.NewTimer(waitTimeout)
	defer timer.Stop()
	select {
	case <-c:
		fmt.Println("运行结束！")
		return nil
	case <-timer.C:
		fmt.Println("超时了老登！")
		return errors.New("Wait timeout!")
	}
}

// 通过计数通道来完成一组退出

func main1() {
	// 创建计数通道（带缓冲区的空结构体通道）
	taskCount := 5
	doneChan := make(chan struct{}, taskCount)

	for i := 0; i < taskCount; i++ {
		go worker(i, doneChan)
	}

	//计数通道
	for i := 0; i < taskCount; i++ {
		<-doneChan // 每次接收相当于 wg.Done()
		fmt.Printf("已完成 %d 个任务\n", i+1)
	}

	fmt.Println("所有任务完成！")
}

func worker(id int, done chan<- struct{}) {
	defer func() {
		done <- struct{}{} // 任务结束时发送信号，相当于 wg.Done()
	}()

	fmt.Printf("Worker %d 开始工作\n", id)
	time.Sleep(time.Duration(id+1) * time.Second) // 模拟不同时长的工作
	fmt.Printf("Worker %d 工作结束\n", id)
}

// 如果使用&wg来传递锁实现一组退出，不是不行：
func main() {
	// 创建计数通道（带缓冲区的空结构体通道）
	taskCount := 5
	var wg sync.WaitGroup

	// for i := 0; i < taskCount; i++ {
	// 	wg.Add(1)         // 外加锁
	// 	go worker(i, &wg) // 内减锁。但是这一步可以优化，因为要充分利用闭包。
	// }
	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go func() { // 都是通过包裹来为goroutine结束完成Done工作。这是为了不破坏原来函数结构。
			worker(i)
			wg.Done()
		}()
	}
	wg.Wait() // 这样的坏处是：没有办法实现超时取消机制。
	// 为什么不能完成超时取消呢？因为程序只会卡在这一处，这一处只有“一个身位”。但是chan实现的卡在某一处，则会有“两个身位”，两个身位的意思其实就是不仅能够接收所有子协程的退出新号，也能够接收timer的退出信号。
	fmt.Println("所有任务完成！")
}

func worker(id int) {

	fmt.Printf("Worker %d 开始工作\n", id)
	time.Sleep(time.Duration(id+1) * time.Second) // 模拟不同时长的工作
	fmt.Printf("Worker %d 工作结束\n", id)
}

// 串行退出
// 总的来说，下面的代码实现了依次关闭组件的效果，且能够近似的保证退出时间不超过waitTimeout。
// type GracefullyShutdowner interface{}

func SequentialShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	start := time.Now()
	var left time.Duration
	timer := time.NewTimer(waitTimeout)

	for _, g := range shutdowners {
		elapsed := time.Since(start)
		left = waitTimeout - elapsed

		c := make(chan struct{})
		go func(shutdowner GracefullyShutdowner) { // 为什么串行退出还是要创建goroutine？因为要实现超时取消机制。超时取消机制一定是多路复用的。
			shutdowner.Shutdown(left)
			c <- struct{}{}
		}(g) // 因为早期版本g是共享的，这么写不会出错。

		timer.Reset(left) // 复用了timer，避免创建多个timer。
		select {
		case <-c:
			// 继续执行
		case <-timer.C:
			return errors.New("wait timeout") // 如果中间某一个组件超时，则本次的整个退出就失败
		}
	}
	return nil
}
