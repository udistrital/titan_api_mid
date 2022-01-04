package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// PreliquidacionhchController operations for Preliquidacioncthch
type PreliquidacionhchController struct {
	beego.Controller
}

func liquidarHCH(contrato models.Contrato) {
	var mesIterativo int              //mes para iterar en el ciclo para liquidar todos los meses de una vez
	var anoIterativo int              //Ano iterativo a la hora de liquidar
	var predicados []models.Predicado //variable para inyectar reglas
	var preliquidacion []models.Preliquidacion
	var contratoPreliquidacion models.ContratoPreliquidacion
	var detallePreliquidacion models.DetallePreliquidacion
	var diasALiquidar string
	var aux map[string]interface{}
	var auxDetalle []models.DetallePreliquidacion
	var reglasAlivios string
	var reglasNuevas string //reglas a usar en cada iteracion
	var semanas_liquidadas int
	cedula, err := strconv.ParseInt(contrato.Documento, 0, 64)

	if err == nil {
		reglasAlivios, contratoPreliquidacion = CargarDatosRetefuente(int(cedula))
	}

	mesIterativo = int(contrato.FechaInicio.Month())
	anoIterativo = contrato.Vigencia

	//Obtener las semanas del contrato

	semanasContrato := int(calcularSemanasContratoHCH(contrato.FechaInicio, contrato.FechaFin))

	predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
	predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
	reglasbase := cargarReglasBase("HCH") + cargarReglasSS() + reglasAlivios + FormatoReglas(predicados)

	for {
		reglasNuevas = ""
		query := "Ano:" + strconv.Itoa(contrato.Vigencia) + ",Mes:" + strconv.Itoa(mesIterativo) + ",Nominaid:415"
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &preliquidacion)
			if preliquidacion[0].Id == 0 {
				preliquidacion[0] = registrarPreliquidacion(contrato.Vigencia, mesIterativo, 476, 415)
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
			//Calcular semanas a liquidar
			if mesIterativo == int(contrato.FechaInicio.Month()) && contrato.Vigencia == anoIterativo {
				//para el mes inicial
			} else if mesIterativo == int(contrato.FechaFin.Month()) && contrato.FechaFin.Year() == anoIterativo {
				//Para el mes final
			} else {
				semanas_liquidadas = 4
			}
			reglasNuevas = reglasNuevas + reglasbase + "periodo(" + strconv.Itoa(contrato.Vigencia) + ")." + "semanas_liquidadas(" + contrato.Documento + "," + strconv.Itoa(semanas_liquidadas) + ")."
			auxDetalle = golog.LiquidarMesHCH(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion)
			for j := 0; j < len(auxDetalle); j++ {
				registrarDetallePreliquidacion(auxDetalle[j])
			}
			if mesIterativo == int(contrato.FechaFin.Month()) && anoIterativo == contrato.FechaFin.Year() {
				break
			} else {
				if mesIterativo == 12 {
					mesIterativo = 1
					anoIterativo = anoIterativo + 1
				} else {
					mesIterativo = mesIterativo + 1
				}
			}
		} else {
			fmt.Println("Error al consultar preliquidaciones")
		}
		preliquidacion[0].Id = 0 //Para evitar errores al obtener la preliquidación del siguiente mes
	}

}
