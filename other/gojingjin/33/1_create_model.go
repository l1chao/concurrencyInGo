package main

type T struct {
}

func spawn(f func()) chan T {
	c := make(chan T)
	go func() {
		// 使用c比那辆进行该goroutine和主goroutine之间的通信。
		// 该goroutine内使用c是通过闭包的方式完成的。

		f()

	}()

	return c
}

func main() {
	c := spawn(func() {})
	// 下面就能够用c和创建了的goroutine进行通信。
}
