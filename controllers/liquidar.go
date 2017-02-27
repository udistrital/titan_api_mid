package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"titan_api_mid/models"

	"github.com/astaxie/beego"
)

// operations for Liquidar
type LiquidarController struct {
	beego.Controller
}

func (c *LiquidarController) URLMapping() {
	c.Mapping("Liquidar", c.Liquidar)
}

func (c *LiquidarController) Liquidar() {
	var idLiquidacion interface{}
	var errores []string
	var v models.DatosLiquidacion
	var d []models.DetallePreliquidacion
	var idL interface{}
	var idDetaLiq interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		liquidacion := models.Liquidacion{Id: v.Preliquidacion.Id, NombreLiquidacion: v.Preliquidacion.Nombre, Nomina: &models.Nomina{Id: v.Preliquidacion.Nomina.Id}, EstadoLiquidacion: "L", FechaLiquidacion: time.Now(), FechaInicio: v.Preliquidacion.FechaInicio, FechaFin: v.Preliquidacion.FechaFin}
		fmt.Println(v.Preliquidacion.Liquidada)
		if v.Preliquidacion.Liquidada == "No" {
			Idpreliquidacion := strconv.Itoa(v.Preliquidacion.Id)
			if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/liquidacion", "POST", &idLiquidacion, &liquidacion); err == nil {
				v.Preliquidacion.Liquidada = "Si"
				if err2 := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/preliquidacion/"+Idpreliquidacion, "PUT", &idL, &v.Preliquidacion); err2 == nil {
					fmt.Println("cambio estado preliquidacion")
					for i := 0; i < len(v.Personas)-1; i++ {
						idPersonaString := strconv.Itoa(v.Personas[i])
						if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=0&query=Preliquidacion:"+Idpreliquidacion+",Persona:"+idPersonaString+"", &d); err == nil {
							for i := 0; i < len(d); i++ {
								detalleliquidacion := models.DetalleLiquidacion{Id: d[i].Id, ValorCalculado: d[i].ValorCalculado, EstadoConcepto: "P", Liquidacion: &models.Liquidacion{Id: v.Preliquidacion.Id}, Persona: d[i].Persona, Concepto: &models.Concepto{Id: d[i].Concepto.Id}, NumeroContrato: &models.ContratoGeneral{Id: d[i].NumeroContrato.Id}}
								if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_liquidacion", "POST", &idDetaLiq, &detalleliquidacion); err == nil {
								} else {
									beego.Debug("error23: ", err)
								}
							}
						}
					}

					/*if err3 := sendJson("http://"+beego.AppConfig.String("Urlmid")+":"+beego.AppConfig.String("Portmid")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_liquidacion/", "POST", &idL, &v.Personas); err3 == nil {
						fmt.Println("cambio estado preliquidacion")
					} else {
						fmt.Print("err3 ")
						fmt.Println(err3)
						errores = append(errores, "error.err3")
					}*/
				} else {
					fmt.Print("err2 ")
					fmt.Println(err2)
					errores = append(errores, "erro.err2")
				}
			} else {
				fmt.Print("err ")
				fmt.Println(err)
				errores = append(errores, "error().err")
			}
		} else {
			fmt.Println("Preliquidacion ya ha sido liquidada!")
			errores = append(errores, "Preliquidacion ya ha sido liquidada")
		}
	}
	fmt.Println(errores)
	if len(errores) > 0 {
		c.Data["json"] = errores
	} else {
		c.Data["json"] = "Ok"
	}
	c.ServeJSON()
}
