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
	c.Mapping("TotalNominaPorFacultad", c.TotalNominaPorFacultad)
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

	var v models.DetallePreliquidacion
	var d []models.DetallePreliquidacion
	var vinculaciones []models.VinculacionDocente
	var total float64;
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		
		ano:= strconv.Itoa(v.Preliquidacion.Ano)
		mes := strconv.Itoa(v.Preliquidacion.Mes)
		id_nomina := strconv.Itoa(v.Preliquidacion.Nomina.Id)

		//query:= de proyecto curricular o facultad o lo que sea
		fmt.Println("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1")
		if err := getJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1", &vinculaciones); err == nil {
			for _, pos := range vinculaciones {
				query := "Preliquidacion.Ano:"+ano+",Preliquidacion.Mes:"+mes+",Preliquidacion.Nomina.Id:"+id_nomina+",NumeroContrato:"+pos.NumeroContrato.String+",VigenciaContrato:"+strconv.Itoa(int(pos.Vigencia.Int64))
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
					if(d != nil){
							total = total + d[0].ValorCalculado
					}
				
					}else{
						fmt.Println("error al traer valor calculado por contrato",err)
					}
			}

			c.Data["json"] = total
		}else{
			fmt.Println("error en vinculaciones",err)
		}
			
		}else{
			c.Data["json"] = err
			fmt.Println("error", err)
		}
	
	c.ServeJSON()

}
