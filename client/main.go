package main

import (
	"fmt"
	"time"
)

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>> 连接服务器失败")
		return
	}
	fmt.Println(">>>> 连接服务器成功")
	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
