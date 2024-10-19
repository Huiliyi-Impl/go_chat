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
			fmt.Println("fmt.Scanln failed, err:", err)
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

// 更新用户名
func (client *Client) updateName() bool {
	go client.dealResponse()
	fmt.Println("请输入用户名:")
	_, err := fmt.Scanln(&client.Name)
	if err != nil {
		fmt.Println("fmt.Scanln failed, err:", err)
		return false
	}
	_, err1 := client.Conn.Write(str2byte("./rename|" + client.Name))
	if err1 != nil {
		fmt.Println("client.Conn.Write failed", err1)
		return false
	}
	return true
}
func (client *Client) selectUser() {
	_, err := client.Conn.Write(str2byte("./who"))
	if err != nil {
		return
	}
}

func (client *Client) privateChat() {
	fmt.Println("请选择在线用户")
	client.selectUser()
	var remoteUser string
	_, err := fmt.Scanln(&remoteUser)
	if err != nil {
		fmt.Println("fmt.Scanln failed, err:", err)
		return
	}
	if remoteUser == client.Name {
		fmt.Println(">>>>can not chat with yourself<<<<")
		return
	}
	if remoteUser == "" {
		fmt.Println(">>>>please enter the correct user name<<<<")
		return
	}
	if remoteUser == "exit" {
		return
	}
	var remoteMsg string
	fmt.Println(">>>>请输入聊天内容(输入exit退出)<<<<")
	for remoteMsg != "exit" {
		_, err := fmt.Scanln(&remoteMsg)
		if err != nil {
			return
		}
		_, err1 := client.Conn.Write(str2byte("./to|" + remoteUser + "|" + remoteMsg))
		if err1 != nil {
			return
		}
		remoteMsg = ""
	}
}

// 公聊模式
func (client *Client) publicChat() {
	go client.dealResponse()
	for {
		var content string
		_, err := fmt.Scanln(&content)
		if err != nil {
			fmt.Println("fmt.Scanln failed, err:", err)
			return
		}
		if content == "exit" {
			_, err2 := client.Conn.Write(str2byte(client.Name + ":exit"))
			if err2 != nil {
				return
			}
			return
		}
		if len(content) != 0 {
			_, err1 := client.Conn.Write(str2byte(content))
			if err1 != nil {
				fmt.Println("client.Conn.Write failed", err1)
				return
			}
		}
	}
}

// 将字符串转化为字节数组并加上换行符
func str2byte(str string) []byte {
	return []byte(str + "\n")
}

// 接收服务器端信息
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
				fmt.Println(">>>>公聊模式(输入exit退出)<<<")
				client.publicChat()
				break
			case 2:
				fmt.Println(">>>私聊模式(按exit退出)<<<")
				client.privateChat()
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
	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial failed, err:", err)
		return nil
	}
	client.Conn = conn
	return client
}
