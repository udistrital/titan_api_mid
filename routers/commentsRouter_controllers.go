package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionContratosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionContratosController"],
		beego.ControllerComments{
			Method: "ListarContratosAgrupadosPorPersona",
			Router: `/listar_contratos_agrupados_por_persona/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

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

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "DesagregadoNominaPorDependencia",
			Router: `/desagregado_nomina_por_dependencia/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "DesagregadoNominaPorFacultad",
			Router: `/desagregado_nomina_por_facultad/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "DesagregadoNominaPorProyectoCurricular",
			Router: `/desagregado_nomina_por_pc/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "TotalNominaPorDependencia",
			Router: `/total_nomina_por_dependencia/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "TotalNominaPorFacultad",
			Router: `/total_nomina_por_facultad/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "TotalNominaPorProyecto",
			Router: `/total_nomina_por_proyecto/`,
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
