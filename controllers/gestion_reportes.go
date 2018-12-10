package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	//"time"
	"github.com/astaxie/beego"
)

// GestionReportesController operations for Preliquidacion
type GestionReportesController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionReportesController) URLMapping() {
	c.Mapping("TotalNominaPorProyecto", c.TotalNominaPorProyecto)
}

// Post ...
// @Title Create
// @Description create TotalNominaPorProyecto
// @Param	body 	models.DetallePreliquidacion	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /total_nomina_por_proyecto/ [post]
func (c *GestionReportesController) TotalNominaPorProyecto() {
	fmt.Println("funcion")

	var v models.ObjetoReporte
	var d []models.DetallePreliquidacion
	var vinculaciones []models.VinculacionDocente
	var total float64;
	var total_descuentos float64;
	var cont_dev = 0;

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		id_nomina := strconv.Itoa(v.Preliquidacion.Nomina.Id)
		proyecto_curricular := strconv.Itoa(v.ProyectoCurricular)

		query:= "IdProyectoCurricular:"+proyecto_curricular
		if err := getJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {
			fmt.Println("hola soy el total de vinculaciones para ese proyecto", len(vinculaciones))
			for _, pos := range vinculaciones {

				IdProveedor := GetIdProveedor(pos.IdPersona)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",Persona:"+strconv.Itoa(IdProveedor)+",Concepto.NaturalezaConcepto.Id:1"
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){
							if(d[0].Concepto.Id == 11) {
									cont_dev = cont_dev + 1;
							}

							total = total + d[0].ValorCalculado

					}

					}else{
						fmt.Println("error al traer valor calculado por devengos",err)
					}

				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",Persona:"+strconv.Itoa(IdProveedor)+",Concepto.NaturalezaConcepto.Id:2"
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query2, &d); err == nil {
					if(d != nil){
							total_descuentos = total_descuentos + d[0].ValorCalculado

					}

					}else{
						fmt.Println("error al traer valor calculado por descuentos",err)
					}
			}

			v.TotalDev = total
			v.TotalDesc =  total_descuentos
			v.TotalDocentes = cont_dev
			c.Data["json"] = v
		}else{
			fmt.Println("error en vinculaciones",err)
		}

		}else{
			c.Data["json"] = err
			fmt.Println("error", err)
		}

	c.ServeJSON()

}


// Post ...
// @Title Create
// @Description create TotalNominaPorFacultad
// @Param	body 	models.DetallePreliquidacion	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /total_nomina_por_facultad/ [post]
func (c *GestionReportesController) TotalNominaPorFacultad() {
	fmt.Println("funcion")

	var v models.ObjetoReporte
	var d []models.DetallePreliquidacion
	var vinculaciones []models.VinculacionDocente
	var total float64;
	var total_descuentos float64;
	var cont_dev = 0;

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		id_nomina := strconv.Itoa(v.Preliquidacion.Nomina.Id)
		facultad := strconv.Itoa(v.Facultad)

		query:= "IdResolucion.IdFacultad:"+facultad
		fmt.Println("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query)
		if err := getJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {
			fmt.Println("hola soy el total de vinculaciones para esa facultad", len(vinculaciones))
			for _, pos := range vinculaciones {

				IdProveedor := GetIdProveedor(pos.IdPersona)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",Persona:"+strconv.Itoa(IdProveedor)+",Concepto.NaturalezaConcepto.Id:1"
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){
							if(d[0].Concepto.Id == 11) {
									cont_dev = cont_dev + 1;
							}
							total = total + d[0].ValorCalculado
					}

					}else{
						fmt.Println("error al traer valor calculado por devengos",err)
					}
				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",Persona:"+strconv.Itoa(IdProveedor)+",Concepto.NaturalezaConcepto.Id:2"
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query2, &d); err == nil {
					if(d != nil){

							total_descuentos = total_descuentos + d[0].ValorCalculado

					}

					}else{
						fmt.Println("error al traer valor calculado por descuentos",err)
					}
			}

			v.TotalDev = total
			v.TotalDesc =  total_descuentos
			v.TotalDocentes = cont_dev
			c.Data["json"] = v
		}else{
			fmt.Println("error en vinculaciones",err)
		}

		}else{
			c.Data["json"] = err
			fmt.Println("error", err)
		}

	c.ServeJSON()

}



// Post ...
// @Title Create
// @Description create DesagregadoNominaPorFacultad
// @Param	body 	models.DetallePreliquidacion	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /desagregado_nomina_por_facultad/ [post]
func (c *GestionReportesController) DesagregadoNominaPorFacultad() {
	fmt.Println("funcion")

	var v models.ObjetoReporte
	var devengos []interface{}
	var descuentos []interface{}
	var res []interface{}
	var vinculaciones []models.VinculacionDocente


	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {


		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		id_nomina := strconv.Itoa(v.Preliquidacion.Nomina.Id)
		facultad := strconv.Itoa(v.Facultad)


		query:= "IdResolucion.IdFacultad:"+facultad
		if err := getJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {
			fmt.Println("hola soy el total de vinculaciones para ese proyecto", len(vinculaciones))
			for _, pos := range vinculaciones {

				IdProveedor := GetIdProveedor(pos.IdPersona)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",Persona:"+strconv.Itoa(IdProveedor)+",Concepto.NaturalezaConcepto.Id:1"
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &devengos); err == nil {
					if(devengos != nil){

						for _, dato := range devengos {
							 aux := models.DetallePreliquidacion{}
							 if err := formatdata.FillStruct(dato, &aux); err == nil{
								 aux.NombreCompleto, _ , aux.Documento, _ = InformacionPersona(v.Preliquidacion.Nomina.TipoNomina.Nombre,aux.NumeroContrato, aux.VigenciaContrato)
								res = append(res, aux)
							 }
						}
					}

					}else{
						fmt.Println("error al traer valor calculado por devengos",err)
					}



				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",Persona:"+strconv.Itoa(IdProveedor)+",Concepto.NaturalezaConcepto.Id:2"

				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query2, &descuentos); err == nil {
					if(descuentos != nil){

						for _, dato := range descuentos {
							aux := models.DetallePreliquidacion{}
							if err := formatdata.FillStruct(dato, &aux); err == nil{
								aux.NombreCompleto, _ , aux.Documento, _ = InformacionPersona(v.Preliquidacion.Nomina.TipoNomina.Nombre,aux.NumeroContrato, aux.VigenciaContrato)
							 res = append(res, aux)
							}
						}
					}

					}else{
						fmt.Println("error al traer valor calculado por descuentos",err)
					}
			}



			c.Data["json"] = res

		}else{
			fmt.Println("error en vinculaciones",err)
		}

		}else{
			c.Data["json"] = err
			fmt.Println("error", err)
		}

	c.ServeJSON()

}


// Post ...
// @Title Create
// @Description create DesagregadoNominaPorProyectoCurricular
// @Param	body 	models.DetallePreliquidacion	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /desagregado_nomina_por_pc/ [post]
func (c *GestionReportesController) DesagregadoNominaPorProyectoCurricular() {
	fmt.Println("funcion")

	var v models.ObjetoReporte
	var devengos []interface{}
	var descuentos []interface{}
	var res []interface{}
	var vinculaciones []models.VinculacionDocente


	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		id_nomina := strconv.Itoa(v.Preliquidacion.Nomina.Id)
		proyecto_curricular := strconv.Itoa(v.ProyectoCurricular)

		query:= "IdProyectoCurricular:"+proyecto_curricular
		if err := getJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {
			fmt.Println("hola soy el total de vinculaciones para ese proyecto", len(vinculaciones))
			for _, pos := range vinculaciones {

				IdProveedor := GetIdProveedor(pos.IdPersona)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",Persona:"+strconv.Itoa(IdProveedor)+",Concepto.NaturalezaConcepto.Id:1"
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &devengos); err == nil {
					if(devengos != nil){
						for _, dato := range devengos {
							aux := models.DetallePreliquidacion{}
							if err := formatdata.FillStruct(dato, &aux); err == nil{
								aux.NombreCompleto, _ , aux.Documento, _ = InformacionPersona(v.Preliquidacion.Nomina.TipoNomina.Nombre,aux.NumeroContrato, aux.VigenciaContrato)
							 res = append(res, aux)
							}
						}
					}

					}else{
						fmt.Println("error al traer valor calculado por devengos",err)
					}



				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",Persona:"+strconv.Itoa(IdProveedor)+",Concepto.NaturalezaConcepto.Id:2"

				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query2, &descuentos); err == nil {
					if(descuentos != nil){

						for _, dato := range descuentos {
							aux := models.DetallePreliquidacion{}
							if err := formatdata.FillStruct(dato, &aux); err == nil{
								aux.NombreCompleto, _ , aux.Documento, _ = InformacionPersona(v.Preliquidacion.Nomina.TipoNomina.Nombre,aux.NumeroContrato, aux.VigenciaContrato)
							 res = append(res, aux)
							}
						}


					}

					}else{
						fmt.Println("error al traer valor calculado por descuentos",err)
					}



			}


			c.Data["json"] = res

		}else{
			fmt.Println("error en vinculaciones",err)
		}

		}else{
			c.Data["json"] = err
			fmt.Println("error", err)
		}

	c.ServeJSON()

}


// Post ...
// @Title Create
// @Description create TotalNominaPorFacultad
// @Param	body 	models.DetallePreliquidacion	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /total_nomina_por_dependencia/ [post]
func (c *GestionReportesController) TotalNominaPorDependencia() {
	fmt.Println("funcion dependencia")

	var v models.ObjetoReporte
	var d []models.DetallePreliquidacion
	var temp_docentes models.ObjetoInformacionContratista
	var temp map[string]interface{}
	arreglo_total := make([]models.DetallePreliquidacion, 0)

	var total float64;
	var total_descuentos float64;
	var cont_dev = 0;

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		id_nomina := "3"
		dependencia := v.Dependencia

				fmt.Println("Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina, "dependencia "+dependencia)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina;
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){
							for _, dato := range d {
								if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/informacion_contrato_contratista/"+dato.NumeroContrato+"/"+strconv.Itoa(dato.VigenciaContrato), &temp); err == nil && temp != nil {
									jsonDocentes, error_json := json.Marshal(temp)
									if error_json == nil {
										json.Unmarshal(jsonDocentes, &temp_docentes)
										if(temp_docentes.InformacionContratista.Dependencia.IdDependencia == dependencia){
											dato.NombreCompleto = temp_docentes.InformacionContratista.NombreCompleto
											dato.Documento = temp_docentes.InformacionContratista.Documento.Numero
											arreglo_total = append(arreglo_total,dato)
										}

								} else {
											c.Data["json"] = error_json
											fmt.Println("error al traer contratos docentes DVE")
									}
								}

							};

							for _, dato := range arreglo_total {
								 if(dato.Concepto.NaturalezaConcepto.CodigoAbreviacion == "NCN_001"){
									 total = 	total + dato.ValorCalculado
									 if(dato.Concepto.Id == 11) {
		 									cont_dev = cont_dev + 1;
		 								}
								 }
								 if(dato.Concepto.NaturalezaConcepto.CodigoAbreviacion == "NCN_002"){
									 	 total_descuentos = 	 total_descuentos + dato.ValorCalculado
								 }
							}



						}

					}else{
						fmt.Println("error al traer preliquidacion",err)
					}


					v.TotalDev = total
					v.TotalDesc =  total_descuentos
					v.TotalDocentes = cont_dev
					c.Data["json"] = v
			}		else{
					c.Data["json"] = err
					fmt.Println("error", err)
		}

	c.ServeJSON()

}


// Post ...
// @Title Create
// @Description create TotalNominaPorFacultad
// @Param	body 	models.DetallePreliquidacion	true		"body for Nomina content"
// @Success 201 {object}
// @Failure 403 body is empty
// @router /desagregado_nomina_por_dependencia/ [post]
func (c *GestionReportesController) DesagregadoNominaPorDependencia() {
	fmt.Println("funcion dependencia")

	var v models.ObjetoReporte
	var d []models.DetallePreliquidacion
	var temp_docentes models.ObjetoInformacionContratista
	var temp map[string]interface{}
	arreglo_total := make([]models.DetallePreliquidacion, 0)


	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		id_nomina := "3"
		dependencia := v.Dependencia

				fmt.Println("Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina, "dependencia "+dependencia)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina;
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){
							for _, dato := range d {
								if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/informacion_contrato_contratista/"+dato.NumeroContrato+"/"+strconv.Itoa(dato.VigenciaContrato), &temp); err == nil && temp != nil {
									jsonDocentes, error_json := json.Marshal(temp)
									if error_json == nil {
										json.Unmarshal(jsonDocentes, &temp_docentes)
										if(temp_docentes.InformacionContratista.Dependencia.IdDependencia == dependencia){
											dato.NombreCompleto = temp_docentes.InformacionContratista.NombreCompleto
											dato.Documento = temp_docentes.InformacionContratista.Documento.Numero
											arreglo_total = append(arreglo_total,dato)
										}

								} else {
											c.Data["json"] = error_json
											fmt.Println("error al traer contratos docentes DVE")
									}
								}

							};

						}

					}else{
						fmt.Println("error al traer preliquidacion",err)
					}

					c.Data["json"] = arreglo_total
			}		else{
					c.Data["json"] = err
					fmt.Println("error", err)
		}

	c.ServeJSON()

}
