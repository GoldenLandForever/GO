package main

import (
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	Ip   string
	Port int
}

//创建一个server的接口

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}
func (this *Server) Handler(conn net.Conn) {
	//..当前链接的业务
	fmt.Println("链接建立成功")
}

// 启动服务器接口
func (this *Server) Run() {
	//socket listen
	listener, err := net.Listen("tcp", this.Ip+":"+strconv.Itoa(this.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	//close linsten socket
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		//do handler
		// go 实现并发
		go this.Handler(conn)
	}
	//accept

	//do handler

}
