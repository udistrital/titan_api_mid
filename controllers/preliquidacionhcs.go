package controllers

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// PreliquidacionHcSController operations for PreliquidacionHcS
type PreliquidacionHcSController struct {
	beego.Controller
}

func liquidarHCS(ctx context.Context, contrato models.Contrato, general bool, porcentaje float64, vigencia_original int, semanas_totales int, valorDia float64, anulacion bool) (mensaje string, err error) {
	var mesIterativo int              //mes para iterar en el ciclo para liquidar todos los meses de una vez
	var anoIterativo int              //Ano iterativo a la hora de liquidar
	var predicados []models.Predicado //variable para inyectar reglas
	var preliquidacion []models.Preliquidacion
	var contratoPreliquidacion models.ContratoPreliquidacion
	var detallePreliquidacion models.DetallePreliquidacion
	var aux map[string]interface{}
	var unico bool = true
	var auxDetalle []models.DetallePreliquidacion
	var reglasAlivios string
	var reglasNuevas string //reglas a usar en cada iteracion
	var semanas_liquidadas int
	var diasALiquidar string
	var porcentaje_ibc float64
	var contratoDVE models.Contrato

	cedula, err := strconv.ParseInt(contrato.Documento, 0, 64)
	var emergencia int //Varibale para evitar loop infinito

	// Buscar si existen contratos vigentes para el docente
	query := "Documento:" + contrato.Documento + ",TipoNominaId:410,Activo:true"
	var contratosDocente []models.Contrato = nil
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratosDocente)

		for i := 0; i < len(contratosDocente); i++ {
			// Si existe algun contrato (que no sea GENERAL) que termine despues de la fecha de inicio del contrato que se va a crear entonces unico es falso
			if contrato.FechaInicio.Before(contratosDocente[i].FechaFin) && !strings.Contains(contratosDocente[i].NumeroContrato, "GENERAL") &&
				contratosDocente[i].NumeroContrato != contrato.NumeroContrato {
				unico = false
			}
		}
		contrato.Unico = unico

	} else {
		logs.Error("Error al buscar contratos adicionales para el docente: ", err)
	}

	if err == nil {
		reglasAlivios, contratoPreliquidacion, err = CargarDatosRetefuente(ctx, int(cedula))
	}

	if err == nil {
		mesIterativo = int(contrato.FechaInicio.Month())
		anoIterativo = contrato.Vigencia

		//Obtener las semanas del contrato
		var semanasContrato int
		if contrato.NumeroSemanas == 0 {
			semanasContrato = int(calcularSemanasContratoDVE(contrato.FechaInicio, contrato.FechaFin))
		} else {
			semanasContrato = contrato.NumeroSemanas
		}

		//Regla para único o general (para apoximar el ibc al tope mínimo)
		if general || contrato.Unico {
			predicados = append(predicados, models.Predicado{Nombre: "general(1)."})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "general(0)."})
		}

		//Si el contrato es completo se tomarán las vacaciones que calcule
		if contrato.Completo {
			predicados = append(predicados, models.Predicado{Nombre: "completo(1)."})
			predicados = append(predicados, models.Predicado{Nombre: "vacaciones(0)."})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "completo(0)."})
			predicados = append(predicados, models.Predicado{Nombre: "vacaciones(" + fmt.Sprintf("%f", contrato.Vacaciones) + ")."})
		}
		if anulacion {
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%v", valorDia*float64(semanas_totales)) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanas_totales) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
		}
		// se deben agregar los predicados de los porcentajes de prestaciones que se van a usar
		predicadosPrestaciones, porcentajesDesagregadoId := ObtenerReglasPrestaciones(ctx, false)
		predicados = append(predicados, predicadosPrestaciones...)

		if !general && !anulacion {
			// 1.) Se guardan los porcentajes de con los que se calcula el desagregado de la preliquidacion
			// 1.1) Obtenemos el contrato DVE desde titan
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato.Id), &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contratoDVE)

				// Se debe modificar para guardar el id de parametros cuando se liquida una vinculacion
				contratoDVE.PorcentajesDesagregadoId = porcentajesDesagregadoId

				// 1.2) Se guardan los porcentajes
				if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato.Id), contratoDVE, &aux); err != nil {
					logs.Error("Error al actualizar porcentajes en el contrato:", err)
				} else {
					fmt.Printf("PUT realizado a contrato %s vigencia %d con datos: %+v\n", contrato.NumeroContrato, contrato.Vigencia, contratoDVE)
				}
			} else {
				logs.Error("Error al obtener el contrato desde Titan CRUD:", err)
			}
		}

		// predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
		// predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})

		for {

			reglasNuevas = ""
			query := "Ano:" + strconv.Itoa(anoIterativo) + ",Mes:" + strconv.Itoa(mesIterativo) + ",Nominaid:416"
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &preliquidacion)
				//En caso de que no exista la preliqudacion la crea
				if preliquidacion[0].Id == 0 {
					preliquidacion[0] = registrarPreliquidacion(ctx, contrato.Vigencia, mesIterativo, 476, 416)
					contratoPreliquidacion = registrarContratoPreliquidacion(ctx, preliquidacion[0].Id, contrato.Id, contratoPreliquidacion)
				} else {
					//En caso contrario únicamente crea el contrato_preliquidación y lo asocia directamente
					contratoPreliquidacion = registrarContratoPreliquidacion(ctx, preliquidacion[0].Id, contrato.Id, contratoPreliquidacion)
				}

				detallePreliquidacion.ContratoPreliquidacionId = &contratoPreliquidacion
				detallePreliquidacion.TipoPreliquidacionId = 397
				detallePreliquidacion.Activo = true
				detallePreliquidacion.EstadoDisponibilidadId = 426
				_, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)

				//Calcular semanas a liquidar
				if contrato.FechaInicio.Month() == contrato.FechaFin.Month() && contrato.FechaInicio.Year() == contrato.FechaFin.Year() {
					//Contratos de un único mes
					//Calcular el numero de días
					var semanas float64
					if contrato.NumeroSemanas == 0 {
						semanasContrato = int(calcularSemanasContratoDVE(contrato.FechaInicio, contrato.FechaFin))
					} else {
						semanasContrato = contrato.NumeroSemanas
					}
					if anulacion {
						semanas = float64(semanas_totales)
					} else {
						semanas = float64(semanasContrato)
					}
					if porcentaje != 0 {
						porcentaje_ibc = porcentaje
					} else {
						porcentaje_ibc = float64(semanas) / 4
					}

					if semanas <= 1 {
						semanas_liquidadas = 1
						detallePreliquidacion.DiasLiquidados = 1
					} else {
						semanas_liquidadas = int(Roundf(semanas))
						detallePreliquidacion.DiasLiquidados = float64(semanas)
					}

					//semanas_liquidadas = semanasContrato
				} else if mesIterativo == int(contrato.FechaInicio.Month()) && contrato.Vigencia == anoIterativo {
					//para el mes inicial
					//Calcular el numero de días
					diasALiquidar, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
					semanas, _ := strconv.ParseFloat(diasALiquidar, 64)

					if porcentaje != 0 {
						porcentaje_ibc = porcentaje
					} else {
						porcentaje_ibc = semanas / 30
					}
					if semanas == 23 {
						semanas -= 1
					}
					semanas = semanas / 7.5

					if semanas <= 1 {
						semanas_liquidadas = 1
						detallePreliquidacion.DiasLiquidados = 1
					} else {
						semanas_liquidadas = int(math.Ceil(semanas))
						detallePreliquidacion.DiasLiquidados = float64(math.Ceil(semanas))
					}

					semanasContrato = semanasContrato - semanas_liquidadas
				} else {
					semanas_liquidadas = 4
					if semanasContrato-semanas_liquidadas <= 0 {
						diasALiquidar, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
						semanas, _ := strconv.ParseFloat(diasALiquidar, 64)

						if porcentaje != 0 {
							porcentaje_ibc = porcentaje
						} else {
							porcentaje_ibc = semanas / 30
						}
						semanas_liquidadas = semanasContrato
						detallePreliquidacion.DiasLiquidados = float64(semanasContrato)
						semanasContrato = 0
						contrato.FechaFin = time.Date(anoIterativo, time.Month(mesIterativo), 30, 12, 0, 0, 0, time.UTC)
					} else {
						semanasContrato = semanasContrato - semanas_liquidadas
						detallePreliquidacion.DiasLiquidados = 4
						porcentaje_ibc = 1
					}
				}

				reglasbase := cargarReglasBase(ctx, "HCS") + reglasAlivios + FormatoReglas(predicados)
				reglasNuevas = reglasNuevas + reglasbase + "porcentaje(" + fmt.Sprintf("%f", porcentaje_ibc) + ").semanas_liquidadas(" + contrato.Documento + "," + strconv.Itoa(semanas_liquidadas) + ")."
				if (mesIterativo == int(contrato.FechaFin.Month()) && anoIterativo == contrato.FechaFin.Year() && !general) || semanasContrato <= 0 {
					reglasNuevas = reglasNuevas + "mesFinal(1)."
					auxDetalle = golog.LiquidarMesHCS(reglasNuevas, contrato, detallePreliquidacion, true)
				} else {
					reglasNuevas = reglasNuevas + "mesFinal(0)."
					auxDetalle = golog.LiquidarMesHCS(reglasNuevas, contrato, detallePreliquidacion, false)
				}

				for j := 0; j < len(auxDetalle); j++ {
					registrarDetallePreliquidacion(ctx, auxDetalle[j])
				}

				if !general {
					LiquidarContratoGeneral(ctx, mesIterativo, anoIterativo, contrato, preliquidacion[0], porcentaje_ibc, "410", vigencia_original, true)
					if !contrato.Unico {
						ReglaDe3(ctx, contrato, mesIterativo, anoIterativo)
					} else {
					}
				}
				if (mesIterativo == int(contrato.FechaFin.Month()) && anoIterativo == contrato.FechaFin.Year()) || semanasContrato <= 0 {
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
				logs.Error("Error al consultar preliquidaciones")
			}
			preliquidacion[0].Id = 0 //Para evitar errores al obtener la preliquidación del siguiente mes
		}
	} else {
		logs.Error("Error al consultar información en Ágora")
		return "Error al consultar información en Ágora: ", err
	}
	return "", nil
}

func liquidarHCSOld(ctx context.Context, contrato models.ContratoOld, general bool, porcentaje float64, vigencia_original int, semanas_totales int, valorDia float64, anulacion bool) (mensaje string, err error) {
	var mesIterativo int              //mes para iterar en el ciclo para liquidar todos los meses de una vez
	var anoIterativo int              //Ano iterativo a la hora de liquidar
	var predicados []models.Predicado //variable para inyectar reglas
	var preliquidacion []models.Preliquidacion
	var contratoPreliquidacion models.ContratoPreliquidacionOld
	var detallePreliquidacion models.DetallePreliquidacionOld
	var aux map[string]interface{}
	var unico bool = true
	var auxDetalle []models.DetallePreliquidacionOld
	var reglasAlivios string
	var reglasNuevas string //reglas a usar en cada iteracion
	var semanas_liquidadas int
	var diasALiquidar string
	var porcentaje_ibc float64

	cedula, err := strconv.ParseInt(contrato.Documento, 0, 64)
	var emergencia int //Varibale para evitar loop infinito

	// Buscar si existen contratos vigentes para el docente
	query := "Documento:" + contrato.Documento + ",TipoNominaId:410" + ",Activo:true"
	var contratosDocente []models.Contrato = nil
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratosDocente)

		for i := 0; i < len(contratosDocente); i++ {
			// Si existe algun contrato (que no sea GENERAL) que termine despues de la fecha de inicio del contrato que se va a crear entonces unico es falso
			if contrato.FechaInicio.Before(contratosDocente[i].FechaFin) && !strings.Contains(contratosDocente[i].NumeroContrato, "GENERAL") &&
				contratosDocente[i].NumeroContrato != contrato.NumeroContrato {
				unico = false
			}
		}
		contrato.Unico = unico

	} else {
		logs.Error("Error al buscar contratos adicionales para el docente: ", err)
	}

	if err == nil {
		reglasAlivios, contratoPreliquidacion, err = CargarDatosRetefuenteOld(ctx, int(cedula))
	}

	if err == nil {
		mesIterativo = int(contrato.FechaInicio.Month())
		anoIterativo = contrato.Vigencia

		//Obtener las semanas del contrato

		semanasContrato := int(calcularSemanasContratoDVE(contrato.FechaInicio, contrato.FechaFin))

		//Regla para único o general (para apoximar el ibc al tope mínimo)
		if general || contrato.Unico {
			predicados = append(predicados, models.Predicado{Nombre: "general(1)."})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "general(0)."})
		}

		//Si el contrato es completo se tomarán las vacaciones que calcule
		if contrato.Completo {
			predicados = append(predicados, models.Predicado{Nombre: "completo(1)."})
			predicados = append(predicados, models.Predicado{Nombre: "vacaciones(0)."})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "completo(0)."})
			predicados = append(predicados, models.Predicado{Nombre: "vacaciones(" + fmt.Sprintf("%f", contrato.Vacaciones) + ")."})
		}
		if anulacion {
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%v", valorDia*float64(semanas_totales)) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanas_totales) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
		}

		for {

			reglasNuevas = ""
			query := "Ano:" + strconv.Itoa(anoIterativo) + ",Mes:" + strconv.Itoa(mesIterativo) + ",Nominaid:416"
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &preliquidacion)
				//En caso de que no exista la preliqudacion la crea
				if preliquidacion[0].Id == 0 {
					preliquidacion[0] = registrarPreliquidacion(ctx, contrato.Vigencia, mesIterativo, 476, 416)
					contratoPreliquidacion = registrarContratoPreliquidacionOld(ctx, preliquidacion[0].Id, contrato.Id, contratoPreliquidacion)
				} else {
					//En caso contrario únicamente crea el contrato_preliquidación y lo asocia directamente
					contratoPreliquidacion = registrarContratoPreliquidacionOld(ctx, preliquidacion[0].Id, contrato.Id, contratoPreliquidacion)
				}

				detallePreliquidacion.ContratoPreliquidacionId = &contratoPreliquidacion
				detallePreliquidacion.TipoPreliquidacionId = 397
				detallePreliquidacion.Activo = true
				detallePreliquidacion.EstadoDisponibilidadId = 426
				_, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)

				//Calcular semanas a liquidar
				if contrato.FechaInicio.Month() == contrato.FechaFin.Month() && contrato.FechaInicio.Year() == contrato.FechaFin.Year() {
					//Contratos de un único mes
					//Calcular el numero de días
					diasALiquidar, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
					// semanas, _ := strconv.ParseFloat(diasALiquidar, 64)
					if porcentaje != 0 {
						porcentaje_ibc = porcentaje
					} else {
						porcentaje_ibc = float64(semanasContrato) / 4
					}

					semanas_liquidadas = semanasContrato
				} else if mesIterativo == int(contrato.FechaInicio.Month()) && contrato.Vigencia == anoIterativo {
					//para el mes inicial
					//Calcular el numero de días
					diasALiquidar, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
					semanas, _ := strconv.ParseFloat(diasALiquidar, 64)

					if porcentaje != 0 {
						porcentaje_ibc = porcentaje
					} else {
						porcentaje_ibc = semanas / 30
					}
					semanas = semanas / 7

					if semanas <= 1 {
						semanas_liquidadas = 1
						detallePreliquidacion.DiasLiquidados = 1
					} else {
						semanas_liquidadas = int(Roundf(semanas))
						detallePreliquidacion.DiasLiquidados = float64(semanas)
					}

					semanasContrato = semanasContrato - semanas_liquidadas
				} else {
					semanas_liquidadas = 4
					if semanasContrato-semanas_liquidadas <= 0 {
						diasALiquidar, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
						semanas, _ := strconv.ParseFloat(diasALiquidar, 64)

						if porcentaje != 0 {
							porcentaje_ibc = porcentaje
						} else {
							porcentaje_ibc = semanas / 30
						}
						semanas_liquidadas = semanasContrato
						detallePreliquidacion.DiasLiquidados = float64(semanasContrato)
						semanasContrato = 0
						contrato.FechaFin = time.Date(anoIterativo, time.Month(mesIterativo), 30, 12, 0, 0, 0, time.UTC)
					} else {
						semanasContrato = semanasContrato - semanas_liquidadas
						detallePreliquidacion.DiasLiquidados = 4
						porcentaje_ibc = 1
					}
				}

				reglasbase := cargarReglasBase(ctx, "HCS") + reglasAlivios + FormatoReglas(predicados)

				reglasNuevas = reglasNuevas + reglasbase + "porcentaje(" + fmt.Sprintf("%f", porcentaje_ibc) + ").semanas_liquidadas(" + contrato.Documento + "," + strconv.Itoa(semanas_liquidadas) + ")."

				if mesIterativo == int(contrato.FechaFin.Month()) && anoIterativo == contrato.FechaFin.Year() && !general {
					reglasNuevas = reglasNuevas + "mesFinal(1)."
					auxDetalle = golog.LiquidarMesHCSOld(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion, true)
				} else {
					reglasNuevas = reglasNuevas + "mesFinal(0)."
					auxDetalle = golog.LiquidarMesHCSOld(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion, false)
				}

				for j := 0; j < len(auxDetalle); j++ {
					registrarDetallePreliquidacionOld(ctx, auxDetalle[j])
				}

				if !general {
					LiquidarContratoGeneralOld(ctx, mesIterativo, anoIterativo, contrato, preliquidacion[0], porcentaje_ibc, "410", vigencia_original, true)
					if !contrato.Unico {
						ReglaDe3Old(ctx, contrato, mesIterativo, anoIterativo)
					} else {
					}
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
				logs.Error("Error al consultar preliquidaciones")
			}
			preliquidacion[0].Id = 0 //Para evitar errores al obtener la preliquidación del siguiente mes
		}
	} else {
		logs.Error("Error al consultar información en Ágora")
		return "Error al consultar información en Ágora: ", err
	}
	return "", nil
}

func ReglaDe3(ctx context.Context, contrato models.Contrato, mesIterativo int, anoIterativo int) {
	var aux map[string]interface{}
	var auxDetalle []models.DetallePreliquidacion
	var contratoGeneral []models.Contrato = nil
	var contratosDocente []models.Contrato = nil
	var contratoPreliquidacionDocente []models.ContratoPreliquidacion = nil
	var auxValor []models.DetallePreliquidacion
	var ibcGeneral float64
	var salarioGeneral float64
	var contratosCambio []int
	var cambioNecesario bool = false
	//Obtener los valores del ibc liquidado para saber si es necesario realizar actualizacion
	query := "Documento:" + contrato.Documento + ",TipoNominaId:410,Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratosDocente)
		if contratosDocente[0].Id != 0 {
			for i := 0; i < len(contratosDocente); i++ {
				query = "ContratoId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
				if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoPreliquidacionDocente)
					if contratoPreliquidacionDocente[0].Id != 0 {
						if contratosDocente[i].NumeroContrato != "GENERAL"+strconv.Itoa(mesIterativo) {
							contratosCambio = append(contratosCambio, contratoPreliquidacionDocente[0].Id)
						} else {
							if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId.Id:"+strconv.Itoa(contratoPreliquidacionDocente[0].Id)+",ConceptoNominaId.Id:521", &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &auxValor)
								if auxValor[0].Id != 0 {
									ibcGeneral = auxValor[0].ValorCalculado
								} else {
								}
							} else {
								logs.Error("Error al obtener el valor del ibc para el contrato: ", contratosDocente[i].NumeroContrato)
							}
							if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId.Id:"+strconv.Itoa(contratoPreliquidacionDocente[0].Id)+",ConceptoNominaId.Id:152", &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &auxValor)
								if auxValor[0].Id != 0 {
									salarioGeneral = auxValor[0].ValorCalculado
								} else {
								}
							} else {
								logs.Error("Error al obtener el valor del ibc para el contrato: ", contratosDocente[i].NumeroContrato)
							}
							if salarioGeneral < ibcGeneral && len(contratosDocente) > 2 {
								cambioNecesario = true
								break
							}
						}
					} else {
					}
				} else {
					logs.Error("Error al obtener el contrato preliquidación para el contrato: ", contratosDocente[i].NumeroContrato)
				}
			}
			//por defecto que se realice la regla de 3
			cambioNecesario = true
			//Hacer regla de 3 en caso de que el cambio sea necesario
			if cambioNecesario {
				//obtener el contrato general
				query = "Documento:" + contrato.Documento + ",TipoNominaId:410,NumeroContrato:GENERAL" + strconv.Itoa(mesIterativo) + ",Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
				if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoGeneral)
					if contratoGeneral[0].Id != 0 {
						//Obtener el contrato preliquidacion del contrato general
						if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query=ContratoId:"+strconv.Itoa(contratoGeneral[0].Id), &aux); err == nil {
							var auxCp []models.ContratoPreliquidacion //Variable auxiliar de contrato preliquidacion
							LimpiezaRespuestaRefactor(aux, &auxCp)
							if auxCp[0].Id != 0 {
								//traer los detalles necesarios para hacer la reglas de tres
								auxDetalle = nil
								if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(auxCp[0].Id), &aux); err == nil {
									LimpiezaRespuestaRefactor(aux, &auxDetalle)
									if auxDetalle[0].Id != 0 {
										var totalHonorarios float64 = 0
										var valorIbc float64 = 0
										var valorSalud float64 = 0
										var valorPension float64 = 0
										var valorArl float64 = 0
										var valorRetefuente float64 = 0
										var valorFondoSol float64 = 0
										var valorFondoSub float64 = 0
										var valorSaludUniversidad float64 = 0
										var valorPensionUniversidad float64 = 0
										var valorMensual float64 = 0
										//obtener los valores totales para realizar la regla de 3
										for i := 0; i < len(auxDetalle); i++ {
											switch auxDetalle[i].ConceptoNominaId.Id {
											case 152:
												totalHonorarios = auxDetalle[i].ValorCalculado
											case 64:
												valorRetefuente = auxDetalle[i].ValorCalculado
											case 170:
												valorFondoSol = auxDetalle[i].ValorCalculado
											case 572:
												valorFondoSub = auxDetalle[i].ValorCalculado
											case 568:
												valorSalud = auxDetalle[i].ValorCalculado
											case 569:
												valorPension = auxDetalle[i].ValorCalculado
											case 570:
												valorArl = auxDetalle[i].ValorCalculado
											case 521:
												valorIbc = auxDetalle[i].ValorCalculado
											case 576:
												valorSaludUniversidad = auxDetalle[i].ValorCalculado
											case 577:
												valorPensionUniversidad = auxDetalle[i].ValorCalculado
											}
										}
										//Obtener los detalles que necesitan cambio
										auxDetalle = nil
										var detalleEnvio models.DetallePreliquidacion
										for i := 0; i < len(contratosCambio); i++ {
											if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contratosCambio[i]), &aux); err == nil {
												LimpiezaRespuestaRefactor(aux, &auxDetalle)
												if auxDetalle[0].Id != 0 {

													for j := 0; j < len(auxDetalle); j++ {
														if auxDetalle[j].ConceptoNominaId.Id == 152 {
															valorMensual = auxDetalle[j].ValorCalculado
														}
													}

													for j := 0; j < len(auxDetalle); j++ {

														switch auxDetalle[j].ConceptoNominaId.Id {
														case 64:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorRetefuente)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {
															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 170:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorFondoSol)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 572:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorFondoSub)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 568:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorSalud)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 569:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorPension)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 570:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorArl)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 521:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorIbc)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 576:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorSaludUniversidad)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 577:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorPensionUniversidad)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {
															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														}
													}
												} else {
												}
											} else {
												logs.Error("Error al obtener detalles para cambio")
											}
										}

									} else {
									}
								} else {
									logs.Error("Error al traer los detalles del contrato general: ", err)
								}
							} else {
							}
						} else {
							logs.Error("Error al obtener el contrato preliquidacion: ", err)
						}
					} else {
					}
				} else {
					logs.Error("Error al obtener el contrato general: ", err)
				}
			} else {
			}
		} else {
		}
	} else {
		logs.Error("Error al intentar obtener contratos del docente: ", err)
	}
}

func ReglaDe3Old(ctx context.Context, contrato models.ContratoOld, mesIterativo int, anoIterativo int) {
	var aux map[string]interface{}
	var auxDetalle []models.DetallePreliquidacionOld
	var contratoGeneral []models.ContratoOld = nil
	var contratosDocente []models.ContratoOld = nil
	var contratoPreliquidacionDocente []models.ContratoPreliquidacionOld = nil
	var auxValor []models.DetallePreliquidacionOld
	var ibcGeneral float64
	var salarioGeneral float64
	var contratosCambio []int
	var cambioNecesario bool = false
	//Obtener los valores del ibc liquidado para saber si es necesario realizar actualizacion
	query := "Documento:" + contrato.Documento + ",TipoNominaId:410,Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratosDocente)
		if contratosDocente[0].Id != 0 {
			for i := 0; i < len(contratosDocente); i++ {
				query = "ContratoId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
				if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoPreliquidacionDocente)
					if contratoPreliquidacionDocente[0].Id != 0 {
						if contratosDocente[i].NumeroContrato != "GENERAL"+strconv.Itoa(mesIterativo) {
							contratosCambio = append(contratosCambio, contratoPreliquidacionDocente[0].Id)
						} else {
							if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId.Id:"+strconv.Itoa(contratoPreliquidacionDocente[0].Id)+",ConceptoNominaId.Id:521", &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &auxValor)
								if auxValor[0].Id != 0 {
									ibcGeneral = auxValor[0].ValorCalculado
								} else {
								}
							} else {
								logs.Error("Error al obtener el valor del ibc para el contrato: ", contratosDocente[i].NumeroContrato)
							}
							if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId.Id:"+strconv.Itoa(contratoPreliquidacionDocente[0].Id)+",ConceptoNominaId.Id:152", &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &auxValor)
								if auxValor[0].Id != 0 {
									salarioGeneral = auxValor[0].ValorCalculado
								} else {
								}
							} else {
								logs.Error("Error al obtener el valor del ibc para el contrato: ", contratosDocente[i].NumeroContrato)
							}
							if salarioGeneral < ibcGeneral && len(contratosDocente) > 2 {
								cambioNecesario = true
								break
							}
						}
					} else {
					}
				} else {
					logs.Error("Error al obtener el contrato preliquidación para el contrato: ", contratosDocente[i].NumeroContrato)
				}
			}
			//por defecto que se realice la regla de 3
			cambioNecesario = true
			//Hacer regla de 3 en caso de que el cambio sea necesario
			if cambioNecesario {
				//obtener el contrato general
				query = "Documento:" + contrato.Documento + ",TipoNominaId:410,NumeroContrato:GENERAL" + strconv.Itoa(mesIterativo) + ",Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
				if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoGeneral)
					if contratoGeneral[0].Id != 0 {
						//Obtener el contrato preliquidacion del contrato general
						if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query=ContratoId:"+strconv.Itoa(contratoGeneral[0].Id), &aux); err == nil {
							var auxCp []models.ContratoPreliquidacionOld //Variable auxiliar de contrato preliquidacion
							LimpiezaRespuestaRefactor(aux, &auxCp)
							if auxCp[0].Id != 0 {
								//traer los detalles necesarios para hacer la reglas de tres
								auxDetalle = nil
								if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(auxCp[0].Id), &aux); err == nil {
									LimpiezaRespuestaRefactor(aux, &auxDetalle)
									if auxDetalle[0].Id != 0 {
										var totalHonorarios float64 = 0
										var valorIbc float64 = 0
										var valorSalud float64 = 0
										var valorPension float64 = 0
										var valorArl float64 = 0
										var valorRetefuente float64 = 0
										var valorFondoSol float64 = 0
										var valorFondoSub float64 = 0
										var valorSaludUniversidad float64 = 0
										var valorPensionUniversidad float64 = 0
										var valorMensual float64 = 0
										//obtener los valores totales para realizar la regla de 3
										for i := 0; i < len(auxDetalle); i++ {
											switch auxDetalle[i].ConceptoNominaId.Id {
											case 152:
												totalHonorarios = auxDetalle[i].ValorCalculado
											case 64:
												valorRetefuente = auxDetalle[i].ValorCalculado
											case 170:
												valorFondoSol = auxDetalle[i].ValorCalculado
											case 572:
												valorFondoSub = auxDetalle[i].ValorCalculado
											case 568:
												valorSalud = auxDetalle[i].ValorCalculado
											case 569:
												valorPension = auxDetalle[i].ValorCalculado
											case 570:
												valorArl = auxDetalle[i].ValorCalculado
											case 521:
												valorIbc = auxDetalle[i].ValorCalculado
											case 576:
												valorSaludUniversidad = auxDetalle[i].ValorCalculado
											case 577:
												valorPensionUniversidad = auxDetalle[i].ValorCalculado
											}
										}
										//Obtener los detalles que necesitan cambio
										auxDetalle = nil
										var detalleEnvio models.DetallePreliquidacionOld
										for i := 0; i < len(contratosCambio); i++ {
											if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contratosCambio[i]), &aux); err == nil {
												LimpiezaRespuestaRefactor(aux, &auxDetalle)
												if auxDetalle[0].Id != 0 {

													for j := 0; j < len(auxDetalle); j++ {
														if auxDetalle[j].ConceptoNominaId.Id == 152 {
															valorMensual = auxDetalle[j].ValorCalculado
														}
													}

													for j := 0; j < len(auxDetalle); j++ {

														switch auxDetalle[j].ConceptoNominaId.Id {
														case 64:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorRetefuente)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {
															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 170:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorFondoSol)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 572:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorFondoSub)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 568:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorSalud)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 569:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorPension)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 570:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorArl)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 521:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorIbc)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 576:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorSaludUniversidad)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {

															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 577:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorPensionUniversidad)
															if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), detalleEnvio, &aux); err == nil {
															} else {
																logs.Error("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														}
													}
												} else {
												}
											} else {
												logs.Error("Error al obtener detalles para cambio")
											}
										}

									} else {
									}
								} else {
									logs.Error("Error al traer los detalles del contrato general: ", err)
								}
							} else {
							}
						} else {
							logs.Error("Error al obtener el contrato preliquidacion: ", err)
						}
					} else {
					}
				} else {
					logs.Error("Error al obtener el contrato general: ", err)
				}
			} else {
			}
		} else {
		}
	} else {
		logs.Error("Error al intentar obtener contratos del docente: ", err)
	}
}
