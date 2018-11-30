package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"
	//"time"
	"github.com/astaxie/beego"
)

// GestionPersonasAPreliquidarController operations for Preliquidacion
type GestionPersonasAPreliquidarController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionPersonasAPreliquidarController) URLMapping() {
	c.Mapping("ListarPersonasAPreliquidar", c.ListarPersonasAPreliquidar)
}

// Post ...
// @Title Create
// @Description create ListarPersonasAPreliquidar: Lista a las personas que tienen vinculaciones activas para ese periodo y que por consiguiente pueden ser preliquidadas
// @Param	body 	models.Preliquidacion	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /listar_personas_a_preliquidar_argo/ [post]
func (c *GestionPersonasAPreliquidarController) ListarPersonasAPreliquidar() {
	var v models.Preliquidacion
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if v.Nomina.TipoNomina.Nombre == "CT" {

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

// Post ...
// @Title Create
// @Description create ListarPersonasAPreliquidar: Lista a las personas pendientes de periodos anteriores para que puedan ser tenidas en cuenta el presente mes
// @Param	body 	models.Preliquidacion	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /listar_personas_a_preliquidar_pendientes/ [post]
func (c *GestionPersonasAPreliquidarController) ListarPersonasAPreliquidarPendientes() {
	var v models.Preliquidacion
	var personas_pend_preliquidacion []models.DetallePreliquidacion
	var error_consulta_informacion_agora error

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion/get_personas_pago_pendiente?idNomina="+strconv.Itoa(v.Nomina.Id), &personas_pend_preliquidacion); err == nil && personas_pend_preliquidacion !=nil {



			for x, dato := range personas_pend_preliquidacion {
				personas_pend_preliquidacion[x].NombreCompleto, _, personas_pend_preliquidacion[x].Documento, error_consulta_informacion_agora= InformacionPersona(v.Nomina.TipoNomina.Nombre,dato.NumeroContrato, dato.VigenciaContrato)
				personas_pend_preliquidacion[x].Preliquidacion = Consultar_datos_preliq(dato.Preliquidacion.Id)
			}

			if(error_consulta_informacion_agora == nil){
				c.Data["json"] = personas_pend_preliquidacion
			}else{
				c.Data["json"] = error_consulta_informacion_agora
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
// @Description create ListarPersonasHCS: Trae las personas que tienen contratos vigentes para ese mes en Hora Cátedra Salarios

func ListarPersonasHCS(objeto_nom models.Preliquidacion)(arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error){

	var temp map[string]interface{}

	var tipo_nom string
	var temp_docentes models.ObjetoFuncionarioContrato
	var temp_docentes_tco models.ObjetoFuncionarioContrato
	var control_error error
	var mes string
	var ano = strconv.Itoa(objeto_nom.Ano);

	if (objeto_nom.Mes >= 1 && objeto_nom.Mes <= 9){
		mes = strconv.Itoa(objeto_nom.Mes);
		mes = "0"+mes
	}else{
		mes = strconv.Itoa(objeto_nom.Mes);
	}



	tipo_nom = "2"

		if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_personas/"+tipo_nom+"/"+ano+"-"+mes+"/"+ano+"-"+mes, &temp); err == nil && temp != nil {
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

		tipo_nom = "18"

		if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_personas/"+tipo_nom+"/"+ano+"-"+mes+"/"+ano+"-"+mes, &temp); err == nil && temp != nil {
			jsonDocentes, error_json := json.Marshal(temp)

			if error_json == nil {

				json.Unmarshal(jsonDocentes, &temp_docentes_tco)

			} else {
				control_error = error_json
				fmt.Println("error al traer contratos docentes DVE")
			}
		} else {
			control_error = err
			fmt.Println("Error al unmarshal datos de nómina",err)


		}

    temp_docentes.ContratosTipo.ContratoTipo = append(temp_docentes.ContratosTipo.ContratoTipo, temp_docentes_tco.ContratosTipo.ContratoTipo...)


	var d []models.DetallePreliquidacion
	for x, dato := range temp_docentes.ContratosTipo.ContratoTipo {
		query := "Preliquidacion.Id:"+strconv.Itoa(objeto_nom.Id)+",Persona:"+dato.Id
		if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
			if len(d) == 0 {
				temp_docentes.ContratosTipo.ContratoTipo[x].Preliquidado = "no"

			}else{
				temp_docentes.ContratosTipo.ContratoTipo[x].Preliquidado = "si"

			}

		}


	}


	return temp_docentes, control_error;

}

// ListarPersonasHCH ...
// @Title  ListarPersonasHCH
// @Description create ListarPersonasHCH: Trae las personas que tienen contratos vigentes para ese mes en Hora Cátedra Honorarios

func ListarPersonasHCH(objeto_nom models.Preliquidacion)(arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error){

	var temp map[string]interface{}

	var tipo_nom string
	var temp_docentes models.ObjetoFuncionarioContrato
  var control_error error
	var mes string
	var ano = strconv.Itoa(objeto_nom.Ano);

	if (objeto_nom.Mes >= 1 && objeto_nom.Mes <= 9){
		mes = strconv.Itoa(objeto_nom.Mes);
		mes = "0"+mes
	}else{
		mes = strconv.Itoa(objeto_nom.Mes);
	}



	tipo_nom = "3"

		if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_personas/"+tipo_nom+"/"+ano+"-"+mes+"/"+ano+"-"+mes, &temp); err == nil && temp != nil {
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

 //SABER SI YA FUE PRELIQUIDADO O NO
 var d []models.DetallePreliquidacion
 for x, dato := range temp_docentes.ContratosTipo.ContratoTipo {
	 query := "Preliquidacion.Id:"+strconv.Itoa(objeto_nom.Id)+",Persona:"+dato.Id
	 if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
		 if len(d) == 0 {
			 temp_docentes.ContratosTipo.ContratoTipo[x].Preliquidado = "no"

		 }else{
			 temp_docentes.ContratosTipo.ContratoTipo[x].Preliquidado = "si"

		 }

	 }


 }


	return temp_docentes, control_error;

}


func  Consultar_datos_preliq(id_pre int)(preliq *models.Preliquidacion){
	var datos_preliquidacion []models.Preliquidacion
	if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/preliquidacion/?query=Id:"+strconv.Itoa(id_pre), &datos_preliquidacion); err == nil && datos_preliquidacion !=nil {
		preliq := &models.Preliquidacion {Id: datos_preliquidacion[0].Id, Descripcion: datos_preliquidacion[0].Descripcion, Mes: datos_preliquidacion[0].Mes,Ano:datos_preliquidacion[0].Ano }

		return preliq

	} else {
		fmt.Println(err)
		fmt.Println("error al consultar preliquidacion")
		return
	}
}
