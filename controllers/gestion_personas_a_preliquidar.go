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
	/*
		var v models.Preliquidacion
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			if v.NominaId == 414 {
				fmt.Println("preliq", v)
				if listaContratos, err := ListaContratosContratistas(v); err == nil {
					c.Ctx.Output.SetStatus(201)
					c.Data["json"] = listaContratos.ContratosTipo.ContratoTipo
				} else {
					c.Data["json"] = err.Error()
					fmt.Println("error : ", err)
				}

			} else if v.NominaId == 416 {
				if listaContratos, err := ListarPersonasHCS(v); err == nil {
					c.Ctx.Output.SetStatus(201)
					c.Data["json"] = listaContratos.ContratosTipo.ContratoTipo
				} else {
					c.Data["json"] = err.Error()
					fmt.Println("error : ", err)
				}
			} else if v.NominaId == 415 {
				if listaContratos, err := ListarPersonasHCH(v); err == nil {
					c.Ctx.Output.SetStatus(201)
					c.Data["json"] = listaContratos.ContratosTipo.ContratoTipo
				} else {
					c.Data["json"] = err.Error()
					fmt.Println("error : ", err)
				}
			} else if v.NominaId == 412 {
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
	*/

}

// ListarPersonasAPreliquidarPendientes ...
// @Title create ListarPersonasAPreliquidarPendientes
// @Description create ListarPersonasAPreliquidar: Lista a las personas pendientes de periodos anteriores para que puedan ser tenidas en cuenta el presente mes
// @Param	body 	body models.Preliquidacion	true		"body for models.Preliquidacion content"
// @Success 201
// @Failure 403 body is empty
// @router /listar_personas_a_preliquidar_pendientes [post]

// ListarPersonasHCS ...
// @Title  ListarPersonasHCS
// @Description ListarPersonasHCS: Trae las personas que tienen contratos vigentes para ese mes en Hora Cátedra Salarios
func ListarPersonasHCS(objeto_nom models.Preliquidacion) (arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error) {

	var temp map[string]interface{}
	var tempDocentes models.ObjetoFuncionarioContrato
	var controlError error
	var mes string

	if objeto_nom.Mes >= 1 && objeto_nom.Mes <= 9 {
		mes = strconv.Itoa(objeto_nom.Mes)
		mes = "0" + mes
	} else {
		mes = strconv.Itoa(objeto_nom.Mes)
	}

	if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/personas_preliquidacion/"+strconv.Itoa(objeto_nom.Id), &temp); err == nil && temp != nil {
		jsonDocentes, errorJSON := json.Marshal(temp)

		if errorJSON == nil {

			json.Unmarshal(jsonDocentes, &tempDocentes)

		} else {
			controlError = errorJSON
			fmt.Println("error al traer contratos docentes DVE")
		}
	} else {
		controlError = err
		fmt.Println("Error al unmarshal datos de nómina", err)

	}

	return tempDocentes, controlError

}

// ListarPersonasHCH ...
// @Title  ListarPersonasHCH
// @Description ListarPersonasHCH: Trae las personas que tienen contratos vigentes para ese mes en Hora Cátedra Honorarios
func ListarPersonasHCH(objeto_nom models.Preliquidacion) (arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error) {

	var temp map[string]interface{}

	var tipoNom string
	var tempDocentes models.ObjetoFuncionarioContrato
	var controlError error
	var mes string
	var ano = strconv.Itoa(objeto_nom.Ano)

	if objeto_nom.Mes >= 1 && objeto_nom.Mes <= 9 {
		mes = strconv.Itoa(objeto_nom.Mes)
		mes = "0" + mes
	} else {
		mes = strconv.Itoa(objeto_nom.Mes)
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
		fmt.Println("Error al unmarshal datos de nómina", err)
	}

	//SABER SI YA FUE PRELIQUIDADO O NO
	var d []models.DetallePreliquidacion
	for x, dato := range tempDocentes.ContratosTipo.ContratoTipo {
		d = nil
		query := "Preliquidacion.Id:" + strconv.Itoa(objeto_nom.Id) + ",Persona:" + dato.Id
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
			if len(d) == 0 || d[0].Id == 0 {
				tempDocentes.ContratosTipo.ContratoTipo[x].Preliquidado = "NO"

			} else {
				tempDocentes.ContratosTipo.ContratoTipo[x].Preliquidado = "SI"

			}

		}

	}

	return tempDocentes, controlError

}

// ConsultarDatosPreliq ...
// @Title ConsultarDatosPreliq
// @Description Trae los datos de una preliquidacion dado su ID
func ConsultarDatosPreliq(id_pre int) (preliq *models.Preliquidacion) {
	var datosPreliquidacion []models.Preliquidacion
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/preliquidacion/?query=Id:"+strconv.Itoa(id_pre), &datosPreliquidacion); err == nil && datosPreliquidacion != nil {
		preliq := &models.Preliquidacion{Id: datosPreliquidacion[0].Id, Descripcion: datosPreliquidacion[0].Descripcion, Mes: datosPreliquidacion[0].Mes, Ano: datosPreliquidacion[0].Ano}

		return preliq

	} else {
		fmt.Println(err)
		fmt.Println("error al consultar preliquidacion")
		return
	}
}
