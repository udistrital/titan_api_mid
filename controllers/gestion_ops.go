package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/astaxie/beego"
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
		fmt.Println("holi", v.Id)
		r := ActualizarEstadoDisponibilidadDetalles(v.Id)
		if r.Type == "success" {
			fmt.Println("resultado de actualización de disponibilidades", r)
			rr := ActualizarEstadoPreliquidacion(v)
			fmt.Println("Resultado de actualizar preliquidaciones",rr)
		}

		c.Data["json"] = r

		//Servicio que dispara job
	} else {
		c.Data["json"] = err.Error()
		fmt.Println("error 2: ", err)
	}
	c.ServeJSON()

}


// ActualizarEstadoDisponibilidadDetalles ...
// @Title ActualizarEstadoDisponibilidadDetalles
// @Description Actualizar estado disponibilidad de los detalles cuyo cumplido fue aprobado (disponible a pagado)
func  ActualizarEstadoDisponibilidadDetalles(id_pre int)(r models.Alert){

	var v models.Alert
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion/update_estado_disponibilidad_detalle?idPreliquidacion="+strconv.Itoa(id_pre), &v); err != nil || v.Type != "success"  {
		fmt.Println("Error:", v.Body)

	}
	return v
}

// ActualizarEstadoPreliquidacion ...
// @Title ActualizarEstadoPreliquidacion
// @Description Actualizar estado estado de la preliquidación: Si la preliquidación tiene personas aún pendientes (estado_disponibilidad = 1 ), se actualiza a estado 4 (Orden de Pago pendientes) o 1 (cerrada)
func  ActualizarEstadoPreliquidacion(mp models.Preliquidacion)(e string){

	var v []models.DetallePreliquidacion
	var respuesta string
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query=Preliquidacion:"+strconv.Itoa(mp.Id)+",EstadoDisponibilidad:1", &v); err != nil || v != nil  {
		fmt.Println("Hay personas pendientes",v)
		mp.EstadoPreliquidacion.Id = 4
		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/preliquidacion/"+strconv.Itoa(mp.Id), "PUT", &respuesta, mp); err == nil  {
			fmt.Println("Estado de preliquidación actualizada")
		} else {
			fmt.Println("Estado de preliquidación actualizada: ", err)
		  respuesta = err.Error()
		}
	}else{
		fmt.Println("No hay personas pendientes. Cerrar preliquidación",v)
		mp.EstadoPreliquidacion.Id = 1
		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/preliquidacion/"+strconv.Itoa(mp.Id), "PUT", &respuesta, mp); err == nil  {
			fmt.Println("Estado de preliquidación actualizada")
		} else {
			fmt.Println("Estado de preliquidación actualizada: ", err)
			respuesta = err.Error()
		}
	}
	return respuesta
}
