package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerPort int
	ServerIP   string
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(ip string, port int) *Client {
	client := &Client{ServerPort: port, ServerIP: ip, flag: 999}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net dial err", err)
		return nil
	}

	client.conn = conn

	return client
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口")
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag < 4 {
		client.flag = flag
		return true
	} else {
		fmt.Println("输入有问题")
		return false
	}

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		case 0:
			return
		}
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名:")
	fmt.Scanln(&client.Name)

	// 没有加上换行符导致乱码
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err", err)
		return false
	}

	return true
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)

	//for {
	//	buff := make([]byte, 4096)
	//	client.conn.Read(buff)
	//	fmt.Println(string(buff))
	//}
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>输入聊天信息，exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write err", err)
			}
		}

		chatMsg = ""
		fmt.Println(">>>>输入聊天信息，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who" + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err", err)
	}
}

func (client *Client) PrivateChat() {
	var chatMsg string
	var remoteUser string
	client.SelectUsers()
	fmt.Println(">>>>输入聊天对象，exit退出")
	fmt.Scanln(&remoteUser)

	for remoteUser != "exit" {

		fmt.Println(">>>>输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteUser + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write err", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>输入聊天内容，exit退出")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println(">>>>输入聊天对象，exit退出")
		fmt.Scanln(&remoteUser)
	}
}

func main() {
	//解析命令行参数
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>连接失败")
		return
	}

	fmt.Println(">>>>>>>>>连接成功")

	go client.DealResponse()

	client.Run()
}
