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

		beego.NSNamespace("/gestion_personas_a_liquidar",
			beego.NSInclude(
				&controllers.GestionPersonasAPreliquidarController{},
			),
		),

		beego.NSNamespace("/gestion_contratos",
			beego.NSInclude(
				&controllers.GestionContratosController{},
			),
		),

		beego.NSNamespace("/gestion_reportes",
			beego.NSInclude(
				&controllers.GestionReportesController{},
			),
		),

		beego.NSNamespace("/concepto_nomina_por_persona",
			beego.NSInclude(
				&controllers.Concepto_nomina_por_personaController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
