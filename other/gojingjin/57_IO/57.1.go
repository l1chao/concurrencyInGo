package io

import (
	"fmt"
	"os"
)

// 理解最基本的IO模式

// writer写出内容。
func directWriteByteSliceToFile(path string, data []byte) (int, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file err: ", err)
		return 0, err
	}

	defer func() {
		f.Sync()
		f.Close()
	}()
	return f.Write(data)
}

// reader接收内容。
func directReadByteSliceFromFile(path string, data []byte) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("open file err :", err)
		return 0, err
	}

	defer f.Close()

	return f.Read(data)
}

func main1() {
	filePath := "./foo.txt"
	text := "hello, gopher!"
	data := make([]byte, 20)

	n, err := directWriteByteSliceToFile(filePath, []byte(text))
	if err != nil {
		fmt.Println("write file err :", err)
		return
	}
	fmt.Printf("write %d byte to file.\n", n)

	n, err = directReadByteSliceFromFile(filePath, data)
	if err != nil {
		fmt.Println("read file err:", err)
		return
	}

	fmt.Printf("read %d byte from file, content:%q \n", n, data)
}
