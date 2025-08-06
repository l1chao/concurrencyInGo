package main

// 串行退出的方法。应当注意，goroutine推出的时候基本都是无顺序推出。需要串行退出的情况很少。

// func main() {
// 	// 使用一个通道来传递退出信号，每个goroutine在退出前通知下一个
// 	ch1 := make(chan struct{})
// 	ch2 := make(chan struct{})
// 	ch3 := make(chan struct{})

// 	var wg sync.WaitGroup
// 	wg.Add(3)

// 	// 第一个goroutine
// 	go func() {
// 		defer wg.Done()
// 		defer close(ch1) // 退出时关闭ch1，通知下一个（这里是第二个）可以退出了
// 		fmt.Println("第一个goroutine开始运行")
// 		// 模拟工作
// 		time.Sleep(2 * time.Second)
// 		fmt.Println("第一个goroutine退出")
// 	}()

// 	// 第二个goroutine
// 	go func() {
// 		defer wg.Done()
// 		defer close(ch2) // 退出时关闭ch2，通知下一个（这里是第三个）可以退出了
// 		// 等待第一个goroutine退出
// 		<-ch1
// 		fmt.Println("第二个goroutine开始运行")
// 		time.Sleep(1 * time.Second)
// 		fmt.Println("第二个goroutine退出")
// 	}()

// 	// 第三个goroutine
// 	go func() {
// 		defer wg.Done()
// 		// 等待第二个goroutine退出
// 		<-ch2
// 		fmt.Println("第三个goroutine开始运行")
// 		time.Sleep(500 * time.Millisecond)
// 		fmt.Println("第三个goroutine退出")
// 		close(ch3) // 这里也可以不关闭，因为已经没有goroutine在等待了
// 	}()

// 	// 如果需要，可以等待第三个goroutine退出
// 	// <-ch3

// 	// 等待所有goroutine退出
// 	wg.Wait()
// 	fmt.Println("所有goroutine已退出")
// }
