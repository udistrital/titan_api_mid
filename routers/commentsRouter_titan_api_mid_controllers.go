package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:LiquidarController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:LiquidarController"],
			beego.ControllerComments{
				"Liquidar",
				`/`,
				[]string{"post"},
				nil})

}
