package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// PreliquidacionctController operations for Preliquidacionct
type PreliquidacionctController struct {
	beego.Controller
}

func liquidarCPS(contrato models.Contrato, diasRestantes int, ano int) {

	var mesIterativo int              //mes para iterar en el ciclo para liquidar todos los meses de una vez
	var predicados []models.Predicado //variable para inyectar reglas
	var preliquidacion []models.Preliquidacion
	var auxNovedades []models.Novedad
	var novedades []models.Novedad //Arreglo para agregar novedades
	var contratoPreliquidacion models.ContratoPreliquidacion
	var detallePreliquidacion models.DetallePreliquidacion
	var diasALiquidar string
	var aux map[string]interface{}
	var auxDetalle []models.DetallePreliquidacion
	var reglasAlivios string
	var reglasNuevas string //reglas a usar en cada iteracion
	var diasContrato float64
	cedula, err := strconv.ParseInt(contrato.Documento, 0, 64)

	if err == nil {
		reglasAlivios, contratoPreliquidacion = CargarDatosRetefuente(int(cedula))
	}

	if ano != contrato.Vigencia { //en caso de que no sea el mismo año en el que se empezó el contrato se procede a liquidar desde el mes uno (cambio de año)
		mesIterativo = 1
	} else {
		mesIterativo = int(contrato.FechaInicio.Month())
	}
	if diasRestantes == 0 {
		diasContrato, _ = CalcularDias(contrato.FechaInicio, contrato.FechaFin)
		diasContrato = diasContrato + 1 //dia inclusive
	} else {
		diasContrato = float64(diasRestantes)
	}

	predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
	predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", diasContrato) + "," + strconv.Itoa(ano) + "). "})
	reglasbase := cargarReglasBase("CT") + cargarReglasSS() + reglasAlivios + FormatoReglas(predicados)

	for mesIterativo <= 12 {
		reglasNuevas = ""

		query := "Ano:" + strconv.Itoa(ano) + ",Mes:" + strconv.Itoa(mesIterativo) + ",Nominaid:414"
		fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/preliquidacion?limit=-1&query=" + query)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &preliquidacion)
			fmt.Println(preliquidacion[0])
			if preliquidacion[0].Id == 0 {
				preliquidacion[0] = registrarPreliquidacion(ano, mesIterativo, 476, 414)
				contratoPreliquidacion = registrarContratoPreliquidacion(preliquidacion[0].Id, contrato.Id, contratoPreliquidacion)
			} else {
				contratoPreliquidacion = registrarContratoPreliquidacion(preliquidacion[0].Id, contrato.Id, contratoPreliquidacion)
			}

			detallePreliquidacion.ContratoPreliquidacionId = &contratoPreliquidacion
			detallePreliquidacion.TipoPreliquidacionId = 397
			detallePreliquidacion.Activo = true
			detallePreliquidacion.EstadoDisponibilidadId = 426
			diasALiquidar, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
			detallePreliquidacion.DiasLiquidados, _ = strconv.ParseFloat(diasALiquidar, 64)
			reglasNuevas = reglasNuevas + reglasbase + "dias_liquidados(" + contrato.Documento + "," + diasALiquidar + ")."

			//Consultar novedades que se aplican para ese contrato en ese mes

			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad?limit=-1&query=ContratoId.Id:"+strconv.Itoa(contrato.Id), &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &auxNovedades)
				if len(auxNovedades) != 0 {
					for i := 0; i < len(auxNovedades); i++ {
						if auxNovedades[i].FechaInicio.Year() == auxNovedades[i].FechaFin.Year() {
							if preliquidacion[0].Mes >= int(auxNovedades[i].FechaInicio.Month()) && preliquidacion[0].Mes <= int(auxNovedades[i].FechaFin.Month()) && preliquidacion[0].Ano == auxNovedades[i].FechaFin.Year() {
								novedades = append(novedades, auxNovedades[i])
							}
						} else {
							if preliquidacion[0].Mes >= int(auxNovedades[i].FechaInicio.Month()) && preliquidacion[0].Mes <= (int(auxNovedades[i].FechaFin.Month())+12) && preliquidacion[0].Ano <= novedades[i].FechaFin.Year() && preliquidacion[0].Ano >= novedades[i].FechaInicio.Year() {
								novedades = append(novedades, auxNovedades[i])
							}
						}
					}
				}
			} else {
				fmt.Println("Error al obtener las novedades")
			}

			auxDetalle = golog.LiquidarMesCPS(reglasNuevas, contrato.Documento, ano, detallePreliquidacion, novedades)
			for j := 0; j < len(auxDetalle); j++ {
				registrarDetallePreliquidacion(auxDetalle[j])
			}
			if mesIterativo == int(contrato.FechaFin.Month()) && ano == contrato.FechaFin.Year() {
				break
			} else {
				mesIterativo = mesIterativo + 1
			}
		} else {
			fmt.Println("Error al consultar preliquidaciones")
		}
		preliquidacion[0].Id = 0 //Para evitar errores al obtener la preliquidación del siguiente mes
	}

}
