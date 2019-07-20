package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:Concepto_nomina_por_personaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:Concepto_nomina_por_personaController"],
		beego.ControllerComments{
			Method: "TrRegistroIncapacidades",
			Router: `/tr_registro_incapacidades`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:Concepto_nomina_por_personaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:Concepto_nomina_por_personaController"],
		beego.ControllerComments{
			Method: "TrRegistroProrrogaIncapacidad",
			Router: `/tr_registro_prorroga_incapacidad`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionContratosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionContratosController"],
		beego.ControllerComments{
			Method: "ListarContratosAgrupadosPorPersona",
			Router: `/listar_contratos_agrupados_por_persona`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionOpsController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionOpsController"],
		beego.ControllerComments{
			Method: "GenerarOrdenPago",
			Router: `/generar_op`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionPersonasAPreliquidarController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionPersonasAPreliquidarController"],
		beego.ControllerComments{
			Method: "ListarPersonasAPreliquidar",
			Router: `/listar_personas_a_preliquidar_argo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionPersonasAPreliquidarController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionPersonasAPreliquidarController"],
		beego.ControllerComments{
			Method: "ListarPersonasAPreliquidarPendientes",
			Router: `/listar_personas_a_preliquidar_pendientes`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "DesagregadoNominaPorDependencia",
			Router: `/desagregado_nomina_por_dependencia`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "DesagregadoNominaPorFacultad",
			Router: `/desagregado_nomina_por_facultad`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "DesagregadoNominaPorProyectoCurricular",
			Router: `/desagregado_nomina_por_pc`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "GetOrdenadoresGasto",
			Router: `/get_ordenadores_gasto`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "TotalNominaPorDependencia",
			Router: `/total_nomina_por_dependencia`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "TotalNominaPorFacultad",
			Router: `/total_nomina_por_facultad`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "TotalNominaPorOrdenador",
			Router: `/total_nomina_por_ordenador`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
		beego.ControllerComments{
			Method: "TotalNominaPorProyecto",
			Router: `/total_nomina_por_proyecto`,
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
			Method: "GetIBCPorNovedad",
			Router: `/get_ibc_novedad`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"],
		beego.ControllerComments{
			Method: "PersonasPorPreliquidacion",
			Router: `/personas_x_preliquidacion`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"],
		beego.ControllerComments{
			Method: "ResumenConceptos",
			Router: `/resumen_conceptos`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:ServicesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:ServicesController"],
		beego.ControllerComments{
			Method: "DesagregacionContratoHCS",
			Router: `/desagregacion_contrato_hcs`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
