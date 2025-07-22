// clockwall.go - 从多个时钟服务器收集时间
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Clock struct {
	Location string
	Addr     string
	Time     string
}

var (
	clocks []*Clock
	mu     sync.Mutex // 保护clocks并发访问
)

func main() {
	// 解析命令行参数：NewYork=localhost:8010 Tokyo=localhost:8020
	flag.Parse()

	// 初始化时钟列表
	for _, arg := range flag.Args() {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "参数格式错误: %s\n", arg)
			os.Exit(1)
		}
		clocks = append(clocks, &Clock{
			Location: parts[0],
			Addr:     parts[1],
		})
	}

	// 启动各时钟服务器的连接
	for _, clock := range clocks {
		go connectToClock(clock)
	}

	// 主循环：每秒刷新显示
	for {
		renderClockWall()
		time.Sleep(1 * time.Second)
	}
}

func connectToClock(clock *Clock) {
	conn, err := net.Dial("tcp", clock.Addr)
	if err != nil {
		log.Printf("无法连接到 %s: %v", clock.Addr, err)
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		mu.Lock()
		clock.Time = scanner.Text()
		mu.Unlock()
	}
}

func renderClockWall() {
	mu.Lock()
	defer mu.Unlock()

	// 清屏并定位到左上角
	fmt.Print("\033[2J\033[0;0H")

	// 打印表头
	fmt.Println("+----------------+--------------+")
	fmt.Println("|     地点       |    时间      |")
	fmt.Println("+----------------+--------------+")

	// 打印各时钟时间
	for _, clock := range clocks {
		fmt.Printf("| %-14s | %-12s |\n", clock.Location, clock.Time)
	}

	// 打印表尾
	fmt.Println("+----------------+--------------+")
	fmt.Printf("更新时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
