package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// NewUser 创建一个用户的方法
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

// Online 用户上线
func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播新用户消息
	u.server.BroadCast(u, "user login")
}

// Offline 用户下线
func (u *User) Offline() {
	//广播新用户消息
	u.server.BroadCast(u, "user login out")
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
}
func (u *User) DoMessage(msg string) {
	if len(msg) > 2 && msg[0:2] == "./" {
		u.Instructions(msg[2:])
	} else {
		u.server.BroadCast(u, msg)
	}

}

// Instructions 指令发送器
func (u *User) Instructions(msg string) {

	suffix := []string{"\r\n", "\n", "\r"}
	for _, s := range suffix {
		if strings.HasSuffix(msg, s) {
			msg = msg[:len(msg)-1]
			break
		}
	}

	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			msg := fmt.Sprintf("[%s] is %s...\n ", user.Name, "online")
			u.SendMsg(msg)
		}
		u.server.mapLock.Unlock()
	} else if msg[:7] == "rename|" {
		//消息格式 rename|newName
		newName := strings.Split(msg, "|")[1]
		//判断name是否存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("The Name Already Exists!")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()
			u.Name = newName
			u.SendMsg("Rename Success!")
		}
	} else if len(msg) > 6 && msg[:3] == "to|" {
		msgArr := strings.Split(msg, "|")
		if len(msgArr) != 3 {
			u.SendMsg("The Directive Does Not Exist!")
		}
		target := msgArr[1]
		content := msgArr[2]
		user, ok := u.server.OnlineMap[target]
		if !ok {
			u.SendMsg("The User Does Not Exist!")
		} else {
			user.SendMsg(fmt.Sprintf("%s say: %s", u.Name, content))
		}
	} else {
		u.SendMsg("The Directive Does Not Exist!")
	}

}
func (u *User) SendMsg(msg string) {
	_, err := u.conn.Write([]byte(msg))
	if err != nil {
		return
	}
	//u.C <- msg
}

// ListenMessage 监听当前用户 channel 的消息
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("login out")
			return
		}
	}
}
