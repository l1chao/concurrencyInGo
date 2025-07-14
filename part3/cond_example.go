package main_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	var (
		mu     sync.Mutex          // 共享数据保护锁
		cond   = sync.NewCond(&mu) // 创建条件变量
		buffer []int               // 共享数据
		done   bool                // 结束标志
	)

	// 消费者（持续消费数据）
	consumer := func(id int) {
		cond.L.Lock() // 1. 获取共享锁
		defer cond.L.Unlock()

		// 循环检查条件：缓冲区有数据或所有生产结束时才继续
		for !(len(buffer) > 0 || done) {
			fmt.Printf("Consumer %d: 等待数据\n", id)
			cond.Wait() // 2. 条件不满足，释放锁并阻塞
		}

		// 消费数据
		if len(buffer) > 0 {
			item := buffer[0]
			buffer = buffer[1:]
			fmt.Printf("Consumer %d: 消费 %d (剩余: %d)\n", id, item, len(buffer))
		}
	}

	// 生产者（持续生产数据）
	producer := func(id int) {
		for i := 0; i < 3; i++ {
			time.Sleep(100 * time.Millisecond) // 模拟生产耗时

			cond.L.Lock() // 3. 获取共享锁
			item := id*100 + i
			buffer = append(buffer, item)
			fmt.Printf("Producer %d: 生产 %d (总量: %d)\n", id, item, len(buffer))

			// 注意，调用 Wait() 前必须持有锁；Wait() 返回后会重新持有锁
			cond.Signal() // 4. 通知一个等待的消费者

			cond.L.Unlock()
		}
	}

	// 启动生产者
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1) // 必须在外面
		go func(id int) {
			defer wg.Done()
			producer(id)
		}(i)
	}

	// 启动消费者
	for i := 0; i < 2; i++ {
		go consumer(i)
	}

	// 等待生产者完成
	wg.Wait()

	// 设置结束标志并广播
	cond.L.Lock()
	done = true
	cond.Broadcast() // 5. 唤醒所有消费者
	cond.L.Unlock()

	time.Sleep(500 * time.Millisecond) // 等待日志输出
}
