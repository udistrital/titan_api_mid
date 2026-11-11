package controllers

import (
	"context"
	"github.com/astaxie/beego/logs"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func cargarReglasBase(ctx context.Context, dominio string) (reglas string) {
	//carga de reglas desde el ruler
	var reglasbase string = ``
	var v []models.Predicado
	var datosConceptos []models.ConceptoNomina
	var aux map[string]interface{}
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlTitanCrud")+"/concepto_nomina?limit=-1", &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &datosConceptos)
		for _, datos := range datosConceptos {
			reglasbase = reglasbase + "codigo_concepto(" + datos.NombreConcepto + "," + strconv.Itoa(datos.Id) + "," + strconv.Itoa(datos.NaturalezaConceptoNominaId) + ")." + "\n"
		}
	} else {
		logs.Error("error al cargar conceptos como reglas", err)
	}
	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlRuler")+"/predicado?limit=-1&query=Dominio.Nombre:"+dominio, &v); err == nil {
		reglasbase = reglasbase + FormatoReglas(v) //funcion general para dar formato a reglas cargadas desde el ruler
	} else {
		logs.Error("error al cargar reglas base: ", err)
	}

	//-----------------------------
	return reglasbase
}

func cargarReglasSS(ctx context.Context) (reglas string) {
	//carga de reglas de SS desde el ruler
	var reglasSS string = ``
	var v []models.Predicado

	if _, err := request.GetWithContext(ctx, beego.AppConfig.String("UrlRuler")+"/predicado?limit=-1&query=Dominio.Nombre:SeguridadSocial", &v); err == nil {
		reglasSS = FormatoReglas(v) //funcion general para dar formato a reglas cargadas desde el ruler
	} else {
		logs.Error("error al cargar reglas base: ", err)
	}

	//-----------------------------
	return reglasSS
}
