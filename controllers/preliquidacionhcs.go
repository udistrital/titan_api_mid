package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"strconv"
	"github.com/udistrital/titan_api_mid/golog"
	"fmt"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"

)

// PreliquidacionHcSController operations for PreliquidacionHcS
type PreliquidacionHcSController struct {
	beego.Controller
}

func (c *PreliquidacionHcSController) Preliquidar(datos models.DatosPreliquidacion , reglasbase string) (res []models.Respuesta) {
	//declaracion de variables
	var resumen_preliqu []models.Respuesta
	var temp_docentes models.ObjetoFuncionarioContrato
	var control_error error
	//-----------------------

	//carga de informacion de los empleados a partir del id de persona Natural (en este momento id proveedor)

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {

		//Si la persona está pendiente, calcula el detalle para el mes que quedó pendiente y lo actualiza
		if(datos.PersonasPreLiquidacion[i].Pendiente == "true"){

						var respuesta string
						var verificacion_pago_pendientes int = 2

						detalles_a_mod := ConsultarDetalleAModificar(datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.PersonasPreLiquidacion[i].Preliquidacion)
						resultado := CrearResultado(detalles_a_mod)
						for _, pos := range detalles_a_mod {

							verificacion_pago_pendientes=verificacion_pago(0,datos.Preliquidacion.Ano, datos.Preliquidacion.Mes,pos.NumeroContrato,  strconv.Itoa(pos.VigenciaContrato),resultado)
							pos.EstadoDisponibilidad = &models.EstadoDisponibilidad{Id: verificacion_pago_pendientes}
							if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion/"+strconv.Itoa(pos.Id), "PUT", &respuesta, pos); err == nil  {

							} else {
								beego.Debug("error al actualizar detalle de preliquidación: ", err)
							}
						}

					}else{



		//BUSCAR CONTRATOS PARA ESA PERSONA

		temp_docentes, control_error = GetContratosPorPersonaHCS(datos,datos.PersonasPreLiquidacion[i])

		//AGRUPAR PARA CALCULAR SOBRE VALORES TOTALES
		if control_error == nil {

			//ELIMINAR REGISTROS SI ESE CONTRATO YA HA SIDO PRELIQUIDADO PARA ESTA PRELIQUIDACION
		if datos.Preliquidacion.Definitiva == true {
		var d []models.DetallePreliquidacion
			query := "Preliquidacion.Id:"+strconv.Itoa(datos.Preliquidacion.Id)+",Persona:"+strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)

							if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
								if len(d) != 0 {
								for _, dato :=  range d {
										urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/detalle_preliquidacion/" + strconv.Itoa(dato.Id)
										var res string
										if err := request.SendJson(urlcrud, "DELETE", &res, nil); err == nil {
											fmt.Println("borrado correctamente")
										}else{
											fmt.Println("error", err)
										}
									}
								}

							}else{
								fmt.Println("error de detalle",err)
							}
	   }

			//temp_agrupar := make(map[string]interface{}) // este mapa tiene la siguiente estructura: temp_agrupar[numero_cedula_docente][id_resolucion][valor_total] (cada resolucion tiene un único tipo de nivel académico, por lo tanto los valores totales se van sumando de acuerdo a la resolución )
			info_resolucion := make(map[string]string)
			info_resoluciones := make(map[string]interface{})

			for _,dato := range temp_docentes.ContratosTipo.ContratoTipo {

							var vinculaciones []models.VinculacionDocente
							query:= "NumeroContrato:"+dato.NumeroContrato+",Vigencia:"+dato.VigenciaContrato
							if err := getJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {

								_, ok := info_resolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)]
								if ok {

												info_resolucion_temp := make(map[string]string)
												temp_valor,_ := strconv.Atoi(info_resolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)])
												temp_valor = temp_valor + int(vinculaciones[0].ValorContrato)
												info_resolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)] =  strconv.Itoa(temp_valor)
												info_resolucion_temp["NumeroContrato"] =  dato.NumeroContrato
												info_resolucion_temp["VigenciaContrato"] =  dato.VigenciaContrato
												info_resolucion_temp["Total"] =  info_resolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)]
												info_resoluciones[strconv.Itoa(vinculaciones[0].IdResolucion.Id)] = info_resolucion_temp

								} else {

												info_resolucion_temp := make(map[string]string)
												temp_valor := int(vinculaciones[0].ValorContrato)
												info_resolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)] =  strconv.Itoa(temp_valor)
												info_resolucion_temp["NumeroContrato"] =  dato.NumeroContrato
												info_resolucion_temp["VigenciaContrato"] =  dato.VigenciaContrato
												info_resolucion_temp["Total"] =  info_resolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)]
												info_resoluciones[strconv.Itoa(vinculaciones[0].IdResolucion.Id)] = info_resolucion_temp

								}
							}

				}

			//CALCULAR PRELIQUIDACIÓN PARA CADA VALOR AGRUPADO
				for key,_ := range info_resoluciones {
					aux := models.ListaContratos{}
				 if err := formatdata.FillStruct(info_resoluciones[key], &aux); err == nil{

					 resumen_preliqu = append(resumen_preliqu,LiquidarContratoHCS(reglasbase,datos.PersonasPreLiquidacion[i].NumDocumento,datos.PersonasPreLiquidacion[i].IdPersona,datos.Preliquidacion,aux)...);
				 }else{
					 fmt.Println("error al guardar información agrupada",err)
				 }
				}

		}

	}

}


		//-----------------------------
		return resumen_preliqu
}

func LiquidarContratoHCS(reglasbase string, NumDocumento,Persona int, preliquidacion *models.Preliquidacion,informacionContrato models.ListaContratos)(res []models.Respuesta){


	var objeto_datos_acta models.ObjetoActaInicio
	var predicados []models.Predicado //variable para inyectar reglas
	var error_consulta_acta error
	var dispo int
	var reglasinyectadas string
  var reglas string
	var predicados_retefuente string
	var idDetaPre interface{}
	var resumen_preliqu []models.Respuesta


	objeto_datos_acta, error_consulta_acta = ActaInicioDVE(informacionContrato.NumeroContrato, informacionContrato.VigenciaContrato)


	  if(error_consulta_acta == nil){

	   datos_acta := objeto_datos_acta
		 vigencia_contrato, _ := strconv.Atoi(informacionContrato.VigenciaContrato)


	  if preliquidacion.Mes == 12 || preliquidacion.Mes == 6 {
	    predicados = append(predicados,models.Predicado{Nombre:"fin_contrato("+strconv.Itoa(Persona)+",si)."} )
	  }else{
	    predicados = append(predicados,models.Predicado{Nombre:"fin_contrato("+strconv.Itoa(Persona)+",no)."} )
	  }


	  predicados = append(predicados,models.Predicado{Nombre:"valor_contrato("+strconv.Itoa(Persona)+","+informacionContrato.Total+")."} )
	  reglasinyectadas = FormatoReglas(predicados)
		reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(Persona, informacionContrato.NumeroContrato, informacionContrato.VigenciaContrato, preliquidacion)
		predicados_retefuente = CargarDatosRetefuente(NumDocumento)
	  reglas =  reglasinyectadas + predicados_retefuente + reglasbase

		temp := golog.CargarReglasHCS(Persona,reglas,preliquidacion, informacionContrato.VigenciaContrato,datos_acta)

	  resultado := temp[len(temp)-1]
	  resultado.NumDocumento = float64(NumDocumento)
		resultado.NumeroContrato = informacionContrato.NumeroContrato
		resultado.VigenciaContrato = informacionContrato.VigenciaContrato
	  dispo=verificacion_pago(NumDocumento,preliquidacion.Ano, preliquidacion.Mes,informacionContrato.NumeroContrato, informacionContrato.VigenciaContrato,resultado)


	  //se guardan los conceptos calculados en la nomina
	  if preliquidacion.Definitiva == true {

	  for _, descuentos := range *resultado.Conceptos{
	    valor, _ := strconv.ParseFloat(descuentos.Valor,64)
	    dias_liquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados,64)
	    tipo_preliquidacion,_ := strconv.Atoi(descuentos.TipoPreliquidacion)
			detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: informacionContrato.NumeroContrato,VigenciaContrato: vigencia_contrato, Persona: Persona, DiasLiquidados: dias_liquidados, TipoPreliquidacion: &models.TipoPreliquidacion {Id: tipo_preliquidacion}, EstadoDisponibilidad: &models.EstadoDisponibilidad {Id: dispo}}

	    if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion","POST",&idDetaPre ,&detallepreliqu); err == nil {

	    }else{
	      beego.Debug("error1: ", err)
	    }
	  }
	}

	  //------------------------------------------------
	  resumen_preliqu = append(resumen_preliqu, resultado)
	  predicados = nil;
  	reglas = ""
	  reglasinyectadas = ""
	  predicados_retefuente = "";

	  }else{
	    fmt.Println("error al traer acta de inicio")
	  }


		return resumen_preliqu



}
