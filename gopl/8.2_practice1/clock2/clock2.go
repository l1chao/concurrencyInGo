// clock2.go - 模拟不同时区的时钟服务器
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// 如果没有传入的port或zone参数，那么这两个就是默认参数。
var port = flag.String("port", "8000", "TCP服务端口号")
var zone = flag.String("zone", "UTC", "时区名称(如Asia/Shanghai)") // 这里是默认UTC

func main() {
	flag.Parse()

	// 设置时区
	if *zone != "UTC" {
		os.Setenv("TZ", *zone) // 如果不是UTC，那么设置环境变量TZ为*zone
	}

	listener, err := net.Listen("tcp", "localhost:"+*port) // 在listen的就是服务器。
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("时钟服务运行中 [时区: %s, 端口: %s]\n", *zone, *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn) // 处理过程应该是并发的，接收完全可以同步阻塞。
	}
}

func handleConn(c net.Conn) {
	defer c.Close() // 对于每一个链接要记得手动关闭。
	for {
		// 获取当前时间并格式化
		current := time.Now().Format("15:04:05\n")
		_, err := io.WriteString(c, current) // 可以直接往tcp里面接里面写东西
		if err != nil {
			return // 客户端断开连接
		}
		time.Sleep(1 * time.Second)
	}
}
