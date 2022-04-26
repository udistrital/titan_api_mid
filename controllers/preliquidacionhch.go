package controllers

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// PreliquidacionhchController operations for Preliquidacioncthch
type PreliquidacionhchController struct {
	beego.Controller
}

func liquidarHCH(contrato models.Contrato, general bool) {
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

	//Para el valor de regla de 3

	//Para el contrato general
	var contratosDocente []models.ContratoPreliquidacion
	var contratoGeneral []models.Contrato //Contrato general mensual para la liquidación general

	if err == nil {
		reglasAlivios, contratoPreliquidacion = CargarDatosRetefuente(int(cedula))
	}

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
					query := "ContratoPreliquidacionId.PreliquidacionId.Ano:" + strconv.Itoa(ano) + ",ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mes) + ",ContratoPreliquidacionId.ContratoId.NumeroContrato:" + contrato.NumeroContrato + ",ContratoPreliquidacionId.ContratoId.Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",ContratoPreliquidacionId.ContratoId.DependenciaId:" + strconv.Itoa(contrato.DependenciaId) + ",ContratoPreliquidacionId.ContratoId.Documento:" + contrato.Documento + ",ContratoPreliquidacionId.ContratoId.Rp:" + strconv.Itoa(contrato.Rp)
					fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/detalle_preliquidacion?limit=-1&query=" + query)
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
			} else {
				semanas_liquidadas = 4
				detallePreliquidacion.DiasLiquidados = 4
			}
			reglasNuevas = reglasNuevas + reglasbase + "semanas_liquidadas(" + contrato.Documento + "," + strconv.Itoa(semanas_liquidadas) + ")."
			auxDetalle = golog.LiquidarMesHCH(reglasNuevas, contrato.Documento, contrato.Vigencia, detallePreliquidacion)
			for j := 0; j < len(auxDetalle); j++ {
				registrarDetallePreliquidacion(auxDetalle[j])
			}

			if !general {

				fmt.Println("Liquidando Contrato General")
				//Buscar el contrato general para este mes para la persona en cuestión, en caso de no existir se crea uno
				query := "NumeroContrato:GENERAL" + strconv.Itoa(mesIterativo) + ",Vigencia:" + strconv.Itoa(anoIterativo) + ",Documento:" + contrato.Documento + ",TipoNominaId:409"
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoGeneral)
					fmt.Println(contratoGeneral[0])
					if contratoGeneral[0].Id == 0 {
						//Crear contrato General
						contratoGeneral[0].NumeroContrato = "GENERAL" + strconv.Itoa(mesIterativo)
						contratoGeneral[0].Vigencia = contrato.Vigencia
						contratoGeneral[0].NombreCompleto = contrato.NombreCompleto
						contratoGeneral[0].Documento = contrato.Documento
						contratoGeneral[0].PersonaId = contrato.PersonaId
						contratoGeneral[0].TipoNominaId = contrato.TipoNominaId
						contratoGeneral[0].Activo = true
						contratoGeneral[0].FechaInicio = time.Date(anoIterativo, time.Month(mesIterativo), 1, 12, 0, 0, 0, time.UTC)
						if mesIterativo == 2 {
							contratoGeneral[0].FechaFin = time.Date(anoIterativo, time.Month(mesIterativo), 28, 12, 0, 0, 0, time.UTC)
						} else {
							contratoGeneral[0].FechaFin = time.Date(anoIterativo, time.Month(mesIterativo), 30, 12, 0, 0, 0, time.UTC)
						}

						//Buscar el valor de los honorarios de los contratos que tiene el docente en ese mes
						query = "PreliquidacionId.Id:" + strconv.Itoa(preliquidacion[0].Id) + ",ContratoId.Documento:" + contrato.Documento + ",ContratoId.TipoNominaId:409"
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &contratosDocente)
							if len(contratosDocente) >= 1 { //Tiene 1 contrato o más
								//Sumar valores de los honorarios para obtener el valor total de ese mes
								contratoGeneral[0].ValorContrato = 0
								for i := 0; i < len(contratosDocente); i++ {
									//Sumar los honorarios de el mes presente para obtener el IBC
									query := "ContratoPreliquidacionId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",ConceptoNominaId.Id:87"
									if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
										LimpiezaRespuestaRefactor(aux, &auxDetalle)
										fmt.Println("Honorarios Encontrados: ", auxDetalle[0].ValorCalculado)
										contratoGeneral[0].ValorContrato = contratoGeneral[0].ValorContrato + auxDetalle[0].ValorCalculado
									} else {
										fmt.Println("Error al obtener los honorarios para el contrato :", contratosDocente[i].ContratoId.NumeroContrato, " ", err)
									}
								}
							}
						} else {
							fmt.Println("Error al obtener los contratos vigentes para el mes actual: ", err)
						}

						//Registrar el contrato nuevo
						contratoGeneral[0], _ = registrarContrato(contratoGeneral[0])

					} else {
						fmt.Println("Contrato Encontrado: ", contratoGeneral[0])
						//Eliminar los detalles del contrato General
						query := "ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contratoGeneral[0].Id) + ",ContratoPreliquidacionId.ContratoId.Vigencia:" + strconv.Itoa(anoIterativo) + ",ContratoPreliquidacionId.ContratoId.Documento:" + contrato.Documento
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &auxDetalle)
							idContratoPeliquidacion := auxDetalle[0].ContratoPreliquidacionId.Id
							for j := 0; j < len(auxDetalle); j++ {
								if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(auxDetalle[j].Id), "DELETE", &aux, nil); err == nil {
									fmt.Println("Detalle Eliminado")
								} else {
									fmt.Println("Error al eliminar detalle: ", err)
								}
							}
							//Eliminar el contrato_preliquidación
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(idContratoPeliquidacion), "DELETE", &aux, nil); err == nil {
								fmt.Println("contrato Preliquidacion Eliminado")
								//Actualizar el valor del contrato general
								//Buscar el valor de los honorarios de los contratos que tiene el docente en ese mes

								query = "PreliquidacionId.Id:" + strconv.Itoa(preliquidacion[0].Id) + ",ContratoId.Documento:" + contrato.Documento
								fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/contrato_preliquidacion?limit=-1&query=" + query)
								if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
									LimpiezaRespuestaRefactor(aux, &contratosDocente)
									if len(contratosDocente) >= 1 { //Tiene más de un contrato
										//Sumar valores de los honorarios para obtener el valor total de ese mes
										contratoGeneral[0].ValorContrato = 0
										for i := 0; i < len(contratosDocente); i++ {
											//Sumar los honorarios de el mes presente para obtener el IBC
											if contratosDocente[i].ContratoId.Id != contratoGeneral[0].Id {
												query := "ContratoPreliquidacionId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",ConceptoNominaId.Id:87"
												fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/detalle_preliquidacion?limit=-1&query=" + query)
												if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
													LimpiezaRespuestaRefactor(aux, &auxDetalle)
													contratoGeneral[0].ValorContrato = contratoGeneral[0].ValorContrato + auxDetalle[0].ValorCalculado
												} else {
													fmt.Println("Error al obtener los honorarios para el contrato :", contratosDocente[i].ContratoId.NumeroContrato, " ", err)
												}
											}
										}
										//Actualizar

										if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contratoGeneral[0].Id), "PUT", &aux, contratoGeneral[0]); err == nil {
											fmt.Println("Valor Actualizado")
										} else {
											fmt.Println("Error al actualizar valor del contrato")
										}
									}
								} else {
									fmt.Println("Error al obtener los contratos vigentes para el mes actual: ", err)
								}
							} else {
								fmt.Println("Error al eliminar contrato_preliquidacion: ", err)
							}
						} else {
							fmt.Println("Error al obtener los detalles para el contrato general del mes")
						}
					}
					//Liquidar el nuevo contrato
					liquidarHCH(contratoGeneral[0], true)
					//Actualizar registros de los valores calculados por regla de 3 para ese mes

					query := "ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",ContratoPreliquidacionId.PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo) + ",ContratoPreliquidacionId.ContratoId.Documento:" + contrato.Documento + ",ContratoPreliquidacionId.ContratoId.TipoNominaId:409"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &auxDetalle)
						if auxDetalle[0].Id != 0 {
							var totalHonorarios float64
							var valorIbc float64
							var valorSalud float64
							var valorPension float64
							var valorArl float64
							var valorReteica float64
							var valorMensual float64
							var valorRetefuente float64
							var valorFondoSol float64
							var valorFondoSub float64
							var detalleEnvio models.DetallePreliquidacion
							var contratosCambio []int //Ids de los contratos que necesitan el ajuste.

							//Obtener los honorarios y el detalle de los valores que necesitan cambio
							for j := 0; j < len(auxDetalle); j++ {
								if auxDetalle[j].ContratoPreliquidacionId.ContratoId.NumeroContrato == "GENERAL"+strconv.Itoa(mesIterativo) {
									if auxDetalle[j].ConceptoNominaId.Id == 87 {
										totalHonorarios = auxDetalle[j].ValorCalculado
									} else if auxDetalle[j].ConceptoNominaId.Id == 64 {
										valorRetefuente = auxDetalle[j].ValorCalculado
									} else if auxDetalle[j].ConceptoNominaId.Id == 170 {
										valorFondoSol = auxDetalle[j].ValorCalculado
									} else if auxDetalle[j].ConceptoNominaId.Id == 572 {
										valorFondoSub = auxDetalle[j].ValorCalculado
									} else if auxDetalle[j].ConceptoNominaId.Id == 568 {
										valorSalud = auxDetalle[j].ValorCalculado
									} else if auxDetalle[j].ConceptoNominaId.Id == 569 {
										valorPension = auxDetalle[j].ValorCalculado
									} else if auxDetalle[j].ConceptoNominaId.Id == 570 {
										valorArl = auxDetalle[j].ValorCalculado
									} else if auxDetalle[j].ConceptoNominaId.Id == 521 {
										valorIbc = auxDetalle[j].ValorCalculado
									} else if auxDetalle[j].ConceptoNominaId.Id == 545 {
										valorReteica = auxDetalle[j].ValorCalculado
									}
								} else if auxDetalle[j].ContratoPreliquidacionId.ContratoId.NumeroContrato != "GENERAL"+strconv.Itoa(mesIterativo) && auxDetalle[j].ConceptoNominaId.Id == 87 {
									contratosCambio = append(contratosCambio, auxDetalle[j].ContratoPreliquidacionId.ContratoId.Id)
								}
							}

							//Recorrer y cambiar los valores teniendo en cuenta su aporte a salud y pension
							for j := 0; j < len(contratosCambio); j++ {
								//Obtener detalles
								query := "ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",ContratoPreliquidacionId.PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo) + ",ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contratosCambio[j])
								if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
									LimpiezaRespuestaRefactor(aux, &auxDetalle)
									if auxDetalle[0].Id != 0 {
										//Obtener los honorarios para hacer regla de 3
										for k := 0; k < len(auxDetalle); k++ {
											if auxDetalle[k].ConceptoNominaId.Id == 87 {
												valorMensual = auxDetalle[k].ValorCalculado
											}
										}
										//Hacer regla de 3
										for k := 0; k < len(auxDetalle); k++ {
											if auxDetalle[k].ConceptoNominaId.Id == 64 {
												detalleEnvio = auxDetalle[k]
												//Actualizar valor
												detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorRetefuente)
												//Enviar a actualización
												if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
													fmt.Println("Rete Fuente actualizada")
												} else {
													fmt.Println("No se pudo actualizar Rete Fuente")
												}

											} else if auxDetalle[k].ConceptoNominaId.Id == 170 {
												detalleEnvio = auxDetalle[k]
												//Actualizar valor
												detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorFondoSol)
												//Enviar a actualización
												if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
													fmt.Println("Fondo sub actualizado")
												} else {
													fmt.Println("No se pudo actualizar fondo sub")
												}
											} else if auxDetalle[k].ConceptoNominaId.Id == 568 {
												detalleEnvio = auxDetalle[k]
												//Actualizar valor
												detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorSalud)
												//Enviar a actualización
												if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
													fmt.Println("Salud actualizada")
												} else {
													fmt.Println("No se pudo actualizar Salud")
												}
											} else if auxDetalle[k].ConceptoNominaId.Id == 569 {
												detalleEnvio = auxDetalle[k]
												//Actualizar valor
												detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorPension)
												//Enviar a actualización
												if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
													fmt.Println("Pensión actualizada")
												} else {
													fmt.Println("No se pudo actualizar Pensión")
												}
											} else if auxDetalle[k].ConceptoNominaId.Id == 570 {
												detalleEnvio = auxDetalle[k]
												//Actualizar valor
												detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorArl)
												//Enviar a actualización
												if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
													fmt.Println("Arl actualizada")
												} else {
													fmt.Println("No se pudo actualizar Arl")
												}
											} else if auxDetalle[k].ConceptoNominaId.Id == 521 {
												detalleEnvio = auxDetalle[k]
												//Actualizar valor
												detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorIbc)
												//Enviar a actualización
												if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
													fmt.Println("Ibc actualizado")
												} else {
													fmt.Println("No se pudo actualizar ibc")
												}
											} else if auxDetalle[k].ConceptoNominaId.Id == 545 {
												detalleEnvio = auxDetalle[k]
												//Actualizar valor
												detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorReteica)
												//Enviar a actualización
												if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
													fmt.Println("Reteica actualizada")
												} else {
													fmt.Println("No se pudo actualizar reteica")
												}
											} else if auxDetalle[k].ConceptoNominaId.Id == 572 {
												detalleEnvio = auxDetalle[k]
												//Actualizar valor
												detalleEnvio.ValorCalculado = math.Round((valorMensual / totalHonorarios) * valorFondoSub)
												//Enviar a actualización
												if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalleEnvio.Id), "PUT", &aux, detalleEnvio); err == nil {
													fmt.Println("Fondo Sol Actualizado actualizada")
												} else {
													fmt.Println("No se pudo actualizar fondo sol")
												}
											}
										}
									} else {
										fmt.Println("No hay conceptos para el contrato: ", contratosCambio[j])
									}
								} else {
									fmt.Println("Error al obtener los detalles del contrato")
								}
							}
						} else {
							fmt.Println("No se encotraron detalles")
						}
					} else {
						fmt.Println("Error al actualizar los conceptos en los contratos")
					}

				} else {
					fmt.Println("Error buscar contrato general mensual: ", err)
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
		preliquidacion[0].Id = 0  //Para evitar errores al obtener la preliquidación del siguiente mes
		contratoGeneral[0].Id = 0 //Para no obtener problemas con el contrato General del siguiente mes
	}
}
