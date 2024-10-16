package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	// 消息队列
	Message chan string
}

// NewServer 创建一个服务器对象
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s] %s", user.Name, msg)
	s.Message <- sendMsg
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// Handler 处理客户端的连接请求
func (s *Server) Handler(conn net.Conn) {
	// 处理客户端业务
	fmt.Println("connected success")
	user := NewUser(conn)
	// 将当前用户加入到onlineMap中
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//广播新用户消息
	s.BroadCast(user, "user login")
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				msgExit := fmt.Sprintf("%s login out", user.Name)
				s.BroadCast(user, msgExit)
				// 删除不在线的用户
				s.mapLock.Lock()
				delete(s.OnlineMap, user.Name)
				s.mapLock.Unlock()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("read failed, err:", err)
				return
			}
			// 去除换行符
			msg := string(buf[:n-1])
			s.BroadCast(user, msg)
		}
	}()
}
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	//close listener
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("close listener failed, err:", err)
		}
	}(listener)

	//启动监听message
	go s.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}
}
