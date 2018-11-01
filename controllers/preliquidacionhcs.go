package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"strconv"
	"github.com/udistrital/titan_api_mid/golog"
	"fmt"
	"time"

)

// PreliquidacionHcSController operations for PreliquidacionHcS
type PreliquidacionHcSController struct {
	beego.Controller
}

func (c *PreliquidacionHcSController) Preliquidar(datos *models.DatosPreliquidacion , reglasbase string) (res []models.Respuesta) {
	//declaracion de variables


	var predicados []models.Predicado //variable para inyectar reglas
	var objeto_datos_contrato models.ObjetoContratoEstado
	var objeto_datos_acta models.ObjetoActaInicio
	var error_consulta_contrato error
	var error_consulta_acta error

	var resumen_preliqu []models.Respuesta
	var meses_contrato float64
	var periodo_liquidacion float64
	var dispo int

	var reglasinyectadas string
	var reglas string
	var predicados_retefuente string

	var idDetaPre interface{}
	var FechaInicioContrato time.Time
	var FechaFinContrato time.Time
	var FechaControl time.Time
	var FechaInicio time.Time
	var FechaFin time.Time

	//-----------------------


	//carga de informacion de los empleados a partir del id de persona Natural (en este momento id proveedor)

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {

		if(datos.PersonasPreLiquidacion[i].Pendiente == "true"){

						var respuesta string
						var verificacion_pago_pendientes int = 2

						detalles_a_mod := ConsultarDetalleAModificar(datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.PersonasPreLiquidacion[i].Preliquidacion)
						resultado := CrearResultado(detalles_a_mod)
						for _, pos := range detalles_a_mod {

							verificacion_pago_pendientes=verificacion_pago(0,datos.Preliquidacion.Ano, datos.Preliquidacion.Mes,pos.NumeroContrato, pos.VigenciaContrato,resultado)
							pos.EstadoDisponibilidad = &models.EstadoDisponibilidad{Id: verificacion_pago_pendientes}
							if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion/"+strconv.Itoa(pos.Id), "PUT", &respuesta, pos); err == nil  {

							} else {
								beego.Debug("error al actualizar detalle de preliquidación: ", err)
							}
						}

					}else{


		objeto_datos_contrato.ContratoEstado.ValorContrato = "3500000"
		objeto_datos_contrato.ContratoEstado.Vigencia = strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)
		objeto_datos_contrato.ContratoEstado.NumeroContrato = datos.PersonasPreLiquidacion[i].NumeroContrato


		//objeto_datos_acta.ActaInicio.FechaInicioTemp = "2018-02-01";
		//objeto_datos_acta.ActaInicio.FechaFinTemp = "2018-06-15";

		objeto_datos_contrato, error_consulta_contrato = ContratosDVE(datos.PersonasPreLiquidacion[i].NumeroContrato,datos.PersonasPreLiquidacion[i].VigenciaContrato )
		objeto_datos_acta, error_consulta_acta = ActaInicioDVE(datos.PersonasPreLiquidacion[i].NumeroContrato,datos.PersonasPreLiquidacion[i].VigenciaContrato )

		if(error_consulta_contrato == nil){
			if(error_consulta_acta == nil){

				datos_contrato := objeto_datos_contrato.ContratoEstado
				datos_acta := objeto_datos_acta.ActaInicio

				layout := "2006-01-02"
				//FechaInicio, _ = time.Parse(layout , "2018-02-01")
				//FechaFin, _ = time.Parse(layout , "2018-06-15")
				FechaInicio, _ = time.Parse(layout , datos_acta.FechaInicioTemp)
				FechaFin, _ = time.Parse(layout , datos_acta.FechaFinTemp)
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

			if datos.Preliquidacion.Mes == 12 || datos.Preliquidacion.Mes == 6 {
				predicados = append(predicados,models.Predicado{Nombre:"fin_contrato("+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)+",si). "} )
			}else{
				predicados = append(predicados,models.Predicado{Nombre:"fin_contrato("+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)+",no). "} )
			}

			predicados = append(predicados,models.Predicado{Nombre:"dias_liquidados("+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)+","+strconv.FormatFloat(periodo_liquidacion, 'f', -1, 64)+"). "} )
			predicados = append(predicados,models.Predicado{Nombre:"valor_contrato("+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)+","+datos_contrato.ValorContrato+"). "} )
			predicados = append(predicados,models.Predicado{Nombre:"duracion_contrato("+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)+","+strconv.FormatFloat(meses_contrato, 'f', -1, 64)+","+vigencia_contrato+"). "} )
			reglasinyectadas = FormatoReglas(predicados)

			reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.Preliquidacion)

			predicados_retefuente = CargarDatosRetefuente(datos.PersonasPreLiquidacion[i].NumDocumento)
			reglas =  reglasinyectadas + reglasbase + predicados_retefuente

			temp := golog.CargarReglasHCS(datos.PersonasPreLiquidacion[i].IdPersona,reglas,vigencia_contrato)

			resultado := temp[len(temp)-1]
			resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].NumDocumento)

			dispo=verificacion_pago(datos.PersonasPreLiquidacion[i].IdPersona,datos.Preliquidacion.Ano, datos.Preliquidacion.Mes,datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato,resultado)

			//se guardan los conceptos calculados en la nomina
			for _, descuentos := range *resultado.Conceptos{
				valor, _ := strconv.ParseFloat(descuentos.Valor,64)
				dias_liquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados,64)
				tipo_preliquidacion,_ := strconv.Atoi(descuentos.TipoPreliquidacion)
				detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: datos.Preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: datos.PersonasPreLiquidacion[i].NumeroContrato,VigenciaContrato: datos.PersonasPreLiquidacion[i].VigenciaContrato, DiasLiquidados: dias_liquidados, TipoPreliquidacion: &models.TipoPreliquidacion {Id: tipo_preliquidacion}, EstadoDisponibilidad: &models.EstadoDisponibilidad {Id: dispo}}

				if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion","POST",&idDetaPre ,&detallepreliqu); err == nil {

				}else{
					beego.Debug("error1: ", err)
				}
			}
			//------------------------------------------------
			resumen_preliqu = append(resumen_preliqu, resultado)
			predicados = nil;
			objeto_datos_contrato = models.ObjetoContratoEstado{}
			reglas = ""
			reglasinyectadas = ""
			predicados_retefuente = "";

			}else{
				fmt.Println("error al traer acta de inicio")
			}

		}else{
			fmt.Println("error al traer valor del contrato")
		}
	}

}


		//-----------------------------
		return resumen_preliqu
}
