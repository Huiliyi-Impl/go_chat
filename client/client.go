package main

import (
	"fmt"
	"io"
	"math"
	"net"
	"os"
)

var modelList = []string{"1:公聊模式", "2:私聊模式", "3:更新用户名", "0:退出"}

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	flag       int
}

// 菜单显示
func (client *Client) menu() bool {
	var flag = math.MaxInt8
	for flag > len(modelList) {
		for _, value := range modelList {
			fmt.Println(value)
		}
		_, err := fmt.Scanln(&flag)
		if err != nil {
			return false
		}
		if flag > len(modelList)-1 {
			fmt.Println(">>>>if you have made a mistake please re enter it<<<<")
			return false
		}
	}
	client.flag = flag
	return true
}
func (client *Client) updateName() bool {
	fmt.Println("请输入用户名:")
	_, err := fmt.Scanln(&client.Name)
	if err != nil {
		fmt.Println("fmt.Scanln failed, err:", err)
		return false
	}
	_, err1 := client.Conn.Write([]byte("./rename|" + client.Name + "\n"))
	if err1 != nil {
		fmt.Println("client.Conn.Write failed", err1)
		return false
	}
	return true
}
func (client *Client) dealResponse() {
	_, err := io.Copy(os.Stdout, client.Conn)
	if err != nil {
		return
	}
}
func (client *Client) run() {
	for client.flag != 0 {
		if client.menu() {
			switch client.flag {
			case 1:
				fmt.Println("公聊模式>>>>")
				break
			case 2:
				fmt.Println("私聊模式>>>>")
				break
			case 3:
				fmt.Println("更新用户名>>>>")
				client.updateName()
				break
			}
		}
	}
	fmt.Println(">>>> exit")
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial failed, err:", err)
		return nil
	}
	client.Conn = conn
	return client
}
