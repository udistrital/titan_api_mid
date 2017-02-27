package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionController"] = append(beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionController"],
		beego.ControllerComments{
			Method:           "Preliquidar",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			Params:           nil})

	beego.GlobalControllerRouter["titan_api_mid/controllers:DetalleLiquidacionController"] = append(beego.GlobalControllerRouter["titan_api_mid/controllers:DetalleLiquidacionController"],
		beego.ControllerComments{
			Method:           "InsertarDetallePreliquidacion",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			Params:           nil})

	beego.GlobalControllerRouter["titan_api_mid/controllers:LiquidarController"] = append(beego.GlobalControllerRouter["titan_api_mid/controllers:LiquidarController"],
		beego.ControllerComments{
			Method:           "Liquidar",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			Params:           nil})

	beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionFpController"] = append(beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionFpController"],
		beego.ControllerComments{
			Method:           "Preliquidar",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			Params:           nil})

	beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionHcController"] = append(beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionHcController"],
		beego.ControllerComments{
			Method:           "Preliquidar",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			Params:           nil})

}
