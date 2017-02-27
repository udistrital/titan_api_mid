package controllers

import (
	"strconv"
	"titan_api_mid/golog"
	"titan_api_mid/models"

	"github.com/astaxie/beego"
	//  "fmt"
	"time"
)

// PreliquidacionHcController operations for PreliquidacionHc
type PreliquidacionFpController struct {
	beego.Controller
}

// Post ...
// @Title Create
// @Description create PreliquidacionHc
// @Param	body		body 	models.PreliquidacionHc	true		"body for PreliquidacionHc content"
// @Success 201 {object} models.PreliquidacionHc
// @Failure 403 body is empty
// @router / [post]
func (c *PreliquidacionFpController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	//declaracion de variables
	var reglasinyectadas string
	var reglas string
	var idDetaPre interface{}
	var resumen_preliqu []models.Respuesta
	var porcentajePT int
	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {
		var informacion_cargo []models.FuncionarioCargo
		filtrodatos := models.FuncionarioCargo{Id: datos.PersonasPreLiquidacion[i].IdPersona, Asignacion_basica: 0}

		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/funcionario_primatec", "POST", &porcentajePT, datos.PersonasPreLiquidacion[i].IdPersona); err != nil {
			porcentajePT = 0
		}

		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/funcionario_cargo", "POST", &informacion_cargo, &filtrodatos); err == nil {
			dias_laborados := CalcularDias(informacion_cargo[0].FechaInicio, informacion_cargo[0].FechaFin)
			esAnual := esAnual(informacion_cargo[0].FechaInicio)
			reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos)
			reglas = reglasinyectadas + reglasbase
			//fmt.Println("reglas: ",reglas)
			temp := golog.CargarReglasFP(reglas, datos.PersonasPreLiquidacion[i].IdPersona, informacion_cargo, dias_laborados, datos.Preliquidacion.Nomina.Periodo, esAnual, porcentajePT)

			resultado := temp[len(temp)-1]
			resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].IdPersona)
			resumen_preliqu = append(resumen_preliqu, resultado)

			for _, descuentos := range *resultado.Conceptos {
				valor, _ := strconv.ParseInt(descuentos.Valor, 10, 64)
				detallepreliqu := models.DetallePreliquidacion{Concepto: &models.Concepto{Id: descuentos.Id}, Persona: datos.PersonasPreLiquidacion[i].IdPersona, Preliquidacion: datos.Preliquidacion.Id, ValorCalculado: valor, NumeroContrato: &models.ContratoGeneral{Id: datos.PersonasPreLiquidacion[i].NumeroContrato}}
				if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

				} else {
					beego.Debug("error1: ", err)
				}
			}

		}

	}

	return resumen_preliqu

}

func CalcularDias(FechaInicio time.Time, FechaFin time.Time) (dias_laborados float64) {
	var a, m, d int
	var meses_contrato float64
	var dias_contrato float64
	if FechaFin.IsZero() {
		var FechaFin2 time.Time
		FechaFin2 = time.Now()
		a, m, d = diff(FechaInicio, FechaFin2)
		meses_contrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
		dias_contrato = meses_contrato * 30

	} else {
		a, m, d = diff(FechaInicio, FechaFin)
		meses_contrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
		dias_contrato = meses_contrato * 30

	}

	return dias_contrato

}

func esAnual(FechaInicio time.Time) (flag int) {
	//Si es uno, es el momento de pagar bonificacion por servicios.
	var esAnual int
	if FechaInicio.Month() == time.Now().Month() {
		esAnual = 1
	} else {
		esAnual = 0
	}

	return esAnual
}
