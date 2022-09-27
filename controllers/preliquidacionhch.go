package controllers

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// PreliquidacionhchController operations for Preliquidacioncthch
type PreliquidacionhchController struct {
	beego.Controller
}

func liquidarHCH(contrato models.Contrato, general bool, porcentaje float64) (mensaje string, err error) {
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
	cedula, err := strconv.ParseInt(contrato.Documento, 0, 64)
	var emergencia int //Varibale para evitar loop infinito

	//Para el valor de regla de 3

	//Para el contrato general
	var contratoGeneral []models.Contrato //Contrato general mensual para la liquidación general

	if err == nil {
		reglasAlivios, contratoPreliquidacion, err = CargarDatosRetefuente(int(cedula))
	}

	if err == nil {
		mesIterativo = int(contrato.FechaInicio.Month())
		anoIterativo = contrato.Vigencia

		//Obtener las semanas del contrato
		semanasContrato := int(calcularSemanasContratoDVE(contrato.FechaInicio, contrato.FechaFin))
		fmt.Println("SemanasContrato: ", semanasContrato)

		//Por si es general o unico
		if general || contrato.Unico {
			predicados = append(predicados, models.Predicado{Nombre: "general(1)."})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "general(0)."})
		}
		predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%f", contrato.ValorContrato) + "). "})
		predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanasContrato) + "," + strconv.Itoa(contrato.Vigencia) + "). "})
		reglasbase := cargarReglasBase("HCH") + reglasAlivios + FormatoReglas(predicados)

		for {

			fmt.Println("Mes: ", mesIterativo)
			fmt.Println("Año: ", anoIterativo)
			reglasNuevas = ""
			query := "Ano:" + strconv.Itoa(anoIterativo) + ",Mes:" + strconv.Itoa(mesIterativo) + ",Nominaid:415"
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
				_, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
				//detallePreliquidacion.DiasLiquidados, _ = strconv.ParseFloat(diasALiquidar, 64)
				//Calcular semanas a liquidar
				//Para cuando son contratos de 1 mes

				if contrato.FechaInicio.Month() == contrato.FechaFin.Month() && contrato.FechaInicio.Year() == contrato.FechaFin.Year() {
					//Calcular el numero de días
					diasALiquidar, detallePreliquidacion.DiasEspecificos = CalcularPeriodoLiquidacion(preliquidacion[0].Ano, preliquidacion[0].Mes, contrato.FechaInicio, contrato.FechaFin)
					semanas, _ := strconv.ParseFloat(diasALiquidar, 64)
					if porcentaje != 0 {
						porcentaje_ibc = porcentaje
					} else {
						porcentaje_ibc = semanas / 30
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
						query := "ContratoPreliquidacionId.PreliquidacionId.Ano:" + strconv.Itoa(ano) + ",ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mes) + ",ContratoPreliquidacionId.ContratoId.NumeroContrato:" + contrato.NumeroContrato + ",ContratoPreliquidacionId.ContratoId.Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",ContratoPreliquidacionId.ContratoId.DependenciaId:" + strconv.Itoa(contrato.DependenciaId) + ",ContratoPreliquidacionId.ContratoId.Documento:" + contrato.Documento + ",ContratoPreliquidacionId.ContratoId.Rp:" + strconv.Itoa(contrato.Rp)
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
							fmt.Println()
							LimpiezaRespuestaRefactor(aux, &semanas)
							for i := 0; i < len(semanas); i++ {
								if semanas[i].ConceptoNominaId.Id == 87 {
									semanas_liquidadas = semanas_liquidadas + int(semanas[i].DiasLiquidados)
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
					semanas_liquidadas = 4
					detallePreliquidacion.DiasLiquidados = 4
					porcentaje_ibc = 1
				}
				reglasNuevas = reglasNuevas + reglasbase + "porcentaje(" + fmt.Sprintf("%f", porcentaje_ibc) + ").semanas_liquidadas(" + contrato.Documento + "," + strconv.Itoa(semanas_liquidadas) + ")."
				auxDetalle = golog.LiquidarMesHCH(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion)
				for j := 0; j < len(auxDetalle); j++ {
					registrarDetallePreliquidacion(auxDetalle[j])
				}

				if !general {
					fmt.Println("Liquidando Contrato General")
					LiquidarContratoGeneral(mesIterativo, anoIterativo, contrato, preliquidacion[0], porcentaje_ibc, "409")
					if !contrato.Unico {
						//nueva lógica
						var contratosDocente []models.Contrato = nil
						var contratoPreliquidacionDocente []models.ContratoPreliquidacion = nil
						//var auxValor []models.DetallePreliquidacion
						//var ibcGeneral float64
						//var salarioGeneral float64
						var contratosCambio []int
						var cambioNecesario bool = true

						//Obtener los valores del ibc liquidado para saber si es necesario realizar actualizacion
						query := "Documento:" + contrato.Documento + ",TipoNominaId:409,Vigencia:" + strconv.Itoa(contrato.Vigencia)
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &contratosDocente)
							if contratosDocente[0].Id != 0 {
								for i := 0; i < len(contratosDocente); i++ {
									query = "ContratoId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
									if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
										contratoPreliquidacionDocente = nil
										LimpiezaRespuestaRefactor(aux, &contratoPreliquidacionDocente)
										if contratoPreliquidacionDocente[0].Id != 0 {
											if contratosDocente[i].NumeroContrato != "GENERAL"+strconv.Itoa(mesIterativo) {
												if !strings.HasPrefix(contratosDocente[i].NumeroContrato, "GENERAL") {
													fmt.Println("Agrego el contrato: ", contratosDocente[i].NumeroContrato)
													contratosCambio = append(contratosCambio, contratoPreliquidacionDocente[0].Id)
												}
											} /*else {
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
												if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId.Id:"+strconv.Itoa(contratoPreliquidacionDocente[0].Id)+",ConceptoNominaId.Id:87", &aux); err == nil {
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
											}*/
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
									query = "Documento:" + contrato.Documento + ",TipoNominaId:409,NumeroContrato:GENERAL" + strconv.Itoa(mesIterativo) + ",Vigencia:" + strconv.Itoa(contrato.Vigencia)
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
															var valorMensual float64 = 0
															//obtener los valores totales para realizar la regla de 3
															for i := 0; i < len(auxDetalle); i++ {
																switch auxDetalle[i].ConceptoNominaId.Id {
																case 87:
																	totalHonorarios = auxDetalle[i].ValorCalculado
																	fmt.Println("Total honorarios:", totalHonorarios)
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
																			if auxDetalle[j].ConceptoNominaId.Id == 87 {
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
								return "El docente no tiene contratos registrados", errors.New("No hay contratos registrados para el docente")
							}
						} else {
							fmt.Println("Error al intentar obtener contratos del docente: ", err)
							return "Error al intentar obtener contratos del docente: ", nil
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
				return "Error al consultar preliquidaciones", err
			}
			preliquidacion[0].Id = 0 //Para evitar errores al obtener la preliquidación del siguiente mes
		}
	} else {
		fmt.Println("Error al consultar información en Ágora")
		return "Error al consultar información en Ágora: ", err
	}
	return "No hubo error", nil
}
