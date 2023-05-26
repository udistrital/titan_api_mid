package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func diff(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	oneDay := time.Hour * 5
	a = a.Add(oneDay)
	b = b.Add(oneDay)
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year = y2 - y1
	month = int(M2 - M1)
	day = d2 - d1

	if day < 0 {

		day = (30 - d1) + d2
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func LimpiezaRespuestaRefactor(respuesta map[string]interface{}, v interface{}) {
	b, err := json.Marshal(respuesta["Data"])
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &v)
}

func FormatoReglas(v []models.Predicado) (reglas string) {
	var arregloReglas = make([]string, len(v))
	reglas = ""
	//var respuesta []models.FormatoPreliqu
	for i := 0; i < len(v); i++ {
		arregloReglas[i] = v[i].Nombre
	}

	for i := 0; i < len(arregloReglas); i++ {
		reglas = reglas + arregloReglas[i] + "\n"
	}
	return
}

func CalcularDias(FechaInicio time.Time, FechaFin time.Time) (diasLaborados float64, meses float64) {

	var a, m, d int
	var mesesContrato float64
	var diasContrato float64
	if FechaFin.IsZero() {
		FechaFin2 := time.Now()
		a, m, d = diff(FechaInicio, FechaFin2)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
		diasContrato = mesesContrato * 30
	} else {
		a, m, d = diff(FechaInicio, FechaFin)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
		diasContrato = mesesContrato * 30
	}
	return diasContrato, mesesContrato

}

func calcularSemanasContratoDVE(FechaInicio time.Time, FechaFin time.Time) (semanas float64) {
	var a, m, d int
	var mesesContrato float64
	if FechaFin.IsZero() {
		FechaFin2 := time.Now()
		a, m, d = diff(FechaInicio, FechaFin2)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)

	} else {
		a, m, d = diff(FechaInicio, FechaFin)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
	}
	if mesesContrato/float64(int(mesesContrato)) != 1 {
		return (mesesContrato * 4) + 1
	} else {
		return (mesesContrato * 4)
	}
}

func registrarPreliquidacion(año, mes, estadoPreliquidacion, nomina int) (preliquidacion models.Preliquidacion) {
	var aux map[string]interface{}
	var nombre string

	if nomina == 412 {
		nombre = "Funcionarios Administrativos-Planta-"
	} else if nomina == 413 {
		nombre = "Docentes de Planta-"
	} else if nomina == 414 {
		nombre = "Contratistas-"
	} else if nomina == 415 {
		nombre = "Hora cátedra honorarios-"
	} else {
		nombre = "Hora cátedra salarios-"
	}
	preliquidacion.Descripcion = nombre + strconv.Itoa(año) + strconv.Itoa(mes)
	preliquidacion.Ano = año
	preliquidacion.Mes = mes
	preliquidacion.NominaId = nomina
	preliquidacion.EstadoPreliquidacionId = estadoPreliquidacion
	preliquidacion.Activo = true
	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion", "POST", &aux, preliquidacion); err == nil {
		LimpiezaRespuestaRefactor(aux, &preliquidacion)
		fmt.Println("Preliquidación creada con éxito")
	} else {
		fmt.Println("Error al guardar preliquidacion: ", err)
	}

	return preliquidacion
}

func registrarContratoPreliquidacion(preliquidacionId int, contratoId int, contratoPreliq models.ContratoPreliquidacion) (contratoPreliquidacion models.ContratoPreliquidacion) {
	var aux map[string]interface{}
	var respuesta models.ContratoPreliquidacion
	contratoPreliq.Id = 0
	contratoPreliquidacion = contratoPreliq
	contratoPreliquidacion.ContratoId = &models.Contrato{Id: contratoId}
	contratoPreliquidacion.PreliquidacionId = &models.Preliquidacion{Id: preliquidacionId}
	contratoPreliquidacion.Cumplido = false
	contratoPreliquidacion.Preliquidado = false
	contratoPreliquidacion.Activo = true

	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion", "POST", &aux, contratoPreliquidacion); err == nil {
		LimpiezaRespuestaRefactor(aux, &respuesta)
		fmt.Println("Contrato_preliquidacion guardado con exito")
	} else {
		fmt.Println("Error al guardar contrato_preliquidacion: ", err)
	}
	return respuesta
}

func registrarDetallePreliquidacion(detallePreliquidacion models.DetallePreliquidacion) {
	var response models.DetallePreliquidacion
	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion", "POST", &response, detallePreliquidacion); err == nil {
		fmt.Println("Detalle guardado con éxito")
	} else {
		fmt.Println("Error al guardar detalle", err)
	}
}

func registrarContrato(contrato models.Contrato) (respuesta models.Contrato, err error) {
	var response map[string]interface{}
	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato", "POST", &response, contrato); err == nil {
		fmt.Println("Contrato guardado con éxito ", response)
		LimpiezaRespuestaRefactor(response, &respuesta)
	} else {
		fmt.Println("Error al guardar contrato", err)
		return contrato, err
	}
	return respuesta, nil
}

func CalcularPeriodoLiquidacion(anoPreliquidacion, mesPreliquidacion int, fechaInicio, fechaFin time.Time) (periodoLiquidacion, periodoEspecifico string) {

	var FechaControl time.Time
	var periodo_liquidacion float64

	//En caso de que sea el mes de inicio y el mes final el mismo
	if fechaInicio.Month() == fechaFin.Month() && fechaInicio.Year() == fechaFin.Year() {
		periodo_liquidacion, _ = CalcularDias(fechaInicio, fechaFin)
		periodo_liquidacion = periodo_liquidacion + 1
		periodoEspecifico = "Del " + strconv.Itoa(fechaInicio.Day()) + " al " + strconv.Itoa(fechaFin.Day()) + " del mes " + strconv.Itoa(mesPreliquidacion) + " del año " + strconv.Itoa(anoPreliquidacion)
		//Para el mes de inicio
	} else if int(fechaInicio.Month()) == mesPreliquidacion && int(fechaInicio.Year()) == anoPreliquidacion {
		//En caso de que sea febrero, si no se coloca 28 tomará días de más
		if mesPreliquidacion == 2 {
			FechaControl = time.Date(anoPreliquidacion, time.Month(mesPreliquidacion), 28, 0, 0, 0, 0, time.UTC)
			periodo_liquidacion, _ = CalcularDias(fechaInicio, FechaControl)
			periodo_liquidacion = periodo_liquidacion + 3 //Dia inclusive y 2 días de desface de febrero
			periodoEspecifico = "Del " + strconv.Itoa(fechaInicio.Day()) + " al " + strconv.Itoa(FechaControl.Day()) + " del mes " + strconv.Itoa(mesPreliquidacion) + " del año " + strconv.Itoa(anoPreliquidacion)
		} else {
			FechaControl = time.Date(anoPreliquidacion, time.Month(mesPreliquidacion), 30, 0, 0, 0, 0, time.UTC)
			periodo_liquidacion, _ = CalcularDias(fechaInicio, FechaControl)
			periodo_liquidacion = periodo_liquidacion + 1 //Dia inclusive
			periodoEspecifico = "Del " + strconv.Itoa(fechaInicio.Day()) + " al " + strconv.Itoa(FechaControl.Day()) + " del mes " + strconv.Itoa(mesPreliquidacion) + " del año " + strconv.Itoa(anoPreliquidacion)
		}
		//Para el mes final
	} else if int(fechaFin.Month()) == mesPreliquidacion && int(fechaFin.Year()) == anoPreliquidacion {
		FechaControl = time.Date(anoPreliquidacion, time.Month(mesPreliquidacion), 1, 0, 0, 0, 0, time.UTC)
		periodo_liquidacion, _ = CalcularDias(FechaControl, fechaFin)
		periodo_liquidacion = periodo_liquidacion + 1 //Dia Inclusivo
		periodoEspecifico = "Del " + strconv.Itoa(FechaControl.Day()) + " al " + strconv.Itoa(fechaFin.Day()) + " del mes " + strconv.Itoa(mesPreliquidacion) + " del año " + strconv.Itoa(anoPreliquidacion)
	} else {
		periodo_liquidacion = 30
		periodoEspecifico = "Del 1 al 30 del mes " + strconv.Itoa(mesPreliquidacion) + " del año " + strconv.Itoa(anoPreliquidacion)
	}

	periodo := strconv.Itoa(int(periodo_liquidacion))

	return periodo, periodoEspecifico
}

func CalcularSemanas(diasLiquidados float64) (semanas int) {
	aux := diasLiquidados / 7

	if aux <= 1 {
		return 1
	} else if aux <= 2 {
		return 2
	} else if aux <= 3 {
		return 3
	} else {
		return 4
	}
}

func Remove(s []models.Contrato, i int) []models.Contrato {
	s = append(s[:i], s[i+1:]...)
	return s
}

func Roundf(x float64) float64 {
	t := math.Trunc(x)
	if math.Abs(x-t) >= 0.5 {
		return t + math.Copysign(1, x)
	}
	return t
}

func LiquidarContratoGeneral(mesIterativo int, anoIterativo int, contrato models.Contrato, preliquidacion models.Preliquidacion, porcentaje float64, nomina string, vigencia_original int, crear bool) {
	var aux map[string]interface{}
	var contratoGeneral []models.Contrato
	var contratosDocente []models.ContratoPreliquidacion
	var auxDetalle []models.DetallePreliquidacion
	var flag bool = true

	query := "NumeroContrato:GENERAL" + strconv.Itoa(mesIterativo) + ",Vigencia:" + strconv.Itoa(anoIterativo) + ",Documento:" + contrato.Documento + ",TipoNominaId:" + nomina + ",Activo:true"
	fmt.Println("QUERY ", query)
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratoGeneral)
		if contratoGeneral[0].Id == 0 {
			//Crear contrato General

			contratoGeneral[0].NumeroContrato = "GENERAL" + strconv.Itoa(mesIterativo)
			contratoGeneral[0].Vigencia = anoIterativo
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

			query = "PreliquidacionId.Id:" + strconv.Itoa(preliquidacion.Id) + ",ContratoId.Documento:" + contrato.Documento + ",ContratoId.TipoNominaId:" + nomina + ",ContratoId.Activo:true"
			fmt.Println("honorarios para general: ", query)
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contratosDocente)
				if len(contratosDocente) >= 1 { //Tiene más de dos contratos
					//Sumar valores de los honorarios para obtener el valor total de ese mes
					contratoGeneral[0].ValorContrato = 0
					var valorContratoGeneralAux = 0
					for i := 0; i < len(contratosDocente); i++ {
						//Sumar los honorarios de el mes presente para obtener el IBC
						if nomina == "410" {
							query = "ContratoPreliquidacionId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",ConceptoNominaId.Id:152,ContratoPreliquidacionId.ContratoId.Activo:true"

						} else {
							query = "ContratoPreliquidacionId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",ConceptoNominaId.Id:87,ContratoPreliquidacionId.ContratoId.Activo:true"
						}
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &auxDetalle)
							valorContratoGeneralAux += int(auxDetalle[0].ValorCalculado)
						} else {
							fmt.Println("Error al obtener los honorarios para el contrato :", contratosDocente[i].ContratoId.NumeroContrato, " ", err)
						}
					}
					contratoGeneral[0].ValorContrato = float64(valorContratoGeneralAux)
				}
			} else {
				fmt.Println("Error al obtener los contratos vigentes para el mes actual: ", err)
			}

			//Registrar el contrato nuevo
			if crear {
				contratoGeneral[0], _ = registrarContrato(contratoGeneral[0])
			}
		} else {
			//Eliminar los detalles del contrato General
			query := "ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contratoGeneral[0].Id) + ",ContratoPreliquidacionId.ContratoId.Vigencia:" + strconv.Itoa(anoIterativo)
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
					contratoGeneral[0].ValorContrato = 0
					query = "PreliquidacionId.Id:" + strconv.Itoa(preliquidacion.Id) + ",ContratoId.Documento:" + contrato.Documento + ",ContratoId.Activo:true"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &contratosDocente)
						if len(contratosDocente) >= 1 && contratosDocente[0].Id != 0 { //Tiene más de un contrato
							//Sumar valores de los honorarios para obtener el valor total de ese mes
							var valorContratoGeneralAux = 0
							for i := 0; i < len(contratosDocente); i++ {
								//Sumar los honorarios de el mes presente para obtener el IBC
								if contratosDocente[i].ContratoId.Id != contratoGeneral[0].Id {
									if nomina == "410" {
										query = "ContratoPreliquidacionId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",ConceptoNominaId.Id:152,ContratoPreliquidacionId.ContratoId.Activo:true"

									} else {
										query = "ContratoPreliquidacionId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",ConceptoNominaId.Id:87,ContratoPreliquidacionId.ContratoId.Activo:true"
									}
									if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
										LimpiezaRespuestaRefactor(aux, &auxDetalle)
										valorContratoGeneralAux += int(auxDetalle[0].ValorCalculado)
									} else {
										fmt.Println("Error al obtener los honorarios para el contrato :", contratosDocente[i].ContratoId.NumeroContrato, " ", err)
									}
								}
							}
							contratoGeneral[0].ValorContrato = float64(valorContratoGeneralAux)
							//Actualizar

							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contratoGeneral[0].Id), "PUT", &aux, contratoGeneral[0]); err == nil {
								fmt.Println("Valor Actualizado")
							} else {
								fmt.Println("Error al actualizar valor del contrato")
							}
						} else {
							flag = false
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

		if nomina == "410" && flag {
			liquidarHCS(contratoGeneral[0], true, porcentaje, vigencia_original, 0, 0, false)
		} else if nomina == "409" && flag {
			liquidarHCH(contratoGeneral[0], true, porcentaje, vigencia_original)
		}
	} else {
		fmt.Println("Error al buscar contrato general:", err)
	}
}

func anularEnGenerales(contrato models.Contrato, fecha_anulacion time.Time, vigencia_original int) {

	var aux map[string]interface{}
	var preliquidacion []models.Preliquidacion

	anio_aux := int(fecha_anulacion.Year())
	mes_aux := int(fecha_anulacion.Month()) + 1

	for {
		var contratosDocente []models.Contrato = nil
		var contratoPreliquidacionDocente []models.ContratoPreliquidacion = nil
		var contratosCambio []int
		query := "Documento:" + contrato.Documento + ",TipoNominaId:" + strconv.Itoa(contrato.TipoNominaId) + ",Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contratosDocente)
			if contratosDocente[0].Id != 0 {
				for i := 0; i < len(contratosDocente); i++ {
					query = "ContratoId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mes_aux) + ",PreliquidacionId.Ano:" + strconv.Itoa(anio_aux)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						contratoPreliquidacionDocente = nil
						LimpiezaRespuestaRefactor(aux, &contratoPreliquidacionDocente)
						if contratoPreliquidacionDocente[0].Id != 0 {
							if contratosDocente[i].NumeroContrato != "GENERAL"+strconv.Itoa(mes_aux) {
								if !strings.HasPrefix(contratosDocente[i].NumeroContrato, "GENERAL") {
									contratosCambio = append(contratosCambio, contratoPreliquidacionDocente[0].Id)
								}
							}
						} else {
							fmt.Println("No se encontraron preliquidaciones asociadas al contrato: ", contratosDocente[i].NumeroContrato)
						}
					} else {
						fmt.Println("Error al obtener el contrato preliquidación para el contrato: ", contratosDocente[i].NumeroContrato)
					}
				}
			} else {
				fmt.Println("El docente no tiene contratos registrados")
			}
		} else {
			fmt.Println("Error al intentar obtener contratos del docente: ", err)
		}
		if contrato.TipoNominaId == 409 {
			query := "Ano:" + strconv.Itoa(anio_aux) + ",Mes:" + strconv.Itoa(mes_aux) + ",Nominaid:415"
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &preliquidacion)
				LiquidarContratoGeneral(mes_aux, anio_aux, contrato, preliquidacion[0], 1, "409", vigencia_original, false)
				borrarContratoGeneral(contrato.Documento, contrato.Vigencia, fecha_anulacion, contrato.FechaFin, contrato.TipoNominaId)
				cambioContrato(true, contrato, mes_aux, contratosCambio)
			}
		} else {
			query := "Ano:" + strconv.Itoa(anio_aux) + ",Mes:" + strconv.Itoa(mes_aux) + ",Nominaid:416"
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &preliquidacion)
				LiquidarContratoGeneral(mes_aux, anio_aux, contrato, preliquidacion[0], 1, "410", vigencia_original, false)
				borrarContratoGeneral(contrato.Documento, contrato.Vigencia, fecha_anulacion, contrato.FechaFin, contrato.TipoNominaId)
			}
		}
		mes_aux += 1
		if mes_aux > 12 {
			mes_aux = 1
			anio_aux += 1
		}
		if anio_aux == int(contrato.FechaFin.Year()) && mes_aux > int(contrato.FechaFin.Month()) {
			break
		}
	}

}

func borrarContratoGeneral(Documento string, Vigencia int, FechaFin time.Time, fechaFinOriginal time.Time, TipoNominaId int) {
	var aux map[string]interface{}
	var contratosDocente []models.Contrato = nil
	var contratos []models.Contrato = nil
	var borrar bool = false
	query := "Documento:" + Documento + ",TipoNominaId:" + strconv.Itoa(TipoNominaId) + ",Vigencia:" + strconv.Itoa(Vigencia) + ",Activo:true"
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratosDocente)
		// Busca si tiene mas de un contrato que no sea general
		for i := 0; i < len(contratosDocente); i++ {
			if !strings.HasPrefix(contratosDocente[i].NumeroContrato, "GENERAL") {
				contratos = append(contratos, contratosDocente[i])
			}
		}
		for j := int(FechaFin.Month()); j <= int(fechaFinOriginal.Month()); j++ {
			for i := 0; i < len(contratos); i++ {
				if int(contratos[i].FechaInicio.Month()) <= j && int(contratos[i].FechaFin.Month()) >= j {
					borrar = false
					break
				} else {
					borrar = true
				}
			}
			if borrar {
				var id int
				for i := 0; i < len(contratosDocente); i++ {
					if strings.HasPrefix(contratosDocente[i].NumeroContrato, "GENERAL"+strconv.Itoa(int(j))) {
						id = contratosDocente[i].Id
						break
					}
				}
				if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(id), "DELETE", &aux, nil); err == nil {
					fmt.Println("Contrato general eliminado con éxito ", int(j))
				} else {
					fmt.Println("Error al Eliminar Contrato General:", err)
				}
				borrar = false
			}
		}
	}
}
