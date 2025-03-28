package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

//创建一个server的接口

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//..当前链接的业务
	//用户上线，加入到onlineMap中
	user := NewUser(conn, this)

	user.Online()

	isLive := make(chan bool)
	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err：", err)
				return
			}

			msg := string(buf[:n-1])

			user.DoMEsssage(msg)

			isLive <- true
		}
	}()
	for {
		select {
		case <-isLive:
			//活跃
			//激活select
		case <-time.After(time.Second * 30):
			//超时
			user.SendMsg("offline")
			close(user.C)
			conn.Close()
			return
		}
	}
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

	go this.ListenMessager()
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		//do handler
		// go 实现并发
		go this.Handler(conn)
	}

}
