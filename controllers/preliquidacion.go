package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"titan_api_mid/models"

	"github.com/astaxie/beego"
)

// PreliquidacionController operations for Preliquidacion
type PreliquidacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *PreliquidacionController) URLMapping() {
	c.Mapping("Preliquidar", c.Preliquidar)
}

// Post ...
// @Title Create
// @Description create Preliquidacion
// @Param	body		body 	models.Preliquidacion	true		"body for Preliquidacion content"
// @Success 201 {object} models.Preliquidacion
// @Failure 403 body is empty
// @router / [post]
func (c *PreliquidacionController) Preliquidar() {
	var v models.DatosPreliquidacion
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//carga de reglas desde el ruler
		reglasbase := CargarReglasBase(v.Preliquidacion.Nomina.TipoNomina.Nombre) //funcion general para dar formato a reglas cargadas desde el ruler

		//-----------------------------
		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "HC" || v.Preliquidacion.Nomina.TipoNomina.Nombre == "HC-SALARIOS" {
			var n *PreliquidacionHcController
			resumen := n.Preliquidar(&v, reglasbase)
			//pr := CargarNovedadesPersona(v[0].PersonasPreLiquidacion[0].IdPersona,&v[0])
			//fmt.Println("prueba: ", pr)
			c.Data["json"] = resumen
			c.ServeJSON()

		}
		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "FP" {

			var n *PreliquidacionFpController
			resumen := n.Preliquidar(&v, reglasbase)

			c.Data["json"] = resumen
			c.ServeJSON()

		}

		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "DP" || v.Preliquidacion.Nomina.TipoNomina.Nombre == "DP-SALARIOS" {
			var n *PreliquidaciondpController
			resumen := n.Preliquidar(&v, reglasbase)
			//pr := CargarNovedadesPersona(v[0].PersonasPreLiquidacion[0].IdPersona,&v[0])
			//fmt.Println("prueba: ", pr)
			c.Data["json"] = resumen
			c.ServeJSON()
		}

		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "PE" {

			var n *PreliquidacionpeController
			resumen := n.Preliquidar(&v, reglasbase)

			c.Data["json"] = resumen
			c.ServeJSON()

		}
		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "CT" || v.Preliquidacion.Nomina.TipoNomina.Nombre == "CT-SALARIOS" {
			var n *PreliquidacionctController //aca se esta creando un objeto del controlador especico
			resumen := n.Preliquidar(&v, reglasbase)
			//pr := CargarNovedadesPersona(v[0].PersonasPreLiquidacion[0].IdPersona,&v[0])
			//fmt.Println("prueba: ", pr)
			c.Data["json"] = resumen
			c.ServeJSON()
		}

	} else {
		fmt.Println("error2: ", err)
	}

}
func CargarReglasBase(dominio string) (reglas string) {
	//carga de reglas desde el ruler
	var reglasbase string = ``
	var v []models.Predicado
	var datos_conceptos []models.Concepto
	if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto?limit=0", &datos_conceptos); err == nil {
		for _, datos := range datos_conceptos {
			reglasbase = reglasbase + `codigo_concepto(` + datos.NombreConcepto + `,` + strconv.Itoa(datos.Id) + `).` + "\n"
		}
	} else {

	}
	fmt.Println(dominio)
	if err := getJson("http://"+beego.AppConfig.String("Urlruler")+":"+beego.AppConfig.String("Portruler")+"/"+beego.AppConfig.String("Nsruler")+"/predicado?limit=0&query=Dominio.Nombre:"+dominio, &v); err == nil {

		reglasbase = reglasbase + FormatoReglas(v) //funcion general para dar formato a reglas cargadas desde el ruler
	} else {
		fmt.Println("err: ", err)
	}

	//-----------------------------
	return reglasbase
}

func FormatoReglas(v []models.Predicado) (reglas string) {
	var arregloReglas = make([]string, len(v))
	reglas = ""
	//var respuesta []models.FormatoPreliqu
	for i := 0; i < len(v); i++ {
		arregloReglas[i] = v[i].Nombre
	}

	for i := 0; i < len(arregloReglas); i++ {
		reglas = reglas + arregloReglas[i] + "\n"
	}
	return
}

func CargarNovedadesPersona(id_persona int, datos_preliqu *models.DatosPreliquidacion) (reglas string) {

	//consulta de la(s) novedades que pueda tener la persona para la pre-liquidacion
	var v []models.ConceptoPorPersona

	reglas = "" //inicializacion de la variable donde se inyectaran las novedades como reglas
	if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_por_persona/novedades_activas/"+strconv.Itoa(id_persona), "POST", &v, &datos_preliqu.Preliquidacion); err == nil {
		if v != nil {

			for i := 0; i < len(v); i++ {
				reglas = reglas + "concepto(" + strconv.Itoa(id_persona) + "," + v[i].Concepto.Naturaleza + ", " + v[i].Tipo + ", " + v[i].Concepto.NombreConcepto + ", " + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + ", " + datos_preliqu.Preliquidacion.Nomina.Periodo + "). " + "\n"
			}

		}

	}
	fmt.Println("novedad: ", reglas)
	//------------------------------------------------------------------------------
	return reglas

}
