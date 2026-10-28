package controllers

import (
	"context"
	"encoding/json"
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
		// dia inclusivo
		d += 1
		if d == 22 {
			d += 1
		}
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
	}
	if mesesContrato/float64(int(mesesContrato)) != 1 {
		return (mesesContrato * 4) + 1
	} else {
		return (mesesContrato * 4)
	}
}

func registrarPreliquidacion(ctx context.Context, año, mes, estadoPreliquidacion, nomina int) (preliquidacion models.Preliquidacion) {
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
	if _, err := request.PostWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion", preliquidacion, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &preliquidacion)
	} else {
		logs.Error("Error al guardar preliquidacion: ", err)
	}

	return preliquidacion
}

func registrarContratoPreliquidacion(ctx context.Context, preliquidacionId int, contratoId int, contratoPreliq models.ContratoPreliquidacion) (contratoPreliquidacion models.ContratoPreliquidacion) {
	var aux map[string]interface{}
	var respuesta models.ContratoPreliquidacion
	contratoPreliq.Id = 0
	contratoPreliquidacion = contratoPreliq
	contratoPreliquidacion.ContratoId = &models.Contrato{Id: contratoId}
	contratoPreliquidacion.PreliquidacionId = &models.Preliquidacion{Id: preliquidacionId}
	contratoPreliquidacion.Cumplido = false
	contratoPreliquidacion.Preliquidado = false
	contratoPreliquidacion.Activo = true

	if _, err := request.PostWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion", contratoPreliquidacion, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &respuesta)
	} else {
		logs.Error("Error al guardar contrato_preliquidacion: ", err)
	}
	return respuesta
}

func registrarContratoPreliquidacionOld(ctx context.Context, preliquidacionId int, contratoId int, contratoPreliq models.ContratoPreliquidacionOld) (contratoPreliquidacion models.ContratoPreliquidacionOld) {
	var aux map[string]interface{}
	var respuesta models.ContratoPreliquidacionOld
	contratoPreliq.Id = 0
	contratoPreliquidacion = contratoPreliq
	contratoPreliquidacion.ContratoId = &models.ContratoOld{Id: contratoId}
	contratoPreliquidacion.PreliquidacionId = &models.Preliquidacion{Id: preliquidacionId}
	contratoPreliquidacion.Cumplido = false
	contratoPreliquidacion.Preliquidado = false
	contratoPreliquidacion.Activo = true

	if _, err := request.PostWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion", contratoPreliquidacion, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &respuesta)
	} else {
		logs.Error("Error al guardar contrato_preliquidacion: ", err)
	}
	return respuesta
}

func registrarDetallePreliquidacion(ctx context.Context, detallePreliquidacion models.DetallePreliquidacion) {
	var response models.DetallePreliquidacion
	if _, err := request.PostWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion", detallePreliquidacion, &response); err == nil {
	} else {
		logs.Error("Error al guardar detalle", err)
	}
}

func registrarDetallePreliquidacionOld(ctx context.Context, detallePreliquidacion models.DetallePreliquidacionOld) {
	var response models.DetallePreliquidacion
	if _, err := request.PostWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion", detallePreliquidacion, &response); err == nil {
	} else {
		logs.Error("Error al guardar detalle", err)
	}
}

func registrarContrato(ctx context.Context, contrato models.Contrato) (respuesta models.Contrato, err error) {
	var response map[string]interface{}
	if _, err := request.PostWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato", contrato, &response); err == nil {
		LimpiezaRespuestaRefactor(response, &respuesta)
	} else {
		logs.Error("Error al guardar contrato", err)
		return contrato, err
	}
	return respuesta, nil
}

func registrarContratoOld(ctx context.Context, contrato models.ContratoOld) (respuesta models.ContratoOld, err error) {
	var response map[string]interface{}
	if _, err := request.PostWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato", contrato, &response); err == nil {
		LimpiezaRespuestaRefactor(response, &respuesta)
	} else {
		logs.Error("Error al guardar contrato", err)
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

func LiquidarContratoGeneral(ctx context.Context, mesIterativo int, anoIterativo int, contrato models.Contrato, preliquidacion models.Preliquidacion, porcentaje float64, nomina string, vigencia_original int, crear bool) {
	var aux map[string]interface{}
	var contratoGeneral []models.Contrato
	var contratosDocente []models.ContratoPreliquidacion
	var auxDetalle []models.DetallePreliquidacion
	var flag bool = true

	query := "NumeroContrato:GENERAL" + strconv.Itoa(mesIterativo) + ",Vigencia:" + strconv.Itoa(anoIterativo) + ",Documento:" + contrato.Documento + ",TipoNominaId:" + nomina + ",Activo:true"
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
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
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
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
						if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &auxDetalle)
							valorContratoGeneralAux += int(auxDetalle[0].ValorCalculado)
						} else {
							logs.Error("Error al obtener los honorarios para el contrato :", contratosDocente[i].ContratoId.NumeroContrato, " ", err)
						}
					}
					contratoGeneral[0].ValorContrato = float64(valorContratoGeneralAux)
				}
			} else {
				logs.Error("Error al obtener los contratos vigentes para el mes actual: ", err)
			}

			//Registrar el contrato nuevo
			if crear {
				contratoGeneral[0], _ = registrarContrato(ctx, contratoGeneral[0])
			}
		} else {
			//Eliminar los detalles del contrato General
			query := "ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contratoGeneral[0].Id) + ",ContratoPreliquidacionId.ContratoId.Vigencia:" + strconv.Itoa(anoIterativo)
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &auxDetalle)
				if auxDetalle[0].Id != 0 {
					idContratoPeliquidacion := auxDetalle[0].ContratoPreliquidacionId.Id
					for j := 0; j < len(auxDetalle); j++ {
						if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(auxDetalle[j].Id), &aux); err == nil {
						} else {
							logs.Error("Error al eliminar detalle: ", err)
						}
					}
					//Eliminar el contrato_preliquidación
					if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(idContratoPeliquidacion), &aux); err == nil {
						//Actualizar el valor del contrato general
						//Buscar el valor de los honorarios de los contratos que tiene el docente en ese mes
						contratoGeneral[0].ValorContrato = 0
						query = "PreliquidacionId.Id:" + strconv.Itoa(preliquidacion.Id) + ",ContratoId.Documento:" + contrato.Documento + ",ContratoId.Activo:true"
						if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
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
										if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
											LimpiezaRespuestaRefactor(aux, &auxDetalle)
											valorContratoGeneralAux += int(auxDetalle[0].ValorCalculado)
										} else {
											logs.Error("Error al obtener los honorarios para el contrato :", contratosDocente[i].ContratoId.NumeroContrato, " ", err)
										}
									}
								}
								contratoGeneral[0].ValorContrato = float64(valorContratoGeneralAux)
								//Actualizar

								if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contratoGeneral[0].Id), contratoGeneral[0], &aux); err == nil {
								} else {
									logs.Error("Error al actualizar valor del contrato")
								}
							} else {
								flag = false
							}
						} else {
							logs.Error("Error al obtener los contratos vigentes para el mes actual: ", err)
						}
					} else {
						logs.Error("Error al eliminar contrato_preliquidacion: ", err)
					}
				}
			} else {
				logs.Error("Error al obtener los detalles para el contrato general del mes")
			}
		}

		if nomina == "410" && flag {
			liquidarHCS(ctx, contratoGeneral[0], true, porcentaje, vigencia_original, 0, 0, false)
		} else if nomina == "409" && flag {
			liquidarHCH(ctx, contratoGeneral[0], true, porcentaje, vigencia_original, 0, 0, false)
		}
	} else {
		logs.Error("Error al buscar contrato general:", err)
	}
}

func LiquidarContratoGeneralOld(ctx context.Context, mesIterativo int, anoIterativo int, contrato models.ContratoOld, preliquidacion models.Preliquidacion, porcentaje float64, nomina string, vigencia_original int, crear bool) {
	var aux map[string]interface{}
	var contratoGeneral []models.ContratoOld
	var contratosDocente []models.ContratoPreliquidacionOld
	var auxDetalle []models.DetallePreliquidacionOld
	var flag bool = true

	query := "NumeroContrato:GENERAL" + strconv.Itoa(mesIterativo) + ",Vigencia:" + strconv.Itoa(anoIterativo) + ",Documento:" + contrato.Documento + ",TipoNominaId:" + nomina + ",Activo:true"
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
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

			query = "PreliquidacionId.Id:" + strconv.Itoa(preliquidacion.Id) + ",ContratoId.Documento:" + contrato.Documento + ",ContratoId.TipoNominaId:" + nomina
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contratosDocente)
				if len(contratosDocente) >= 1 { //Tiene más de dos contratos
					//Sumar valores de los honorarios para obtener el valor total de ese mes
					contratoGeneral[0].ValorContrato = 0
					for i := 0; i < len(contratosDocente); i++ {
						//Sumar los honorarios de el mes presente para obtener el IBC
						if nomina == "410" {
							query = "ContratoPreliquidacionId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",ConceptoNominaId.Id:152"

						} else {
							query = "ContratoPreliquidacionId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",ConceptoNominaId.Id:87"
						}
						if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &auxDetalle)
							contratoGeneral[0].ValorContrato = contratoGeneral[0].ValorContrato + auxDetalle[0].ValorCalculado
						} else {
							logs.Error("Error al obtener los honorarios para el contrato :", contratosDocente[i].ContratoId.NumeroContrato, " ", err)
						}
					}
				}
			} else {
				logs.Error("Error al obtener los contratos vigentes para el mes actual: ", err)
			}

			//Registrar el contrato nuevo
			if crear {
				contratoGeneral[0], _ = registrarContratoOld(ctx, contratoGeneral[0])
			}
		} else {
			//Eliminar los detalles del contrato General
			query := "ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contratoGeneral[0].Id) + ",ContratoPreliquidacionId.ContratoId.Vigencia:" + strconv.Itoa(anoIterativo)
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &auxDetalle)
				idContratoPeliquidacion := auxDetalle[0].ContratoPreliquidacionId.Id
				for j := 0; j < len(auxDetalle); j++ {
					if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(auxDetalle[j].Id), &aux); err == nil {
					} else {
						logs.Error("Error al eliminar detalle: ", err)
					}
				}
				//Eliminar el contrato_preliquidación
				if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(idContratoPeliquidacion), &aux); err == nil {
					//Actualizar el valor del contrato general
					//Buscar el valor de los honorarios de los contratos que tiene el docente en ese mes
					contratoGeneral[0].ValorContrato = 0
					query = "PreliquidacionId.Id:" + strconv.Itoa(preliquidacion.Id) + ",ContratoId.Documento:" + contrato.Documento
					if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &contratosDocente)
						if len(contratosDocente) >= 1 && contratosDocente[0].Id != 0 { //Tiene más de un contrato
							//Sumar valores de los honorarios para obtener el valor total de ese mes
							for i := 0; i < len(contratosDocente); i++ {
								//Sumar los honorarios de el mes presente para obtener el IBC
								if contratosDocente[i].ContratoId.Id != contratoGeneral[0].Id {
									if nomina == "410" {
										query = "ContratoPreliquidacionId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",ConceptoNominaId.Id:152"

									} else {
										query = "ContratoPreliquidacionId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",ConceptoNominaId.Id:87"
									}
									if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
										LimpiezaRespuestaRefactor(aux, &auxDetalle)
										contratoGeneral[0].ValorContrato = contratoGeneral[0].ValorContrato + auxDetalle[0].ValorCalculado
									} else {
										logs.Error("Error al obtener los honorarios para el contrato :", contratosDocente[i].ContratoId.NumeroContrato, " ", err)
									}
								}
							}
							//Actualizar

							if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contratoGeneral[0].Id), contratoGeneral[0], &aux); err == nil {
							} else {
								logs.Error("Error al actualizar valor del contrato")
							}
						} else {
							flag = false
						}
					} else {
						logs.Error("Error al obtener los contratos vigentes para el mes actual: ", err)
					}
				} else {
					logs.Error("Error al eliminar contrato_preliquidacion: ", err)
				}
			} else {
				logs.Error("Error al obtener los detalles para el contrato general del mes")
			}
		}

		if nomina == "410" && flag {
			liquidarHCSOld(ctx, contratoGeneral[0], true, porcentaje, vigencia_original, 0, 0, false)
		} else if nomina == "409" && flag {
			liquidarHCHOld(ctx, contratoGeneral[0], true, porcentaje, vigencia_original)
		}
	} else {
		logs.Error("Error al buscar contrato general:", err)
	}
}

func anularEnGenerales(ctx context.Context, contrato models.Contrato, fecha_anulacion time.Time, vigencia_original int, reduccionTotal bool) {

	var aux map[string]interface{}
	var preliquidacion []models.Preliquidacion

	anio_aux := int(fecha_anulacion.Year())
	mes_aux := int(fecha_anulacion.Month()) + 1
	count := 1

	for {
		var contratosDocente []models.Contrato = nil
		var contratoPreliquidacionDocente []models.ContratoPreliquidacion = nil
		var contratosCambio []int
		query := "Documento:" + contrato.Documento + ",TipoNominaId:" + strconv.Itoa(contrato.TipoNominaId) + ",Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
		if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contratosDocente)
			if contratosDocente[0].Id != 0 {
				for i := 0; i < len(contratosDocente); i++ {
					query = "ContratoId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mes_aux) + ",PreliquidacionId.Ano:" + strconv.Itoa(anio_aux)
					if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						contratoPreliquidacionDocente = nil
						LimpiezaRespuestaRefactor(aux, &contratoPreliquidacionDocente)
						if contratoPreliquidacionDocente[0].Id != 0 {
							if contratosDocente[i].NumeroContrato != "GENERAL"+strconv.Itoa(mes_aux) {
								if !strings.HasPrefix(contratosDocente[i].NumeroContrato, "GENERAL") {
									contratosCambio = append(contratosCambio, contratoPreliquidacionDocente[0].Id)
								}
							}
						} else {
						}
					} else {
						logs.Error("Error al obtener el contrato preliquidación para el contrato: ", contratosDocente[i].NumeroContrato)
					}
				}
			} else {
			}
		} else {
			logs.Error("Error al intentar obtener contratos del docente: ", err)
		}
		if contrato.TipoNominaId == 409 {
			query := "Ano:" + strconv.Itoa(anio_aux) + ",Mes:" + strconv.Itoa(mes_aux) + ",Nominaid:415"
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &preliquidacion)
				LiquidarContratoGeneral(ctx, mes_aux, anio_aux, contrato, preliquidacion[0], 1, "409", vigencia_original, false)
				if reduccionTotal {
					borrarContratoGeneral(ctx, contrato.Documento, contrato.Vigencia, fecha_anulacion.AddDate(0, count, 0), contrato.FechaFin, contrato.TipoNominaId)
				} else {
					borrarContratoGeneral(ctx, contrato.Documento, contrato.Vigencia, fecha_anulacion, contrato.FechaFin, contrato.TipoNominaId)
				}
				cambioContrato(ctx, true, contrato, mes_aux, contratosCambio)
			}
		} else {
			query := "Ano:" + strconv.Itoa(anio_aux) + ",Mes:" + strconv.Itoa(mes_aux) + ",Nominaid:416"
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &preliquidacion)
				LiquidarContratoGeneral(ctx, mes_aux, anio_aux, contrato, preliquidacion[0], 1, "410", vigencia_original, false)
				if reduccionTotal {
					borrarContratoGeneral(ctx, contrato.Documento, contrato.Vigencia, fecha_anulacion.AddDate(0, count, 0), contrato.FechaFin, contrato.TipoNominaId)
				} else {
					borrarContratoGeneral(ctx, contrato.Documento, contrato.Vigencia, fecha_anulacion, contrato.FechaFin, contrato.TipoNominaId)
				}
			}
		}
		mes_aux += 1
		count += 1
		if mes_aux > 12 {
			mes_aux = 1
			anio_aux += 1
		}
		if anio_aux > int(contrato.FechaFin.Year()) {
			break
		} else if anio_aux == int(contrato.FechaFin.Year()) {
			if mes_aux > int(contrato.FechaFin.Month()) {
				break
			}
		}

	}

}

func anularEnGeneralesOld(ctx context.Context, contrato models.ContratoOld, fecha_anulacion time.Time, vigencia_original int) {

	var aux map[string]interface{}
	var preliquidacion []models.Preliquidacion

	anio_aux := int(fecha_anulacion.Year())
	mes_aux := int(fecha_anulacion.Month()) + 1

	for {
		var contratosDocente []models.ContratoOld = nil
		var contratoPreliquidacionDocente []models.ContratoPreliquidacionOld = nil
		var contratosCambio []int
		query := "Documento:" + contrato.Documento + ",TipoNominaId:" + strconv.Itoa(contrato.TipoNominaId) + ",Vigencia:" + strconv.Itoa(contrato.Vigencia) + ",Activo:true"
		if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contratosDocente)
			if contratosDocente[0].Id != 0 {
				for i := 0; i < len(contratosDocente); i++ {
					query = "ContratoId.Id:" + strconv.Itoa(contratosDocente[i].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mes_aux) + ",PreliquidacionId.Ano:" + strconv.Itoa(anio_aux)
					if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						contratoPreliquidacionDocente = nil
						LimpiezaRespuestaRefactor(aux, &contratoPreliquidacionDocente)
						if contratoPreliquidacionDocente[0].Id != 0 {
							if contratosDocente[i].NumeroContrato != "GENERAL"+strconv.Itoa(mes_aux) {
								if !strings.HasPrefix(contratosDocente[i].NumeroContrato, "GENERAL") {
									contratosCambio = append(contratosCambio, contratoPreliquidacionDocente[0].Id)
								}
							}
						} else {
						}
					} else {
						logs.Error("Error al obtener el contrato preliquidación para el contrato: ", contratosDocente[i].NumeroContrato)
					}
				}
			} else {
			}
		} else {
			logs.Error("Error al intentar obtener contratos del docente: ", err)
		}
		if contrato.TipoNominaId == 409 {
			query := "Ano:" + strconv.Itoa(anio_aux) + ",Mes:" + strconv.Itoa(mes_aux) + ",Nominaid:415"
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &preliquidacion)
				LiquidarContratoGeneralOld(ctx, mes_aux, anio_aux, contrato, preliquidacion[0], 1, "409", vigencia_original, false)
				borrarContratoGeneralOld(ctx, contrato.Documento, contrato.Vigencia, fecha_anulacion, contrato.FechaFin, contrato.TipoNominaId)
				cambioContratoOld(ctx, true, contrato, mes_aux, contratosCambio)
			}
		} else {
			query := "Ano:" + strconv.Itoa(anio_aux) + ",Mes:" + strconv.Itoa(mes_aux) + ",Nominaid:416"
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &preliquidacion)
				LiquidarContratoGeneralOld(ctx, mes_aux, anio_aux, contrato, preliquidacion[0], 1, "410", vigencia_original, false)
				borrarContratoGeneralOld(ctx, contrato.Documento, contrato.Vigencia, fecha_anulacion, contrato.FechaFin, contrato.TipoNominaId)
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

func borrarContratoGeneral(ctx context.Context, Documento string, Vigencia int, FechaFin time.Time, fechaFinOriginal time.Time, TipoNominaId int) {
	var aux map[string]interface{}
	var contratosDocente []models.Contrato = nil
	var contratos []models.Contrato = nil
	var borrar bool = true
	query := "Documento:" + Documento + ",TipoNominaId:" + strconv.Itoa(TipoNominaId) + ",Vigencia:" + strconv.Itoa(Vigencia) + ",Activo:true"
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratosDocente)
		// Busca si tiene mas de un contrato que no sea general
		for i := 0; i < len(contratosDocente); i++ {
			if !strings.HasPrefix(contratosDocente[i].NumeroContrato, "GENERAL") {
				contratos = append(contratos, contratosDocente[i])
			}
		}
		for j := int(FechaFin.Month()); j <= int(fechaFinOriginal.Month()); j++ {
			for i := 0; i < len(contratos); i++ {
				if int(contratos[i].FechaInicio.Month()) <= j && int(contratos[i].FechaFin.Month()) >= j && contratos[i].ValorContrato != 0 {
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
				if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(id), &aux); err == nil {
				} else {
					logs.Error("Error al Eliminar Contrato General:", err)
				}
				borrar = false
			}
		}
	}
}

func borrarContratoGeneralOld(ctx context.Context, Documento string, Vigencia int, FechaFin time.Time, fechaFinOriginal time.Time, TipoNominaId int) {
	var aux map[string]interface{}
	var contratosDocente []models.ContratoOld = nil
	var contratos []models.ContratoOld = nil
	var borrar bool = false
	query := "Documento:" + Documento + ",TipoNominaId:" + strconv.Itoa(TipoNominaId) + ",Vigencia:" + strconv.Itoa(Vigencia) + ",Activo:true"
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
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
				if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(id), &aux); err == nil {
				} else {
					logs.Error("Error al Eliminar Contrato General:", err)
				}
				borrar = false
			}
		}
	}
}

func Preliquidacion(ctx context.Context, contrato models.Contrato) (mensaje string, codigo string, contratoReturn *models.Contrato, err error) {
	var aux map[string]interface{}

	if contrato.FechaInicio.Before(contrato.FechaFin) {
		if _, err := request.PostWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato", contrato, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato)
			if contrato.TipoNominaId == 411 {
				mensaje, err = liquidarCPS(ctx, contrato)
				if err == nil {
					return "Successful", "201", &contrato, nil
				} else {
					codigo = "404"
					// Titan CRUD deletes by path ID and returns a JSON success envelope; the rollback body is unused.
					if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato.Id), &aux); err == nil {
					} else {
						mensaje = "Error al eliminar el contrato que no se liquidó: "
						codigo = "404"
						return mensaje, codigo, nil, err
					}
					return mensaje, codigo, nil, err
				}
			} else if contrato.TipoNominaId == 409 {
				mensaje, err = liquidarHCH(ctx, contrato, false, 0, contrato.Vigencia, 0, 0, false)
				if err == nil {
					return "Successful", "201", &contrato, nil
				} else {
					codigo = "404"
					if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato.Id), &aux); err == nil {
					} else {
						mensaje = "Error al eliminar el contrato que no se liquidó: "
						codigo = "404"
						return mensaje, codigo, nil, err
					}
					return mensaje, codigo, nil, err
				}
			} else if contrato.TipoNominaId == 410 {
				mensaje, err = liquidarHCS(ctx, contrato, false, 0, contrato.Vigencia, 0, 0, false)
				var contratoActualizado models.Contrato
				if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato.Id), &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoActualizado)
					contratoReturn = &contratoActualizado
					contratoReturn.Desagregado = contrato.Desagregado
					contratoReturn.NumeroSemanas = contrato.NumeroSemanas
					contratoReturn.Completo = contrato.Completo
					contratoReturn.Unico = contrato.Unico
				}
				if err == nil {
					return "Successful", "201", contratoReturn, nil
				} else {
					codigo = "404"
					if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato.Id), &aux); err == nil {
					} else {
						mensaje = "Error al eliminar el contrato que no se liquidó: "
						codigo = "404"
						return mensaje, codigo, nil, err
					}
					return mensaje, codigo, nil, err
				}
			}
		} else {
			logs.Error("No se pudo guardar el contrato", err)
			mensaje = "No se pudo guardar el contrato"
			codigo = "404"
			return mensaje, codigo, nil, err
		}
	} else {
		mensaje = "La fecha inicio no puede estar después de la fecha fin"
		codigo = "404"
		return mensaje, codigo, nil, err
	}
	return

}

func Anulacion(ctx context.Context, anulacion models.Anulacion, valorContrato float64, semanas int, semanasAnteriores int, anulacionReduccion bool) (mensaje string, codigo string, contratoReturn *models.Contrato, err error, fechaOriginal time.Time, completo bool) {

	var aux map[string]interface{}
	var contrato []models.Contrato
	var contratoOriginal models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var valorDia float64
	var detalles []models.DetallePreliquidacion
	var semanasContrato int
	var semanasTotales int
	var conceptoNominaId int
	var mismoMes bool = false
	var codAux int
	var DetallesAux []models.DetallePreliquidacion
	var sumaContratosTemp float64
	var anulacion_completa bool
	var valorNuevo float64
	var semanasAnulacion int
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+anulacion.NumeroContrato+",Vigencia:"+strconv.Itoa(anulacion.Vigencia)+",Documento:"+anulacion.Documento+",Activo:true", &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contrato)
		if contrato[0].Id != 0 {

			//Ordenar los contratos para tomar el más reciente
			for i := 0; i < len(contrato); i++ {
				if contrato[0].Id < contrato[i].Id {
					auxContrato := contrato[0]
					contrato[0] = contrato[i]
					contrato[i] = auxContrato
				}
			}

			contratoOriginal = contrato[0]
			contrato[0].Desagregado = anulacion.Desagregado
			semanasAnulacion = semanas
			if anulacion.FechaAnulacion.Equal(contrato[0].FechaInicio) {
				anulacion.FechaAnulacion = contrato[0].FechaInicio
				contrato[0].ValorContrato = 0
				anulacion_completa = true
			} else {
				anulacion_completa = false
			}

			anoIterativo := anulacion.FechaAnulacion.Year()
			mesIterativo := int(anulacion.FechaAnulacion.Month())

			//Eliminar los detalles y los contratos_preliquidacion
			for {
				//Obtener contrato_preliquidacion para ese mes
				query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
				if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
					if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &detalles)
						for j := 0; j < len(detalles); j++ {
							if contrato[0].TipoNominaId == 410 {
								if detalles[j].ConceptoNominaId.Id == 152 {
									valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
								}
							} else {
								if detalles[j].ConceptoNominaId.Id == 87 {
									valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
								}
							}
							if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalles[j].Id), &aux); err == nil {
							} else {
								logs.Error("Error al Eliminar Detalles:", err)
								mensaje = "Error al Eliminar Detalles: "
								codigo = "400"
								return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
							}
						}
						if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
						} else {
							logs.Error("Error al eliminar contrato preliquidacion: ", err)
							mensaje = "Error al Eliminar contrato preliquidación: "
							codigo = "400"
							return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
						}
					} else {
						logs.Error("Error al obtener detalles")
						mensaje = "Error al obtener detalles: "
						codigo = "400"
						return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
					}
				} else {
					logs.Error("Error al obtener contrato_preliquidacion")
					mensaje = "Error al obtener contrato_preliquidacion"
					codigo = "400"
					return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
				}
				if mesIterativo == int(contrato[0].FechaFin.Month()) && anoIterativo == contrato[0].FechaFin.Year() {
					break
				} else {
					if mesIterativo == 12 {
						mesIterativo = 1
						anoIterativo = anoIterativo + 1
					} else {
						mesIterativo = mesIterativo + 1
					}
				}
			}
			contrato[0].FechaFin = anulacion.FechaAnulacion
			//Actualizar fecha de finalización del contrato
			contratoOriginal.Activo = false
			contratoOriginal.Id = 0

			if contrato[0].TipoNominaId == 409 {
				conceptoNominaId = 87
			} else if contrato[0].TipoNominaId == 410 {
				conceptoNominaId = 152
			}

			query := "ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contrato[0].Id) + ",ConceptoNominaId.Id:" + strconv.Itoa(conceptoNominaId)
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &detalles)
				for j := 0; j < len(detalles); j++ {
					valorNuevo = valorNuevo + detalles[j].ValorCalculado
				}
			}
			// CREA CONTRATO NUEVO CON LA INFORMACIÓN DEL CONTRATO ORIGINAL CON CAMPO ACTIVO FALSE
			if _, err := request.PostWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato", contratoOriginal, &aux); err == nil {
				fechaOriginal = contratoOriginal.FechaFin
				if anulacionReduccion && semanasAnteriores == semanasAnulacion {
					anulacion_completa = true
				}
				if contrato[0].FechaInicio.Month() != anulacion.FechaAnulacion.Month() || contrato[0].FechaInicio.Year() != anulacion.FechaAnulacion.Year() {
					//semanasTotales = int(calcularSemanasContratoDVE(contrato[0].FechaInicio, anulacion.FechaAnulacion))
					var semanasTranscurridas = 0
					for i := int(contratoOriginal.FechaInicio.Month()); i < int(anulacion.FechaAnulacion.Month()); i++ {
						var contrato_preliquidacion_aux []models.ContratoPreliquidacion
						var detalle_preliquidacion_aux []models.DetallePreliquidacion
						if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?query=ContratoId:"+strconv.Itoa(contrato[0].Id)+",PreliquidacionId__Mes:"+strconv.Itoa(i)+",PreliquidacionId__Ano:"+strconv.Itoa(contratoOriginal.FechaInicio.Year()), &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion_aux)
							if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion_aux[0].Id), &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &detalle_preliquidacion_aux)
								semanasTranscurridas += int(detalle_preliquidacion_aux[0].DiasLiquidados)
							}
						}
					}
					if semanasAnulacion == 0 {
						semanasContrato = int(calcularSemanasContratoDVE(time.Date(anulacion.FechaAnulacion.Year(), anulacion.FechaAnulacion.Month(), 1, 12, 0, 0, 0, time.UTC), anulacion.FechaAnulacion))
					} else {
						semanasContrato = semanasAnulacion - semanasTranscurridas
					}
					semanasTotales = semanasContrato
					contrato[0].FechaInicio = time.Date(anulacion.FechaAnulacion.Year(), anulacion.FechaAnulacion.Month(), 1, 12, 0, 0, 0, time.UTC)
					if valorContrato == 0 {
						contrato[0].ValorContrato = valorDia * float64(semanasContrato)
					} else {
						contrato[0].ValorContrato = valorContrato
					}
				} else if anulacion_completa {
					semanasContrato = 0
					semanasTotales = 0
					contrato[0].ValorContrato = 0
					contrato[0].Activo = false
				} else {
					mismoMes = true
					diaAux := contrato[0].FechaInicio.AddDate(0, 0, 1)
					semanasContrato = int(calcularSemanasContratoDVE(diaAux, contrato[0].FechaFin))
					/*if !anulacionReduccion {
						semanasContrato -= 1
					}*/
					semanasTotales = semanasContrato
					contrato[0].NumeroSemanas = semanasTotales
					if valorContrato == 0 {
						contrato[0].ValorContrato = valorDia * float64(semanasContrato)
					} else {
						contrato[0].ValorContrato = math.Floor(valorContrato)
					}
				}

				if contrato[0].TipoNominaId == 409 {
					codAux = 87
				} else if contrato[0].TipoNominaId == 410 {
					codAux = 152
				}

				// Trae todos los detalles preliquidación no eliminados para calcular el nuevo valor del contrato
				query := "ContratoPreliquidacionId__ContratoId__Id:" + strconv.Itoa(contrato[0].Id) + ",ConceptoNominaId:" + strconv.Itoa(codAux)
				if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil && valorContrato == 0 {
					LimpiezaRespuestaRefactor(aux, &DetallesAux)
					for i := 0; i < len(DetallesAux); i++ {
						sumaContratosTemp += DetallesAux[i].ValorCalculado
					}
				} else {
					logs.Error("Error al obtener detalles preliquidacion ", err)
					mensaje = "Error al obtener detalles preliquidacion"
					codigo = "400"
				}

				contratoAux := contrato[0]
				contratoAux.FechaInicio = contratoOriginal.FechaInicio
				sumaContratosTemp += contrato[0].ValorContrato
				contratoAux.ValorContrato = sumaContratosTemp
				if mismoMes {
					valorNuevo = 0
				}
				var valorContratoCalculo float64
				if valorContrato == 0 {
					valorContratoCalculo = contratoAux.ValorContrato
				} else {
					valorContratoCalculo = valorContrato
				}

				if semanasTotales > 0 {
					valorDia = (valorContratoCalculo - valorNuevo) / float64(semanasTotales)
				} else {
					valorDia = 0
				}
				// Actualiza los datos del contrato: Fecha fin y valor
				if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), contratoAux, &aux); err == nil {
					if anulacion.FechaAnulacion.Day() != 30 {
						contrato[0].ValorContrato = Roundf(contrato[0].ValorContrato)
						if contrato[0].TipoNominaId == 409 && semanasTotales > 0 {
							// REVISAR VALOR DEL DIA Y SI ES O NO ANULACIÓN
							mensaje, err = liquidarHCH(ctx, contrato[0], false, 0, contrato[0].Vigencia, semanasTotales, valorDia, true)
							anularEnGenerales(ctx, contratoOriginal, anulacion.FechaAnulacion, anulacion.Vigencia, false)
						} else if contrato[0].TipoNominaId == 410 && semanasTotales > 0 {
							mensaje, err = liquidarHCS(ctx, contrato[0], false, 0, contrato[0].Vigencia, semanasTotales, valorDia, true)
							anularEnGenerales(ctx, contratoOriginal, anulacion.FechaAnulacion, anulacion.Vigencia, false)
						} else if contrato[0].TipoNominaId == 410 && semanasTotales == 0 && !anulacion_completa {
							contrato[0].NumeroSemanas = semanasAnulacion
							registrarPrestaciones(ctx, contrato[0])
						}

						if err == nil {
							mensaje = "Registration successful"
							codigo = "201"
							return mensaje, codigo, &contrato[0], nil, fechaOriginal, anulacion_completa
						} else {
							mensaje = "Error al cancelar contrato"
							codigo = "400"
							logs.Error("Error al cancelar contrato: ", err)
						}
					}
				} else {
					logs.Error("Error al crear el contrato nuevo ", err)
					mensaje = "Error al crear el contrato nuevo"
					codigo = "400"
				}

			} else {
				logs.Error("Error al actualizar el contrato: ", err)
				mensaje = "Error al actualizar el contrato"
				codigo = "400"
				return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
			}
		} else {
			logs.Error("Error al obtener el contrato porque no tiene id: ", err)
			mensaje = "Error al obtener el contrato"
			codigo = "400"
			return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
		}
	} else {
		logs.Error("Error al obtener el contrato: ", err)
		mensaje = "Error al obtener el contrato"
		codigo = "400"
		return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
	}
	return
}

func AnulacionPosgrado(ctx context.Context, anulacion models.Anulacion, valorContrato float64, semanas int, semanasAnteriores int, anulacionReduccion bool) (mensaje string, codigo string, contratoReturn *models.Contrato, err error, fechaOriginal time.Time, completo bool) {
	var aux map[string]interface{}
	var contrato []models.Contrato
	var contratoOriginal models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var valorDia float64
	var detalles []models.DetallePreliquidacion
	var semanasContrato int
	var semanasTotales int
	var conceptoNominaId int
	var mismoMes bool = false
	var codAux int
	var DetallesAux []models.DetallePreliquidacion
	var sumaContratosTemp float64
	var anulacion_completa bool
	var valorNuevo float64
	var semanasAnulacion int
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+anulacion.NumeroContrato+",Vigencia:"+strconv.Itoa(anulacion.Vigencia)+",Documento:"+anulacion.Documento+",Activo:true", &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contrato)
		if contrato[0].Id != 0 {

			//Ordenar los contratos para tomar el más reciente
			for i := 0; i < len(contrato); i++ {
				if contrato[0].Id < contrato[i].Id {
					auxContrato := contrato[0]
					contrato[0] = contrato[i]
					contrato[i] = auxContrato
				}
			}

			contratoOriginal = contrato[0]
			contrato[0].Desagregado = anulacion.Desagregado
			if anulacionReduccion {
				semanasAnulacion = semanasAnteriores - semanas
			} else {
				semanasAnulacion = semanas
			}
			if anulacion.FechaAnulacion.Equal(contrato[0].FechaInicio) {
				anulacion.FechaAnulacion = contrato[0].FechaInicio
				contrato[0].ValorContrato = 0
				anulacion_completa = true
			} else {
				anulacion_completa = false
			}

			anoIterativo := int(contratoOriginal.FechaInicio.Year())
			mesIterativo := int(contratoOriginal.FechaInicio.Month())

			//Eliminar los detalles y los contratos_preliquidacion
			for {
				//Obtener contrato_preliquidacion para ese mes
				query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
				if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
					if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &detalles)
						for j := 0; j < len(detalles); j++ {
							if contrato[0].TipoNominaId == 410 {
								if detalles[j].ConceptoNominaId.Id == 152 {
									valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
								}
							} else {
								if detalles[j].ConceptoNominaId.Id == 87 {
									valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
								}
							}
							if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalles[j].Id), &aux); err == nil {
							} else {
								logs.Error("Error al Eliminar Detalles:", err)
								mensaje = "Error al Eliminar Detalles: "
								codigo = "400"
								return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
							}
						}
						if _, err := request.DeleteWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
						} else {
							logs.Error("Error al eliminar contrato preliquidacion: ", err)
							mensaje = "Error al Eliminar contrato preliquidación: "
							codigo = "400"
							return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
						}
					} else {
						logs.Error("Error al obtener detalles")
						mensaje = "Error al obtener detalles: "
						codigo = "400"
						return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
					}
				} else {
					logs.Error("Error al obtener contrato_preliquidacion")
					mensaje = "Error al obtener contrato_preliquidacion"
					codigo = "400"
					return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
				}

				if mesIterativo == int(contrato[0].FechaFin.Month()) && anoIterativo == contrato[0].FechaFin.Year() {
					break
				} else {
					if mesIterativo == 12 {
						mesIterativo = 1
						anoIterativo = anoIterativo + 1
					} else {
						mesIterativo = mesIterativo + 1
					}
				}
			}

			if valorContrato > 0 {
				valorDia = valorContrato / float64(semanasAnulacion)
			}

			contrato[0].FechaFin = anulacion.FechaAnulacion
			//Actualizar fecha de finalización del contrato
			contratoOriginal.Activo = false
			contratoOriginal.Id = 0

			if contrato[0].TipoNominaId == 409 {
				conceptoNominaId = 87
			} else if contrato[0].TipoNominaId == 410 {
				conceptoNominaId = 152
			}

			query := "ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contrato[0].Id) + ",ConceptoNominaId.Id:" + strconv.Itoa(conceptoNominaId)
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &detalles)
				for j := 0; j < len(detalles); j++ {
					valorNuevo = valorNuevo + detalles[j].ValorCalculado
				}
			}
			// CREA CONTRATO NUEVO CON LA INFORMACIÓN DEL CONTRATO ORIGINAL CON CAMPO ACTIVO FALSE
			if _, err := request.PostWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato", contratoOriginal, &aux); err == nil {
				fechaOriginal = contratoOriginal.FechaFin
				if anulacionReduccion && semanasAnulacion == 0 {
					anulacion_completa = true
				}
				if contrato[0].FechaInicio.Month() != anulacion.FechaAnulacion.Month() || contrato[0].FechaInicio.Year() != anulacion.FechaAnulacion.Year() {
					//semanasTotales = int(calcularSemanasContratoDVE(contrato[0].FechaInicio, anulacion.FechaAnulacion))
					if semanas == 0 {
						semanasContrato = int(calcularSemanasContratoDVE(time.Date(anulacion.FechaAnulacion.Year(), anulacion.FechaAnulacion.Month(), 1, 12, 0, 0, 0, time.UTC), anulacion.FechaAnulacion))
					} else {
						semanasContrato = semanasAnulacion
					}
					semanasTotales = semanasContrato
					contrato[0].FechaInicio = contratoOriginal.FechaInicio
					if valorContrato == 0 {
						contrato[0].ValorContrato = valorDia * float64(semanasContrato)
					} else {
						contrato[0].ValorContrato = valorContrato
					}
				} else if anulacion_completa {
					semanasContrato = 0
					semanasTotales = 0
					contrato[0].ValorContrato = 0
					contrato[0].Activo = false
				} else {
					mismoMes = true
					diaAux := contrato[0].FechaInicio.AddDate(0, 0, 1)
					semanasContrato = int(calcularSemanasContratoDVE(diaAux, contrato[0].FechaFin))
					if !anulacionReduccion {
						semanasContrato -= 1
					}
					semanasTotales = semanasContrato
					contrato[0].NumeroSemanas = semanasTotales
					if valorContrato == 0 {
						contrato[0].ValorContrato = valorDia * float64(semanasContrato)
					} else {
						contrato[0].ValorContrato = valorContrato
					}
				}

				if contrato[0].TipoNominaId == 409 {
					codAux = 87
				} else if contrato[0].TipoNominaId == 410 {
					codAux = 152
				}

				// Trae todos los detalles preliquidación no eliminados para calcular el nuevo valor del contrato
				query := "ContratoPreliquidacionId__ContratoId__Id:" + strconv.Itoa(contrato[0].Id) + ",ConceptoNominaId:" + strconv.Itoa(codAux)
				if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil && valorContrato == 0 {
					LimpiezaRespuestaRefactor(aux, &DetallesAux)
					for i := 0; i < len(DetallesAux); i++ {
						sumaContratosTemp += DetallesAux[i].ValorCalculado
					}
				} else {
					logs.Error("Error al obtener detalles preliquidacion ", err)
					mensaje = "Error al obtener detalles preliquidacion"
					codigo = "400"
				}

				contratoAux := contrato[0]
				contratoAux.FechaInicio = contratoOriginal.FechaInicio
				sumaContratosTemp += contrato[0].ValorContrato
				contratoAux.ValorContrato = sumaContratosTemp
				if mismoMes {
					valorNuevo = 0
				}
				// Actualiza los datos del contrato: Fecha fin y valor
				if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), contratoAux, &aux); err == nil {
					if anulacion.FechaAnulacion.Day() != 30 {
						contrato[0].ValorContrato = Roundf(contrato[0].ValorContrato)
						if contrato[0].TipoNominaId == 409 && semanasTotales > 0 {
							// REVISAR VALOR DEL DIA Y SI ES O NO ANULACIÓN
							anularEnGenerales(ctx, contratoOriginal, contratoOriginal.FechaInicio, anulacion.Vigencia, false)
							mensaje, err = liquidarHCH(ctx, contrato[0], false, 0, contrato[0].Vigencia, semanasTotales, valorDia, true)
						} else if contrato[0].TipoNominaId == 410 && semanasTotales > 0 {
							anularEnGenerales(ctx, contratoOriginal, anulacion.FechaAnulacion, anulacion.Vigencia, false)
							contrato[0].NumeroSemanas = semanasAnulacion
							mensaje, err = liquidarHCS(ctx, contrato[0], false, 0, contrato[0].Vigencia, semanasTotales, valorDia, true)
						}

						if err == nil {
							mensaje = "Registration successful"
							codigo = "201"
							return mensaje, codigo, &contrato[0], nil, fechaOriginal, anulacion_completa
						} else {
							mensaje = "Error al cancelar contrato"
							codigo = "400"
							logs.Error("Error al cancelar contrato: ", err)
						}
					}
				} else {
					logs.Error("Error al crear el contrato nuevo ", err)
					mensaje = "Error al crear el contrato nuevo"
					codigo = "400"
				}

			} else {
				logs.Error("Error al actualizar el contrato: ", err)
				mensaje = "Error al actualizar el contrato"
				codigo = "400"
				return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
			}
		} else {
			logs.Error("Error al obtener el contrato porque no tiene id: ", err)
			mensaje = "Error al obtener el contrato"
			codigo = "400"
			return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
		}
	} else {
		logs.Error("Error al obtener el contrato: ", err)
		mensaje = "Error al obtener el contrato"
		codigo = "400"
		return mensaje, codigo, nil, err, fechaOriginal, anulacion_completa
	}
	return
}

func registrarPrestaciones(ctx context.Context, contrato models.Contrato) {
	var predicados []models.Predicado
	var reglasNuevas string = ""
	var reglasAlivios string
	var auxDetalle []models.DetallePreliquidacion
	var aux map[string]interface{}
	//var detalles []models.DetallePreliquidacion
	var detallePreliquidacion []models.DetallePreliquidacion

	//Regla para único o general (para apoximar el ibc al tope mínimo)
	predicados = append(predicados, models.Predicado{Nombre: "general(1)."})
	//Si el contrato es completo se tomarán las vacaciones que calcule
	predicados = append(predicados, models.Predicado{Nombre: "completo(1)."})
	predicados = append(predicados, models.Predicado{Nombre: "vacaciones(0)."})
	valorDia := contrato.ValorContrato / float64(contrato.NumeroSemanas)
	semanas_totales := contrato.NumeroSemanas
	predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + contrato.Documento + "," + fmt.Sprintf("%v", valorDia*float64(semanas_totales)) + "). "})
	predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + contrato.Documento + "," + strconv.Itoa(semanas_totales) + "," + strconv.Itoa(contrato.Vigencia) + "). "})

	cedula, err := strconv.ParseInt(contrato.Documento, 0, 64)

	if err == nil {
		reglasAlivios, _, err = CargarDatosRetefuente(ctx, int(cedula))
	}

	query := "ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contrato.Id) + ",ConceptoNominaId.Id:152"
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &detallePreliquidacion)
	}

	reglasbase := cargarReglasBase(ctx, "HCS") + reglasAlivios + FormatoReglas(predicados)
	reglasNuevas = reglasNuevas + reglasbase + "porcentaje(1).semanas_liquidadas(" + contrato.Documento + "," + strconv.Itoa(contrato.NumeroSemanas) + ")."
	reglasNuevas = reglasNuevas + "mesFinal(1)."
	auxDetalle = golog.LiquidarMesHCS(reglasNuevas, contrato, detallePreliquidacion[len(detallePreliquidacion)-1], true)
	conceptos := []string{"primaNavidad", "cesantias", "priServ", "primaVacaciones", "vacaciones", "interesCesantias", "bonServ"}
	for j := 0; j < len(auxDetalle); j++ {
		if contains(conceptos, auxDetalle[j].ConceptoNominaId.NombreConcepto) {
			registrarDetallePreliquidacion(ctx, auxDetalle[j])
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func CambioCumplido(ctx context.Context, ano string, mes string, numeroContrato string, contrato []models.Contrato) (contrato_general_preliquidacion models.ContratoPreliquidacion, mensaje string, codigo string) {
	var aux map[string]interface{}
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var id int
	cumplidoCompleto := true

	query := "PreliquidacionId.Ano:" + ano + ",PreliquidacionId.Mes:" + mes + ",ContratoId.NumeroContrato:" + numeroContrato
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
		//actualiar cumplido
		contrato_preliquidacion[0].Cumplido = true
		if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), contrato_preliquidacion[0], &aux); err == nil {
			query := "ContratoId.Documento:" + contrato[0].Documento + ",PreliquidacionId.Mes:" + mes + ",PreliquidacionId.Ano:" + ano
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
				if contrato_preliquidacion[0].Id != 0 {
					for i := 0; i < len(contrato_preliquidacion); i++ {
						if contrato_preliquidacion[i].ContratoId.NumeroContrato != "GENERAL"+mes {
							if !contrato_preliquidacion[i].Cumplido {
								cumplidoCompleto = false
							}
						} else {
							id = i
						}
					}
					if cumplidoCompleto {
						//Actualizar el cumplido del contrato General
						contrato_preliquidacion[id].Cumplido = true
						if _, err := request.PutWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[id].Id), contrato_preliquidacion[id], &aux); err == nil {
							contrato_general_preliquidacion = contrato_preliquidacion[0]
							mensaje = "Cumplido actualizado"
							codigo = "200"
							return contrato_general_preliquidacion, mensaje, codigo
						} else {
							logs.Error("Error al actualizar cumplido general: ", err)
							mensaje = "Error al actualizar cumplido general: " + err.Error()
							codigo = "404"
							return contrato_general_preliquidacion, mensaje, codigo
						}
					} else {
						contrato_general_preliquidacion = contrato_preliquidacion[0]
						mensaje = "Cumplido actualizado"
						codigo = "200"
						return contrato_general_preliquidacion, mensaje, codigo
					}
				} else {
					logs.Error("Error al obetener cumplidos: ", err)
					mensaje = "Error al actualizar cumplido general: " + err.Error()
					codigo = "404"
					return contrato_general_preliquidacion, mensaje, codigo
				}
			} else {
				logs.Error("Error al actualizar cumplido: ", err)
				mensaje = "Error al actualizar cumplido general: " + err.Error()
				codigo = "404"
				return contrato_general_preliquidacion, mensaje, codigo
			}
		} else {
			logs.Error("Error al actualizar cumplido: ", err)
			mensaje = "Error al actualizar cumplido general: " + err.Error()
			codigo = "404"
			return contrato_general_preliquidacion, mensaje, codigo
		}

	} else {
		logs.Error("Error al obtener contrato preliquidación")
		mensaje = "Error al actualizar cumplido general: " + err.Error()
		codigo = "404"
		return contrato_general_preliquidacion, mensaje, codigo
	}
}

// Funcionalidad para saber la cantidad de dias de un mes
func daysInMonth(month, year int) int {
	switch time.Month(month) {
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) { // leap year
			return 29
		}
		return 28
	default:
		return 31
	}
}

func ConstruirReglasDesagregado(ctx context.Context, vinculacion models.DatosVinculacion, numSemanas int, contratoOriginal ...models.Contrato) (reglas string, lowCategoria string, lowDedicacion string) {
	var predicados []models.Predicado
	var predicadosPrestaciones []models.Predicado
	if vinculacion.Dedicacion == "HCP" {
		if vinculacion.NivelAcademico == "POSGRADO" {
			lowDedicacion = "hcpos"
		} else {
			lowDedicacion = "hcpre"
		}
		predicados = append(predicados, models.Predicado{Nombre: "aplica_prima(0)."})
	} else {
		if vinculacion.Dedicacion == "MTO" {
			predicados = append(predicados, models.Predicado{Nombre: "aplica_prima(0)."})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "aplica_prima(1)."})
		}
		lowDedicacion = strings.ToLower(vinculacion.Dedicacion)
	}
	lowCategoria = strings.ToLower(vinculacion.Categoria)
	predicados = append(predicados, models.Predicado{Nombre: "horas_semanales(" + strconv.Itoa(vinculacion.HorasSemanales) + ")."})
	predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + vinculacion.Documento + "," + strconv.Itoa(numSemanas) + "," + strconv.Itoa(vinculacion.Vigencia) + ")."})
	predicados = append(predicados, models.Predicado{Nombre: "valor_punto(" + strconv.Itoa(vinculacion.Vigencia) + "," + strconv.Itoa(int(vinculacion.PuntoSalarial)) + ")."})
	if len(contratoOriginal) > 0 {
		switch vinculacion.ObjetoNovedad.TipoResolucion {
		case "RCAN":
			predicadosPrestaciones, _ = ObtenerReglasPrestaciones(ctx, false, contratoOriginal[0])
		case "RADD", "RRED":
			contratoOriginal[0].NumeroSemanas = vinculacion.NumeroSemanas + vinculacion.ObjetoNovedad.SemanasNuevas
			predicadosPrestaciones, _ = ObtenerReglasPrestaciones(ctx, true, contratoOriginal[0])
		default:
			predicadosPrestaciones, _ = ObtenerReglasPrestaciones(ctx, false)
		}
	} else {
		predicadosPrestaciones, _ = ObtenerReglasPrestaciones(ctx, false)
	}
	predicados = append(predicados, predicadosPrestaciones...)
	reglas = cargarReglasBase(ctx, "HCS") + FormatoReglas(predicados)
	return reglas, lowCategoria, lowDedicacion
}

func ObtenerReglasPrestaciones(ctx context.Context, novedad bool, contratoOriginal ...models.Contrato) (predicados []models.Predicado, porcentajesDesagregadoIdNew int) {
	var aux map[string]interface{}
	anoActual := time.Now().Year()
	if len(contratoOriginal) > 0 {
		// en este caso se deben cargar las reglas con el id de parametroPeriodo obtenido (que es para alguna novedad)
		var parametroPeriodo []models.ParametroPeriodo
		query := "id:" + strconv.Itoa(contratoOriginal[0].PorcentajesDesagregadoId)
		if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlParametrosCrud")+"/parametro_periodo?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &parametroPeriodo)
			// Construir reglas dinámicas de porcentaje según los parámetros obtenidos
			porcentajesDesagregadoIdNew = parametroPeriodo[0].Id
			for _, pp := range parametroPeriodo {
				var valores map[string]map[string]float64
				json.Unmarshal([]byte(pp.Valor), &valores)
				for concepto, porcentajes := range valores {
					if novedad {
						// Cuando son adiciones o reducciones
						semanasOriginales := contratoOriginal[0].NumeroSemanas
						predicados = append(predicados, models.Predicado{Nombre: "semanas_contrato_original(" + strconv.Itoa(semanasOriginales) + ")."})
					} else {
						// Cuando son cancelaciones
						predicados = append(predicados, models.Predicado{Nombre: "semanas_contrato_original(0)."})
					}
					if mayor, ok := porcentajes["porcentaje_mayor"]; ok {
						predicados = append(predicados, models.Predicado{Nombre: "porcentaje_mayor(" + strconv.Itoa(anoActual) + "," + strings.ToLower(concepto) + "," + fmt.Sprintf("%.5f", mayor) + ")."})
					}
					if menor, ok := porcentajes["porcentaje_menor"]; ok {
						predicados = append(predicados, models.Predicado{Nombre: "porcentaje_menor(" + strconv.Itoa(anoActual) + "," + strings.ToLower(concepto) + "," + fmt.Sprintf("%.5f", menor) + ")."})
					}
				}
			}
		} else {
			logs.Error("Error al obtener parametro", err)
		}

	} else {
		// se debe obtener desde parametros los valores de porcentaje de prestaciones y cargar los predicados dinamicos
		// obtener el periodo vigente para app de resoluciones
		// contemplar agregar el aplicacion_id para crear periodos exclusivos para resoluciones
		var periodo []models.Periodo
		predicados = append(predicados, models.Predicado{Nombre: "semanas_contrato_original(0)."})
		query := "year:" + strconv.Itoa(anoActual) + ",codigo_abreviacion:PAR,aplicacion_id:30,activo:true"
		if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlParametrosCrud")+"/periodo?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &periodo)
			// obtener el id de parametro de porcentajes de prestaciones
			var parametro []models.Parametro
			query2 := "codigo_abreviacion:PDVE,activo:true"
			if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlParametrosCrud")+"/parametro?limit=-1&query="+query2, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &parametro)
				// finalmente obtener los valores de parametro_periodo
				var parametroPeriodo []models.ParametroPeriodo
				query3 := "parametro_id:" + strconv.Itoa(parametro[0].Id) + ",periodo_id:" + strconv.Itoa(periodo[0].Id) + ",activo:true"
				if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlParametrosCrud")+"/parametro_periodo?limit=-1&query="+query3, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &parametroPeriodo)
					// Construir reglas dinámicas de porcentaje según los parámetros obtenidos
					porcentajesDesagregadoIdNew = parametroPeriodo[0].Id
					for _, pp := range parametroPeriodo {
						var valores map[string]map[string]float64
						json.Unmarshal([]byte(pp.Valor), &valores)
						for concepto, porcentajes := range valores {
							if mayor, ok := porcentajes["porcentaje_mayor"]; ok {
								predicados = append(predicados, models.Predicado{Nombre: "porcentaje_mayor(" + strconv.Itoa(anoActual) + "," + strings.ToLower(concepto) + "," + fmt.Sprintf("%.5f", mayor) + ")."})
							}
							if menor, ok := porcentajes["porcentaje_menor"]; ok {
								predicados = append(predicados, models.Predicado{Nombre: "porcentaje_menor(" + strconv.Itoa(anoActual) + "," + strings.ToLower(concepto) + "," + fmt.Sprintf("%.5f", menor) + ")."})
							}
						}
					}
				} else {
					logs.Error("Error al obtener parametro_periodo", err)
				}
			} else {
				logs.Error("Error al obtener parametro", err)
			}
		} else {
			logs.Error("Error al obtener periodo", err)
		}
	}
	return predicados, porcentajesDesagregadoIdNew
}
