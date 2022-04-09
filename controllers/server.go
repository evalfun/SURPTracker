package controllers

import (
	"SURPTracker/connectionmanager"
	"SURPTracker/utils"
	"encoding/json"
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var pushAddrListRateLimit *utils.SimpleRateLimiter
var serverWSConnectRateLimit *utils.SimpleRateLimiter

func init() {
	ServerPushAddrListRateDuration, err := beego.AppConfig.Int("ServerPushAddrListRateDuration")
	if err != nil {
		log.Panicln("ServerPushAddrListRateDuration: " + err.Error())
	}
	ServerPushAddrListRateLimit, err := beego.AppConfig.Int("ServerPushAddrListRateLimit")
	if err != nil {
		log.Panicln("ServerPushAddrListRateLimit: " + err.Error())
	}
	pushAddrListRateLimit = utils.NewSimpleRateLimiter(ServerPushAddrListRateDuration, ServerPushAddrListRateLimit)
	ServerWSConnectRateDuration, err := beego.AppConfig.Int("ServerWSConnectRateDuration")
	if err != nil {
		log.Panicln("ServerWSConnectRateDuration: " + err.Error())
	}
	ServerWSConnectRateLimit, err := beego.AppConfig.Int("ServerWSConnectRateLimit")
	if err != nil {
		log.Panicln("ServerWSConnectRateLimit: " + err.Error())
	}
	serverWSConnectRateLimit = utils.NewSimpleRateLimiter(ServerWSConnectRateDuration, ServerWSConnectRateLimit)
}

type ServerWSController struct {
	beego.Controller
}

func (this *ServerWSController) Get() {
	if serverWSConnectRateLimit.BelowRate(this.Ctx.Input.IP()) == false {
		this.Ctx.ResponseWriter.WriteHeader(503)
		this.Ctx.WriteString("excceed rate limit")
		log.Printf("user %s - websocket client - connect error: excceed rate limit", this.Ctx.Input.IP())
		return
	}
	server_id := this.GetString("server_id")
	user_id := this.GetString("user_id")

	if !(utils.IsValidUUID(server_id) && utils.IsValidUUID(user_id)) {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.WriteString("device_id and user_id must be uuid format")
		return
	}
	rest, reason := connectionmanager.PermCheckServerWebSocketConnect(user_id, server_id)
	if !rest {
		this.Ctx.ResponseWriter.WriteHeader(403)
		this.Ctx.WriteString(reason)
		log.Printf("user %s %s websocket client %s connect fail: %s", this.Ctx.Input.IP(), user_id, server_id, reason)
		return
	}

	log.Printf("user %s %s websocket client %s connect to server", this.Ctx.Input.IP(), user_id, server_id)
	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if err != nil {
		log.Printf("WS connection error: %s", err.Error())
		return
	}
	defer func() {
		connectionmanager.ServerOffline(server_id, "Connection closed")
	}()
	connectionmanager.ServerOnline(server_id, &connectionmanager.ServerConnection{
		WSConnection: ws,
		WSAddr:       this.Ctx.Input.IP(),
		UserID:       user_id,
		CreateTime:   time.Now().Unix(),
	})
	var data connectionmanager.Message
	this.Ctx.WriteString("")
	for {
		_, raw_data, err := ws.ReadMessage()
		if err != nil {
			log.Printf("user %s %s websocket client %s connect error: %s", this.Ctx.Input.IP(), user_id, server_id, err.Error())
			return
		}
		err = json.Unmarshal(raw_data, &data)
		if err != nil {
			log.Printf("user %s %s websocket client %s parse json error: %s", this.Ctx.Input.IP(), user_id, server_id, err.Error())
			continue
		}
		if data.MessageType == connectionmanager.MSG_TYPE_UPDATEIP {
			if pushAddrListRateLimit.BelowRate(server_id) == false {
				continue
			}
			var addrList []string
			if len(data.AddrList) > 10 {
				data.AddrList = data.AddrList[:10]
			}
			for i := range data.AddrList {
				if len(data.AddrList[i]) < 128 && len(data.AddrList[i]) > 3 {
					addrList = append(addrList, data.AddrList[i])
				}
			}
			connectionmanager.AssociatServerIDUserID(user_id, server_id)
			log.Printf("user %s %s websocket client %s push addr: %v", this.Ctx.Input.IP(), user_id, server_id, addrList)
			connectionmanager.SetServerAddrList(server_id, addrList)

		}
	}
}
