package controllers

import (
	"fmt"
	"titan_api_mid/golog"
	"titan_api_mid/models"
	"strconv"

	"github.com/astaxie/beego"
)

// PreliquidacionpeController operations for Preliquidacionpe
type PreliquidacionpeController struct {
	beego.Controller
}

func (c *PreliquidacionpeController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	//var predicados []models.Predicado //variable para inyectar reglas
	var resumen_preliqu []models.Respuesta
	var idDetaPre interface{}
	var pensionados []models.InformacionPensionado // arreglo de informacion_pensionado
	var beneficiarios string

	var reglasinyectadas string
	var reglas string

	//	var reglasinyectadas string
	//var reglas string
	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {
		filtrodatos := models.InformacionPensionado{Id: datos.PersonasPreLiquidacion[i].IdPersona}

		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/informacion_pensionado", "POST", &pensionados, &filtrodatos); err == nil {
id := strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)
			if err2 := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/beneficiario/beneficiarioDatos", "POST", &beneficiarios, &id); err == nil {
fmt.Println("beneficiarios")
fmt.Println(beneficiarios)
		reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos)
		reglas = reglasbase + reglasinyectadas
		fmt.Println(reglas)
			temp := golog.CargarReglasPE(reglas, pensionados[i], beneficiarios)


			resultado := temp[len(temp)-1]
			resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].IdPersona)
			resumen_preliqu = append(resumen_preliqu, resultado)
			fmt.Println("resultado Preliquidacion")
			fmt.Println(resumen_preliqu[0].Conceptos)

			for _, descuentos := range *resultado.Conceptos {
				valor, _ := strconv.ParseInt(descuentos.Valor, 10, 64)
				//fmt.Println("asdfg"+datos.PersonasPreLiquidacion[i].NumeroContrato)
				detallepreliqu := models.DetallePreliquidacion{Concepto: &models.Concepto{Id: descuentos.Id}, Persona: datos.PersonasPreLiquidacion[i].IdPersona, Preliquidacion: datos.Preliquidacion.Id, ValorCalculado: valor, NumeroContrato: &models.ContratoGeneral{Id: "DVE8000"}}
				if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

					} else {
						beego.Debug("error1: ", err)
					}
				}
			 }else {
				 fmt.Println(err2)
			}
				}else {
					fmt.Println(err)
				}
			}
			return resumen_preliqu
		}
