package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// PreliquidacionHcSController operations for PreliquidacionHcS
type PreliquidacionHcSController struct {
	beego.Controller
}

func liquidarHCS(contrato models.Contrato) {
	var mesIterativo int              //mes para iterar en el ciclo para liquidar todos los meses de una vez
	var anoIterativo int              //Ano iterativo a la hora de liquidar
	var predicados []models.Predicado //variable para inyectar reglas
	var preliquidacion []models.Preliquidacion
	var contratoPreliquidacion models.ContratoPreliquidacion
	var detallePreliquidacion models.DetallePreliquidacion
	var aux map[string]interface{}
	var auxDetalle []models.DetallePreliquidacion
	var reglasAlivios string
	var reglasNuevas string //reglas a usar en cada iteracion
	var semanas_liquidadas int
	var diasALiquidar string
	cedula, err := strconv.ParseInt(contrato.Documento, 0, 64)
	var emergencia int //Varibale para evitar loop infinito

	if err == nil {
		reglasAlivios, contratoPreliquidacion = CargarDatosRetefuente(int(cedula))
	}

	mesIterativo = int(contrato.FechaInicio.Month())
	anoIterativo = contrato.Vigencia

	//Obtener las semanas del contrato

	semanasContrato := int(calcularSemanasContratoDVE(contrato.FechaInicio, contrato.FechaFin))
	fmt.Println("SemanasContrato: ", semanasContrato)

	predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
	predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
	reglasbase := cargarReglasBase("HCS") + reglasAlivios + FormatoReglas(predicados)

	for {

		fmt.Println("Mes: ", mesIterativo)
		fmt.Println("Año: ", anoIterativo)
		reglasNuevas = ""
		query := "Ano:" + strconv.Itoa(anoIterativo) + ",Mes:" + strconv.Itoa(mesIterativo) + ",Nominaid:416"
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &preliquidacion)
			//En caso de que no exista la preliqudacion la crea
			if preliquidacion[0].Id == 0 {
				preliquidacion[0] = registrarPreliquidacion(contrato.Vigencia, mesIterativo, 476, 416)
				contratoPreliquidacion = registrarContratoPreliquidacion(preliquidacion[0].Id, contrato.Id, contratoPreliquidacion)
			} else {
				//En caso contrario únicamente crea el contrato_preliquidación y lo asocia directamente
				contratoPreliquidacion = registrarContratoPreliquidacion(preliquidacion[0].Id, contrato.Id, contratoPreliquidacion)
			}

			detallePreliquidacion.ContratoPreliquidacionId = &contratoPreliquidacion
			detallePreliquidacion.TipoPreliquidacionId = 397
			detallePreliquidacion.Activo = true
			detallePreliquidacion.EstadoDisponibilidadId = 426
			_, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
			//Calcular semanas a liquidar
			if mesIterativo == int(contrato.FechaInicio.Month()) && contrato.Vigencia == anoIterativo {
				//para el mes inicial

				//Calcular el numero de días
				diasALiquidar, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
				semanas, _ := strconv.ParseFloat(diasALiquidar, 64)
				semanas = semanas / 7

				if semanas <= 1 {
					semanas_liquidadas = 1
					detallePreliquidacion.DiasLiquidados = 1
					fmt.Println("Semanas: ", semanas)
				} else {
					semanas_liquidadas = int(Roundf(semanas))
					detallePreliquidacion.DiasLiquidados = float64(semanas)
					fmt.Println("Semanas: ", semanas)
				}

			} else if mesIterativo == int(contrato.FechaFin.Month()) && contrato.FechaFin.Year() == anoIterativo {
				//Para el mes final
				//Contar las semanas liquidadas
				var aux map[string]interface{}
				var semanas []models.DetallePreliquidacion
				var mes = int(contrato.FechaInicio.Month())
				var ano = contrato.FechaFin.Year()
				semanas_liquidadas = 0
				for {
					if mes == int(contrato.FechaFin.Month()) && ano == contrato.FechaFin.Year() {
						break
					}
					query := "ContratoPreliquidacionId.PreliquidacionId.Ano:" + strconv.Itoa(ano) + ",ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mes) + ",ContratoPreliquidacionId.ContratoId.NumeroContrato:" + contrato.NumeroContrato + ",ContratoPreliquidacionId.ContratoId.Vigencia:" + strconv.Itoa(contrato.Vigencia)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						fmt.Println()
						LimpiezaRespuestaRefactor(aux, &semanas)
						for i := 0; i < len(semanas); i++ {
							if semanas[i].ConceptoNominaId.Id == 87 {
								semanas_liquidadas = semanas_liquidadas + int(semanas[i].DiasLiquidados)
								fmt.Println("Semanas Liquidadas: ", semanas_liquidadas)
							}
						}
					} else {
						fmt.Println("Error al conseguir las semanas liquidadas: ", err)
					}

					if mes == 12 {
						mes = 1
						ano = ano + 1
					} else {
						mes = mes + 1
					}
				}

				semanas_liquidadas = semanasContrato - semanas_liquidadas
				fmt.Println("Semanas Liquidadas: ", semanas_liquidadas)
				detallePreliquidacion.DiasLiquidados = float64(semanas_liquidadas)
			} else {
				semanas_liquidadas = 4
				detallePreliquidacion.DiasLiquidados = 4
			}
			reglasNuevas = reglasNuevas + reglasbase + "periodo(" + strconv.Itoa(contrato.Vigencia) + ")." + "semanas_liquidadas(" + contrato.Documento + "," + strconv.Itoa(semanas_liquidadas) + ")."

			if mesIterativo == int(contrato.FechaFin.Month()) && anoIterativo == contrato.FechaFin.Year() {
				auxDetalle = golog.LiquidarMesHCS(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion, true)
			} else {
				auxDetalle = golog.LiquidarMesHCS(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion, false)
			}

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
				emergencia = emergencia + 1
			}
			if emergencia == 12 {
				break
			}
		} else {
			fmt.Println("Error al consultar preliquidaciones")
		}
		preliquidacion[0].Id = 0 //Para evitar errores al obtener la preliquidación del siguiente mes
	}
}
