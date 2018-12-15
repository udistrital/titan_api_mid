package controllers

import (
	"fmt"
	"strconv"

	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
	"time"
	"github.com/astaxie/beego"
	"encoding/json"
)

// operations for Preliquidacionct
type PreliquidacionctController struct {
	beego.Controller
}

func (c *PreliquidacionctController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	//declaracion de variables

	var predicados []models.Predicado //variable para inyectar reglas
	var objeto_datos_contrato models.ObjetoContratoEstado
	var objeto_datos_acta models.ObjetoActaInicio
	var error_consulta_contrato error
	var error_consulta_acta error


	var resumen_preliqu []models.Respuesta

	var reglasinyectadas string
	var reglas string
	var disp int

	var idDetaPre interface{}

	var FechaInicio time.Time
	var FechaFin time.Time



	//var al, ml, dl int
	//-----------------------

	//carga de informacion de los empleados a partir del id de persona Natural (en este momento id proveedor)

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {

		if(datos.PersonasPreLiquidacion[i].Pendiente == "true"){

			var respuesta string
			var verificacion_pago_pendientes int = 2

			detalles_a_mod := ConsultarDetalleAModificar(datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.PersonasPreLiquidacion[i].Preliquidacion)
			resultado := CrearResultado(detalles_a_mod)
			for _, pos := range detalles_a_mod {

				verificacion_pago_pendientes=verificacion_pago(0,datos.Preliquidacion.Ano, datos.Preliquidacion.Mes,pos.NumeroContrato, strconv.Itoa(pos.VigenciaContrato),resultado)
				pos.EstadoDisponibilidad = &models.EstadoDisponibilidad{Id: verificacion_pago_pendientes}
				if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion/"+strconv.Itoa(pos.Id), "PUT", &respuesta, pos); err == nil  {
					fmt.Println("preliquidaciones actualizadas")
				} else {
					beego.Debug("error al actualizar detalle de preliquidación: ", err)
				}
			}

		}else{

			if datos.Preliquidacion.Definitiva == true {
			var d []models.DetallePreliquidacion
			query := "Preliquidacion.Id:"+strconv.Itoa(datos.Preliquidacion.Id)+",NumeroContrato:"+datos.PersonasPreLiquidacion[i].NumeroContrato+",VigenciaContrato:"+strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)

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




			objeto_datos_contrato, error_consulta_contrato = ContratosContratistas(datos.PersonasPreLiquidacion[i].NumeroContrato,datos.PersonasPreLiquidacion[i].VigenciaContrato )

			objeto_datos_acta, error_consulta_acta = ActaInicioContratistas(datos.PersonasPreLiquidacion[i].NumeroContrato,datos.PersonasPreLiquidacion[i].VigenciaContrato )



		if(error_consulta_contrato == nil){
			if(error_consulta_acta == nil){
				datos_contrato := objeto_datos_contrato.ContratoEstado
				datos_acta := objeto_datos_acta.ActaInicio

				layout := "2006-01-02"

				FechaInicio, _ = time.Parse(layout , datos_acta.FechaInicioTemp)
				FechaFin, _ = time.Parse(layout , datos_acta.FechaFinTemp)

				dias_contrato := CalcularDias(FechaInicio, FechaFin) + 1  //Suma uno para día inclusive
				fmt.Println("valor del contarto ", datos_contrato.ValorContrato)
				fmt.Println("días contrato ",dias_contrato)

				vigencia_contrato := strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)
				predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona) + "," + datos_contrato.ValorContrato+ "). "})
				predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona) + "," + strconv.FormatFloat(dias_contrato, 'f', -1, 64) + "," + vigencia_contrato + "). "})
				predicados = append(predicados, models.Predicado{Nombre: "pensionado(no)."})

				reglasinyectadas = FormatoReglas(predicados)

				reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato), datos.Preliquidacion)
				reglas = reglasinyectadas + reglasbase

				temp := golog.CargarReglasCT(datos.PersonasPreLiquidacion[i].IdPersona, reglas, datos.Preliquidacion , vigencia_contrato, objeto_datos_acta)

				resultado := temp[len(temp)-1]
				resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].NumDocumento)
				resultado.NumeroContrato = datos.PersonasPreLiquidacion[i].NumeroContrato
				resultado.VigenciaContrato = strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)
				resultado.TotalDevengos, resultado.TotalDescuentos, resultado.TotalAPagar = CalcularTotalesPorPersona(*resultado.Conceptos);

				disp=verificacion_pago(datos.PersonasPreLiquidacion[i].IdPersona,datos.Preliquidacion.Ano, datos.Preliquidacion.Mes,datos.PersonasPreLiquidacion[i].NumeroContrato,  strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato),resultado)


				//ELIMINAR REGISTROS SI ESE CONTRATO YA HA SIDO PRELIQUIDADO PARA ESTA PRELIQUIDACION
				if datos.Preliquidacion.Definitiva == true {
				for _, descuentos := range *resultado.Conceptos{
					valor, _ := strconv.ParseFloat(descuentos.Valor,64)
					dias_liquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados,64)
					tipo_preliquidacion,_ := strconv.Atoi(descuentos.TipoPreliquidacion)
					detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: datos.Preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: datos.PersonasPreLiquidacion[i].NumeroContrato,VigenciaContrato: datos.PersonasPreLiquidacion[i].VigenciaContrato,Persona: datos.PersonasPreLiquidacion[i].IdPersona, DiasLiquidados: dias_liquidados, TipoPreliquidacion: &models.TipoPreliquidacion {Id: tipo_preliquidacion}, EstadoDisponibilidad: &models.EstadoDisponibilidad {Id: disp}}

					if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion","POST",&idDetaPre ,&detallepreliqu); err == nil {

					}else{
						beego.Debug("error1: ", err)
					}
				}
			}

				//
				resumen_preliqu = append(resumen_preliqu, resultado)
				predicados = nil
				objeto_datos_contrato = models.ObjetoContratoEstado{}
				reglas = ""
				reglasinyectadas = ""

			}else{
				fmt.Println("error al traer acta de inicio")
			}

		}else{
			fmt.Println("error al traer valor del contrato")
		}
	}
}

	return resumen_preliqu
}

func ContratosContratistas(id_contrato string, vigencia int)(datos models.ObjetoContratoEstado,  err error){

	var temp map[string]interface{}
	var temp_docentes models.ObjetoContratoEstado
	var control_error error
	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contrato_estado/"+id_contrato+"/"+strconv.Itoa(vigencia), &temp); err == nil && temp != nil {
		jsonDocentes, error_json := json.Marshal(temp)

		if error_json == nil {

			json.Unmarshal(jsonDocentes, &temp_docentes)

		} else {
			control_error = error_json
			fmt.Println("error al traer contratos docentes DVE")
		}
	} else {
		control_error = err
		fmt.Println("Error al unmarshal datos de nómina",err)


	}

		return temp_docentes, control_error;
}

func ActaInicioContratistas(id_contrato string, vigencia int)(datos models.ObjetoActaInicio,  err error){

	var temp map[string]interface{}
	var temp_docentes models.ObjetoActaInicio
	var control_error error

	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/acta_inicio/"+id_contrato+"/"+strconv.Itoa(vigencia), &temp); err == nil && temp != nil {
		jsonDocentes, error_json := json.Marshal(temp)

		if error_json == nil {

			json.Unmarshal(jsonDocentes, &temp_docentes)

		} else {
			control_error = error_json
			fmt.Println("error al traer contratos docentes DVE")
		}
	} else {
		control_error = err
		fmt.Println("Error al unmarshal datos de nómina",err)


	}

		return temp_docentes, control_error;
}

func ConsultarDetalleAModificar(id_contrato string, vigencia, preliquidacion int)(det []models.DetallePreliquidacion){
	var v []models.DetallePreliquidacion
	query := "NumeroContrato:"+id_contrato+",VigenciaContrato:"+strconv.Itoa(vigencia)+",Preliquidacion.Id:"+strconv.Itoa(preliquidacion)
	if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?query="+query, &v); err == nil && v != nil{

	}else{
		fmt.Println("error al consultar preliquidacion a modificar ",err)
	}

	return v

}
