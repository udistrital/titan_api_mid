package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
	//"time"
	"github.com/astaxie/beego"
)

// GestionPersonasAPreliquidarController operations for GestionPersonasAPreliquidar
type GestionPersonasAPreliquidarController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionPersonasAPreliquidarController) URLMapping() {
	c.Mapping("ListarPersonasAPreliquidar", c.ListarPersonasAPreliquidar)
}

// ListarPersonasAPreliquidar ...
// @Title create ListarPersonasAPreliquidar
// @Description create ListarPersonasAPreliquidar: Lista a las personas que tienen vinculaciones activas para ese periodo y que por consiguiente pueden ser preliquidadas
// @Param	body 	body    models.Preliquidacion	true		"body for models.Preliquidacion content"
// @Success 201
// @Failure 403 body is empty
// @router /listar_personas_a_preliquidar_argo [post]
func (c *GestionPersonasAPreliquidarController) ListarPersonasAPreliquidar() {
	var v models.Preliquidacion
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if v.Nomina.TipoNomina.Nombre == "CT" {
			fmt.Println("preliq",v)
			if listaContratos, err := ListaContratosContratistas(v); err == nil {
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = listaContratos.ContratosTipo.ContratoTipo
			} else {
				c.Data["json"] = err.Error()
				fmt.Println("error : ", err)
			}


		} else if v.Nomina.TipoNomina.Nombre == "HCS" {
				if listaContratos, err := ListarPersonasHCS(v); err == nil {
					c.Ctx.Output.SetStatus(201)
					c.Data["json"] = listaContratos.ContratosTipo.ContratoTipo
				} else {
					c.Data["json"] = err.Error()
					fmt.Println("error : ", err)
				}
			}	else if v.Nomina.TipoNomina.Nombre == "HCH" {
					if listaContratos, err := ListarPersonasHCH(v); err == nil {
						c.Ctx.Output.SetStatus(201)
						c.Data["json"] = listaContratos.ContratosTipo.ContratoTipo
					} else {
						c.Data["json"] = err.Error()
						fmt.Println("error : ", err)
					}
				}else if v.Nomina.TipoNomina.Nombre == "FP" {
					fmt.Println("Planta")
					if listaContratos, err := ListaContratosFuncionariosPlanta(); err == nil {
						c.Ctx.Output.SetStatus(201)
						c.Data["json"] = listaContratos
					} else {
						c.Data["json"] = err.Error()
						fmt.Println("error : ", err)
					}
				}

	} else {
		c.Data["json"] = err.Error()
		fmt.Println("error 2: ", err)
	}
	c.ServeJSON()

}

// ListarPersonasAPreliquidarPendientes ...
// @Title create ListarPersonasAPreliquidarPendientes
// @Description create ListarPersonasAPreliquidar: Lista a las personas pendientes de periodos anteriores para que puedan ser tenidas en cuenta el presente mes
// @Param	body 	body models.Preliquidacion	true		"body for models.Preliquidacion content"
// @Success 201 {object} models.DetallePreliquidacion
// @Failure 403 body is empty
// @router /listar_personas_a_preliquidar_pendientes [post]
func (c *GestionPersonasAPreliquidarController) ListarPersonasAPreliquidarPendientes() {
	var v models.Preliquidacion
	var personasPendPreliquidacion []models.DetallePreliquidacion
	var errorConsultaInformacionAgora error

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion/get_personas_pago_pendiente?idNomina="+strconv.Itoa(v.Nomina.Id), &personasPendPreliquidacion); err == nil && personasPendPreliquidacion !=nil {



			for x, dato := range personasPendPreliquidacion {
				personasPendPreliquidacion[x].NombreCompleto, _, personasPendPreliquidacion[x].Documento, errorConsultaInformacionAgora= InformacionPersona(v.Nomina.TipoNomina.Nombre,dato.NumeroContrato, dato.VigenciaContrato)
				personasPendPreliquidacion[x].Preliquidacion = ConsultarDatosPreliq(dato.Preliquidacion.Id)
			}

			if(errorConsultaInformacionAgora == nil){
				c.Data["json"] = personasPendPreliquidacion
			}else{
				c.Data["json"] = errorConsultaInformacionAgora
				fmt.Println("error al consultar información en Agora")
			}


		} else {
			fmt.Println(err)
			fmt.Println("error al traer pendientes")
			c.Data["json"] = err
		}

	} else {
		c.Data["json"] = err.Error()
		fmt.Println("error al leer json de nómina: ", err)
	}
	c.ServeJSON()

}

// ListarPersonasHCS ...
// @Title  ListarPersonasHCS
// @Description ListarPersonasHCS: Trae las personas que tienen contratos vigentes para ese mes en Hora Cátedra Salarios
func ListarPersonasHCS(objeto_nom models.Preliquidacion)(arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error){

	var temp map[string]interface{}

	var tipoNom string
	var tempDocentes models.ObjetoFuncionarioContrato
	var tempDocentesTco models.ObjetoFuncionarioContrato
	var controlError error
	var mes string
	var ano = strconv.Itoa(objeto_nom.Ano);

	if (objeto_nom.Mes >= 1 && objeto_nom.Mes <= 9){
		mes = strconv.Itoa(objeto_nom.Mes);
		mes = "0"+mes
	}else{
		mes = strconv.Itoa(objeto_nom.Mes);
	}



	tipoNom = "2"

		if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_personas/"+tipoNom+"/"+ano+"-"+mes+"/"+ano+"-"+mes, &temp); err == nil && temp != nil {
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

		tipoNom = "18"

		if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_personas/"+tipoNom+"/"+ano+"-"+mes+"/"+ano+"-"+mes, &temp); err == nil && temp != nil {
			jsonDocentes, errorJSON := json.Marshal(temp)

			if errorJSON == nil {

				json.Unmarshal(jsonDocentes, &tempDocentesTco)

			} else {
				controlError = errorJSON
				fmt.Println("error al traer contratos docentes DVE")
			}
		} else {
			controlError = err
			fmt.Println("Error al unmarshal datos de nómina",err)


		}

    tempDocentes.ContratosTipo.ContratoTipo = append(tempDocentes.ContratosTipo.ContratoTipo, tempDocentesTco.ContratosTipo.ContratoTipo...)


	var d []models.DetallePreliquidacion
	for x, dato := range tempDocentes.ContratosTipo.ContratoTipo {
		query := "Preliquidacion.Id:"+strconv.Itoa(objeto_nom.Id)+",Persona:"+dato.Id
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
			if len(d) == 0 {
				tempDocentes.ContratosTipo.ContratoTipo[x].Preliquidado = "No"

			}else{
				tempDocentes.ContratosTipo.ContratoTipo[x].Preliquidado = "Sí"

			}

		}


	}


	return tempDocentes, controlError;

}

// ListarPersonasHCH ...
// @Title  ListarPersonasHCH
// @Description ListarPersonasHCH: Trae las personas que tienen contratos vigentes para ese mes en Hora Cátedra Honorarios
func ListarPersonasHCH(objeto_nom models.Preliquidacion)(arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error){

	var temp map[string]interface{}

	var tipoNom string
	var tempDocentes models.ObjetoFuncionarioContrato
  var controlError error
	var mes string
	var ano = strconv.Itoa(objeto_nom.Ano);

	if (objeto_nom.Mes >= 1 && objeto_nom.Mes <= 9){
		mes = strconv.Itoa(objeto_nom.Mes);
		mes = "0"+mes
	}else{
		mes = strconv.Itoa(objeto_nom.Mes);
	}



	tipoNom = "3"

		if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_personas/"+tipoNom+"/"+ano+"-"+mes+"/"+ano+"-"+mes, &temp); err == nil && temp != nil {
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

 //SABER SI YA FUE PRELIQUIDADO O NO
 var d []models.DetallePreliquidacion
 for x, dato := range tempDocentes.ContratosTipo.ContratoTipo {
	 query := "Preliquidacion.Id:"+strconv.Itoa(objeto_nom.Id)+",Persona:"+dato.Id
	 if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
		 if len(d) == 0 {
			 tempDocentes.ContratosTipo.ContratoTipo[x].Preliquidado = "No"

		 }else{
			 tempDocentes.ContratosTipo.ContratoTipo[x].Preliquidado = "Sí"

		 }

	 }


 }


	return tempDocentes, controlError;

}


// ConsultarDatosPreliq ...
// @Title ConsultarDatosPreliq
// @Description Trae los datos de una preliquidacion dado su ID
func  ConsultarDatosPreliq(id_pre int)(preliq *models.Preliquidacion){
	var datosPreliquidacion []models.Preliquidacion
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/preliquidacion/?query=Id:"+strconv.Itoa(id_pre), &datosPreliquidacion); err == nil && datosPreliquidacion !=nil {
		preliq := &models.Preliquidacion {Id: datosPreliquidacion[0].Id, Descripcion: datosPreliquidacion[0].Descripcion, Mes: datosPreliquidacion[0].Mes,Ano:datosPreliquidacion[0].Ano }

		return preliq

	} else {
		fmt.Println(err)
		fmt.Println("error al consultar preliquidacion")
		return
	}
}
