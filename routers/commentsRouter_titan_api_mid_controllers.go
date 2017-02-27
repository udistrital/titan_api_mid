package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionController"] = append(beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionController"],
		beego.ControllerComments{
			"Preliquidar",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionFpController"] = append(beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionFpController"],
		beego.ControllerComments{
			"Preliquidar",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionHcController"] = append(beego.GlobalControllerRouter["titan_api_mid/controllers:PreliquidacionHcController"],
		beego.ControllerComments{
			"Preliquidar",
			`/`,
			[]string{"post"},
			nil})

}
