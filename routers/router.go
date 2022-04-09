package routers

import (
	"SURPTracker/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/server", &controllers.ServerWSController{})
	beego.Router("/server/list", &controllers.ListServerConnectionController{})
	beego.Router("/client", &controllers.ClientController{})
	beego.Router("/server/asso", &controllers.ListServerIDAssociationController{})
}
