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

// GestionContratosController operations for GestionContratos
type GestionContratosController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionContratosController) URLMapping() {
	c.Mapping("ListarContratosAgrupados", c.ListarContratosAgrupadosPorPersona)
}

// ListarContratosAgrupadosPorPersona ...
// @Title Create ListarContratosAgrupadosPorPersona
// @Description Lista por persona los contratos que tiene vigentes. Para el caso de los docentes HC, agrupará los que sean de la misma resolución
// @Param	body		body  models.DatosPreliquidacion	true		"body for models.DatosPreliquidacion content"
// @Success 201
// @Failure 403 body is empty
// @router /listar_contratos_agrupados_por_persona [post]
func (c *GestionContratosController) ListarContratosAgrupadosPorPersona() {
	var v models.DatosPreliquidacion
	var controlError error
	var tempDocentes models.ObjetoFuncionarioContrato

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "CT" {
			tempDocentes, controlError = GetContratosPorPersonaCT(v, v.PersonasPreLiquidacion[0].NumDocumento)
		}
		//Buscar contratos vigentes en ese periodo para esa persona
		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCH" {
			tempDocentes, controlError = GetContratosPorPersonaHCH(v, v.PersonasPreLiquidacion[0])
		}

		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCS" {
			tempDocentes, controlError = GetContratosPorPersonaHCS(v, v.PersonasPreLiquidacion[0])
		}

		if controlError == nil {
			if v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCH" || v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCS" {
				tempAgrupar := make(map[string]interface{}) // este mapa tiene la siguiente estructura: tempAgrupar[numero_cedula_docente][id_resolucion][valor_total] (cada resolucion tiene un único tipo de nivel académico, por lo tanto los valores totales se van sumando de acuerdo a la resolución )
				infoContratos := make(map[string]interface{})

				for x, dato := range tempDocentes.ContratosTipo.ContratoTipo {

					var vinculaciones []models.VinculacionDocente
					query := "NumeroContrato:" + dato.NumeroContrato + ",Vigencia:" + dato.VigenciaContrato
					if err := request.GetJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {

						infoContrato := make(map[string]interface{})
						infoContrato["VigenciaContrato"] = dato.VigenciaContrato
						infoContrato["NumeroContrato"] = dato.NumeroContrato
						infoContrato["NivelAcademico"] = vinculaciones[0].IdResolucion.NivelAcademico
						infoContrato["Resolucion"] = vinculaciones[0].IdResolucion.Id
						infoContratos[strconv.Itoa(x)] = infoContrato

					}
				}
				tempAgrupar["Contratos"] = infoContratos
				c.Data["json"] = tempAgrupar
			}

			if v.Preliquidacion.Nomina.TipoNomina.Nombre == "CT" {
				tempAgrupar := make(map[string]interface{}) // este mapa tiene la siguiente estructura: tempAgrupar[numero_cedula_docente][id_resolucion][valor_total] (cada resolucion tiene un único tipo de nivel académico, por lo tanto los valores totales se van sumando de acuerdo a la resolución )
				infoContratos := make(map[string]interface{})
				for x, dato := range tempDocentes.ContratosTipo.ContratoTipo {
					infoContrato := make(map[string]interface{})
					infoContrato["VigenciaContrato"] = dato.VigenciaContrato
					infoContrato["NumeroContrato"] = dato.NumeroContrato
					infoContratos[strconv.Itoa(x)] = infoContrato
				}

				tempAgrupar["Contratos"] = infoContratos
				c.Data["json"] = tempAgrupar
			}

		} else {
			c.Data["json"] = controlError
		}

		c.ServeJSON()

	} else {
		fmt.Println("Error al leer datos", controlError)
	}
}

// ListaContratosContratistas ...
// @Title ListaContratosContratistas
// @Description Lista los contratos vigentes para esa preliquidaciones para contratistas
func ListaContratosContratistas(objeto_nom models.Preliquidacion) (arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error) {

	var temp map[string]interface{}
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

	if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo/6/"+ano+"-"+mes+"/"+ano+"-"+mes, &temp); err == nil && temp != nil {
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

	for x, dato := range tempDocentes.ContratosTipo.ContratoTipo {
		var d []models.DetallePreliquidacion
		query := "Preliquidacion.Id:" + strconv.Itoa(objeto_nom.Id) + ",Persona:" + dato.Id
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
			if d[0].Id == 0 || len(d) == 0 {
				tempDocentes.ContratosTipo.ContratoTipo[x].Preliquidado = "No"
				tempDocentes.ContratosTipo.ContratoTipo[x].EstadoPago = "No liquidado"
			} else {
				tempDocentes.ContratosTipo.ContratoTipo[x].Preliquidado = "Sí"
				tempDocentes.ContratosTipo.ContratoTipo[x].EstadoPago = d[0].EstadoDisponibilidad.Nombre
			}

		}
	}

	return tempDocentes, controlError

}

// ListaContratosFuncionariosPlanta ...
// @Title ListaContratosFuncionariosPlanta
// @Description Lista a los contratistas que se van a liquidar para la preliquidacion
func ListaContratosFuncionariosPlanta() (arreglo_contratos []models.Funcionario_x_Proveedor, e error) {
	var err error
	var datosPlanta []models.Funcionario_x_Proveedor
	fmt.Println("listar personas de planta")
	if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/funcionario_proveedor/get_funcionarios_planta", &datosPlanta); err == nil && datosPlanta != nil {
		fmt.Println("funcionario")

	} else {
		fmt.Println("error al traer personas de planta", err)
	}

	return datosPlanta, err
}

// GetContratosPorPersonaHCH ...
// @Title GetContratosPorPersonaHCH
// @Description Trae los contratos por docente de hora cátedra honorarios vigentes para esa preliquidacion
func GetContratosPorPersonaHCH(v models.DatosPreliquidacion, docente models.PersonasPreliquidacion) (arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error) {

	var temp map[string]interface{}

	var tipoNom string
	var tempDocentes models.ObjetoFuncionarioContrato
	var controlError error
	var mes string
	var ano string
	var persona string

	ano = strconv.Itoa(v.Preliquidacion.Ano)
	persona = strconv.Itoa(docente.NumDocumento)

	if v.Preliquidacion.Mes >= 1 && v.Preliquidacion.Mes <= 9 {
		mes = strconv.Itoa(v.Preliquidacion.Mes)
		mes = "0" + mes
	} else {
		mes = strconv.Itoa(v.Preliquidacion.Mes)
	}

	tipoNom = "3"
	if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_persona/"+tipoNom+"/"+ano+"-"+mes+"/"+ano+"-"+mes+"/"+persona, &temp); err == nil && temp != nil {
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

// GetContratosPorPersonaCT ...
// @Title GetContratosPorPersonaCT
// @Description Trae los contratos por persona de los contratistas
func GetContratosPorPersonaCT(v models.DatosPreliquidacion, docente int) (arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error) {

	var temp map[string]interface{}

	var tipoNom string
	var tempDocentes models.ObjetoFuncionarioContrato
	var controlError error
	var mes string
	var ano string
	var persona string

	ano = strconv.Itoa(v.Preliquidacion.Ano)
	persona = strconv.Itoa(docente)

	if v.Preliquidacion.Mes >= 1 && v.Preliquidacion.Mes <= 9 {
		mes = strconv.Itoa(v.Preliquidacion.Mes)
		mes = "0" + mes
	} else {
		mes = strconv.Itoa(v.Preliquidacion.Mes)
	}

	tipoNom = "6"
	if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_persona/"+tipoNom+"/"+ano+"-"+mes+"/"+ano+"-"+mes+"/"+persona, &temp); err == nil && temp != nil {
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

// GetContratosPorPersonaHCS ...
// @Title GetContratosPorPersonaHCS
// @Description Trae los contratos por docente de hora cátedra salarios vigentes para esa preliquidacion
func GetContratosPorPersonaHCS(v models.DatosPreliquidacion, docente models.PersonasPreliquidacion) (arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error) {

	var temp map[string]interface{}

	var tipoNom string
	var tempDocentes models.ObjetoFuncionarioContrato
	var tempDocentesTco models.ObjetoFuncionarioContrato
	var controlError error
	var mes string
	var ano string
	var persona string

	ano = strconv.Itoa(v.Preliquidacion.Ano)
	persona = strconv.Itoa(docente.NumDocumento)

	if v.Preliquidacion.Mes >= 1 && v.Preliquidacion.Mes <= 9 {
		mes = strconv.Itoa(v.Preliquidacion.Mes)
		mes = "0" + mes
	} else {
		mes = strconv.Itoa(v.Preliquidacion.Mes)
	}

	tipoNom = "2"
	if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_persona/"+tipoNom+"/"+ano+"-"+mes+"/"+ano+"-"+mes+"/"+persona, &temp); err == nil && temp != nil {
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

	tipoNom = "18"
	if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_persona/"+tipoNom+"/"+ano+"-"+mes+"/"+ano+"-"+mes+"/"+persona, &temp); err == nil && temp != nil {
		jsonDocentes, errorJSON := json.Marshal(temp)

		if errorJSON == nil {

			json.Unmarshal(jsonDocentes, &tempDocentesTco)

		} else {
			controlError = errorJSON
			fmt.Println("error al traer contratos docentes DVE")
		}
	} else {
		controlError = err
		fmt.Println("Error al unmarshal datos de nómina", err)

	}

	tempDocentes.ContratosTipo.ContratoTipo = append(tempDocentes.ContratosTipo.ContratoTipo, tempDocentesTco.ContratosTipo.ContratoTipo...)

	return tempDocentes, controlError
}
