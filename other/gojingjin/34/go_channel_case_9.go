package main

// 用nil channel避免一直读空channel的零值。
import (
	"fmt"
	"time"
)

// 下面是要解决的问题代码：在关闭了c1之后，主协程将会一直读取被关闭的协程，以至于出现很多读空值的操作。
// 解决办法：1. 每一次读取的时候验证ok 2.利用nil channel将读取 closed channel的动作阻塞（对于nil channel进行读取的case将不再被选中）。
func main() {
	c1, c2 := make(chan int), make(chan int)

	go func() {
		time.Sleep(100 * time.Millisecond)
		c1 <- 1
		close(c1)
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		c2 <- 2
		close(c2)
	}()

	ok1, ok2 := false, false
	for {
		select {
		case k := <-c1:
			ok1 = true
			fmt.Println(k)
		case k := <-c2:
			ok2 = true
			fmt.Println(k)
		}

		if ok1 && ok2 {
			break
		}
	}

}
