package main

// 下面协程模型的关键是：这是一个排号器，多个协程并发拿到自己的票号。
// 对于像下面这种，生产者不应该分散到各个消费者进程里面。v := Increase_sync()的调用就是在消费者协程里面完成的。用chan来替代这种情况下的锁是更好的，不仅更加符合go的设计哲学，而且生产者是一个单独的协程，消费者想要使用生产出来的编号只需要通过chan接收。

// type counter_sync struct {
// 	sync.Mutex // 内嵌结构体，方法直接提升！
// 	i int
// }

// var cter_sync counter_sync

// func Increase_sync() int {
// 	cter_sync.Lock()
// 	defer cter_sync.Unlock()

// 	cter_sync.i++
// 	return cter_sync.i
// }

// func main() {
// 	for i := range 10 {
// 		go func(i int) {
// 			v := Increase_sync()
// 			fmt.Printf("goroutine-%d: current counter value is %d\n", i, v)
// 		}(i)
// 	}
// 	time.Sleep(5 * time.Second)
// }

// 利用channel替代上面的锁机制
// type counter_chan struct {
// 	c chan int
// 	i int
// }

// var cter_chan counter_chan

// func InitCounter() {
// 	cter_chan = counter_chan{
// 		c: make(chan int),
// 	}
// 	go func() {
// 		for {
// 			cter_chan.i++
// 			cter_chan.c <- cter_chan.i
// 		}
// 	}()
// 	fmt.Println("counter init ok.")
// }

// func Increase() int {
// 	return <-cter_chan.c
// }

// func init() {
// 	InitCounter()
// }

// func main() {
// 	for i := range 10 {
// 		go func(i int) {
// 			v := Increase()
// 			fmt.Printf("goroutine-%d,current counter value is %d\n", i, v)
// 		}(i)
// 	}
// 	time.Sleep(3 * time.Second)
// }
