// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"titan_api_mid/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/preliquidacion",
			beego.NSInclude(
				&controllers.PreliquidacionController{},
			),
		),

		beego.NSNamespace("/liquidacion",
			beego.NSInclude(
				&controllers.LiquidarController{},
			),
		),

		beego.NSNamespace("/detalle_liquidacion",
			beego.NSInclude(
				&controllers.DetalleLiquidacionController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
