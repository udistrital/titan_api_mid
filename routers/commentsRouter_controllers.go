package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:Concepto_nomina_por_personaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:Concepto_nomina_por_personaController"],
        beego.ControllerComments{
            Method: "TrRegistroIncapacidades",
            Router: "/tr_registro_incapacidades",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:Concepto_nomina_por_personaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:Concepto_nomina_por_personaController"],
        beego.ControllerComments{
            Method: "TrRegistroProrrogaIncapacidad",
            Router: "/tr_registro_prorroga_incapacidad",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:CumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:CumplidoController"],
        beego.ControllerComments{
            Method: "ActualizarCumplido",
            Router: "/obtener_detalle_CT/:ano/:mes/:contrato",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:DetallePreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:DetallePreliquidacionController"],
        beego.ControllerComments{
            Method: "ObtenerDetalleCT",
            Router: "/obtener_detalle_CT/:ano/:mes/:contrato/:documento",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:DetallePreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:DetallePreliquidacionController"],
        beego.ControllerComments{
            Method: "ObtenerDetalleHCH",
            Router: "/obtener_detalle_HCH/:ano/:mes/:documento",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionContratosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionContratosController"],
        beego.ControllerComments{
            Method: "ListarContratosAgrupadosPorPersona",
            Router: "/listar_contratos_agrupados_por_persona",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionOpsController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionOpsController"],
        beego.ControllerComments{
            Method: "GenerarOrdenPago",
            Router: "/generar_op",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionPersonasAPreliquidarController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionPersonasAPreliquidarController"],
        beego.ControllerComments{
            Method: "ListarPersonasAPreliquidar",
            Router: "/listar_personas_a_preliquidar_argo",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
        beego.ControllerComments{
            Method: "DesagregadoNominaPorDependencia",
            Router: "/desagregado_nomina_por_dependencia",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
        beego.ControllerComments{
            Method: "DesagregadoNominaPorFacultad",
            Router: "/desagregado_nomina_por_facultad",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
        beego.ControllerComments{
            Method: "DesagregadoNominaPorProyectoCurricular",
            Router: "/desagregado_nomina_por_pc",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
        beego.ControllerComments{
            Method: "GetOrdenadoresGasto",
            Router: "/get_ordenadores_gasto",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
        beego.ControllerComments{
            Method: "TotalNominaPorDependencia",
            Router: "/total_nomina_por_dependencia",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
        beego.ControllerComments{
            Method: "TotalNominaPorFacultad",
            Router: "/total_nomina_por_facultad",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
        beego.ControllerComments{
            Method: "TotalNominaPorOrdenador",
            Router: "/total_nomina_por_ordenador",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionReportesController"],
        beego.ControllerComments{
            Method: "TotalNominaPorProyecto",
            Router: "/total_nomina_por_proyecto",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"],
        beego.ControllerComments{
            Method: "AgregarNovedad",
            Router: "/agregar_novedad",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"],
        beego.ControllerComments{
            Method: "CancelarContrato",
            Router: "/cancelar_contrato/:NumeroContrato",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"],
        beego.ControllerComments{
            Method: "CederContrato",
            Router: "/ceder_contrato",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"],
        beego.ControllerComments{
            Method: "EliminarNovedad",
            Router: "/eliminar_novedad/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"],
        beego.ControllerComments{
            Method: "AplicarOtrosi",
            Router: "/otrosi_contrato",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadController"],
        beego.ControllerComments{
            Method: "SuspenderContrato",
            Router: "/suspender_contrato",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"],
        beego.ControllerComments{
            Method: "Preliquidar",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:PreliquidacionController"],
        beego.ControllerComments{
            Method: "ObtenerResumenPreliquidacion",
            Router: "/obtener_resumen_preliquidacion/:mes/:ano/:nomina",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:ServicesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:ServicesController"],
        beego.ControllerComments{
            Method: "DesagregacionContratoHCS",
            Router: "/desagregacion_contrato_hcs",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
