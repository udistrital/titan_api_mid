package controllers

import (
	"fmt"
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

func liquidarHCS(contrato models.Contrato, general bool, porcentaje float64, vigencia_original int, semanas_totales int, valorDia float64, anulacion bool) (mensaje string, err error) {
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
	var porcentajesDesagregado models.PorcentajeDesagregado

	cedula, err := strconv.ParseInt(contrato.Documento, 0, 64)
	var emergencia int //Varibale para evitar loop infinito
	// Validamos que solo se guarden los porcentajes en los DVE, no en los generales
	if general == false {
		// 1.) Se guardan los porcentajes de con los que se calcula el desagregado de la preliquidacion
		// 1.1) Obtenemos el contrato DVE desde titan
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"contrato/"+strconv.Itoa(contrato.Id), &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contratoDVE)
		} else {
			fmt.Println("Error al obtener el contrato desde Titan CRUD:", err)
		}

		porcentajesDesagregado = golog.ObtenerPorcentajesDesagregado(cargarReglasBase("HCS"), contrato)
		contratoDVE.PorcentajeCesantias = porcentajesDesagregado.PorcentajeCesantias
		contratoDVE.PorcentajePrimaNavidad = porcentajesDesagregado.PorcentajePrimaNavidad
		contratoDVE.PorcentajePrimaVacaciones = porcentajesDesagregado.PorcentajePrimaVacaciones
		contratoDVE.PorcentajeVacaciones = porcentajesDesagregado.PorcentajeVacaciones
		contratoDVE.PorcentajePrimaServicios = porcentajesDesagregado.PorcentajePrimaServicios

		// 1.2) Se guardan los porcentajes
		if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"contrato/"+strconv.Itoa(contrato.Id), "PUT", &aux, contratoDVE); err != nil {
			fmt.Println("Error al actualizar porcentajes en el contrato:", err)
		}
	}

	// Buscar si existen contratos vigentes para el docente
	query := "Documento:" + contrato.Documento + ",TipoNominaId:410,Activo:true"
	var contratosDocente []models.Contrato = nil
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
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
		fmt.Println("Error al buscar contratos adicionales para el docente: ", err)
	}

	if err == nil {
		reglasAlivios, contratoPreliquidacion, err = CargarDatosRetefuente(int(cedula))
	}

	if err == nil {
		mesIterativo = int(contrato.FechaInicio.Month())
		anoIterativo = contrato.Vigencia

		//Obtener las semanas del contrato
		fmt.Println("CONTRato ", contrato.NumeroSemanas)
		var semanasContrato int
		if contrato.NumeroSemanas == 0 {
			semanasContrato = int(calcularSemanasContratoDVE(contrato.FechaInicio, contrato.FechaFin))
		} else {
			semanasContrato = contrato.NumeroSemanas
		}
		fmt.Println("SemanasContrato: ", semanasContrato)

		//Regla para único o general (para apoximar el ibc al tope mínimo)
		if general || contrato.Unico {
			predicados = append(predicados, models.Predicado{Nombre: "general(1)."})
			fmt.Println("El contrato es general o único, se carga regla")
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "general(0)."})
			fmt.Println("El docente tiene varios contratos, no se carga regla de único")
		}

		//Si el contrato es completo se tomarán las vacaciones que calcule
		if contrato.Completo {
			predicados = append(predicados, models.Predicado{Nombre: "completo(1)."})
			predicados = append(predicados, models.Predicado{Nombre: "vacaciones(0)."})
			fmt.Println("El contrato es completo, no requiere vacaciones")
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "completo(0)."})
			predicados = append(predicados, models.Predicado{Nombre: "vacaciones(" + fmt.Sprintf("%f", contrato.Vacaciones) + ")."})
			fmt.Println("El contrato no es completo, requiere de las vacaciones")
		}
		if anulacion {
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%v", valorDia*float64(semanas_totales)) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanas_totales) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
		}

		// predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
		// predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})

		for {

			fmt.Println("Mes: ", mesIterativo)
			fmt.Println("Año: ", anoIterativo)
			fmt.Println("+++ ", contrato.FechaFin.Month())
			fmt.Println("+++ ", contrato.FechaFin.Year())
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
				if contrato.FechaInicio.Month() == contrato.FechaFin.Month() && contrato.FechaInicio.Year() == contrato.FechaFin.Year() {
					fmt.Println("CONTRATOS UNICO MES ")
					fmt.Println(contrato.FechaInicio)
					fmt.Println(contrato.FechaFin)
					fmt.Println(semanas_totales)
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
					fmt.Println("CONTRATOS MES INICIAL ")
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
					fmt.Println("CONTRATOS VARIOS MESES ")
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
				predicados = append(predicados, models.Predicado{Nombre: "cancelacion(0)."})
				reglasbase := cargarReglasBase("HCS") + reglasAlivios + FormatoReglas(predicados)
				reglasNuevas = reglasNuevas + reglasbase + "porcentaje(" + fmt.Sprintf("%f", porcentaje_ibc) + ").semanas_liquidadas(" + contrato.Documento + "," + strconv.Itoa(semanas_liquidadas) + ")."
				if (mesIterativo == int(contrato.FechaFin.Month()) && anoIterativo == contrato.FechaFin.Year() && !general) || semanasContrato <= 0 {
					reglasNuevas = reglasNuevas + "mesFinal(1)."
					auxDetalle = golog.LiquidarMesHCS(reglasNuevas, contrato, detallePreliquidacion, true)
				} else {
					reglasNuevas = reglasNuevas + "mesFinal(0)."
					auxDetalle = golog.LiquidarMesHCS(reglasNuevas, contrato, detallePreliquidacion, false)
				}

				for j := 0; j < len(auxDetalle); j++ {
					registrarDetallePreliquidacion(auxDetalle[j])
				}

				if !general {
					fmt.Println("Liquidando Contrato General")
					LiquidarContratoGeneral(mesIterativo, anoIterativo, contrato, preliquidacion[0], porcentaje_ibc, "410", vigencia_original, true)
					if !contrato.Unico {
						fmt.Println("Realizando Regla de 3 con los conceptos de ibc")
						ReglaDe3(contrato, mesIterativo, anoIterativo)
					} else {
						fmt.Println("El contrato es único, no requiere de actualización")
					}
				}
				fmt.Println("Mes: ", mesIterativo)
				fmt.Println("Año: ", anoIterativo)
				fmt.Println("--- ", contrato.FechaFin.Month())
				fmt.Println("--- ", contrato.FechaFin.Year())
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
				fmt.Println("Error al consultar preliquidaciones")
			}
			preliquidacion[0].Id = 0 //Para evitar errores al obtener la preliquidación del siguiente mes
		}
	} else {
		fmt.Println("Error al consultar información en Ágora")
		return "Error al consultar información en Ágora: ", err
	}
	return "", nil
}

func liquidarHCSOld(contrato models.ContratoOld, general bool, porcentaje float64, vigencia_original int, semanas_totales int, valorDia float64, anulacion bool) (mensaje string, err error) {
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
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
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
		fmt.Println("Error al buscar contratos adicionales para el docente: ", err)
	}

	if err == nil {
		reglasAlivios, contratoPreliquidacion, err = CargarDatosRetefuenteOld(int(cedula))
	}

	if err == nil {
		mesIterativo = int(contrato.FechaInicio.Month())
		anoIterativo = contrato.Vigencia

		//Obtener las semanas del contrato

		semanasContrato := int(calcularSemanasContratoDVE(contrato.FechaInicio, contrato.FechaFin))
		fmt.Println("SemanasContrato: ", semanasContrato)

		//Regla para único o general (para apoximar el ibc al tope mínimo)
		if general || contrato.Unico {
			predicados = append(predicados, models.Predicado{Nombre: "general(1)."})
			fmt.Println("El contrato es general o único, se carga regla")
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "general(0)."})
			fmt.Println("El docente tiene varios contratos, no se carga regla de único")
		}

		//Si el contrato es completo se tomarán las vacaciones que calcule
		if contrato.Completo {
			predicados = append(predicados, models.Predicado{Nombre: "completo(1)."})
			predicados = append(predicados, models.Predicado{Nombre: "vacaciones(0)."})
			fmt.Println("El contrato es completo, no requiere vacaciones")
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "completo(0)."})
			predicados = append(predicados, models.Predicado{Nombre: "vacaciones(" + fmt.Sprintf("%f", contrato.Vacaciones) + ")."})
			fmt.Println("El contrato no es completo, requiere de las vacaciones")
		}
		if anulacion {
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%v", valorDia*float64(semanas_totales)) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanas_totales) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
		}

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
					contratoPreliquidacion = registrarContratoPreliquidacionOld(preliquidacion[0].Id, contrato.Id, contratoPreliquidacion)
				} else {
					//En caso contrario únicamente crea el contrato_preliquidación y lo asocia directamente
					contratoPreliquidacion = registrarContratoPreliquidacionOld(preliquidacion[0].Id, contrato.Id, contratoPreliquidacion)
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

				reglasbase := cargarReglasBase("HCS") + reglasAlivios + FormatoReglas(predicados)

				reglasNuevas = reglasNuevas + reglasbase + "porcentaje(" + fmt.Sprintf("%f", porcentaje_ibc) + ").semanas_liquidadas(" + contrato.Documento + "," + strconv.Itoa(semanas_liquidadas) + ")."

				if mesIterativo == int(contrato.FechaFin.Month()) && anoIterativo == contrato.FechaFin.Year() && !general {
					reglasNuevas = reglasNuevas + "mesFinal(1)."
					auxDetalle = golog.LiquidarMesHCSOld(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion, true)
				} else {
					reglasNuevas = reglasNuevas + "mesFinal(0)."
					auxDetalle = golog.LiquidarMesHCSOld(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion, false)
				}

				for j := 0; j < len(auxDetalle); j++ {
					registrarDetallePreliquidacionOld(auxDetalle[j])
				}

				if !general {
					fmt.Println("Liquidando Contrato General")
					LiquidarContratoGeneralOld(mesIterativo, anoIterativo, contrato, preliquidacion[0], porcentaje_ibc, "410", vigencia_original, true)
					if !contrato.Unico {
						fmt.Println("Realizando Regla de 3 con los conceptos de ibc")
						ReglaDe3Old(contrato, mesIterativo, anoIterativo)
					} else {
						fmt.Println("El contrato es único, no requiere de actualización")
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
				fmt.Println("Error al consultar preliquidaciones")
			}
			preliquidacion[0].Id = 0 //Para evitar errores al obtener la preliquidación del siguiente mes
		}
	} else {
		fmt.Println("Error al consultar información en Ágora")
		return "Error al consultar información en Ágora: ", err
	}
	return "", nil
}

func ReglaDe3(contrato models.Contrato, mesIterativo int, anoIterativo int) {
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
	fmt.Println("Ingreso a regla de 3")
	//Obtener los valores del ibc liquidado para saber si es necesario realizar actualizacion
	query := "Documento:" + contrato.Documento + ",TipoNominaId:410,Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratosDocente)
		if contratosDocente[0].Id != 0 {
			fmt.Println("Tamaño arreglo contratos: ", len(contratosDocente))
			for i := 0; i < len(contratosDocente); i++ {
				fmt.Println("iteracion: ", i)
				fmt.Println(contratosDocente[i].NumeroContrato)
				query = "ContratoId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoPreliquidacionDocente)
					if contratoPreliquidacionDocente[0].Id != 0 {
						if contratosDocente[i].NumeroContrato != "GENERAL"+strconv.Itoa(mesIterativo) {
							fmt.Println("Agrego el contrato: ", contratosDocente[i].NumeroContrato)
							contratosCambio = append(contratosCambio, contratoPreliquidacionDocente[0].Id)
						} else {
							if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId.Id:"+strconv.Itoa(contratoPreliquidacionDocente[0].Id)+",ConceptoNominaId.Id:521", &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &auxValor)
								if auxValor[0].Id != 0 {
									ibcGeneral = auxValor[0].ValorCalculado
								} else {
									fmt.Println("No se encontró ibc para el contrato: ", contratosDocente[i].NumeroContrato)
								}
							} else {
								fmt.Println("Error al obtener el valor del ibc para el contrato: ", contratosDocente[i].NumeroContrato)
							}
							if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId.Id:"+strconv.Itoa(contratoPreliquidacionDocente[0].Id)+",ConceptoNominaId.Id:152", &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &auxValor)
								if auxValor[0].Id != 0 {
									salarioGeneral = auxValor[0].ValorCalculado
								} else {
									fmt.Println("No se encontraron salarion para el contrato: ", contratosDocente[i].NumeroContrato)
								}
							} else {
								fmt.Println("Error al obtener el valor del ibc para el contrato: ", contratosDocente[i].NumeroContrato)
							}
							// fmt.Println("salarioGeneral: ", salarioGeneral)
							// fmt.Println("ibcGeneral: ", ibcGeneral)
							// fmt.Println("contratosDocente: ", contratosDocente)
							if salarioGeneral < ibcGeneral && len(contratosDocente) > 2 {
								cambioNecesario = true
								break
							}
						}
					} else {
						fmt.Println("No se encontraron preliquidaciones asociadas al contrato: ", contratosDocente[i].NumeroContrato)
					}
				} else {
					fmt.Println("Error al obtener el contrato preliquidación para el contrato: ", contratosDocente[i].NumeroContrato)
				}
			}
			//por defecto que se realice la regla de 3
			cambioNecesario = true
			//Hacer regla de 3 en caso de que el cambio sea necesario
			if cambioNecesario {
				//obtener el contrato general
				query = "Documento:" + contrato.Documento + ",TipoNominaId:410,NumeroContrato:GENERAL" + strconv.Itoa(mesIterativo) + ",Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoGeneral)
					if contratoGeneral[0].Id != 0 {
						//Obtener el contrato preliquidacion del contrato general
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query=ContratoId:"+strconv.Itoa(contratoGeneral[0].Id), &aux); err == nil {
							var auxCp []models.ContratoPreliquidacion //Variable auxiliar de contrato preliquidacion
							LimpiezaRespuestaRefactor(aux, &auxCp)
							if auxCp[0].Id != 0 {
								//traer los detalles necesarios para hacer la reglas de tres
								auxDetalle = nil
								if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(auxCp[0].Id), &aux); err == nil {
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
												fmt.Println("Total honorarios:", totalHonorarios)
												fmt.Println("------------------------------------------------------------")
											case 64:
												valorRetefuente = auxDetalle[i].ValorCalculado
												fmt.Println("Total retefuente:", valorRetefuente)
												fmt.Println("------------------------------------------------------------")
											case 170:
												valorFondoSol = auxDetalle[i].ValorCalculado
												fmt.Println("Total fondo sol:", valorFondoSol)
												fmt.Println("------------------------------------------------------------")
											case 572:
												valorFondoSub = auxDetalle[i].ValorCalculado
												fmt.Println("Total fondo sub:", valorFondoSub)
												fmt.Println("------------------------------------------------------------")
											case 568:
												valorSalud = auxDetalle[i].ValorCalculado
												fmt.Println("Total Salud:", valorSalud)
												fmt.Println("------------------------------------------------------------")
											case 569:
												valorPension = auxDetalle[i].ValorCalculado
												fmt.Println("Total Pension:", valorPension)
												fmt.Println("------------------------------------------------------------")
											case 570:
												valorArl = auxDetalle[i].ValorCalculado
												fmt.Println("Total Arl:", valorArl)
												fmt.Println("------------------------------------------------------------")
											case 521:
												valorIbc = auxDetalle[i].ValorCalculado
												fmt.Println("Total ibc:", valorIbc)
												fmt.Println("------------------------------------------------------------")
											case 576:
												valorSaludUniversidad = auxDetalle[i].ValorCalculado
												fmt.Println("Total salud Universidad:", valorSaludUniversidad)
												fmt.Println("------------------------------------------------------------")
											case 577:
												valorPensionUniversidad = auxDetalle[i].ValorCalculado
												fmt.Println("Total Pensión universidad:", valorPensionUniversidad)
												fmt.Println("------------------------------------------------------------")
											}
										}
										//Obtener los detalles que necesitan cambio
										auxDetalle = nil
										var detalleEnvio models.DetallePreliquidacion
										for i := 0; i < len(contratosCambio); i++ {
											if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contratosCambio[i]), &aux); err == nil {
												LimpiezaRespuestaRefactor(aux, &auxDetalle)
												if auxDetalle[0].Id != 0 {

													for j := 0; j < len(auxDetalle); j++ {
														if auxDetalle[j].ConceptoNominaId.Id == 152 {
															valorMensual = auxDetalle[j].ValorCalculado
															fmt.Println("Honorarios para el contrato: ", valorMensual)
														}
													}

													for j := 0; j < len(auxDetalle); j++ {

														switch auxDetalle[j].ConceptoNominaId.Id {
														case 64:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorRetefuente)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																//fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)
															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 170:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorFondoSol)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																//fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 572:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorFondoSub)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																//fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 568:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorSalud)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																//fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 569:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorPension)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																//fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 570:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorArl)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																//fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 521:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorIbc)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																//fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 576:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorSaludUniversidad)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																//fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 577:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorPensionUniversidad)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																//fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)
															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														}
													}
												} else {
													fmt.Println("No se encontraron detalles que requieran cambio")
												}
											} else {
												fmt.Println("Error al obtener detalles para cambio")
											}
										}

									} else {
										fmt.Println("No se encontraron los detalles del contrato general")
									}
								} else {
									fmt.Println("Error al traer los detalles del contrato general: ", err)
								}
							} else {
								fmt.Println("no se encontró contrato preliquidación para el contrato general")
							}
						} else {
							fmt.Println("Error al obtener el contrato preliquidacion: ", err)
						}
					} else {
						fmt.Println("No se encontró el conrato general")
					}
				} else {
					fmt.Println("Error al obtener el contrato general: ", err)
				}
			} else {
				fmt.Println("No se requiere actualización de valores")
			}
		} else {
			fmt.Println("El docente no tiene contratos registrados")
		}
	} else {
		fmt.Println("Error al intentar obtener contratos del docente: ", err)
	}
}

func ReglaDe3Old(contrato models.ContratoOld, mesIterativo int, anoIterativo int) {
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
	fmt.Println("Ingreso a regla de 3")
	//Obtener los valores del ibc liquidado para saber si es necesario realizar actualizacion
	query := "Documento:" + contrato.Documento + ",TipoNominaId:410,Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratosDocente)
		if contratosDocente[0].Id != 0 {
			fmt.Println("Tamaño arreglo contratos: ", len(contratosDocente))
			for i := 0; i < len(contratosDocente); i++ {
				fmt.Println("iteracion: ", i)
				fmt.Println(contratosDocente[i].NumeroContrato)
				query = "ContratoId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoPreliquidacionDocente)
					if contratoPreliquidacionDocente[0].Id != 0 {
						if contratosDocente[i].NumeroContrato != "GENERAL"+strconv.Itoa(mesIterativo) {
							fmt.Println("Agrego el contrato: ", contratosDocente[i].NumeroContrato)
							contratosCambio = append(contratosCambio, contratoPreliquidacionDocente[0].Id)
						} else {
							if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId.Id:"+strconv.Itoa(contratoPreliquidacionDocente[0].Id)+",ConceptoNominaId.Id:521", &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &auxValor)
								if auxValor[0].Id != 0 {
									ibcGeneral = auxValor[0].ValorCalculado
								} else {
									fmt.Println("No se encontró ibc para el contrato: ", contratosDocente[i].NumeroContrato)
								}
							} else {
								fmt.Println("Error al obtener el valor del ibc para el contrato: ", contratosDocente[i].NumeroContrato)
							}
							if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId.Id:"+strconv.Itoa(contratoPreliquidacionDocente[0].Id)+",ConceptoNominaId.Id:152", &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &auxValor)
								if auxValor[0].Id != 0 {
									salarioGeneral = auxValor[0].ValorCalculado
								} else {
									fmt.Println("No se encontraron salarion para el contrato: ", contratosDocente[i].NumeroContrato)
								}
							} else {
								fmt.Println("Error al obtener el valor del ibc para el contrato: ", contratosDocente[i].NumeroContrato)
							}
							// fmt.Println("salarioGeneral: ", salarioGeneral)
							// fmt.Println("ibcGeneral: ", ibcGeneral)
							// fmt.Println("contratosDocente: ", contratosDocente)
							if salarioGeneral < ibcGeneral && len(contratosDocente) > 2 {
								cambioNecesario = true
								break
							}
						}
					} else {
						fmt.Println("No se encontraron preliquidaciones asociadas al contrato: ", contratosDocente[i].NumeroContrato)
					}
				} else {
					fmt.Println("Error al obtener el contrato preliquidación para el contrato: ", contratosDocente[i].NumeroContrato)
				}
			}
			//por defecto que se realice la regla de 3
			cambioNecesario = true
			//Hacer regla de 3 en caso de que el cambio sea necesario
			if cambioNecesario {
				//obtener el contrato general
				query = "Documento:" + contrato.Documento + ",TipoNominaId:410,NumeroContrato:GENERAL" + strconv.Itoa(mesIterativo) + ",Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoGeneral)
					if contratoGeneral[0].Id != 0 {
						//Obtener el contrato preliquidacion del contrato general
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query=ContratoId:"+strconv.Itoa(contratoGeneral[0].Id), &aux); err == nil {
							var auxCp []models.ContratoPreliquidacionOld //Variable auxiliar de contrato preliquidacion
							LimpiezaRespuestaRefactor(aux, &auxCp)
							if auxCp[0].Id != 0 {
								//traer los detalles necesarios para hacer la reglas de tres
								auxDetalle = nil
								if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(auxCp[0].Id), &aux); err == nil {
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
												fmt.Println("Total honorarios:", totalHonorarios)
												fmt.Println("------------------------------------------------------------")
											case 64:
												valorRetefuente = auxDetalle[i].ValorCalculado
												fmt.Println("Total retefuente:", valorRetefuente)
												fmt.Println("------------------------------------------------------------")
											case 170:
												valorFondoSol = auxDetalle[i].ValorCalculado
												fmt.Println("Total fondo sol:", valorFondoSol)
												fmt.Println("------------------------------------------------------------")
											case 572:
												valorFondoSub = auxDetalle[i].ValorCalculado
												fmt.Println("Total fondo sub:", valorFondoSub)
												fmt.Println("------------------------------------------------------------")
											case 568:
												valorSalud = auxDetalle[i].ValorCalculado
												fmt.Println("Total Salud:", valorSalud)
												fmt.Println("------------------------------------------------------------")
											case 569:
												valorPension = auxDetalle[i].ValorCalculado
												fmt.Println("Total Pension:", valorPension)
												fmt.Println("------------------------------------------------------------")
											case 570:
												valorArl = auxDetalle[i].ValorCalculado
												fmt.Println("Total Arl:", valorArl)
												fmt.Println("------------------------------------------------------------")
											case 521:
												valorIbc = auxDetalle[i].ValorCalculado
												fmt.Println("Total ibc:", valorIbc)
												fmt.Println("------------------------------------------------------------")
											case 576:
												valorSaludUniversidad = auxDetalle[i].ValorCalculado
												fmt.Println("Total salud Universidad:", valorSaludUniversidad)
												fmt.Println("------------------------------------------------------------")
											case 577:
												valorPensionUniversidad = auxDetalle[i].ValorCalculado
												fmt.Println("Total Pensión universidad:", valorPensionUniversidad)
												fmt.Println("------------------------------------------------------------")
											}
										}
										//Obtener los detalles que necesitan cambio
										auxDetalle = nil
										var detalleEnvio models.DetallePreliquidacionOld
										for i := 0; i < len(contratosCambio); i++ {
											if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contratosCambio[i]), &aux); err == nil {
												LimpiezaRespuestaRefactor(aux, &auxDetalle)
												if auxDetalle[0].Id != 0 {

													for j := 0; j < len(auxDetalle); j++ {
														if auxDetalle[j].ConceptoNominaId.Id == 152 {
															valorMensual = auxDetalle[j].ValorCalculado
															fmt.Println("Honorarios para el contrato: ", valorMensual)
														}
													}

													for j := 0; j < len(auxDetalle); j++ {

														switch auxDetalle[j].ConceptoNominaId.Id {
														case 64:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorRetefuente)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)
															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 170:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorFondoSol)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 572:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorFondoSub)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 568:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorSalud)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 569:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorPension)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 570:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorArl)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 521:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorIbc)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 576:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorSaludUniversidad)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)

															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														case 577:
															detalleEnvio = auxDetalle[j]
															//Actualizar valor
															detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorPensionUniversidad)
															if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
																fmt.Println("Se ha actualizado: ", detalleEnvio.ConceptoNominaId.AliasConcepto, " con el valor de: ", detalleEnvio.ValorCalculado)
															} else {
																fmt.Println("Error al actualizar el valor de: ", detalleEnvio.ConceptoNominaId.AliasConcepto)
															}
														}
													}
												} else {
													fmt.Println("No se encontraron detalles que requieran cambio")
												}
											} else {
												fmt.Println("Error al obtener detalles para cambio")
											}
										}

									} else {
										fmt.Println("No se encontraron los detalles del contrato general")
									}
								} else {
									fmt.Println("Error al traer los detalles del contrato general: ", err)
								}
							} else {
								fmt.Println("no se encontró contrato preliquidación para el contrato general")
							}
						} else {
							fmt.Println("Error al obtener el contrato preliquidacion: ", err)
						}
					} else {
						fmt.Println("No se encontró el conrato general")
					}
				} else {
					fmt.Println("Error al obtener el contrato general: ", err)
				}
			} else {
				fmt.Println("No se requiere actualización de valores")
			}
		} else {
			fmt.Println("El docente no tiene contratos registrados")
		}
	} else {
		fmt.Println("Error al intentar obtener contratos del docente: ", err)
	}
}
