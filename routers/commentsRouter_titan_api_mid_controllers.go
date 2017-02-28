package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"],
		beego.ControllerComments{
			"Preliquidar",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionFpController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionFpController"],
		beego.ControllerComments{
			"Preliquidar",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionHcController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionHcController"],
		beego.ControllerComments{
			"Preliquidar",
			`/`,
			[]string{"post"},
			nil})

}
