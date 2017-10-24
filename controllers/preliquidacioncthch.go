package controllers

import (
	"fmt"
	"strconv"

	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"

	"time"
	"github.com/astaxie/beego"
	"encoding/json"
)

// operations for Preliquidacioncthch
type PreliquidacioncthchController struct {
	beego.Controller
}

func (c *PreliquidacioncthchController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	//declaracion de variables

	var predicados []models.Predicado //variable para inyectar reglas
	var datos_contrato models.ContratoEstado
	var datos_acta models.ActaInicioP
	var datos_pruebas []models.DatosPruebas
	//var datos_novedades []models.ConceptoPorPersona
	var resumen_preliqu []models.Respuesta
	var periodo_liquidacion float64

	var reglasinyectadas string
	var reglas string
	var url_consulta string
	var filtrodatos_acta string
	var idDetaPre interface{}
	var FechaControl time.Time
	var FechaInicioContrato time.Time
	var FechaFinContrato time.Time
	var FechaInicio time.Time
	var FechaFin time.Time



	var arreglo_pruebas []models.PruebaGo
	arreglo_pruebas = make([]models.PruebaGo, len(datos.PersonasPreLiquidacion))
	var informacion_cargo []models.FuncionarioCargo

	if(datos.Preliquidacion.Nomina.TipoNomina.Nombre == "CT"){
		url_consulta = "ContratistasPruebas"
	}else{
		url_consulta = "HonorariosPruebas"

	}
	//var al, ml, dl int
	//-----------------------

	//carga de informacion de los empleados a partir del id de persona Natural (en este momento id proveedor)

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {

		consulta_contratos := models.ContratoGeneral{Id: datos.PersonasPreLiquidacion[i].NumeroContrato,Vigencia:datos.PersonasPreLiquidacion[i].VigenciaContrato}
		//consulta_contratos := models.ContratoGeneral{Id: "658",Vigencia:2017}

		filtrodatos_acta = "NumeroContrato:"+(datos.PersonasPreLiquidacion[i].NumeroContrato)+",Vigencia:"+strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)



		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/contrato_general/contratos"+url_consulta, "POST", &datos_contrato, &consulta_contratos); err == nil {

			if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/acta_inicio/actaInicio"+url_consulta, "POST", &datos_acta, &consulta_contratos); err == nil {

				layout := "2006-01-02"
				FechaInicio, err = time.Parse(layout , "2017-02-01")
				FechaFin, err = time.Parse(layout , "2017-06-15")
				//FechaInicio, err = time.Parse(layout , datos_acta.FechaInicioTemp)
				//FechaFin, err = time.Parse(layout , datos_acta.FechaFinTemp)

			FechaInicioContrato = time.Date(FechaInicio.Year(), FechaInicio.Month(), FechaInicio.Day(), 0, 0, 0, 0, time.UTC)
			FechaFinContrato = time.Date(FechaFin.Year(), FechaFin.Month(), FechaFin.Day(), 0, 0, 0, 0, time.UTC)

			dias_contrato := CalcularDias(FechaInicio, FechaFin)


			if int(FechaInicioContrato.Month()) == datos.Preliquidacion.Mes && int(FechaInicioContrato.Year()) == datos.Preliquidacion.Ano {
				FechaControl = time.Date(datos.Preliquidacion.Ano, time.Month(datos.Preliquidacion.Mes), 30, 0, 0, 0, 0, time.UTC)
				periodo_liquidacion = CalcularDias(FechaInicioContrato, FechaControl)


			} else if int(FechaFinContrato.Month()) == datos.Preliquidacion.Mes && int(FechaFinContrato.Year()) == datos.Preliquidacion.Ano {
				FechaControl = time.Date(datos.Preliquidacion.Ano, time.Month(datos.Preliquidacion.Mes), 1, 0, 0, 0, 0, time.UTC)
				periodo_liquidacion = CalcularDias(FechaControl, FechaFinContrato)

			} else {
				periodo_liquidacion = 30


			}
			fmt.Println("periodo de liquidacion")
			fmt.Println(periodo_liquidacion)

			vigencia_contrato := strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)
			predicados = append(predicados, models.Predicado{Nombre: "dias_liquidados(" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona) + "," + strconv.FormatFloat(periodo_liquidacion, 'f', -1, 64) + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona) + "," + datos_contrato.ValorContrato+ "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona) + "," + strconv.FormatFloat(dias_contrato, 'f', -1, 64) + "," + vigencia_contrato + "). "})
			fmt.Println(predicados)
			reglasinyectadas = FormatoReglas(predicados)

			reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.Preliquidacion)
			reglas = reglasinyectadas + reglasbase

			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/datos_pruebas?limit=-1&query=MesPreliq:"+strconv.Itoa(datos.Preliquidacion.Mes)+",AnoPreliq:"+strconv.Itoa(datos.Preliquidacion.Ano)+",NumDocumento:"+strconv.Itoa(datos.PersonasPreLiquidacion[i].NumDocumento), &datos_pruebas); err == nil && datos_pruebas != nil{
				arreglo_pruebas[i] = models.PruebaGo{informacion_cargo, "",datos.Preliquidacion.FechaRegistro, datos_pruebas[0].ValorSalario,datos_pruebas[0].ValorReteica,datos_pruebas[0].ValorEstampillaUD,datos_pruebas[0].ValorProCultura,datos_pruebas[0].ValorAdultoMayor,"","","","",datos.PersonasPreLiquidacion[i].IdPersona,datos.PersonasPreLiquidacion[i].NumDocumento,0,datos.Preliquidacion.Mes, datos.Preliquidacion.Ano, 0, 0}
			}else{
				fmt.Println(err)
			}
			
			temp := golog.CargarReglasCT(datos.PersonasPreLiquidacion[i].IdPersona, reglas, vigencia_contrato)

			resultado := temp[len(temp)-1]
			resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].NumDocumento)

			disponiblidad:= calcular_disponibilidad(2439,2017,resultado)
			fmt.Println("disponibilidad")
			fmt.Println(disponiblidad)
			//calcular_disponibilidad(datos.PersonasPreLiquidacion[i].NumDocumento, 	datos.PersonasPreLiquidacion[i].VigenciaContrato, resultado)

			for _, descuentos := range *resultado.Conceptos {
				valor, _ := strconv.ParseFloat(descuentos.Valor,64)
				dias_liquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados,64)
				tipo_preliquidacion,_ := strconv.Atoi(descuentos.TipoPreliquidacion)
				detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: datos.Preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: datos.PersonasPreLiquidacion[i].NumeroContrato,VigenciaContrato: datos.PersonasPreLiquidacion[i].VigenciaContrato, DiasLiquidados: dias_liquidados, TipoPreliquidacion: &models.TipoPreliquidacion {Id: tipo_preliquidacion}, EstadoDisponibilidad: &models.EstadoDisponibilidad {Id: disponiblidad}}

				if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

				} else {
					beego.Debug("error1: ", err)
				}
			}
			//------------------------------------------------

		//

			fmt.Println(resumen_preliqu)
			predicados = nil
			datos_contrato = models.ContratoEstado{}
			reglas = ""
			reglasinyectadas = ""
		} else {
			fmt.Println(filtrodatos_acta)
			fmt.Println("errorsote: ", err)
		}
		}else{

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
	if err := WriteStringToFile("pruebaContratistas"+ano+mes+".txt", str); err != nil {
			panic(err)
	}
	//-----------------------------
	return resumen_preliqu
}

func consultar_rp (num_documento, vigencia int) (saldo float64){
		var registro_presupuestal []models.RegistroPresupuestal
		var saldo_rp float64
		var num_documento_string = strconv.Itoa(num_documento)
		var vigencia_string = strconv.Itoa(vigencia)
		if err := getJson("http://"+beego.AppConfig.String("Urlkronos")+":"+beego.AppConfig.String("Portkronos")+"/"+beego.AppConfig.String("Nskronos")+"/registro_presupuestal?limit=-1&query=Beneficiario:"+num_documento_string+",Vigencia:"+vigencia_string, &registro_presupuestal); err == nil {
			var id_registro_pre = strconv.Itoa(registro_presupuestal[0].Id)
			if err := getJson("http://"+beego.AppConfig.String("Urlkronos")+":"+beego.AppConfig.String("Portkronos")+"/"+beego.AppConfig.String("Nskronos")+"/registro_presupuestal/ValorActualRp/"+id_registro_pre, &saldo_rp); err == nil {
				fmt.Println("saldo rp")
				fmt.Println(saldo_rp)
			}else{
				fmt.Println("error al consultar saldo de rp")
				fmt.Println(err)
			}



		}else{
			fmt.Println("error en consulta de rp")
			fmt.Println(err)
		}

		return saldo_rp
}


func total_a_pagar(respuesta models.Respuesta)(total float64){
	var total_dev float64
	for _, descuentos := range *respuesta.Conceptos {
		if(descuentos.NaturalezaConcepto == 1){
			valor, _ := strconv.ParseFloat(descuentos.Valor,64)
			total_dev = total_dev + valor
		}


}
fmt.Println("total a pagar")
fmt.Println(total_dev)
return total_dev
}

func calcular_disponibilidad(num_documento, vigencia int,respuesta models.Respuesta)(disp int){
	var valor_a_pagar float64
	var saldo_rp float64
	var disponibilidad int
	saldo_rp = consultar_rp(num_documento, vigencia)
	valor_a_pagar = total_a_pagar(respuesta)
	if(valor_a_pagar > saldo_rp){
		disponibilidad = 1;
		fmt.Println("no hay dinero")
	}else{
		disponibilidad = 2;
		fmt.Println("si hay dinero ")
	}

	return disponibilidad
}
