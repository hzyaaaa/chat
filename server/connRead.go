package server

import (
	"net"
)

func connRead(conn net.Conn, overTime chan bool) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf) //阻塞等待别人写数据
		//与服务器通信的客户端用户名
		thisUser := onlineUsers[conn.RemoteAddr().String()].name
		if n == 0 { //当手动关闭客户端的时候，会向tcp发送FIN，此时n==0
			for _, v := range onlineUsers {
				if thisUser != "" { //往每个已经登录的用户里写入“xxx已退出”，真正退出的时候拿来用！
					v.AddUser <- []byte("用户[" + thisUser + "]已退出\n")
				}
			}
			delete(onlineUsers, conn.RemoteAddr().String())
			return
		}

		if err != nil {
			ccF.Println("连接读取失败:", err)
			return
		}
		//处理消息内容
		var msg []byte
		//不等于"\n"
		if buf[0] != 10 {
			//控制台客户端：buf[n-1]->13,buf[n]->10
			//ccF.Println("buf[:n]:", buf[:n])
			//buf[:n-2]去除\n换行符10和enter回车符13
			msg = append([]byte("["+thisUser+"]说>:"), buf[:n-2]...)
		} else {
			msg = nil
		}
		//发送消息到信道
		overTime <- true //只要read到数据，overTime就是true
		message <- msg
		LogToDb(string(msg), thisUser)
	}
}
