// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/titan_api_mid/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/preliquidacion",
			beego.NSInclude(
				&controllers.PreliquidacionController{},
			),
		),

		beego.NSNamespace("/gestion_ops",
			beego.NSInclude(
				&controllers.GestionOpsController{},
			),
		),

		beego.NSNamespace("/detalle_preliquidacion",
			beego.NSInclude(
				&controllers.DetallePreliquidacionController{},
			),
		),

		beego.NSNamespace("/novedad",
			beego.NSInclude(
				&controllers.NovedadController{},
			),
		),
		beego.NSNamespace("/contrato_preliquidacion",
			beego.NSInclude(
				&controllers.CumplidoController{},
			),
		),
		beego.NSNamespace("/notificaciones",
			beego.NSInclude(
				&controllers.NotificacionesController{},
			),
		),

		beego.NSNamespace("/contratos",
			beego.NSInclude(
				&controllers.ContratosController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
