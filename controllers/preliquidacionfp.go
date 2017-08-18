package controllers

import (
	"strconv"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"encoding/json"
	"github.com/astaxie/beego"
  "fmt"
	"time"
)

// PreliquidacionHcController operations for PreliquidacionHc
type PreliquidacionFpController struct {
	beego.Controller
}

func (c *PreliquidacionFpController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	//declaracion de variables
	var reglasinyectadas string
	var reglas string
	var idDetaPre interface{}
	var resumen_preliqu []models.Respuesta
	var porcentajePT int
	var tipoNom int;
	var arreglo_pruebas []models.PruebaGo
	arreglo_pruebas = make([]models.PruebaGo, len(datos.PersonasPreLiquidacion))

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {


		var informacion_cargo []models.FuncionarioCargo
		filtrodatos := models.FuncionarioCargo{Id: datos.PersonasPreLiquidacion[i].IdPersona, Asignacion_basica: 0}
		tipoNom = 2

		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/funcionario_primatec", "POST", &porcentajePT, datos.PersonasPreLiquidacion[i].IdPersona); err != nil {
			porcentajePT = 0
		}

		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/funcionario_cargo", "POST", &informacion_cargo, &filtrodatos); err == nil {

			dias_laborados := CalcularDias(informacion_cargo[0].FechaInicio, informacion_cargo[0].FechaFin)
			//reglasNominasEspeciales := crearHechosNominasEspeciales(datos.Preliquidacion,datos.PersonasPreLiquidacion[i].IdPersona)
			esAnual := esAnual(datos.Preliquidacion.Mes, informacion_cargo[0].FechaInicio)
			reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.Preliquidacion)
			reglas = reglasinyectadas + reglasbase // + reglasNominasEspeciales

			//fmt.Println(datos.Preliquidacion.Fecha, datos.PersonasPreLiquidacion[i].IdPersona, informacion_cargo, dias_laborados, datos.Preliquidacion.Nomina.Periodo, esAnual, porcentajePT,tipoNom)

			arreglo_pruebas[i] = models.PruebaGo{informacion_cargo, "",datos.Preliquidacion.FechaRegistro, "",datos.PersonasPreLiquidacion[i].IdPersona,0,dias_laborados,datos.Preliquidacion.Mes, datos.Preliquidacion.Ano,esAnual, porcentajePT, tipoNom}

			temp := golog.CargarReglasFP(datos.Preliquidacion.Mes, datos.Preliquidacion.Ano,reglas, datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, informacion_cargo, dias_laborados, esAnual, porcentajePT,tipoNom)

			resultado := temp[len(temp)-1]
			resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].IdPersona)
			resumen_preliqu = append(resumen_preliqu, resultado)

			for _, descuentos := range *resultado.Conceptos {
				valor, _ := strconv.ParseFloat(descuentos.Valor,64)
				dias_liquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados,64)
				tipo_preliquidacion,_ := strconv.Atoi(descuentos.TipoPreliquidacion)
				detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: datos.Preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: datos.PersonasPreLiquidacion[i].NumeroContrato,VigenciaContrato: datos.PersonasPreLiquidacion[i].VigenciaContrato, DiasLiquidados: dias_liquidados, TipoPreliquidacion: &models.TipoPreliquidacion {Id: tipo_preliquidacion}}
				if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

				} else {
					beego.Debug("error1: ", err)
				}
			}

		}
		reglasinyectadas = "";
	}

	data, err := json.Marshal(arreglo_pruebas)
	if err != nil {
			fmt.Println("error en json")
		}
	str := fmt.Sprintf("%s", data)
	mes := strconv.Itoa(datos.Preliquidacion.Mes)
	if err := WriteStringToFile("prueba"+mes+".txt", str); err != nil {
			panic(err)
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

func esAnual(MesPreliquidacion int, FechaIngreso time.Time) (flag int) {
	//Si es uno, es el momento de pagar bonificacion por servicios.
	var esAnual int

	if MesPreliquidacion == int(FechaIngreso.Month()) {

		esAnual = 1
	} else {
		esAnual = 0
	}

	return esAnual
}
