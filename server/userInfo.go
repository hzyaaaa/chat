package server

//userInfo用于存储用户信息
type userInfo struct {
	name    string
	perC    chan []byte //除了上线的消息通知外，其他的消息都是通过信道perC完成的。
	AddUser chan []byte //广播用户进入或退出
}
