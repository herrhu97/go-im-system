package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// BroadCast 将用户信息加msg，发送到Message进行广播
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

// Handler 处理连接conn
func (this *Server) Handler(conn net.Conn) {
	//处理的业务逻辑
	//fmt.Println("connected...")

	user := NewUser(conn, this)

	user.Online()

	isAlive := make(chan bool)

	//接受客户端发来的消息
	go func() {
		buff := make([]byte, 4096)
		for {
			n, err := conn.Read(buff)
			//ctrl + D
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}

			//去除'\n'，乱码问题
			msg := string(buff[:n-1])
			// 将消息加上用户信息广播
			user.DoMsg(msg)

			isAlive <- true
		}

	}()

	// 当前goroutine阻塞
	for {
		select {
		case <-isAlive:
		//	啥也不做，为了激活，为了更新计时器

		case <-time.After(time.Hour * 5):
			//计时器到点，触发此case，执行踢人逻辑
			user.SendMsg("你被踢了")

			close(isAlive)

			conn.Close()

			return
		}
	}
}

func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

// ListenMessage 监听Message，并发给所有User
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// Start 启动服务器的接口
func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		fmt.Println("net listen err:", err)
	}

	defer listener.Close()

	go this.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
		}

		go this.Handler(conn)
	}

}
