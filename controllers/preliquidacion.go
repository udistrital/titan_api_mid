package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"
	"time"
	"github.com/astaxie/beego"
)

// PreliquidacionController operations for Preliquidacion
type PreliquidacionController struct {
	beego.Controller
}

var reglas_nov_dev string = ``
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

		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCS" {
			var n *PreliquidacionHcSController
			resumen := n.Preliquidar(&v, reglasbase)

			c.Data["json"] = resumen
			c.ServeJSON()

		}

		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "CT" || v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCH" {
			var n *PreliquidacioncthchController //aca se esta creando un objeto del controlador especico
			resumen := n.Preliquidar(&v, reglasbase)
			c.Data["json"] = resumen
			c.ServeJSON()
		}

		/*
		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "FP" {

			var n *PreliquidacionFpController
			resumen := n.Preliquidar(&v, reglasbase)
	 		c.Data["json"] = resumen
			c.ServeJSON()

		}

		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "DP"  {
			var n *PreliquidaciondpController
			resumen := n.Preliquidar(&v, reglasbase)
			c.Data["json"] = resumen
			c.ServeJSON()
		}

		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "PE" {
			var n *PreliquidacionpeController
			resumen := n.Preliquidar(&v, reglasbase)
			c.Data["json"] = resumen
			c.ServeJSON()
*/





	}else {
		fmt.Println("errorrrr2: ", err)
	}

}
func CargarReglasBase(dominio string) (reglas string) {
	//carga de reglas desde el ruler
	var reglasbase string = ``

	var v []models.Predicado
	var datos_conceptos []models.ConceptoNomina
	if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina?limit=-1", &datos_conceptos); err == nil {
		for _, datos := range datos_conceptos {
			reglasbase = reglasbase + "codigo_concepto("+datos.NombreConcepto + "," + strconv.Itoa(datos.Id) + "," + strconv.Itoa(datos.NaturalezaConcepto.Id)+")." + "\n"
		}
	} else {
		fmt.Println("error al cargar conceptos como reglas")
	}
	fmt.Println(dominio)
	if err := getJson("http://"+beego.AppConfig.String("Urlruler")+":"+beego.AppConfig.String("Portruler")+"/"+beego.AppConfig.String("Nsruler")+"/predicado?limit=-1&query=Dominio.Nombre:"+dominio, &v); err == nil {

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

func CargarNovedadesPersona(id_persona int, numero_contrato string, vigencia int, datos_preliqu *models.Preliquidacion) (reglas string) {

	//consulta de la(s) novedades que pueda tener la persona para la pre-liquidacion
	var v []models.ConceptoNominaPorPersona
	reglas_nov_dev = ""
	reglas = ""//inicializacion de la variable donde se inyectaran las novedades como reglas
	query := "Activo:true,Persona:"+strconv.Itoa(id_persona)
	if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina_por_persona?limit=-1&query="+query, &v); err == nil {
		if v != nil {
			for i := 0; i < len(v); i++ {
				esActiva := validarNovedades_segSocial(datos_preliqu.Mes,datos_preliqu.Ano , v[i].FechaDesde, v[i].FechaHasta)
				if esActiva == 1 {
					//todo debe estar dentro de este if para que se verifiquen fechas siempre
					if (v[i].Concepto.NaturalezaConcepto.Nombre == "seguridad_social"){

						year, month, day := v[i].FechaDesde.Date()
						year2, month2, day2 := v[i].FechaHasta.Date()

						reglas = reglas + "seg_social("+v[i].Concepto.NombreConcepto+","+strconv.Itoa(year)+","+strconv.Itoa(int(month))+","+strconv.Itoa(day + 1)+","+strconv.Itoa(year2)+","+strconv.Itoa(int(month2))+","+strconv.Itoa(day2 + 1)+")." + "\n"
						reglas = reglas + "concepto(" + strconv.Itoa(id_persona) + "," + v[i].Concepto.NaturalezaConcepto.Nombre + ", " + v[i].Concepto.TipoConcepto.Nombre + ", " + v[i].Concepto.NombreConcepto + ", " + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + ", " + strconv.Itoa(datos_preliqu.Ano) + "). " + "\n"

						}
         }

				 if(v[i].NumCuotas != 999 && v[i].NumCuotas != 0 ){
					 fmt.Println("cuotas")
					 numCuotas := cuotasPagas(numero_contrato, vigencia,v[i].Concepto.Id)
					 if(numCuotas == int(v[i].NumCuotas)){
						 fmt.Println("se pagó tota la novedad")
						 v[i].Activo = false
						 desactivarNovedad(v[i].Concepto.Id, v[i])
						 //inhabilitar cuotas
					 }else{
						 fmt.Println("no pagadas")
						 reglas = reglas + "concepto(" + strconv.Itoa(id_persona) + "," + v[i].Concepto.NaturalezaConcepto.Nombre + ", " + v[i].Concepto.TipoConcepto.Nombre + ", " + v[i].Concepto.NombreConcepto + ", " + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + ", " + strconv.Itoa(datos_preliqu.Ano) + "). " + "\n"
						 if (v[i].Concepto.NaturalezaConcepto.Nombre == "devengo"){
							 reglas = reglas + "devengo("+strconv.FormatFloat(v[i].ValorNovedad,'f', -1, 64)+","+v[i].Concepto.NombreConcepto+")." + "\n"
						 }
					 }
				 }

				 if(v[i].NumCuotas == 999 && v[i].Concepto.NaturalezaConcepto.Nombre != "seguridad_social"){
					  fmt.Println("no cuotas")
						reglas = reglas + "concepto(" + strconv.Itoa(id_persona) + "," + v[i].Concepto.NaturalezaConcepto.Nombre + ", " + v[i].Concepto.TipoConcepto.Nombre + ", " + v[i].Concepto.NombreConcepto + ", " + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + ", " + strconv.Itoa(datos_preliqu.Ano) + "). " + "\n"
					 if (v[i].Concepto.NaturalezaConcepto.Nombre == "devengo"){
							 reglas = reglas + "devengo("+strconv.FormatFloat(v[i].ValorNovedad,'f', -1, 64)+","+v[i].Concepto.NombreConcepto+")." + "\n"

						 }
				 }
			 }
		 }
	 }else{
		 fmt.Println(err)
	 }
	 fmt.Println(reglas)
	//------------------------------------------------------------------------------
	return reglas


}

func validarNovedades_segSocial(Mes, Ano int, FechaDesde, FechaHasta time.Time) (flag int) {

	if FechaDesde.Month() == time.Month(Mes) && FechaDesde.Year() == Ano {
		flag = 1

	} else if FechaHasta.Month() == time.Month(Mes) && FechaHasta.Year() == Ano {
		flag = 1

	} else if FechaHasta.Month() == time.Month(Mes) && FechaHasta.Year() == Ano {
		flag = 1

	} else{
		flag = 0

	}

	return flag

}

func cuotasPagas(numero_contrato string, vigencia,idConcepto int)(cuotas_pagas int){

	var vigencia_string string
	var idConcepto_string string

	vigencia_string = strconv.Itoa(vigencia)
	idConcepto_string = strconv.Itoa(idConcepto)

	var numero_cuotas_pagas int

	var detalle_preliquidacion []models.DetallePreliquidacion
	fmt.Println(vigencia_string)
	fmt.Println(numero_contrato)
	fmt.Println(idConcepto)
	if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query=Preliquidacion.EstadoPreliquidacion.Nombre:Cerrada,VigenciaContrato:"+vigencia_string+",NumeroContrato:"+numero_contrato+",Concepto.Id:"+idConcepto_string+"", &detalle_preliquidacion); err == nil {
		numero_cuotas_pagas = len(detalle_preliquidacion)
  }

	return numero_cuotas_pagas
}

func desactivarNovedad(idNovedad int, v models.ConceptoNominaPorPersona){
		var idNovedad_string string
		var idCPP interface{}
		idNovedad_string = strconv.Itoa(idNovedad)
		if err2 := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina_por_persona/"+idNovedad_string, "PUT", &idCPP, &v); err2 == nil {

	}
}
