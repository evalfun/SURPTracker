package controllers

import (
	"SURPTracker/connectionmanager"
	"SURPTracker/utils"
	"encoding/json"
	"log"

	"github.com/astaxie/beego"
)

type ClientController struct {
	beego.Controller
}

var getAddrLimit *utils.SimpleRateLimiter
var reqConnectLimit *utils.SimpleRateLimiter

func init() {
	ClientGetAddrListRateDuration, err := beego.AppConfig.Int("ClientGetAddrListRateDuration")
	if err != nil {
		log.Panicln("ClientGetAddrListRateDuration: " + err.Error())
	}
	ClientGetAddrListRateLimit, err := beego.AppConfig.Int("ClientGetAddrListRateLimit")
	if err != nil {
		log.Panicln("ClientGetAddrListRateLimit: " + err.Error())
	}
	getAddrLimit = utils.NewSimpleRateLimiter(ClientGetAddrListRateDuration, ClientGetAddrListRateLimit)
	ClientRequestConnectRateDuration, err := beego.AppConfig.Int("ClientRequestConnectRateDuration")
	if err != nil {
		log.Panicln("ClientRequestConnectRateDuration: " + err.Error())
	}
	ClientRequestConnectRateLimit, err := beego.AppConfig.Int("ClientRequestConnectRateLimit")
	if err != nil {
		log.Panicln("ClientRequestConnectRateLimit: " + err.Error())
	}
	reqConnectLimit = utils.NewSimpleRateLimiter(ClientRequestConnectRateDuration, ClientRequestConnectRateLimit)
}

//获取服务器地址
func (this *ClientController) Get() {
	if getAddrLimit.BelowRate(this.Ctx.Input.IP()) == false {
		this.Ctx.ResponseWriter.WriteHeader(503)
		this.Ctx.WriteString("excceed rate limit")
		log.Printf("user %s - Client: query server - fail: excceed rate limit", this.Ctx.Input.IP())
		return
	}
	server_id := this.GetString("server_id")
	user_id := this.GetString("user_id")
	if user_id != "" && !utils.IsValidUUID(user_id) {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.WriteString("invaild user_id")
		log.Printf("user %s %s Client: query server %s fail: invalid user id", this.Ctx.Input.IP(), user_id, server_id)
		return
	}
	if user_id == "" {
		user_id = "anonymous"
	}
	if !utils.IsValidUUID(server_id) {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.WriteString("invaild server_id")
		log.Printf("user %s %s Client: query server %s fail: invalid server id", this.Ctx.Input.IP(), user_id, server_id)
		return
	}

	server := connectionmanager.GetServerConnectionByUUID(server_id)
	if server == nil {
		this.Ctx.ResponseWriter.WriteHeader(404)
		this.Ctx.WriteString("server offline or did not exist")
		log.Printf("user %s %s Client: query server %s fail: server did not found", this.Ctx.Input.IP(), user_id, server_id)
		return
	}
	log.Printf("user %s %s Client: query server %s", this.Ctx.Input.IP(), user_id, server_id)
	this.Ctx.Output.JSON(server.ServerAddrList, true, false)
}

//请求服务器发起主动连接
func (this *ClientController) Post() {
	if reqConnectLimit.BelowRate(this.Ctx.Input.IP()) == false {
		this.Ctx.ResponseWriter.WriteHeader(503)
		log.Printf("user %s - req server connect - fail: excceed rate limit", this.Ctx.Input.IP())
		this.Ctx.WriteString("excceed rate limit")
		return
	}
	type Request struct {
		ServerID string
		UserID   string
		AddrList []string //客户端地址列表
	}
	var req Request
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
		log.Printf("user %s - Client: request server connect - fail: %s", this.Ctx.Input.IP(), err.Error())
		this.Ctx.WriteString("error json")
		return
	}
	if !utils.IsValidUUID(req.ServerID) {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.WriteString("invalid ServerID")
		log.Printf("user %s %s Client: request server connect %s fail: invalid server id", this.Ctx.Input.IP(), req.UserID, req.ServerID)
		return
	}
	if !utils.IsValidUUID(req.UserID) && req.UserID != "" {
		this.Ctx.ResponseWriter.WriteHeader(400)
		log.Printf("user %s %s Client: request server connect %s fail: invalid user id", this.Ctx.Input.IP(), req.UserID, req.ServerID)
		this.Ctx.WriteString("invaild UserID")
		return
	}
	if req.UserID == "" {
		req.UserID = "anonymous"
	}

	if len(req.AddrList) > 10 {
		req.ServerID = req.ServerID[:10]
	}
	var addrList []string
	for i := range req.AddrList {
		if len(req.AddrList) < 128 {
			addrList = append(addrList, req.AddrList[i])
		}
	}
	if len(addrList) == 0 {
		this.Ctx.ResponseWriter.WriteHeader(400)
		log.Printf("user %s %s Client: request server connect %s fail: invaild AddrList", this.Ctx.Input.IP(), req.UserID, req.ServerID)
		this.Ctx.WriteString("invaild AddrList")
		return
	}

	err = connectionmanager.ConnectToClient(req.ServerID, addrList)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.WriteString(err.Error())
		log.Printf("user %s %s Client: request server connect %s fail: %s", this.Ctx.Input.IP(), req.UserID, req.ServerID, err.Error())
		return
	}
	this.Ctx.ResponseWriter.WriteHeader(200)
	this.Ctx.WriteString("success")
	log.Printf("user %s %s Client: request server connect %s success", this.Ctx.Input.IP(), req.UserID, req.ServerID)
	return
}
