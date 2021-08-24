package controllers

import (
	"strconv"

	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"

	//"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
)

// PreliquidacionFpController operations for PreliquidacionHc
type PreliquidacionFpController struct {
	beego.Controller
}

// Preliquidar ...
// @Title Preliquidar
// @Description Preliquidacion para administrativos
func (c *PreliquidacionFpController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {

	//declaracion de variables
	var reglasinyectadas string
	var reglas string
	var idDetaPre interface{}
	var resumenPreliqu []models.Respuesta

	var porcentajePT int
	var tipoNom int

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {

		var informacionCargo []models.FuncionarioCargo
		filtrodatos := models.FuncionarioCargo{Id: datos.PersonasPreLiquidacion[i].IdPersona, Asignacion_basica: 0}
		tipoNom = 2

		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/funcionario_primatec", "POST", &porcentajePT, datos.PersonasPreLiquidacion[i].IdPersona); err != nil {
			porcentajePT = 0
		}

		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/funcionario_cargo/get_asignacion_basica", "POST", &informacionCargo, &filtrodatos); err == nil {

			diasLaborados := CalcularDias(informacionCargo[0].FechaInicio, informacionCargo[0].FechaFin)
			//reglasNominasEspeciales := crearHechosNominasEspeciales(datos.Preliquidacion,datos.PersonasPreLiquidacion[i].IdPersona)
			esAnual := esAnual(datos.Preliquidacion.Mes, informacionCargo[0].FechaInicio)
			reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato), datos.Preliquidacion)
			reglas = reglasinyectadas + reglasbase + esAnual // + reglasNominasEspeciales

			//fmt.Println(datos.Preliquidacion.Fecha, datos.PersonasPreLiquidacion[i].IdPersona, informacionCargo, diasLaborados, datos.Preliquidacion.Nomina.Periodo, esAnual, porcentajePT,tipoNom)

			temp := golog.CargarReglasFP("", datos.Preliquidacion.Mes, datos.Preliquidacion.Ano, reglas, datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, informacionCargo, diasLaborados, porcentajePT, tipoNom)

			resultado := temp[len(temp)-1]
			resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].IdPersona)
			//fmt.Println(resultado.Conceptos)
			resultado.TotalDevengos, resultado.TotalDescuentos, resultado.TotalAPagar = CalcularTotalesPorPersona(*resultado.Conceptos)

			if datos.Preliquidacion.Definitiva == true {

				for _, descuentos := range *resultado.Conceptos {
					valor, _ := strconv.ParseFloat(descuentos.Valor, 64)
					diasLiquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados, 64)
					tipoPreliquidacion, _ := strconv.Atoi(descuentos.TipoPreliquidacion)
					detallepreliqu := models.DetallePreliquidacion{ConceptoNominaId: &models.ConceptoNomina{Id: descuentos.Id}, PreliquidacionId: &models.Preliquidacion{Id: datos.Preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: datos.PersonasPreLiquidacion[i].NumeroContrato, VigenciaContrato: datos.PersonasPreLiquidacion[i].VigenciaContrato, PersonaId: datos.PersonasPreLiquidacion[i].IdPersona, DiasLiquidados: diasLiquidados, TipoPreliquidacionId: &models.TipoPreliquidacion{Id: tipoPreliquidacion}}

					//detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: informacionContrato.NumeroContrato,VigenciaContrato: vigenciaContrato, Persona: Persona, DiasLiquidados: diasLiquidados, TipoPreliquidacion: &models.TipoPreliquidacion {Id: tipoPreliquidacion}, EstadoDisponibilidad: &models.EstadoDisponibilidad {Id: dispo}}

					if err := request.SendJson(beego.AppConfig.String("UrlCrudTitan")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {
						//fmt.Println(idDetaPre)
					} else {
						fmt.Println("error1: ", err)
					}
				}
			}

			resumenPreliqu = append(resumenPreliqu, resultado)
			reglas = ""
			reglasinyectadas = ""

		} else {
			fmt.Println("error al consultar Asignacion_basica", err)
		}
		reglasinyectadas = ""
	}

	return resumenPreliqu

}

// CalcularDias ...
// @Title CalcularDias
// @Description Dada una fecha inicio y una fecha fin, indicar los días transcurridos
func CalcularDias(FechaInicio time.Time, FechaFin time.Time) (diasLaborados float64) {
	var a, m, d int
	var mesesContrato float64
	var diasContrato float64
	if FechaFin.IsZero() {
		var FechaFin2 time.Time
		FechaFin2 = time.Now()
		a, m, d = diff(FechaInicio, FechaFin2)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
		diasContrato = mesesContrato * 30

	} else {
		a, m, d = diff(FechaInicio, FechaFin)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
		diasContrato = mesesContrato * 30

	}

	//fmt.Println("dias de contFP", diasContrato)

	return diasContrato

}

func esAnual(MesPreliquidacion int, FechaIngreso time.Time) (regla_anual string) {
	//Si es uno, es el momento de pagar bonificacion por servicios.
	var esAnual string

	if MesPreliquidacion == int(FechaIngreso.Month()) {

		esAnual = "esAnual(si)."
	} else {
		esAnual = "esAnual(no)."
	}

	return esAnual
}
