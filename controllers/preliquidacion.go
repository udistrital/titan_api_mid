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

// Resumen ...
// @Title Create
// @Description create Resumen
// @Param	body		body 	var v models.Preliquidacion	true		"body for Preliquidacion content"
// @Success 201 {object}  []models.InformePreliquidacion
// @Failure 403 body is empty
// @router /resumen/ [post]
func (c *PreliquidacionController) Resumen() {
	var v models.Preliquidacion
	var datos_preliquidacion []models.InformePreliquidacion
	var error_consulta_informacion_agora error
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/preliquidacion/resumen", "POST", &datos_preliquidacion, &v); err == nil {

			for x, dato := range datos_preliquidacion {

				datos_preliquidacion[x].NombreCompleto, datos_preliquidacion[x].NumeroContrato, datos_preliquidacion[x].Documento, error_consulta_informacion_agora= InformacionPersona(v.Nomina.TipoNomina.Nombre,dato.NumeroContrato, dato.Vigencia)

			}

			if(error_consulta_informacion_agora == nil){

				c.Data["json"] = datos_preliquidacion
			}else{
				c.Data["json"] = error_consulta_informacion_agora
				fmt.Println("error al consultar información en Agora")
			}


		} else {

			fmt.Println("error al traer resumen de preliquidacion")
			c.Data["json"] = err
		}

	}else{
		c.Data["json"] = err

		fmt.Println("error al leer datos de preliquidación a listar")
	}

	c.ServeJSON()

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


		if v.Preliquidacion.Nomina.TipoNomina.Nombre == "FP" {

			var n *PreliquidacionFpController
			resumen := n.Preliquidar(&v, reglasbase)
	 		c.Data["json"] = resumen
			c.ServeJSON()

		}
		/*
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
		fmt.Println("error al leer datos de preliquidacion ", err)
	}

}
func CargarReglasBase(dominio string) (reglas string) {
	//carga de reglas desde el ruler
	var reglasbase string = ``
	fmt.Println("dominio", dominio)
	var v []models.Predicado
	var datos_conceptos []models.ConceptoNomina

	if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina?limit=-1", &datos_conceptos); err == nil {
		for _, datos := range datos_conceptos {
			reglasbase = reglasbase + "codigo_concepto("+datos.NombreConcepto + "," + strconv.Itoa(datos.Id) + "," + strconv.Itoa(datos.NaturalezaConcepto.Id)+")." + "\n"
		}
	} else {
		fmt.Println("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina?limit=-1")
		fmt.Println("error al cargar conceptos como reglas")
	}

	if err := getJson("http://"+beego.AppConfig.String("Urlruler")+":"+beego.AppConfig.String("Portruler")+"/"+beego.AppConfig.String("Nsruler")+"/predicado?limit=-1&query=Dominio.Nombre:"+dominio, &v); err == nil {

		reglasbase = reglasbase + FormatoReglas(v) //funcion general para dar formato a reglas cargadas desde el ruler
	} else {
		fmt.Println("error al cargar reglas base: ", err)
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
	query := "Activo:true,NumeroContrato:"+numero_contrato+",VigenciaContrato:"+strconv.Itoa(vigencia)
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

					 numCuotas := cuotasPagas(numero_contrato, vigencia,v[i].Concepto.Id)
					 if(numCuotas == int(v[i].NumCuotas)){

						 v[i].Activo = false
						 desactivarNovedad(v[i].Concepto.Id, v[i])
						 //inhabilitar cuotas
					 }else{

						 reglas = reglas + "concepto(" + strconv.Itoa(id_persona) + "," + v[i].Concepto.NaturalezaConcepto.Nombre + ", " + v[i].Concepto.TipoConcepto.Nombre + ", " + v[i].Concepto.NombreConcepto + ", " + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + ", " + strconv.Itoa(datos_preliqu.Ano) + "). " + "\n"
						 if (v[i].Concepto.NaturalezaConcepto.Nombre == "devengo"){
							 reglas = reglas + "devengo("+strconv.FormatFloat(v[i].ValorNovedad,'f', -1, 64)+","+v[i].Concepto.NombreConcepto+")." + "\n"
						 }
					 }
				 }

				 if(v[i].NumCuotas == 999 && v[i].Concepto.NaturalezaConcepto.Nombre != "seguridad_social"){

						reglas = reglas + "concepto(" + strconv.Itoa(id_persona) + "," + v[i].Concepto.NaturalezaConcepto.Nombre + ", " + v[i].Concepto.TipoConcepto.Nombre + ", " + v[i].Concepto.NombreConcepto + ", " + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + ", " + strconv.Itoa(datos_preliqu.Ano) + "). " + "\n"
					 if (v[i].Concepto.NaturalezaConcepto.Nombre == "devengo"){
							 reglas = reglas + "devengo("+strconv.FormatFloat(v[i].ValorNovedad,'f', -1, 64)+","+v[i].Concepto.NombreConcepto+")." + "\n"

						 }
				 }
			 }
		 }
	 }else{
		 fmt.Println("Error al traer novedades",err)
	 }
	 fmt.Println("reglas de novedades",reglas)
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


func CargarDatosRetefuente(cedula int) (reglas string) {

	var v []models.InformacionPersonaNatural
	reglas = ""
	query := "Id:"+strconv.Itoa(cedula)
	if err := getJson("https://"+beego.AppConfig.String("Urlargoamazon")+"/"+beego.AppConfig.String("Nsargoamazon")+"/informacion_persona_natural?limit=-1&query="+query, &v); err == nil {
		if v != nil {

				if(v[0].PersonasACargo == true){
					reglas = reglas + "dependiente(si)."
				}else{
					reglas = reglas + "dependiente(no)."
				}

				if(v[0].DeclaranteRenta == true){
					reglas = reglas + "declarante(si)."
				}else{
					reglas = reglas + "declarante(no)."
				}

				if(v[0].MedicinaPrepagada == true){
					reglas = reglas + "medicina_prepagada(si)."
				}else{
					reglas = reglas + "medicina_prepagada(no)."
				}

				if(v[0].IdFondoPension == 119){
					reglas = reglas + "pensionado(si)."
				}else{
					reglas = reglas + "pensionado(no)."
				}

				reglas = reglas + "intereses_vivienda("+strconv.Itoa(int(v[0].InteresViviendaAfc))+")."
		}
	}

	return reglas
}
