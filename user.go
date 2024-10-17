package main

import (
	"fmt"
	"net"
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
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//广播新用户消息
	u.server.BroadCast(u, "user login out")
}
func (u *User) DoMessage(msg string) {
	u.server.BroadCast(u, msg)
}

// ListenMessage 监听当前用户 channel 的消息
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("login out")
		}
	}
}
