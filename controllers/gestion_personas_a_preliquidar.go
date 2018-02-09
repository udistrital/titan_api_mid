package controllers

import (
	"encoding/json"
	"fmt"
	//"strconv"
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
// @router /listar_personas_a_preliquidar [post]
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
