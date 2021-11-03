package main

import (
	"bufio"
	"bytes"
	"flag"
	"os"
	"sync"
)

var (
	file   = flag.String("f", "access.log", "输入文件")
	emaill = flag.String("e", "", "email")
	out    = flag.String("o", "out.log", "输出文件")
	ipp    = flag.Bool("ip", false, "获取使用者ip")
	urll   = flag.Bool("url", false, "获取访问路径")
	o      *os.File
)

func main() {
	flag.Parse()
	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	o, err = os.OpenFile(*out, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer o.Close()
	buf := bufio.NewReader(f)
	ips, urls := sync.Map{}, sync.Map{}
	w := sync.WaitGroup{}
	ch := make(chan struct{}, 10) // 多协程处理
	for {
		line, next, err := buf.ReadLine()
		if err != nil {
			panic(err)
		}
		if !next {
			break
		}
		ch <- struct{}{}
		w.Add(1)
		go func(line []byte) {
			defer func() {
				<-ch
				w.Done()
			}()
			arr := bytes.Split(line, []byte(" "))
			if len(arr) == 7 && string(arr[6]) == *emaill {
				//time := bytes.Join(arr[:2], []byte(" "))
				ip := bytes.Split(arr[2], []byte(":"))[0]
				uri := bytes.Split(arr[4], []byte(":"))
				//typee := uri[0]
				url := uri[1]
				//port := uri[2]

				// todo 完成一个ip对应多个路径
				if *ipp {
					ips.Store(ip, struct{}{})
				}
				if *urll {
					urls.Store(url, struct{}{})
				}
			}
		}(line)
	}
	w.Wait()
	if *ipp {
		ips.Range(func(key, value interface{}) bool {
			o.Write(key.([]byte))
			o.Write([]byte("\n"))
			return true
		})
	}
	if *urll {
		urls.Range(func(key, value interface{}) bool {
			o.Write(key.([]byte))
			o.Write([]byte("\n"))
			return true
		})
	}
}
