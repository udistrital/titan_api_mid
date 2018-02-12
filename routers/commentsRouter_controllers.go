package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionPersonasAPreliquidarController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionPersonasAPreliquidarController"],
		beego.ControllerComments{
			Method: "ListarPersonasAPreliquidar",
			Router: `/listar_personas_a_preliquidar_argo/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionPersonasAPreliquidarController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionPersonasAPreliquidarController"],
		beego.ControllerComments{
			Method: "ListarPersonasAPreliquidarPendientes",
			Router: `/listar_personas_a_preliquidar_pendientes/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"],
		beego.ControllerComments{
			Method: "Preliquidar",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"],
		beego.ControllerComments{
			Method: "Resumen",
			Router: `/resumen/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
