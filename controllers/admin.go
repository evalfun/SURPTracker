package controllers

import (
	"SURPTracker/connectionmanager"

	"github.com/astaxie/beego"
)

var adminPasswd string

func init() {

	adminPasswd = beego.AppConfig.String("adminPasswd")
}

type ListServerConnectionController struct {
	beego.Controller
}

func (this *ListServerConnectionController) Get() {
	password := this.GetString("password")
	if password != adminPasswd {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.WriteString("403")
		return
	}
	this.Ctx.Output.JSON(connectionmanager.GetServers(), true, false)
	return
}

type ListServerIDAssociationController struct {
	beego.Controller
}

func (this *ListServerIDAssociationController) Get() {
	password := this.GetString("password")
	if password != adminPasswd {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.WriteString("403")
		return
	}
	type Resp struct {
		UserID          string
		ServerID        string
		AssociationTime int64
	}
	var resp []Resp
	connectionmanager.ServerID2userIDAssociationMap.Range(func(key, value interface{}) bool {
		association := value.(*connectionmanager.ServerID2UserIDAssociation)
		resp = append(resp, Resp{
			AssociationTime: association.AssociationTime,
			UserID:          association.UserID,
			ServerID:        key.(string),
		})
		return true
	})
	this.Ctx.Output.JSON(resp, true, false)
	return
}
