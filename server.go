package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// 创建server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Message广播消息channel的goroutine，一旦有消息就发给所有在线的User
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

// 广播消息的方法
func (this *Server) BroudCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//...当前连接的任务
	//fmt.Println("连接建立成功")

	user := NewUser(conn, this)
	//用户上线, 将用户加入到OnlineMap中
	user.Online()

	//监听用户是否活跃的channel
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
				fmt.Println("Conn Read err", err)
				return
			}

			//提取用户的消息，去除'\n'
			msg := string(buf[:n-1])

			//将得到的消息进行广播
			user.DoMessage(msg)

			//用户的任意消息，代表当前用户是一个活跃的
			isLive <- true
		}
	}()
	//当前handler阻塞,不能让Handler的goroutine死亡，防止这里面的子功能全死了（用户就没了）
	//select {}
	for {
		select {
		case <-isLive:
			//当前用户是活跃的，应该重置定时器
			//激活select,更新定时器
			//抢先执行，保证下面一条case超时处理不执行，但执行case语句计算（也就是执行time.After来重置定时器）
			//这是因为select的case语句一定会计算，而满足条件的case里的处理内容会随机执行
		case <-time.After(time.Second * 300):
			//已经超时
			//将当前的User强制关闭

			user.SendMsg("你超时啦，被踢咯\n")

			//销毁用的资源
			close(user.C)

			//关闭连接
			conn.Close()

			//退出当前的Handler
			return //runtime.Goexit()
		}
	}
}

// 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	//第二个参数是由ip和端口用冒号组成的地址，可以用Sprintf拼接成字符串
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listen accept err: ", err)
			continue
		}

		//do handler
		go this.Handler(conn)

	}

}
