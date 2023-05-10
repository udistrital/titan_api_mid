package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:ContratosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:ContratosController"],
        beego.ControllerComments{
            Method: "ObtenerContratosDVE",
            Router: "/docentesDVE/:nomina/:mes:/:ano",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:CumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:CumplidoController"],
        beego.ControllerComments{
            Method: "ActualizarCumplido",
            Router: "/cumplido/:ano/:mes/:contrato/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:CumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:CumplidoController"],
        beego.ControllerComments{
            Method: "ActualizarCumplidoRp",
            Router: "/cumplido_rp/:ano/:mes/:contrato/:vigencia/:rp",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:CumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:CumplidoController"],
        beego.ControllerComments{
            Method: "ActualizarPreliquidado",
            Router: "/preliquidado/:ano/:mes/:contrato/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:DesagregadoHCSController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:DesagregadoHCSController"],
        beego.ControllerComments{
            Method: "ObtenerDesagregado",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:DetallePreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:DetallePreliquidacionController"],
        beego.ControllerComments{
            Method: "ObtenerDetalleCT",
            Router: "/obtener_detalle_CT/:ano/:mes/:contrato/:vigencia/:documento",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:DetallePreliquidacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:DetallePreliquidacionController"],
        beego.ControllerComments{
            Method: "ObtenerDetalleDVE",
            Router: "/obtener_detalle_DVE/:ano/:mes/:documento/:nomina",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionOpsController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:GestionOpsController"],
        beego.ControllerComments{
            Method: "GenerarOrdenPago",
            Router: "/generar_op/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NotificacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NotificacionesController"],
        beego.ControllerComments{
            Method: "EnviarNotificacion",
            Router: "/enviar_notificacion/:dependencia/:mes/:ano",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadCPSController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadCPSController"],
        beego.ControllerComments{
            Method: "CancelarContrato",
            Router: "/cancelar_contrato",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadCPSController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadCPSController"],
        beego.ControllerComments{
            Method: "CederContrato",
            Router: "/ceder_contrato",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadCPSController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadCPSController"],
        beego.ControllerComments{
            Method: "AplicarOtrosi",
            Router: "/otrosi_contrato",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadCPSController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadCPSController"],
        beego.ControllerComments{
            Method: "ReiniciarContrato",
            Router: "/reiniciar_contrato",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadCPSController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadCPSController"],
        beego.ControllerComments{
            Method: "SuspenderContrato",
            Router: "/suspender_contrato",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"],
        beego.ControllerComments{
            Method: "AgregarNovedad",
            Router: "/agregar_novedad",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"],
        beego.ControllerComments{
            Method: "AplicarAnulacion",
            Router: "/aplicar_anulacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"],
        beego.ControllerComments{
            Method: "AplicarReduccion",
            Router: "/aplicar_reduccion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"],
        beego.ControllerComments{
            Method: "EliminarNovedad",
            Router: "/eliminar_novedad/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"],
        beego.ControllerComments{
            Method: "GenerarAdicion",
            Router: "/generar_adicion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"] = append(beego.GlobalControllerRouter["github.com/udistrital/titan_api_mid/controllers:NovedadVEController"],
        beego.ControllerComments{
            Method: "VerificarDescuentos",
            Router: "/verificar_descuentos",
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

}
