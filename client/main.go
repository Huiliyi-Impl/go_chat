package main

import (
	"flag"
	"fmt"
)

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "i", "127.0.0.1", "设置服务器IP地址(默认127.0.0.1)")
	flag.IntVar(&serverPort, "p", 8888, "设置服务器端口(默认8888)")
}
func main() {
	// 命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>> 连接服务器失败")
		return
	}
	fmt.Println(">>>> 连接服务器成功")
	client.run()
}
