package controllers

import (
	"fmt"
	"strconv"

	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"

	"time"

	"github.com/astaxie/beego"
)

// operations for Preliquidacionct
type PreliquidacionctController struct {
	beego.Controller
}

func (c *PreliquidacionctController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	//declaracion de variables

	var predicados []models.Predicado //variable para inyectar reglas
	var datos_contrato []models.ActaInicio
	//var datos_novedades []models.ConceptoPorPersona
	var resumen_preliqu []models.Respuesta
	var periodo_liquidacion float64

	var reglasinyectadas string
	var reglas string
	var filtrodatos string
	var idDetaPre interface{}
	var FechaControl time.Time
	var FechaInicioContrato time.Time
	var FechaFinContrato time.Time
	var FechaPreliq time.Time

	//var al, ml, dl int
	//-----------------------

	//carga de informacion de los empleados a partir del id de persona Natural (en este momento id proveedor)

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {
		filtrodatos = "NumeroContrato.Id:" + (datos.PersonasPreLiquidacion[i].NumeroContrato) + ",Vigencia:" + datos.Preliquidacion.Nomina.Periodo
		//fmt.Println("Reglas: ", reglasbase)
		if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/acta_inicio?limit=1&query="+filtrodatos, &datos_contrato); err == nil && datos_contrato != nil {

			FechaInicioContrato = time.Date(datos_contrato[0].FechaInicio.Year(), datos_contrato[0].FechaInicio.Month(), datos_contrato[0].FechaInicio.Day()+1, 0, 0, 0, 0, time.UTC)
			FechaFinContrato = time.Date(datos_contrato[0].FechaFin.Year(), datos_contrato[0].FechaFin.Month(), datos_contrato[0].FechaFin.Day()+1, 0, 0, 0, 0, time.UTC)
			FechaPreliq = time.Date(datos.Preliquidacion.Fecha.Year(), datos.Preliquidacion.Fecha.Month(), datos.Preliquidacion.Fecha.Day()+1, 0, 0, 0, 0, time.UTC)

			dias_contrato := CalcularDias(datos_contrato[0].FechaInicio, datos_contrato[0].FechaFin)

			fmt.Println(periodo_liquidacion)
			fmt.Println(FechaInicioContrato)
			fmt.Println(FechaFinContrato)
			fmt.Println(FechaPreliq)
			if FechaInicioContrato.Month() == FechaPreliq.Month() && FechaInicioContrato.Year() == FechaPreliq.Year() {
				FechaControl = time.Date(FechaPreliq.Year(), FechaPreliq.Month(), 30, 0, 0, 0, 0, time.UTC)
				periodo_liquidacion = CalcularDias(FechaInicioContrato, FechaControl) + 1

				fmt.Println("Prueba")
				fmt.Println(periodo_liquidacion)
			} else if FechaFinContrato.Month() == FechaPreliq.Month() && FechaFinContrato.Year() == FechaPreliq.Year() {
				FechaControl = time.Date(FechaPreliq.Year(), FechaPreliq.Month(), 1, 0, 0, 0, 0, time.UTC)
				periodo_liquidacion = CalcularDias(FechaControl, FechaFinContrato) + 1
				fmt.Println("Prueba2")
				fmt.Println(periodo_liquidacion)
			} else {
				periodo_liquidacion = 30
				fmt.Println("Prueba3")

			}

			
			predicados = append(predicados, models.Predicado{Nombre: "dias_liquidados(" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona) + "," + strconv.FormatFloat(periodo_liquidacion, 'f', -1, 64) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona) + "," + strconv.FormatFloat(datos_contrato[0].NumeroContrato.ValorContrato, 'f', -1, 64) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona) + "," + strconv.FormatFloat(dias_contrato, 'f', -1, 64) + "," + datos.Preliquidacion.Nomina.Periodo + "). "})
			reglasinyectadas = FormatoReglas(predicados)

			reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos)
			reglas = reglasinyectadas + reglasbase
			//fmt.Println("Reglas: ", reglasbase)
			temp := golog.CargarReglasCT(reglas, datos.Preliquidacion.Nomina.Periodo)

			resultado := temp[len(temp)-1]
			resultado.NumDocumento = datos_contrato[0].NumeroContrato.Contratista.NumDocumento
			//se guardan los conceptos calculados en la nomina
			for _, descuentos := range *resultado.Conceptos {
				valor, _ := strconv.ParseInt(descuentos.Valor, 10, 64)
				detallepreliqu := models.DetallePreliquidacion{Concepto: &models.Concepto{Id: descuentos.Id}, Persona: datos.PersonasPreLiquidacion[i].IdPersona, Preliquidacion: datos.Preliquidacion.Id, ValorCalculado: valor, NumeroContrato: &models.ContratoGeneral{Id: datos_contrato[0].NumeroContrato.Id}}
				if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

				} else {
					beego.Debug("error1: ", err)
				}
			}
			//------------------------------------------------
			resumen_preliqu = append(resumen_preliqu, resultado)
			predicados = nil
			datos_contrato = nil
			reglas = ""
			reglasinyectadas = ""
		} else {
			fmt.Println(filtrodatos)
			fmt.Println("error3: ", err)
		}

	}
	//-----------------------------
	return resumen_preliqu
}
