package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"strconv"
	"github.com/udistrital/titan_api_mid/golog"
	"fmt"
	"time"
	"encoding/json"
)

// PreliquidacionHcSController operations for PreliquidacionHcS
type PreliquidacionHcSController struct {
	beego.Controller
}

func (c *PreliquidacionHcSController) Preliquidar(datos *models.DatosPreliquidacion , reglasbase string) (res []models.Respuesta) {
	//declaracion de variables


	var predicados []models.Predicado //variable para inyectar reglas
	var datos_contrato models.ContratoEstado
	var datos_acta models.ActaInicioP
	var datos_pruebas []models.DatosPruebas
	//var datos_novedades []models.ConceptoPorPersona
	var resumen_preliqu []models.Respuesta
	var meses_contrato float64
	var periodo_liquidacion float64

	var reglasinyectadas string
	var reglas string
	var filtrodatos string
	var filtrodatos_acta string
	var idDetaPre interface{}
	var FechaInicioContrato time.Time
	var FechaFinContrato time.Time
	var FechaControl time.Time
	var FechaInicio time.Time
	var FechaFin time.Time
	var arreglo_pruebas []models.PruebaGo
	arreglo_pruebas = make([]models.PruebaGo, len(datos.PersonasPreLiquidacion))
	var informacion_cargo []models.FuncionarioCargo
//	var al,ml,dl int
	//-----------------------


	//carga de informacion de los empleados a partir del id de persona Natural (en este momento id proveedor)

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {

		consulta_contratos := models.ContratoGeneral{Id: datos.PersonasPreLiquidacion[i].NumeroContrato,Vigencia:datos.PersonasPreLiquidacion[i].VigenciaContrato}
		//consulta_contratos := models.ContratoGeneral{Id: "658",Vigencia:2017}

		filtrodatos_acta = "NumeroContrato:"+(datos.PersonasPreLiquidacion[i].NumeroContrato)+",Vigencia:"+strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)

		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/contrato_general/contratosHonorariosPruebas", "POST", &datos_contrato, &consulta_contratos); err == nil {

			if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/acta_inicio/actaInicioHonorariosPruebas", "POST", &datos_acta, &consulta_contratos); err == nil {
				layout := "2006-01-02"
				//FechaInicio, err = time.Parse(layout , "2017-02-01")
				//FechaFin, err = time.Parse(layout , "2017-06-15")
				FechaInicio, err = time.Parse(layout , datos_acta.FechaInicioTemp)
				FechaFin, err = time.Parse(layout , datos_acta.FechaFinTemp)
				a,m,d := diff(FechaInicio,FechaFin)

			FechaInicioContrato = time.Date(FechaInicio.Year(), FechaInicio.Month(), FechaInicio.Day(), 0, 0, 0, 0, time.UTC)
			FechaFinContrato = time.Date(FechaFin.Year(), FechaFin.Month(), FechaFin.Day(), 0, 0, 0, 0, time.UTC)

			if int(FechaInicioContrato.Month()) == datos.Preliquidacion.Mes && int(FechaInicioContrato.Year()) == datos.Preliquidacion.Ano {
				FechaControl = time.Date(datos.Preliquidacion.Ano, time.Month(datos.Preliquidacion.Mes), 30, 0, 0, 0, 0, time.UTC)
				periodo_liquidacion = CalcularDias(FechaInicioContrato, FechaControl)


			} else if int(FechaFinContrato.Month()) == datos.Preliquidacion.Mes && int(FechaFinContrato.Year()) == datos.Preliquidacion.Ano {
				FechaControl = time.Date(datos.Preliquidacion.Ano, time.Month(datos.Preliquidacion.Mes), 1, 0, 0, 0, 0, time.UTC)
				periodo_liquidacion = CalcularDias(FechaControl, FechaFinContrato)

			} else {
				periodo_liquidacion = 30


			}

			vigencia_contrato := strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)
			meses_contrato = (float64(a*12))+float64(m)+(float64(d)/30)
		
			fmt.Println(a)
			fmt.Println(m)
			fmt.Println(d)
			fmt.Println("meses: ",meses_contrato)
			fmt.Println("dias: ",periodo_liquidacion)
			if int(FechaFinContrato.Month()) == datos.Preliquidacion.Mes && int(FechaFinContrato.Year()) == datos.Preliquidacion.Ano {
				predicados = append(predicados,models.Predicado{Nombre:"fin_contrato("+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)+",si). "} )
			}else{
				predicados = append(predicados,models.Predicado{Nombre:"fin_contrato("+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)+",no). "} )
			}

			predicados = append(predicados,models.Predicado{Nombre:"dias_liquidados("+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)+","+strconv.FormatFloat(periodo_liquidacion, 'f', -1, 64)+"). "} )
			predicados = append(predicados,models.Predicado{Nombre:"valor_contrato("+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)+","+datos_contrato.ValorContrato+"). "} )
			predicados = append(predicados,models.Predicado{Nombre:"duracion_contrato("+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)+","+strconv.FormatFloat(meses_contrato, 'f', -1, 64)+","+vigencia_contrato+"). "} )
			reglasinyectadas = FormatoReglas(predicados)

			reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.Preliquidacion)
			reglas =  reglasinyectadas + reglasbase

			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/datos_pruebas?limit=-1&query=MesPreliq:"+strconv.Itoa(datos.Preliquidacion.Mes)+",AnoPreliq:"+strconv.Itoa(datos.Preliquidacion.Ano)+",NumDocumento:"+strconv.Itoa(datos.PersonasPreLiquidacion[i].NumDocumento), &datos_pruebas); err == nil && datos_pruebas != nil{
				arreglo_pruebas[i] = models.PruebaGo{informacion_cargo, "",datos.Preliquidacion.FechaRegistro, datos_pruebas[0].ValorSalario,"","","","","","",datos_pruebas[0].ValorSalud,datos_pruebas[0].ValorPension,datos.PersonasPreLiquidacion[i].IdPersona,datos.PersonasPreLiquidacion[i].NumDocumento,0,datos.Preliquidacion.Mes, datos.Preliquidacion.Ano, 0, 0}
			}else{
				fmt.Println(err)
			}

			temp := golog.CargarReglasHCS(datos.PersonasPreLiquidacion[i].IdPersona,reglas,vigencia_contrato)

			resultado := temp[len(temp)-1]
			resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].NumDocumento)
			//se guardan los conceptos calculados en la nomina
			for _, descuentos := range *resultado.Conceptos{
				valor, _ := strconv.ParseFloat(descuentos.Valor,64)
				dias_liquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados,64)
				tipo_preliquidacion,_ := strconv.Atoi(descuentos.TipoPreliquidacion)
				detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: datos.Preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: datos.PersonasPreLiquidacion[i].NumeroContrato,VigenciaContrato: datos.PersonasPreLiquidacion[i].VigenciaContrato, DiasLiquidados: dias_liquidados, TipoPreliquidacion: &models.TipoPreliquidacion {Id: tipo_preliquidacion}}

				if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion","POST",&idDetaPre ,&detallepreliqu); err == nil {

				}else{
					beego.Debug("error1: ", err)
				}
			}
			//------------------------------------------------
			resumen_preliqu = append(resumen_preliqu, resultado)
			predicados = nil;
			datos_contrato = models.ContratoEstado{}
			reglas = ""
			reglasinyectadas = ""
		}else{
			fmt.Println(filtrodatos_acta)
			fmt.Println("error3: ", err)
		}

	}else{
		fmt.Println(filtrodatos)
		fmt.Println("error2: ", err)
	}
}

			data, err := json.Marshal(arreglo_pruebas)
			if err != nil {
					fmt.Println("error en json")
				}
			str := fmt.Sprintf("%s", data)
			mes := strconv.Itoa(datos.Preliquidacion.Mes)
			ano := strconv.Itoa(datos.Preliquidacion.Ano)
			if err := WriteStringToFile("pruebaHC"+ano+mes+".txt", str); err != nil {
					panic(err)
			}
		//-----------------------------
		return resumen_preliqu
}