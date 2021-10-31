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

func liquidarCPS(contrato models.Contrato) {

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
	cedula, err := strconv.ParseInt(contrato.Documento, 0, 64)

	if err == nil {
		reglasAlivios, contratoPreliquidacion = CargarDatosRetefuente(int(cedula))
	}

	mesIterativo = int(contrato.FechaInicio.Month())

	diasContrato, mesesContrato := calcularDiasContratoCPS(contrato.FechaInicio, contrato.FechaFin)

	predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
	predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(diasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
	reglasbase := cargarReglasBase("CT") + cargarReglasSS() + reglasAlivios + FormatoReglas(predicados)

	if contrato.FechaInicio.Day() != 1 {
		mesesContrato = mesesContrato + 1
	}

	for i := 0; i < int(mesesContrato); i++ {
		reglasNuevas = ""
		query := "Ano:" + strconv.Itoa(contrato.Vigencia) + ",Mes:" + strconv.Itoa(mesIterativo) + ",Nominaid:414"
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &preliquidacion)
			if preliquidacion[0].Id == 0 {
				preliquidacion[0] = registrarPreliquidacion(contrato.Vigencia, mesIterativo, 476, 414)
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
				if auxNovedades[0].Id != 0 {
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

			auxDetalle = golog.LiquidarMesCPS(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion, novedades)
			for j := 0; j < len(auxDetalle); j++ {
				registrarDetallePreliquidacion(auxDetalle[j])
			}
			if mesIterativo == 12 {
				break
			} else {
				mesIterativo = mesIterativo + 1
			}
		} else {
			fmt.Println("Error al consultar preliquidaciones")
		}
		preliquidacion[0].Id = 0 //PAra evitar errores al obtener la preliquidación del siguiente mes
	}

}
