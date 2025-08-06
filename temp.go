package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 创建计数通道（带缓冲区的空结构体通道）
	taskCount := 5
	var wg sync.WaitGroup

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait()
	fmt.Println("所有任务完成！")
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d 开始工作\n", id)
	time.Sleep(time.Duration(id+1) * time.Second) // 模拟不同时长的工作
	fmt.Printf("Worker %d 工作结束\n", id)
}
