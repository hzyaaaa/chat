package server

import (
	"net"
	"time"
)

// HandleConnect 处理客户端请求
func HandleConnect(conn net.Conn) {
	defer conn.Close()
	//信道overTime处理超时
	overTime := make(chan bool)
	bufUName := make([]byte, 4096)
	n, err := conn.Read(bufUName) //读取用户名，此时conn里没有其他消息，然后只用来connWrite协程写数据，和connRead读数据。
	if err != nil {
		ccF.Println("连接读取失败:", err)
		return
	}
	userName := string(bufUName[:n])
	perC := make(chan []byte)
	perAddUser := make(chan []byte)
	//!!!未判断用户名重复问题
	user := userInfo{name: userName, perC: perC, AddUser: perAddUser}
	onlineUsers[conn.RemoteAddr().String()] = user
	//新客户端连接后广播
	go broadcast(userName)

	//监听客户端自己的信道，conn是每个客户端独有的
	go connWrite(conn, user)

	//循环读取客户端发来的消息
	go connRead(conn, overTime)

	for {
		select {
		case <-overTime: //只要该用户read到数据，overtime就是true，只是走这里的空，而每次循环time.After重新计时，只要read阻塞300秒则执行time.After后的代码。
		case <-time.After(time.Second * 300):
			_, _ = conn.Write([]byte("已被系统踢出\n")) //一直read不到数据就write“xx被踢出，然后xx下线”
			thisUser := onlineUsers[conn.RemoteAddr().String()].name
			for _, v := range onlineUsers { //把xx下线的消息发给在线用户集合onlineUsers
				if thisUser != "" {
					v.AddUser <- []byte("用户[" + thisUser + "]已被踢出\n") //告诉其他人，某某已经被踢出！
				}
			}
			delete(onlineUsers, conn.RemoteAddr().String())
			return
		}
	}
}
