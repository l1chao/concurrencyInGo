package main

// 实际上就是，spawn函数将某个chan处理之后又传递给另一个spawn继续处理。

func newNumGenerator(start, count int) <-chan int {
	c := make(chan int)

	go func() {
		for i := start; i < start+count; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}

func filterOdd(in int) (int, bool) {
	if in%2 != 0 {
		return 0, false
	}
	return in, true
}

func square(in int) (int, bool) {
	return in * in, true
}

// spawn: To ​​generate​​, ​​cause the creation of
func spawn1(f func(int) (int, bool), in <-chan int) <-chan int { //直观：一个用来处理的函数 + 一个待处理的通道。
	out := make(chan int)

	go func() {
		for v := range in {
			r, ok := f(v)
			if ok {
				out <- r
			}
		}
		close(out)
	}()
	return out
}

func main() {
	in := newNumGenerator(1, 20)
	out := spawn1(square, spawn1(filterOdd, in)) // 多层管道嵌套
	for v := range out {
		println(v)
	}
}
