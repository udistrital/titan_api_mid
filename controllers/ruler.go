package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func cargarReglasBase(dominio string) (reglas string) {
	//carga de reglas desde el ruler
	var reglasbase string = ``
	var v []models.Predicado
	var datosConceptos []models.ConceptoNomina
	var aux map[string]interface{}
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/concepto_nomina?limit=-1", &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &datosConceptos)
		for _, datos := range datosConceptos {
			reglasbase = reglasbase + "codigo_concepto(" + datos.NombreConcepto + "," + strconv.Itoa(datos.Id) + "," + strconv.Itoa(datos.NaturalezaConceptoNominaId) + ")." + "\n"
		}
	} else {
		fmt.Println("error al cargar conceptos como reglas", err)
	}
	fmt.Println(beego.AppConfig.String("UrlRuler") + "/predicado?limit=-1&query=Dominio.Nombre:" + dominio)
	if err := request.GetJson(beego.AppConfig.String("UrlRuler")+"/predicado?limit=-1&query=Dominio.Nombre:"+dominio, &v); err == nil {
		reglasbase = reglasbase + FormatoReglas(v) //funcion general para dar formato a reglas cargadas desde el ruler
	} else {
		fmt.Println("error al cargar reglas base: ", err)
	}

	//-----------------------------
	return reglasbase
}

func cargarReglasSS() (reglas string) {
	//carga de reglas de SS desde el ruler
	var reglasSS string = ``
	var v []models.Predicado

	if err := request.GetJson(beego.AppConfig.String("UrlRuler")+"/predicado?limit=-1&query=Dominio.Nombre:SeguridadSocial", &v); err == nil {
		reglasSS = FormatoReglas(v) //funcion general para dar formato a reglas cargadas desde el ruler
	} else {
		fmt.Println("error al cargar reglas base: ", err)
	}

	//-----------------------------
	return reglasSS
}
