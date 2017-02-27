package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"titan_api_mid/models"

	"github.com/astaxie/beego"
)

// operations for Liquidar
type DetalleLiquidacionController struct {
	beego.Controller
}

func (c *DetalleLiquidacionController) URLMapping() {
	c.Mapping("detalle_liquidacion", c.InsertarDetallePreliquidacion)
}

func (c *DetalleLiquidacionController) InsertarDetallePreliquidacion() {
	fmt.Println("detalle")
	var v []int
	var tam int
	var IdPreliquidacion int
	var idPersonaString string
	var idPreliquidacionString string
	var d []models.DetallePreliquidacion
	var idDetaLiq interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		tam = len(v)
		IdPreliquidacion = v[tam-1]
		for i := 0; i < len(v)-1; i++ {
			idPersonaString = strconv.Itoa(v[i])
			idPreliquidacionString = strconv.Itoa(IdPreliquidacion)
			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=0&query=Preliquidacion:"+idPreliquidacionString+",Persona:"+idPersonaString+"", &d); err == nil {
				for i := 0; i < len(d); i++ {
					detalleliquidacion := models.DetalleLiquidacion{Id: d[i].Id, ValorCalculado: d[i].ValorCalculado, EstadoConcepto: "P", Liquidacion: &models.Liquidacion{Id: IdPreliquidacion}, Persona: d[i].Persona, Concepto: &models.Concepto{Id: d[i].Concepto.Id}, NumeroContrato: &models.ContratoGeneral{Id: d[i].NumeroContrato.Id}}
					if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_liquidacion", "POST", &idDetaLiq, &detalleliquidacion); err == nil {
					} else {
						beego.Debug("error1: ", err)
					}
				}
			} else {

			}
		}

		//http://localhost:8082/v1/detalle_preliquidacion?limit=0&query=Preliquidacion:7,Persona:184
	}
}
