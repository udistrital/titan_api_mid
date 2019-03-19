package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
	"time"
	"github.com/astaxie/beego"
)

// GestionReportesController operations for GestionReportes
type GestionReportesController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionReportesController) URLMapping() {
	c.Mapping("TotalNominaPorProyecto", c.TotalNominaPorProyecto)
	c.Mapping("TotalNominaPorFacultad", c.TotalNominaPorFacultad)
	c.Mapping("DesagregadoNominaPorFacultad", c.DesagregadoNominaPorFacultad)
	c.Mapping("DesagregadoNominaPorProyectoCurricular", c.DesagregadoNominaPorProyectoCurricular)
	c.Mapping("TotalNominaPorDependencia", c.TotalNominaPorDependencia)
	c.Mapping("GetOrdenadoresGasto", c.GetOrdenadoresGasto)


}

// GetOrdenadoresGasto ...
// @Title create GetOrdenadoresGasto
// @Description create GetOrdenadoresGasto
// @Success 201
// @Failure 403 body is empty
// @router /get_ordenadores_gasto [post]
func (c *GestionReportesController) GetOrdenadoresGasto() {
	fmt.Println("funcion ordenadores")

	var temp interface{}

	FechaActual := time.Now();
	var ano = strconv.Itoa(FechaActual.Year());
	var dia = strconv.Itoa(FechaActual.Day());
	var mes string;

	if (int(FechaActual.Month()) >= 1 && int(FechaActual.Month()) <= 9){
		mes = strconv.Itoa(int(FechaActual.Month()));
		mes = "0"+mes
	}else{
		mes = strconv.Itoa(int(FechaActual.Month()));
	}

	if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/lista_ordenadores/"+ano+"-"+mes+"-"+dia+"/"+ano+"-"+mes+"-"+dia, &temp); err == nil  {

		c.Data["json"]= temp.(map[string]interface{})["ListaOrdenadores"].(map[string]interface{})["Ordenadores"]
		fmt.Println(temp.(map[string]interface{})["ListaOrdenadores"].(map[string]interface{})["Ordenadores"])
	} else {
		c.Data["json"] = err
		fmt.Println("Error al traer lista de ordenadores",err)


	}


	c.ServeJSON()

}
// TotalNominaPorProyecto ...
// @Title create TotalNominaPorProyecto
// @Description create TotalNominaPorProyecto
// @Param	body 	body models.ObjetoReporte	true		"body for models.ObjetoReporte	true content"
// @Success 201
// @Failure 403 body is empty
// @router /total_nomina_por_proyecto [post]
func (c *GestionReportesController) TotalNominaPorProyecto() {
	fmt.Println("funcion")

	var v models.ObjetoReporte
	var d []models.DetallePreliquidacion
	var vinculaciones []models.VinculacionDocente
	var total float64;
	var totalDescuentos float64;
	var contDev = 0;

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		IDNomina := strconv.Itoa(v.Preliquidacion.Nomina.Id)
		proyectoCurricular := strconv.Itoa(v.ProyectoCurricular)

		query:= "IdProyectoCurricular:"+proyectoCurricular
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {
			fmt.Println("hola soy el total de vinculaciones para ese proyecto", len(vinculaciones))
			for _, pos := range vinculaciones {

				IDProveedor := GetIDProveedor(pos.IdPersona)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina+",Persona:"+strconv.Itoa(IDProveedor)+",Concepto.NaturalezaConcepto.Id:1"
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){
							if(d[0].Concepto.Id == 11 || d[0].Concepto.Id == 10) {

									contDev = contDev + 1;
							}

							total = total + d[0].ValorCalculado

					}

					}else{
						fmt.Println("error al traer valor calculado por devengos",err)
					}

				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina+",Persona:"+strconv.Itoa(IDProveedor)+",Concepto.NaturalezaConcepto.Id:2"
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query2, &d); err == nil {
					if(d != nil){
							totalDescuentos = totalDescuentos + d[0].ValorCalculado

					}

					}else{
						fmt.Println("error al traer valor calculado por descuentos",err)
					}
			}

			v.TotalDev = total
			v.TotalDesc =  totalDescuentos
			v.TotalDocentes = contDev
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


// TotalNominaPorFacultad ...
// @Title create TotalNominaPorFacultad
// @Description create TotalNominaPorFacultad
// @Param	body 	body models.ObjetoReporte	true		"body for models.ObjetoReporte content"
// @Success 201
// @Failure 403 body is empty
// @router /total_nomina_por_facultad [post]
func (c *GestionReportesController) TotalNominaPorFacultad() {
	fmt.Println("funcion")

	var v models.ObjetoReporte
	var d []models.DetallePreliquidacion
	var vinculaciones []models.VinculacionDocente
	var total float64;
	var totalDescuentos float64;
	var contDev = 0;

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		IDNomina := strconv.Itoa(v.Preliquidacion.Nomina.Id)
		facultad := strconv.Itoa(v.Facultad)

		query:= "IdResolucion.IdFacultad:"+facultad
		fmt.Println("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query)
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {
			fmt.Println("hola soy el total de vinculaciones para esa facultad", len(vinculaciones))
			for _, pos := range vinculaciones {

				IDProveedor := GetIDProveedor(pos.IdPersona)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina+",Persona:"+strconv.Itoa(IDProveedor)+",Concepto.NaturalezaConcepto.Id:1"
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){
							if(d[0].Concepto.Id == 11 || d[0].Concepto.Id == 10) {
									contDev = contDev + 1;
							}
							total = total + d[0].ValorCalculado
					}

					}else{
						fmt.Println("error al traer valor calculado por devengos",err)
					}
				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina+",Persona:"+strconv.Itoa(IDProveedor)+",Concepto.NaturalezaConcepto.Id:2"
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query2, &d); err == nil {
					if(d != nil){

							totalDescuentos = totalDescuentos + d[0].ValorCalculado

					}

					}else{
						fmt.Println("error al traer valor calculado por descuentos",err)
					}
			}

			v.TotalDev = total
			v.TotalDesc =  totalDescuentos
			v.TotalDocentes = contDev
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



// DesagregadoNominaPorFacultad ...
// @Title create DesagregadoNominaPorFacultad
// @Description create DesagregadoNominaPorFacultad
// @Param	body  body models.ObjetoReporte	true		"body for models.ObjetoReporte content"
// @Success 201
// @Failure 403 body is empty
// @router /desagregado_nomina_por_facultad [post]
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
		IDNomina := strconv.Itoa(v.Preliquidacion.Nomina.Id)
		facultad := strconv.Itoa(v.Facultad)


		query:= "IdResolucion.IdFacultad:"+facultad
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {
			fmt.Println("hola soy el total de vinculaciones para ese proyecto", len(vinculaciones))
			for _, pos := range vinculaciones {

				IDProveedor := GetIDProveedor(pos.IdPersona)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina+",Persona:"+strconv.Itoa(IDProveedor)+",Concepto.NaturalezaConcepto.Id:1"
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &devengos); err == nil {
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



				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina+",Persona:"+strconv.Itoa(IDProveedor)+",Concepto.NaturalezaConcepto.Id:2"

				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query2, &descuentos); err == nil {
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


// DesagregadoNominaPorProyectoCurricular ...
// @Title create DesagregadoNominaPorProyectoCurricular
// @Description create DesagregadoNominaPorProyectoCurricular
// @Param	body 	body models.ObjetoReporte	true		"body for models.ObjetoReporte content"
// @Success 201
// @Failure 403 body is empty
// @router /desagregado_nomina_por_pc [post]
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
		IDNomina := strconv.Itoa(v.Preliquidacion.Nomina.Id)
		proyectoCurricular := strconv.Itoa(v.ProyectoCurricular)

		query:= "IdProyectoCurricular:"+proyectoCurricular
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {
			fmt.Println("hola soy el total de vinculaciones para ese proyecto", len(vinculaciones))
			for _, pos := range vinculaciones {

				IDProveedor := GetIDProveedor(pos.IdPersona)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina+",Persona:"+strconv.Itoa(IDProveedor)+",Concepto.NaturalezaConcepto.Id:1"
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &devengos); err == nil {
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



				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina+",Persona:"+strconv.Itoa(IDProveedor)+",Concepto.NaturalezaConcepto.Id:2"

				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query2, &descuentos); err == nil {
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


// TotalNominaPorDependencia ...
// @Title create TotalNominaPorDependencia
// @Description create TotalNominaPorFacultad
// @Param	body 	body models.ObjetoReporte	true		"body for models.ObjetoReporte content"
// @Success 201
// @Failure 403 body is empty
// @router /total_nomina_por_dependencia [post]
func (c *GestionReportesController) TotalNominaPorDependencia() {
	fmt.Println("funcion dependencia")

	var v models.ObjetoReporte
	var d []models.DetallePreliquidacion
	var tempDocentes models.ObjetoInformacionContratista
	var temp map[string]interface{}
	arregloTotal := make([]models.DetallePreliquidacion, 0)

	var total float64;
	var totalDescuentos float64;
	var contDev = 0;

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		IDNomina := "3"
		dependencia := v.Dependencia

				fmt.Println("Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina, "dependencia "+dependencia)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina;
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){
							for _, dato := range d {
								if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/informacion_contrato_contratista/"+dato.NumeroContrato+"/"+strconv.Itoa(dato.VigenciaContrato), &temp); err == nil {
									jsonDocentes, errorJSON := json.Marshal(temp)
									if errorJSON == nil {
										json.Unmarshal(jsonDocentes, &tempDocentes)
										if(tempDocentes.InformacionContratista.Dependencia.IdDependencia == dependencia){
											dato.NombreCompleto = tempDocentes.InformacionContratista.NombreCompleto
											dato.Documento = tempDocentes.InformacionContratista.Documento.Numero
											arregloTotal = append(arregloTotal,dato)
										}

								} else {
											c.Data["json"] = errorJSON
											fmt.Println("error al traer contratos docentes DVE")
									}
								}

							};

							for _, dato := range arregloTotal {
								 if(dato.Concepto.NaturalezaConcepto.CodigoAbreviacion == "NCN_001"){
									 total = 	total + dato.ValorCalculado
									 if(dato.Concepto.Id == 11 || d[0].Concepto.Id == 10) {
		 									contDev = contDev + 1;
		 								}
								 }
								 if(dato.Concepto.NaturalezaConcepto.CodigoAbreviacion == "NCN_002"){
									 	 totalDescuentos = 	 totalDescuentos + dato.ValorCalculado
								 }
							}



						}

					}else{
						fmt.Println("error al traer preliquidacion",err)
					}


					v.TotalDev = total
					v.TotalDesc =  totalDescuentos
					v.TotalDocentes = contDev
					c.Data["json"] = v
			}		else{
					c.Data["json"] = err
					fmt.Println("error", err)
		}

	c.ServeJSON()

}


// DesagregadoNominaPorDependencia ...
// @Title create DesagregadoNominaPorDependencia
// @Description create TotalNominaPorFacultad
// @Param	body 	body models.ObjetoReporte	true		"body for models.ObjetoReporte content"
// @Success 201
// @Failure 403 body is empty
// @router /desagregado_nomina_por_dependencia [post]
func (c *GestionReportesController) DesagregadoNominaPorDependencia() {
	fmt.Println("funcion dependencia")

	var v models.ObjetoReporte
	var d []models.DetallePreliquidacion
	var tempDocentes models.ObjetoInformacionContratista
	var temp map[string]interface{}
	arregloTotal := make([]models.DetallePreliquidacion, 0)


	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		IDNomina := "3"
		dependencia := v.Dependencia

				fmt.Println("Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina, "dependencia "+dependencia)
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina;
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){
							for _, dato := range d {
								if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/informacion_contrato_contratista/"+dato.NumeroContrato+"/"+strconv.Itoa(dato.VigenciaContrato), &temp); err == nil {
									jsonDocentes, errorJSON := json.Marshal(temp)
									if errorJSON == nil {
										json.Unmarshal(jsonDocentes, &tempDocentes)
										if(tempDocentes.InformacionContratista.Dependencia.IdDependencia == dependencia){
											dato.NombreCompleto = tempDocentes.InformacionContratista.NombreCompleto
											dato.Documento = tempDocentes.InformacionContratista.Documento.Numero
											arregloTotal = append(arregloTotal,dato)
										}

								} else {
											c.Data["json"] = errorJSON
											fmt.Println("error al traer contratos docentes DVE")
									}
								}

							};

						}

					}else{
						fmt.Println("error al traer preliquidacion",err)
					}

					c.Data["json"] = arregloTotal
			}		else{
					c.Data["json"] = err
					fmt.Println("error", err)
		}

	c.ServeJSON()

}



// TotalNominaPorOrdenador ...
// @Title create TotalNominaPorOrdenador
// @Description create TotalNominaPorDependencia
// @Param	body 	body models.ObjetoReporte	true		"body for models.ObjetoReporte content"
// @Success 201 
// @Failure 403 body is empty
// @router /total_nomina_por_ordenador [post]
func (c *GestionReportesController) TotalNominaPorOrdenador() {
	fmt.Println("funcion odenador total")

	var v models.ObjetoReporte
	var d []models.DetallePreliquidacion
	var temp interface{}
	arregloTotal := make([]models.DetallePreliquidacion, 0)

	var total float64;
	var totalDescuentos float64;
	var contDev = 0;

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		IDNomina := "3"
		ordenador := v.Ordenador

				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+IDNomina;
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){
							for _, dato := range d {
								//Buscar contrato elaborado o suscrito ??
								if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contrato/"+dato.NumeroContrato+"/"+strconv.Itoa(dato.VigenciaContrato), &temp); err == nil {
										contrato := temp.(map[string]interface{})["contrato"].(map[string]interface{})["ordenador_gasto"]

										if contrato != nil{
											IDOrdenador:= temp.(map[string]interface{})["contrato"].(map[string]interface{})["ordenador_gasto"].(map[string]interface{})["id"]
											if(IDOrdenador == ordenador){
													NombreCompleto, Documento, _ := InformacionPersonaProveedor(dato.Persona)
													dato.NombreCompleto = NombreCompleto
													dato.Documento = strconv.Itoa(Documento)
													arregloTotal = append(arregloTotal,dato)
												}
										}

								}else{
									  fmt.Println("error al traer id del ordenador",err)
								}

								}

								fmt.Println("arreglo totaaaal", arregloTotal)
								for _, dato := range arregloTotal {
									 if(dato.Concepto.NaturalezaConcepto.CodigoAbreviacion == "NCN_001"){
										 total = 	total + dato.ValorCalculado
										 if(dato.Concepto.Id == 11 || d[0].Concepto.Id == 10) {
												contDev = contDev + 1;
											}
									 }
									 if(dato.Concepto.NaturalezaConcepto.CodigoAbreviacion == "NCN_002"){
											 totalDescuentos = 	 totalDescuentos + dato.ValorCalculado
									 }
								}

							}else{
								fmt.Println("hola no tengo datos de prelis", d)
								};

					}else{
						fmt.Println("error al traer preliquidacion",err)
					}

					v.TotalDev = total
					v.TotalDesc =  totalDescuentos
					v.TotalDocentes = contDev

					c.Data["json"] = v
			}		else{
					c.Data["json"] = err
					fmt.Println("error", err)
		}

	c.ServeJSON()

}
