package controllers

import (
	"fmt"
	"math"
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

func liquidarHCS(contrato models.Contrato, general bool, porcentaje float64) {
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
	var porcentaje_ibc float64

	//Para el contrato general
	var contratoGeneral []models.Contrato //Contrato general mensual para la liquidación general

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
	if general || contrato.Unico {
		predicados = append(predicados, models.Predicado{Nombre: "general(1)."})
		fmt.Println("El contrato es general o único, se carga regla")
	} else {
		predicados = append(predicados, models.Predicado{Nombre: "general(0)."})
		fmt.Println("El docente tiene varios contratos, no se carga regla de único")
	}
	predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
	predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})

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
				predicados = append(predicados, models.Predicado{Nombre: "vacaciones(" + fmt.Sprintf("%f", contrato.Vacaciones) + ")."})
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

			} else if mesIterativo == int(contrato.FechaFin.Month()) && contrato.FechaFin.Year() == anoIterativo {
				//Para el mes final
				predicados = append(predicados, models.Predicado{Nombre: "vacaciones(" + fmt.Sprintf("%f", contrato.Vacaciones) + ")."})

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
					query := "ContratoPreliquidacionId.PreliquidacionId.Ano:" + strconv.Itoa(ano) + ",ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mes) + ",ContratoPreliquidacionId.ContratoId.NumeroContrato:" + contrato.NumeroContrato + ",ContratoPreliquidacionId.ContratoId.Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",ContratoPreliquidacionId.ContratoId.Documento:" + contrato.Documento + ",ContratoPreliquidacionId.ContratoId.DependenciaId:" + strconv.Itoa(contrato.DependenciaId) + ",ContratoPreliquidacionId.ContratoId.Rp:" + strconv.Itoa(contrato.Rp)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						fmt.Println()
						LimpiezaRespuestaRefactor(aux, &semanas)
						for i := 0; i < len(semanas); i++ {
							if semanas[i].ConceptoNominaId.Id == 152 {
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
				detallePreliquidacion.DiasLiquidados = float64(semanas_liquidadas)

				if porcentaje != 0 {
					porcentaje_ibc = porcentaje
				} else {
					diasALiquidar, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
					dias, _ := strconv.ParseFloat(diasALiquidar, 64)
					porcentaje_ibc = dias / 30
				}

			} else {
				if general {
					predicados = append(predicados, models.Predicado{Nombre: "vacaciones(" + fmt.Sprintf("%f", contrato.Vacaciones) + ")."})
				} else {
					predicados = append(predicados, models.Predicado{Nombre: "vacaciones(0)."})

				}
				semanas_liquidadas = 4
				detallePreliquidacion.DiasLiquidados = 4
				porcentaje_ibc = 1
			}

			reglasbase := cargarReglasBase("HCS") + reglasAlivios + FormatoReglas(predicados)

			reglasNuevas = reglasNuevas + reglasbase + "porcentaje(" + fmt.Sprintf("%f", porcentaje_ibc) + ").semanas_liquidadas(" + contrato.Documento + "," + strconv.Itoa(semanas_liquidadas) + ")."

			if mesIterativo == int(contrato.FechaFin.Month()) && anoIterativo == contrato.FechaFin.Year() && !general {
				auxDetalle = golog.LiquidarMesHCS(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion, true)
			} else {
				auxDetalle = golog.LiquidarMesHCS(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion, false)
			}

			for j := 0; j < len(auxDetalle); j++ {
				registrarDetallePreliquidacion(auxDetalle[j])
			}

			if !general {
				fmt.Println("Liquidando Contrato General")
				LiquidarContratoGeneral(mesIterativo, anoIterativo, contrato, preliquidacion[0], porcentaje_ibc, "410")
				if !contrato.Unico {
					var contratosDocente []models.Contrato = nil
					var contratoPreliquidacionDocente []models.ContratoPreliquidacion = nil
					var auxValor []models.DetallePreliquidacion
					var ibcGeneral float64
					var salarioGeneral float64
					var contratosCambio []int
					var cambioNecesario bool = false

					//Obtener los valores del ibc liquidado para saber si es necesario realizar actualizacion
					query := "Documento:" + contrato.Documento + ",TipoNominaId:410,Vigencia:" + strconv.Itoa(contrato.Vigencia)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &contratosDocente)
						if contratosDocente[0].Id != 0 {
							fmt.Println("Tamaño arreglo: ", len(contratosDocente))
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

							//Hacer regla de 3 en caso de que el cambio sea necesario
							if cambioNecesario {
								//obtener el contrato general
								query = "Documento:" + contrato.Documento + ",TipoNominaId:410,NumeroContrato:GENERAL" + strconv.Itoa(mesIterativo) + ",Vigencia:" + strconv.Itoa(contrato.Vigencia)
								if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
									contratoGeneral = nil
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
}
