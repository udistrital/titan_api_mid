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
// @Description create ListarPersonasAPreliquidar
// @Param	body 	models.Nomina	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /listar_personas_a_preliquidar_argo/ [post]
func (c *GestionPersonasAPreliquidarController) ListarPersonasAPreliquidar() {
	var v models.Nomina
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if v.TipoNomina.Nombre == "CT" {

			if listaContratos, err := ListaContratosContratistas(v); err == nil {
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = listaContratos.ContratosTipo.ContratoTipo
			} else {
				c.Data["json"] = err.Error()
				fmt.Println("error : ", err)
			}


		} else if v.TipoNomina.Nombre == "HCS" || v.TipoNomina.Nombre == "HCH" {
				if listaContratos, err := ListaContratosDocentesDVE(v); err == nil {
					c.Ctx.Output.SetStatus(201)
					c.Data["json"] = listaContratos.ContratosTipo.ContratoTipo
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
// @Description create ListarPersonasAPreliquidar
// @Param	body 	models.Nomina	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /listar_personas_a_preliquidar_pendientes/ [post]
func (c *GestionPersonasAPreliquidarController) ListarPersonasAPreliquidarPendientes() {
	var v models.Nomina
	var personas_pend_preliquidacion []models.DetallePreliquidacion
	var error_consulta_informacion_agora error

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion/get_personas_pago_pendiente?idNomina="+strconv.Itoa(v.Id), &personas_pend_preliquidacion); err == nil && personas_pend_preliquidacion !=nil {

			fmt.Println("personas pendientes", personas_pend_preliquidacion)

			for x, dato := range personas_pend_preliquidacion {
				personas_pend_preliquidacion[x].NombreCompleto, _, personas_pend_preliquidacion[x].Documento, error_consulta_informacion_agora= InformacionContratista(dato.NumeroContrato, dato.VigenciaContrato)
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

func ListaContratosDocentesDVE(objeto_nom models.Nomina)(arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error){

	var temp map[string]interface{}
	var tipo_nom string
	var temp_docentes models.ObjetoFuncionarioContrato
	var control_error error

	if(objeto_nom.TipoNomina.Nombre == "HCH") {
		tipo_nom = "3"
	}else {
		tipo_nom = "2"
	}

	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo/"+tipo_nom, &temp); err == nil && temp != nil {
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

func ListaContratosContratistas(objeto_nom models.Nomina)(arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error){
	fmt.Println("contratistas")
	var temp map[string]interface{}
	var temp_docentes models.ObjetoFuncionarioContrato
	var control_error error

	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_tipo/6", &temp); err == nil && temp != nil {
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
