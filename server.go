package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer 创建一个服务器对象
func NewServer(ip string, port int) *Server {
	server := &Server{Ip: ip, Port: port}
	return server
}

// Handler 处理客户端的连接请求
func (s *Server) Handler(conn net.Conn) {
	// 处理客户端业务
	fmt.Println("connected success")
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
