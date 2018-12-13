package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"
	//"time"
	"github.com/astaxie/beego"
)

// GestionContratosController operations for Preliquidacion
type GestionContratosController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionContratosController) URLMapping() {
	c.Mapping("ListarContratosAgrupados", c.ListarContratosAgrupadosPorPersona)
}



// Post ...
// @Title Create
// @Description create ListarContratosAgrupados: Lista por persona los contratos que tiene vigentes. Para el caso de los docentes HC, agrupará los que sean de la misma resolución
// @Param	body 	models.DatosPreliquidacion	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /listar_contratos_agrupados_por_persona/ [post]
func (c *GestionContratosController) ListarContratosAgrupadosPorPersona() {
	var v models.DatosPreliquidacion;
	var control_error error
	var temp_docentes models.ObjetoFuncionarioContrato


	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		if(v.Preliquidacion.Nomina.TipoNomina.Nombre == "CT") {
			temp_docentes, control_error = GetContratosPorPersonaCT(v, v.PersonasPreLiquidacion[0])
		}
		//Buscar contratos vigentes en ese periodo para esa persona
		if(v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCH") {
			temp_docentes, control_error = GetContratosPorPersonaHCH(v, v.PersonasPreLiquidacion[0])
		}

		if(v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCS") {
			temp_docentes, control_error = GetContratosPorPersonaHCS(v,v.PersonasPreLiquidacion[0])
		}

		if (control_error == nil){
			if (v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCH" || v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCS"){
			temp_agrupar := make(map[string]interface{}) // este mapa tiene la siguiente estructura: temp_agrupar[numero_cedula_docente][id_resolucion][valor_total] (cada resolucion tiene un único tipo de nivel académico, por lo tanto los valores totales se van sumando de acuerdo a la resolución )
			info_contratos := make(map[string]interface{})

				for x,dato := range temp_docentes.ContratosTipo.ContratoTipo {

							var vinculaciones []models.VinculacionDocente
							query:= "NumeroContrato:"+dato.NumeroContrato+",Vigencia:"+dato.VigenciaContrato
							if err := getJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {


								info_contrato := make(map[string]interface{})
								info_contrato["VigenciaContrato"] = dato.VigenciaContrato
								info_contrato["NumeroContrato"] = dato.NumeroContrato
								info_contrato["NivelAcademico"] = vinculaciones[0].IdResolucion.NivelAcademico
								info_contrato["Resolucion"] = vinculaciones[0].IdResolucion.Id
								info_contratos[strconv.Itoa(x)] = info_contrato;

							}
				}
				temp_agrupar["Contratos"] =  info_contratos
				c.Data["json"] = temp_agrupar
			}

			if (v.Preliquidacion.Nomina.TipoNomina.Nombre == "CT"){
				temp_agrupar := make(map[string]interface{}) // este mapa tiene la siguiente estructura: temp_agrupar[numero_cedula_docente][id_resolucion][valor_total] (cada resolucion tiene un único tipo de nivel académico, por lo tanto los valores totales se van sumando de acuerdo a la resolución )
				info_contratos := make(map[string]interface{})
				for x, dato := range temp_docentes.ContratosTipo.ContratoTipo {
					info_contrato := make(map[string]interface{})
					info_contrato["VigenciaContrato"] = dato.VigenciaContrato
					info_contrato["NumeroContrato"] = dato.NumeroContrato
					info_contratos[strconv.Itoa(x)] = info_contrato;
				}

				temp_agrupar["Contratos"] =  info_contratos
				c.Data["json"] = temp_agrupar
			}


		}else{
			c.Data["json"] = control_error
		}

			c.ServeJSON()


}else{
	fmt.Println("Error al leer datos", control_error)
}
}


func ListaContratosContratistas(objeto_nom models.Preliquidacion)(arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error){

	var temp map[string]interface{}
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


	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo/6/"+ano+"-"+mes+"/"+ano+"-"+mes, &temp); err == nil && temp != nil {
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

func ListaContratosFuncionariosPlanta()(arreglo_contratos []models.Funcionario_x_Proveedor, e error){
	var err error
	var datos_planta []models.Funcionario_x_Proveedor
	if err = getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/funcionario_proveedor/get_funcionarios_planta", &datos_planta); err == nil && datos_planta !=nil {
		fmt.Println("funcionario")
		fmt.Println(datos_planta)
	}

	return datos_planta, err
}

func GetContratosPorPersonaHCH(v models.DatosPreliquidacion,docente models.PersonasPreliquidacion) (arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error){

	var temp map[string]interface{}

	var tipo_nom string
	var temp_docentes models.ObjetoFuncionarioContrato
	var control_error error
	var mes string
	var ano string
	var persona string

	ano = strconv.Itoa(v.Preliquidacion.Ano);
  persona = strconv.Itoa(docente.NumDocumento)

	if (v.Preliquidacion.Mes >= 1 && v.Preliquidacion.Mes <= 9){
				mes = strconv.Itoa(v.Preliquidacion.Mes);
				mes = "0"+mes
			}else{
				mes = strconv.Itoa(v.Preliquidacion.Mes);
			}

	tipo_nom = "3"
	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_persona/"+tipo_nom+"/"+ano+"-"+mes+"/"+ano+"-"+mes+"/"+persona, &temp); err == nil && temp != nil {
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

func GetContratosPorPersonaCT(v models.DatosPreliquidacion,docente models.PersonasPreliquidacion) (arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error){

	var temp map[string]interface{}

	var tipo_nom string
	var temp_docentes models.ObjetoFuncionarioContrato
	var control_error error
	var mes string
	var ano string
	var persona string

	ano = strconv.Itoa(v.Preliquidacion.Ano);
  persona = strconv.Itoa(docente.NumDocumento)

	if (v.Preliquidacion.Mes >= 1 && v.Preliquidacion.Mes <= 9){
				mes = strconv.Itoa(v.Preliquidacion.Mes);
				mes = "0"+mes
			}else{
				mes = strconv.Itoa(v.Preliquidacion.Mes);
			}

	tipo_nom = "6"
	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_persona/"+tipo_nom+"/"+ano+"-"+mes+"/"+ano+"-"+mes+"/"+persona, &temp); err == nil && temp != nil {
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


func GetContratosPorPersonaHCS(v models.DatosPreliquidacion,docente models.PersonasPreliquidacion) (arreglo_contratos models.ObjetoFuncionarioContrato, cont_error error){

	var temp map[string]interface{}

	var tipo_nom string
	var temp_docentes models.ObjetoFuncionarioContrato
	var temp_docentes_tco models.ObjetoFuncionarioContrato
	var control_error error
	var mes string
	var ano string
	var persona string

	ano = strconv.Itoa(v.Preliquidacion.Ano);
  persona = strconv.Itoa(docente.NumDocumento)

	if (v.Preliquidacion.Mes >= 1 && v.Preliquidacion.Mes <= 9){
				mes = strconv.Itoa(v.Preliquidacion.Mes);
				mes = "0"+mes
			}else{
				mes = strconv.Itoa(v.Preliquidacion.Mes);
			}

			tipo_nom = "2"
			if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_persona/"+tipo_nom+"/"+ano+"-"+mes+"/"+ano+"-"+mes+"/"+persona, &temp); err == nil && temp != nil {
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
			if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contratos_elaborado_tipo_persona/"+tipo_nom+"/"+ano+"-"+mes+"/"+ano+"-"+mes+"/"+persona, &temp); err == nil && temp != nil {
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

		return temp_docentes, control_error;
}
