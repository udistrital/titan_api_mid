package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:LiquidarController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:LiquidarController"],
			beego.ControllerComments{
				Method: "Liquidar",
				Router: `/`,
				AllowHTTPMethods: []string{"post"},
				Params: nil})
}
