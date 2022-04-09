package connectionmanager

import (
	"log"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

type ServerID2UserIDAssociation struct {
	UserID          string
	AssociationTime int64
}

var ServerID2userIDAssociationMap sync.Map
var KeepServerIDUserIDAssociation int64

func init() {
	var err error
	KeepServerIDUserIDAssociation, err = beego.AppConfig.Int64("KeepServerIDUserIDAssociation")
	if err != nil {
		log.Panicln("KeepServerIDUserIDAssociation: " + err.Error())
	}
	go cleanAssociationMapProc()
}

func cleanAssociationMapProc() {

	for {
		time.Sleep(60 * time.Second)
		currentTime := time.Now().Unix()
		ServerID2userIDAssociationMap.Range(func(key, value interface{}) bool {
			association := value.(*ServerID2UserIDAssociation)
			if currentTime > association.AssociationTime+KeepServerIDUserIDAssociation {
				ServerID2userIDAssociationMap.Delete(key)
			}
			return true
		})
	}

}

func AssociatServerIDUserID(userID string, serverID string) {
	_savedServer, ok := ServerID2userIDAssociationMap.Load(serverID)
	if !ok {
		ServerID2userIDAssociationMap.Store(serverID, &ServerID2UserIDAssociation{
			UserID:          userID,
			AssociationTime: time.Now().Unix(),
		})
		return
	}
	savedServer := _savedServer.(*ServerID2UserIDAssociation)
	savedServer.AssociationTime = time.Now().Unix()
}

//websocket客户端连接检查
func PermCheckServerWebSocketConnect(userID string, serverID string) (bool, string) {
	_savedServer, ok := ServerID2userIDAssociationMap.Load(serverID)
	if !ok {
		AssociatServerIDUserID(userID, serverID)
		return true, ""
	}
	savedServer := _savedServer.(*ServerID2UserIDAssociation)
	if savedServer.UserID != userID {
		return false, "This server id is already used by other user"
	} else {
		savedServer.AssociationTime = time.Now().Unix()
	}
	return true, ""
}

func PermCheckQueryServerAddrList(userID string, serverID string) (bool, string) {
	return true, ""
}

func PermCheckRequestServerConnect(userID string, serverID string) (bool, string) {

	return true, ""
}
