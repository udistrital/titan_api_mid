package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"
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
		fmt.Println("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query)
		if err := getJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {
			fmt.Println("hola soy el total de vinculaciones para ese proyecto", len(vinculaciones))
			for _, pos := range vinculaciones {
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",NumeroContrato:"+pos.NumeroContrato.String+",VigenciaContrato:"+strconv.Itoa(int(pos.Vigencia.Int64))+",Concepto.NaturalezaConcepto.Id:1"
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

				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",NumeroContrato:"+pos.NumeroContrato.String+",VigenciaContrato:"+strconv.Itoa(int(pos.Vigencia.Int64))+",Concepto.NaturalezaConcepto.Id:2"

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
			fmt.Println("hola soy el total de vinculaciones para ese proyecto", len(vinculaciones))
			for _, pos := range vinculaciones {
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",NumeroContrato:"+pos.NumeroContrato.String+",VigenciaContrato:"+strconv.Itoa(int(pos.Vigencia.Int64))+",Concepto.NaturalezaConcepto.Id:1"
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

				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",NumeroContrato:"+pos.NumeroContrato.String+",VigenciaContrato:"+strconv.Itoa(int(pos.Vigencia.Int64))+",Concepto.NaturalezaConcepto.Id:2"

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
	var d []models.DetallePreliquidacion
	var res []models.DetallePreliquidacion
	var vinculaciones []models.VinculacionDocente


	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		fmt.Println("objeto",v)
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		id_nomina := strconv.Itoa(v.Preliquidacion.Nomina.Id)
		facultad := strconv.Itoa(v.Facultad)

		query:= "IdResolucion.IdFacultad:"+facultad
		fmt.Println("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query)
		if err := getJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {
			fmt.Println("hola soy el total de vinculaciones para ese proyecto", len(vinculaciones))
			for _, pos := range vinculaciones {
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",NumeroContrato:"+pos.NumeroContrato.String+",VigenciaContrato:"+strconv.Itoa(int(pos.Vigencia.Int64))+",Concepto.NaturalezaConcepto.Id:1"
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){

							res = append(res,d...)
							fmt.Println(res)
					}

					}else{
						fmt.Println("error al traer valor calculado por devengos",err)
					}

				query2 := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",NumeroContrato:"+pos.NumeroContrato.String+",VigenciaContrato:"+strconv.Itoa(int(pos.Vigencia.Int64))+",Concepto.NaturalezaConcepto.Id:2"

				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query2, &d); err == nil {
					if(d != nil){

						res = append(res,d...)


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
