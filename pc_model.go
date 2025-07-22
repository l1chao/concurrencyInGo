package main

// 生产者消费者模型啊

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 订单结构体
type Order struct {
	ID     int
	Amount float64
	Items  []string
}

// 生成随机订单数据
func generateOrder(id int) Order {
	items := []string{"T-Shirt", "Laptop", "Book", "Headphones", "Camera"}
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	itemCount := rand.Intn(3) + 1 // 1-3个商品
	return Order{
		ID:     id,
		Amount: float64(rand.Intn(500)+50) + rand.Float64(), // 50.00-549.99
		Items:  items[:itemCount],
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 使用带缓冲的channel作为订单队列 (容量100)
	orderQueue := make(chan Order, 100)

	// WaitGroup用于等待所有消费者完成
	var wg sync.WaitGroup

	// 创建10个消费者协程
	const numConsumers = 10
	wg.Add(numConsumers)

	fmt.Println("🚀 启动订单处理系统...")
	fmt.Printf("🛒 生产者开始生成订单 | 👥 创建%d个消费者\n", numConsumers)

	// 生产者：生成并发送订单
	go func() {
		const totalOrders = 50 // 总共生成50个订单

		for orderID := 1; orderID <= totalOrders; orderID++ {
			order := generateOrder(orderID)

			// 模拟随机订单到达间隔
			interval := time.Duration(rand.Intn(150)) * time.Millisecond
			time.Sleep(interval)

			fmt.Printf("📦 生产者: 创建订单 #%d (%.2f) - %v | 队列状态: %d/%d\n",
				order.ID, order.Amount, order.Items, len(orderQueue), cap(orderQueue))

			// 将订单发送到队列(阻塞操作直到有可用空间)
			orderQueue <- order
		}

		fmt.Printf("\n🛑 生产者已创建所有%d个订单，关闭订单队列\n", totalOrders)
		close(orderQueue) // 关闭通道以通知消费者
	}()

	time.Sleep(1500 * time.Millisecond)

	// 启动消费者协程。先创建消费者协程，再创建生产者协程吗？哈吉并，你这家伙。
	for i := 1; i <= numConsumers; i++ {
		go func(consumerID int) {
			defer wg.Done() // 一个consumer的range关闭了，就Done。所有consumer的range关闭了，就表示orderQueue真的没有了！

			for order := range orderQ ueue { // orderQueue通道关闭后，for range​​不会立即终止​​：for range会​​继续读取通道中剩余的所有数据​​，直到通道被完全清空。
				fmt.Printf("👷 消费者%d 开始处理订单 #%d (金额: $%.2f)\n",
					consumerID, order.ID, order.Amount)

				// 模拟订单处理时间
				processTime := time.Duration(rand.Intn(800)+200) * time.Millisecond
				time.Sleep(processTime)

				// 模拟支付处理
				if rand.Float32() < 0.92 { // 92%支付成功率
					fmt.Printf("✅ 消费者%d 成功处理订单 #%d | 耗时: %v\n",
						consumerID, order.ID, processTime.Round(time.Millisecond))
				} else {
					fmt.Printf("❌ 消费者%d 支付失败 #%d | 耗时: %v\n",
						consumerID, order.ID, processTime.Round(time.Millisecond))
				}
			}

			fmt.Printf("🛑 消费者%d 停止工作\n", consumerID)
		}(i)
	}

	// 主协程监控队列状态
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			fmt.Printf("📊 监控: 当前队列长度 %d/%d | 活跃消费者: %d\n",
				len(orderQueue), cap(orderQueue), numConsumers)
		}
	}()

	// 等待所有消费者完成工作
	wg.Wait()
	fmt.Println("\n🔚 所有订单处理完成，系统关闭")
}
