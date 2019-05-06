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

// PreliquidacionctController operations for Preliquidacionct
type PreliquidacionctController struct {
	beego.Controller
}

// Preliquidar ...
// @Title Preliquidar
// @Description Preliquidacion para contratistas
func (c *PreliquidacionctController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	//declaracion de variables

	var predicados []models.Predicado //variable para inyectar reglas
	var objetoDatosContrato models.ObjetoContratoEstado
	var objetoDatosActa models.ObjetoActaInicio
	var errorConsultaContrato error
	var errorConsultaActa error


	var resumenPreliqu []models.Respuesta

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

		if(datos.PersonasPreLiquidacion[i].EstadoDisponibilidad == 1){

			fmt.Println("soy una persona pendiente", datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.PersonasPreLiquidacion[i].Preliquidacion)
			var respuesta string
			var verificacionPagoPendientes = 2

			detallesAMod := ConsultarDetalleAModificar(datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.PersonasPreLiquidacion[i].Preliquidacion)
			resultado := CrearResultado(detallesAMod)
			for _, pos := range detallesAMod {

				verificacionPagoPendientes=verificacionPago(0,datos.Preliquidacion.Ano, datos.Preliquidacion.Mes,pos.NumeroContrato, strconv.Itoa(pos.VigenciaContrato),resultado)
				pos.EstadoDisponibilidad = &models.EstadoDisponibilidad{Id: verificacionPagoPendientes}
				if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion/"+strconv.Itoa(pos.Id), "PUT", &respuesta, pos); err == nil  {
					fmt.Println("preliquidaciones actualizadas")
				} else {
					fmt.Println("error al actualizar detalle de preliquidación: ", err)
				}
			}

		}else{

			//eliminar los registros ya existentes en caso de ser definitiva y no solo consulta
			if datos.Preliquidacion.Definitiva == true {
			var d []models.DetallePreliquidacion
			query := "Preliquidacion.Id:"+strconv.Itoa(datos.Preliquidacion.Id)+",NumeroContrato:"+datos.PersonasPreLiquidacion[i].NumeroContrato+",VigenciaContrato:"+strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)

							if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
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

			objetoDatosContrato, errorConsultaContrato = ContratosContratistas(datos.PersonasPreLiquidacion[i].NumeroContrato,datos.PersonasPreLiquidacion[i].VigenciaContrato )

			objetoDatosActa, errorConsultaActa = ActaInicioContratistas(datos.PersonasPreLiquidacion[i].NumeroContrato,datos.PersonasPreLiquidacion[i].VigenciaContrato )



		if(errorConsultaContrato == nil){
			if(errorConsultaActa == nil){
				datosContrato := objetoDatosContrato.ContratoEstado
				datosActa := objetoDatosActa.ActaInicio

				layout := "2006-01-02"

				FechaInicio, _ = time.Parse(layout , datosActa.FechaInicioTemp)
				FechaFin, _ = time.Parse(layout , datosActa.FechaFinTemp)

				diasContrato := CalcularDias(FechaInicio, FechaFin) + 1  //Suma uno para día inclusive
				fmt.Println("valor del contarto ", datosContrato.ValorContrato)
				fmt.Println("días contrato ",diasContrato)

				vigenciaContrato := strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)
				predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona) + "," + datosContrato.ValorContrato+ "). "})
				predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona) + "," + strconv.FormatFloat(diasContrato, 'f', -1, 64) + "," + vigenciaContrato + "). "})
				predicados = append(predicados, models.Predicado{Nombre: "pensionado(no)."})

				reglasinyectadas = FormatoReglas(predicados)

				reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato), datos.Preliquidacion)
				reglas = reglasinyectadas + reglasbase

				temp := golog.CargarReglasCT(datos.PersonasPreLiquidacion[i].IdPersona, reglas, datos.Preliquidacion , vigenciaContrato, objetoDatosActa)

				resultado := temp[len(temp)-1]
				resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].NumDocumento)
				resultado.NumeroContrato = datos.PersonasPreLiquidacion[i].NumeroContrato
				resultado.VigenciaContrato = strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)
				resultado.TotalDevengos, resultado.TotalDescuentos, resultado.TotalAPagar = CalcularTotalesPorPersona(*resultado.Conceptos);

				disp=verificacionPago(datos.PersonasPreLiquidacion[i].IdPersona,datos.Preliquidacion.Ano, datos.Preliquidacion.Mes,datos.PersonasPreLiquidacion[i].NumeroContrato,  strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato),resultado)


				//INSERTAR LOS REGISTROS SI LA PRELIQUIDACIÓN ES DEFINITIVA
				if datos.Preliquidacion.Definitiva == true {
				for _, descuentos := range *resultado.Conceptos{
					valor, _ := strconv.ParseFloat(descuentos.Valor,64)
					diasLiquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados,64)
					tipoPreliquidacion,_ := strconv.Atoi(descuentos.TipoPreliquidacion)
					detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: datos.Preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: datos.PersonasPreLiquidacion[i].NumeroContrato,VigenciaContrato: datos.PersonasPreLiquidacion[i].VigenciaContrato,Persona: datos.PersonasPreLiquidacion[i].IdPersona, DiasLiquidados: diasLiquidados, TipoPreliquidacion: &models.TipoPreliquidacion {Id: tipoPreliquidacion}, EstadoDisponibilidad: &models.EstadoDisponibilidad {Id: disp}}

					if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion","POST",&idDetaPre ,&detallepreliqu); err == nil {

					}else{
						fmt.Println("error1: ", err)
					}
				}
			}

				//
				resumenPreliqu = append(resumenPreliqu, resultado)
				predicados = nil
				objetoDatosContrato = models.ObjetoContratoEstado{}
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

	return resumenPreliqu
}

// ContratosContratistas ...
// @Title ContratosContratistas
// @Description Trae de Argo la informacion del contrato por su número y su vigencia
func ContratosContratistas(id_contrato string, vigencia int)(datos models.ObjetoContratoEstado,  err error){

	var temp map[string]interface{}
	var tempDocentes models.ObjetoContratoEstado
	var controlError error
	if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contrato_estado/"+id_contrato+"/"+strconv.Itoa(vigencia), &temp); err == nil && temp != nil {
		jsonDocentes, errorJSON := json.Marshal(temp)

		if errorJSON == nil {

			json.Unmarshal(jsonDocentes, &tempDocentes)

		} else {
			controlError = errorJSON
			fmt.Println("error al traer contratos docentes DVE")
		}
	} else {
		controlError = err
		fmt.Println("Error al unmarshal datos de nómina",err)


	}

		return tempDocentes, controlError;
}

// ActaInicioContratistas ...
// @Title ActaInicioContratistas
// @Description Trae el acta de inicio por contrato y vigencia
func ActaInicioContratistas(id_contrato string, vigencia int)(datos models.ObjetoActaInicio,  err error){

	var temp map[string]interface{}
	var tempDocentes models.ObjetoActaInicio
	var controlError error

	if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/acta_inicio/"+id_contrato+"/"+strconv.Itoa(vigencia), &temp); err == nil && temp != nil {
		jsonDocentes, errorJSON := json.Marshal(temp)

		if errorJSON == nil {

			json.Unmarshal(jsonDocentes, &tempDocentes)

		} else {
			controlError = errorJSON
			fmt.Println("error al traer contratos docentes DVE")
		}
	} else {
		controlError = err
		fmt.Println("Error al unmarshal datos de nómina",err)


	}

		return tempDocentes, controlError;
}

// ConsultarDetalleAModificar ...
// @Title ConsultarDetalleAModificar
// @Description Consultar los detalles por contrato, vigencia y preliquidacion
func ConsultarDetalleAModificar(id_contrato string, vigencia, preliquidacion int)(det []models.DetallePreliquidacion){
	var v []models.DetallePreliquidacion
	query := "NumeroContrato:"+id_contrato+",VigenciaContrato:"+strconv.Itoa(vigencia)+",Preliquidacion.Id:"+strconv.Itoa(preliquidacion)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?query="+query, &v); err == nil && v != nil{

	}else{
		fmt.Println("error al consultar preliquidacion a modificar ",err)
	}

	return v

}
