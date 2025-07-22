package main

import (
	"CInG/gopl"
	"log"
	"net"
)

// 也就是说，导入路径是完全的模块名+目录名！这里的包名，就是该文件夹对应的包名，这个包名会在下面的代码里面使用。

// "CInG/gopl"
// func main() {
// 	go gopl.Spinner(100 * time.Millisecond) // 在代码中则用package名来调用导入包部件。
// 	const n = 45
// 	fibN := gopl.Fib(n) // slow
// 	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
// }

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go gopl.HandleConn(conn) // handle one connection at a time
	}
}
