package connectionmanager

import (
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var serverConnectionMap sync.Map

func GetServers() []*ServerInfo {
	var respList []*ServerInfo

	serverConnectionMap.Range(func(key, value interface{}) bool {
		serverID := key.(string)
		serverInfo := value.(*ServerConnection)
		respList = append(respList, &ServerInfo{
			ServerID:       serverID,
			UserID:         serverInfo.UserID,
			WSAddr:         serverInfo.WSAddr,
			CreateTime:     serverInfo.CreateTime,
			ServerAddrList: serverInfo.ServerAddrList,
			LastReport:     serverInfo.LastReport,
		})
		return true
	})
	return respList
}

type ServerConnection struct {
	ServerAddrList []string
	UserID         string
	WSAddr         string
	WSConnection   *websocket.Conn
	CreateTime     int64
	LastReport     int64
}

type ServerInfo struct {
	ServerID       string
	UserID         string
	WSAddr         string
	CreateTime     int64
	LastReport     int64
	ServerAddrList []string
}

func GetServerConnectionByUUID(uuid string) *ServerConnection {
	_connection, ok := serverConnectionMap.Load(uuid)
	if !ok {
		return nil
	}
	return _connection.(*ServerConnection)
}

func ConnectToClient(uuid string, addrList []string) error {
	_connection, ok := serverConnectionMap.Load(uuid)
	if !ok {
		return errors.New("server offline or not exist")
	}
	connection := _connection.(*ServerConnection)
	err := connection.WSConnection.WriteJSON(Message{
		MessageType: MSG_TYPE_CONNECT,
		AddrList:    addrList,
	})
	if err != nil {
		return err
	}
	return nil
}

func ServerOnline(uuid string, connection *ServerConnection) {
	_conn, ok := serverConnectionMap.Load(uuid)
	if ok {
		conn := _conn.(*ServerConnection)
		conn.WSConnection.Close()
		serverConnectionMap.Delete(uuid)
		ServerOffline(uuid, "another client online")
	}

	serverConnectionMap.Store(uuid, connection)
}
func ServerOffline(uuid string, reason string) {
	_conn, ok := serverConnectionMap.Load(uuid)
	if !ok {
		return
	}
	conn := _conn.(*ServerConnection)
	conn.WSConnection.Close()
	serverConnectionMap.Delete(uuid)
}

func SetServerAddrList(uuid string, addrList []string) {
	_conn, ok := serverConnectionMap.Load(uuid)
	if !ok {
		return
	}
	conn := _conn.(*ServerConnection)
	conn.LastReport = time.Now().Unix()
	conn.ServerAddrList = addrList
}
