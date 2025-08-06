package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func fetchAPI(ctx context.Context, url string) (string, error) {
	// 1. 创建绑定到上下文的HTTP请求
	// http.NewRequestWithContext将传入的ctx与请求绑定
	// 当ctx被取消或超时时，HTTP客户端会主动中断请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	log.Println("发送请求...")

	// 2. 执行HTTP请求
	// http.DefaultClient.Do()执行请求时持续监听ctx.Done()通道
	// 一旦ctx被取消，请求会被立即中断
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// 此处的err可能是正常的网络错误，也可能是ctx中断造成的
		return "", fmt.Errorf("请求执行失败: %w", err)
	}

	// 确保响应体被关闭（防止资源泄漏）
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	log.Println("读取响应...")

	// 3. 读取响应体（此操作也可能被ctx取消）
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	return string(data), nil
}

func Test1(t *testing.T) {
	// 4. 创建500ms超时的上下文
	// context.WithTimeout创建：
	//   ctx - 带有超时的上下文
	//   cancel - 取消函数（可用于手动取消）
	//
	// 在后台会启动一个计时器：
	//   500ms后自动关闭ctx.Done()通道
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// 5. 确保在任何情况下调用cancel释放资源
	// defer确保无论函数如何退出（正常返回或panic）都会调用cancel()
	defer cancel()

	// 6. 执行API调用，传入超时上下文
	result, err := fetchAPI(ctx, "https://httpbin.org/delay/1") // 模拟延迟1秒的API

	// 7. 错误处理
	if err != nil {
		// 检查是否是超时错误
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("请求超时: 服务器未在指定时间内响应")
			return
		}

		// 其他类型的错误
		log.Fatal("请求失败:", err)
	}

	fmt.Println("请求成功:", result)
}

func Test2(t *testing.T) {
	// ================= 示例1：从字符串读取 =================
	// 创建字符串读取器
	stringReader := strings.NewReader("Hello, 世界! 🌍")

	// 一次性读取全部内容
	data, err := io.ReadAll(stringReader)
	if err != nil {
		log.Fatalf("字符串读取失败: %v", err)
	}

	// 打印结果
	fmt.Printf("示例1 - 字符串读取: \n\t内容: %s\n\t长度: %d 字节\n\t原始字节: %v\n\n",
		string(data), len(data), data)

	// ================= 示例2：从字节缓冲区读取 =================
	// 创建包含二进制数据的缓冲区
	buf := bytes.NewBuffer([]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}) // ASCII 的 "Hello"
	buf.WriteByte(0x20)                                          // 空格
	buf.WriteString("World")                                     // 字符串

	// 读取缓冲区全部内容
	bufData, err := io.ReadAll(buf)
	if err != nil {
		log.Fatalf("缓冲区读取失败: %v", err)
	}

	fmt.Printf("示例2 - 缓冲区读取: \n\t内容: %s\n\t十六进制: % x\n\n",
		string(bufData), bufData)

	// ================= 示例3：从文件读取 =================
	// 创建临时文件（实践中应使用真实文件路径）
	tmpFile, err := os.CreateTemp("", "readall-example-*.txt")
	if err != nil {
		log.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // 程序结束时删除临时文件

	// 向文件写入多行文本
	content := "第一行\n第二行\n第三行"
	if _, err := tmpFile.WriteString(content); err != nil {
		log.Fatalf("写入文件失败: %v", err)
	}

	// 重置文件指针到开头（重要！）
	if _, err := tmpFile.Seek(0, 0); err != nil {
		log.Fatalf("重置文件指针失败: %v", err)
	}

	// 读取整个文件内容
	fileData, err := io.ReadAll(tmpFile)
	if err != nil {
		log.Fatalf("文件读取失败: %v", err)
	}

	fmt.Printf("示例3 - 文件读取: \n\t内容: \n%s\n\t行数: %d\n\n",
		string(fileData), bytes.Count(fileData, []byte{'\n'}))

	// ================= 示例4：有限读取器 =================
	// 创建一个限制长度的读取器（最多读10字节）
	limitedReader := io.LimitReader(strings.NewReader("这段内容将被截断"), 10)

	limitedData, err := io.ReadAll(limitedReader)
	if err != nil {
		log.Fatalf("有限读取失败: %v", err)
	}

	fmt.Printf("示例4 - 有限读取(10字节): \n\t结果: %s\n\n", limitedData)

	// ================= 示例5：错误处理 =================
	// 创建自定义错误读取器
	errorReader := &ErrorReader{Msg: "模拟读取错误"}

	_, err = io.ReadAll(errorReader)
	if err != nil {
		fmt.Printf("示例5 - 错误处理: \n\t错误信息: %v\n\t错误类型: %T\n",
			err, err)
	}
}

// ================= 自定义错误读取器 =================
// 实现 io.Reader 接口但始终返回错误
type ErrorReader struct {
	Msg string
}

func (er *ErrorReader) Read(p []byte) (n int, err error) {
	// 返回自定义错误
	return 0, fmt.Errorf("自定义读取错误: %s", er.Msg)
}
