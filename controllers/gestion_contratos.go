package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

type ContratosController struct {
	beego.Controller
}

func (c *ContratosController) URLMapping() {
	c.Mapping("ObtenerContratosDVE", c.ObtenerContratosDVE)
}

// Get ...
// @Title Obtener contratos DVE
// @Description Retorna todos los docnetes de vinvulación especial
// @Param	nomina		path 	string	true		"Nomina de la preliquidacion"
// @Param	mes		path 	string	true		"Mes de la preliquidación"
// @Param	ano		path 	string	true		"Año de la preliquidación"
// @Success 201 {object} []models.ContratoDVE
// @Failure 403 body is empty
// @router /docentesDVE/:nomina/:mes:/:ano [get]
func (c *ContratosController) ObtenerContratosDVE() {

	ano := c.Ctx.Input.Param(":ano")
	mes := c.Ctx.Input.Param(":mes")
	nomina := c.Ctx.Input.Param(":nomina")

	var aux map[string]interface{}
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var contratoDVE []models.ContratoDVE
	var auxContratoDVE models.ContratoDVE
	//Obtener todos los contratos pertenecientes a esa preliquidación
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query=PreliquidacionId.Mes:"+mes+",PreliquidacionId.Ano:"+ano+",PreliquidacionId.NominaId:"+nomina, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)

		//Añadir cada persona al arreglo
		for i := 0; i < len(contrato_preliquidacion); i++ {
			if i == 0 {
				auxContratoDVE.NombreCompleto = contrato_preliquidacion[0].ContratoId.NombreCompleto
				auxContratoDVE.Documento = contrato_preliquidacion[0].ContratoId.Documento
				auxContratoDVE.Cumplido = false
				auxContratoDVE.Preliquidado = false
				contratoDVE = append(contratoDVE, auxContratoDVE)
			} else {
				if !BuscarDocumento(contratoDVE, contrato_preliquidacion[i].ContratoId.Documento) {
					auxContratoDVE.NombreCompleto = contrato_preliquidacion[i].ContratoId.NombreCompleto
					auxContratoDVE.Documento = contrato_preliquidacion[i].ContratoId.Documento
					auxContratoDVE.Cumplido = true
					auxContratoDVE.Preliquidado = true
					contratoDVE = append(contratoDVE, auxContratoDVE)
				}
			}
		}

		//Verficar los cumplidos

		for i := 0; i < len(contratoDVE); i++ {
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query=PreliquidacionId.Mes:"+mes+",PreliquidacionId.Ano:"+ano+",ContratoId.Documento:"+contratoDVE[i].Documento, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)

				for j := 0; j < len(contrato_preliquidacion); j++ {
					if !contrato_preliquidacion[j].Cumplido {
						contratoDVE[i].Cumplido = false
						contratoDVE[i].Preliquidado = false
					}
				}

			} else {
				fmt.Println("Error al obtener los contratos del docente ", err)
				c.Data["mesaage"] = "Error al generar orden de pago " + err.Error()
				c.Abort("404")
			}
		}

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": contratoDVE}

	} else {
		fmt.Println("Error al obtener contrato preliquidacion: ", err)
		c.Data["mesaage"] = "Error al obtener contratos" + err.Error()
		c.Abort("404")
	}

	c.ServeJSON()

}

func BuscarDocumento(contratos []models.ContratoDVE, documento string) bool {
	for i := 0; i < len(contratos); i++ {
		if documento == contratos[i].Documento {
			return true
		}
	}
	return false
}
