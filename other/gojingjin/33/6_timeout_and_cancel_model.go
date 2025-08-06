package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

type result struct {
	value string
}

// 第一版：
func first1(servers ...*httptest.Server) (result, error) {
	c := make(chan result, len(servers))

	queryFunc := func(i int, server *httptest.Server) {
		url := server.URL
		resp, err := http.Get(url) // 直接请求
		if err != nil {
			log.Printf("http get error:%s\n", err)

		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		c <- result{
			value: string(data),
		}
	}
	for i, serv := range servers {
		go queryFunc(i, serv)
	}

	return <-c, nil
}

func fakeWeatherServer1(name string, interval int) *httptest.Server {
	// 相当于开了一个http服务器。可以通过某些方法获取该服务器的url等以供访问。
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s recerve a http request.\n", name)
		time.Sleep(time.Duration(interval) * time.Millisecond)
		w.Write([]byte(name + ":ok"))
	}))
}

func main() {
	result, err := first1(fakeWeatherServer1("open-weather-1", 200),
		fakeWeatherServer1("open-weather-2", 1000),
		fakeWeatherServer1("open-weather-3", 600))
	if err != nil {
		log.Println("invoke first error:", err)
		return
	}

	log.Println(result)
}

// 第二版：
func first2(servers ...*httptest.Server) (result, error) {
	c := make(chan result, len(servers))

	queryFunc := func(i int, server *httptest.Server) {
		url := server.URL
		resp, err := http.Get(url) // 直接请求
		if err != nil {
			log.Printf("http get error:%s\n", err)
			return
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		c <- result{
			value: string(data),
		}
	}
	for i, serv := range servers {
		go queryFunc(i, serv)
	}

	select {
	case r := <-c:
		return r, nil
	case <-time.After(500 * time.Millisecond):
		return result{}, errors.New("wait timeout.")
	}
}

func fakeWeatherServer2(name string, interval int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s recerve a http request.\n", name)
		time.Sleep(time.Duration(interval) * time.Millisecond)
		w.Write([]byte(name + ":ok"))
	}))
}

func main() {
	result, err := first2(fakeWeatherServer2("open-weather-1", 200),
		fakeWeatherServer2("open-weather-2", 1000),
		fakeWeatherServer2("open-weather-3", 600))
	if err != nil {
		log.Println("invoke first error:", err)
		return
	}

	log.Println(result)
}

// 第三版：
func first(servers ...*httptest.Server) (result, error) {
	c := make(chan result)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queryFunc := func(i int, server *httptest.Server) {
		url := server.URL
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil) // 这里看得出，http包天然支持ctx机制。
		// 上面一行的等价写法如下：
		// req, err := http.NewRequest("GET", url, nil)
		// req = req.WithContext(ctx)
		if err != nil {
			log.Printf("query goroutine-%d: http NewRequest error: %s\n", i, err)
			return
		}

		log.Printf("query goroutine-%d: send request ...\n", i)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("query goroutine-%d: get return error: %s\n", i, err)
			return
		}
		log.Printf("query goroutine-%d: get response\n", i)
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)

		c <- result{
			value: string(data),
		}
		return
	}

	for i, serv := range servers {
		go queryFunc(i, serv)
	}

	select {
	case r := <-c:
		return r, nil
	case <-time.After(500 * time.Millisecond):
		return result{}, errors.New("timeout")
	}
}

func fakeWeatherServer(name string, interval int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s recerve a http request.\n", name)
		time.Sleep(time.Duration(interval) * time.Millisecond)
		w.Write([]byte(name + ":ok"))
	}))
}

func main() {
	result, err := first(fakeWeatherServer("open-weather-1", 200),
		fakeWeatherServer("open-weather-2", 1000),
		fakeWeatherServer("open-weather-3", 600))
	if err != nil {
		log.Println("invoke first error:", err)
		return
	}

	fmt.Println(result)
	time.Sleep(10 * time.Second)
}

func foo() {
	c := make(chan int, 5)

	value, ok := <-c

	for {
		value, ok := <-c
	}

	select {
	case value, ok := <-c:
		// ...
	case time.After(5 * time.Millisecond):
		// 超时之后的动作。
	}

	for v := range c {
		//
	}
}
