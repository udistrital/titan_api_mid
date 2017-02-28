package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"],
		beego.ControllerComments{
			Method: "Preliquidar",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionFpController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionFpController"],
		beego.ControllerComments{
			Method: "Preliquidar",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionHcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionHcController"],
		beego.ControllerComments{
			Method: "Preliquidar",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}
