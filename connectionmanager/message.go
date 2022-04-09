package connectionmanager

const (
	MSG_TYPE_UPDATEIP = "UPDATEIP"
	MSG_TYPE_OFFLINE  = "OFFLINE"
	MSG_TYPE_CONNECT  = "CONNECT"
)

type Message struct {
	MessageType string
	AddrList    []string
	Message     string
}

// 消息类型
// UPDATEIP 更新服务端ip列表
// OFFLINE  服务端连接断开
// CONNECT  服务端注重发起连接到客户端
