package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GestionOpsController operations for GestionOps
type GestionOpsController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionOpsController) URLMapping() {
	c.Mapping("GenerarOrdenPago", c.GenerarOrdenPago)

}

// GenerarOrdenPago ...
// @Title create GenerarOrdenPago
// @Description Lanzar Job para crear órdenes de pago y actualizar estados de disponibilidad de detalles de preliquidación
// @Success 201
// @Failure 403 body is empty
// @router /generar_op [post]
func (c *GestionOpsController) GenerarOrdenPago() {
	fmt.Println("generar órden de pagooooooooooooooo")

	var v models.Preliquidacion
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//fmt.Println("holi", v.Id)
		r := ActualizarEstadoDisponibilidadDetalles(v.Id)
		if r.Type == "success" {
			//fmt.Println("resultado de actualización de disponibilidades", r)
			rr := ActualizarEstadoPreliquidacion(v)
			fmt.Println("Resultado de actualizar preliquidaciones", rr)
		}

		c.Data["json"] = r

		var anno string

		if v.Ano >= 10 {

			anno = strconv.Itoa(v.Ano)

		} else {

			anno = "0" + strconv.Itoa(v.Ano)
		}

		var temp interface{}
		//Servicio que dispara job

		if v.NominaId.TipoNominaId.Nombre == "CT" {
			if err := request.GetJsonWSO2(beego.AppConfig.String("UrlArgoPruebas")+"/liquidacion/"+anno+"/"+strconv.Itoa(v.Mes), &temp); err == nil {

			} else {

				fmt.Println("error job", err)
			}

		}

	} else {
		c.Data["json"] = err.Error()
		fmt.Println("error 2: ", err)
	}
	c.ServeJSON()

}

// ActualizarEstadoDisponibilidadDetalles ...
// @Title ActualizarEstadoDisponibilidadDetalles
// @Description Actualizar estado disponibilidad de los detalles cuyo cumplido fue aprobado (disponible a pagado)
func ActualizarEstadoDisponibilidadDetalles(id_pre int) (r models.Alert) {
	var aux map[string]interface{}
	var v models.Alert
	if err := request.GetJson(beego.AppConfig.String("UrlCrudTitan")+"/detalle_preliquidacion/update_estado_disponibilidad_detalle?idPreliquidacion="+strconv.Itoa(id_pre), &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &v)
	} else {
		fmt.Println("error: ", err)
	}
	return v
}

// ActualizarEstadoPreliquidacion ...
// @Title ActualizarEstadoPreliquidacion
// @Description Actualizar estado estado de la preliquidación: Si la preliquidación tiene personas aún pendientes (estado_disponibilidad = 1 ), se actualiza a estado 4 (Orden de Pago pendientes) o 1 (cerrada)
func ActualizarEstadoPreliquidacion(mp models.Preliquidacion) (e string) {
	var aux map[string]interface{}
	var v []models.DetallePreliquidacion
	var respuesta string
	if err := request.GetJson(beego.AppConfig.String("UrlCrudTitan")+"/detalle_preliquidacion?limit=-1&query=PreliquidacionId:"+strconv.Itoa(mp.Id)+",EstadoDisponibilidadId:1", &aux); err != nil {
		LimpiezaRespuestaRefactor(aux, &v)
		//fmt.Println("Hay personas pendientes", v)
		mp.EstadoPreliquidacionId.Id = 4
		if err := request.SendJson(beego.AppConfig.String("UrlCrudTitan")+"/preliquidacion/"+strconv.Itoa(mp.Id), "PUT", &respuesta, mp); err == nil {
			fmt.Println("Estado de preliquidación actualizada")
		} else {
			fmt.Println("Estado de preliquidación actualizada: ", err)
			respuesta = err.Error()
		}
	} else {
		//fmt.Println("No hay personas pendientes. Cerrar preliquidación", v)
		mp.EstadoPreliquidacionId.Id = 1
		if err := request.SendJson(beego.AppConfig.String("UrlCrudTitan")+"/preliquidacion/"+strconv.Itoa(mp.Id), "PUT", &respuesta, mp); err == nil {
			fmt.Println("Estado de preliquidación actualizada")
		} else {
			fmt.Println("Estado de preliquidación actualizada: ", err)
			respuesta = err.Error()
		}
	}
	return respuesta
}
