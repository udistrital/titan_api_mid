package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// PreliquidacionController operations for Preliquidacion
type PreliquidacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *PreliquidacionController) URLMapping() {
	c.Mapping("Preliquidar", c.Preliquidar)
	c.Mapping("ObtenerResumenPreliquidacion", c.ObtenerResumenPreliquidacion)
}

// Preliquidar ...
// @Title Preliquidar Contrato
// @Description Preliquida todos los meses del contrato que se le pase como parámetro
// @Param	body		body 	models.Contrato		true		"body for DatosPreliquidacion content"
// @Success 201 object models.Contrato
// @Failure 404 not found resource
// @router / [post]
func (c *PreliquidacionController) Preliquidar() {
	var contrato models.Contrato
	//var aux map[string]interface{}
	//var mensaje string
	//var codigo string
	fmt.Println("Hola Mundo producción")
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &contrato); err == nil {
		fmt.Println(contrato)
		fmt.Println("Contrato Recibido: ")
		fmt.Println("Numero de contrato: ", contrato.NumeroContrato)
		fmt.Println("Documento: ", contrato.Documento)
		fmt.Println("Dependencia : ", contrato.DependenciaId)
		fmt.Println("Fecha Inicio: ", contrato.FechaInicio)
		fmt.Println("Fecha Fin: ", contrato.FechaFin)
		fmt.Println("Nombre completo: ", contrato.NombreCompleto)
		fmt.Println("Valor Contrato: ", contrato.ValorContrato)
		fmt.Println("Vacaciones: ", contrato.Vacaciones)
		fmt.Println("Unico: ", contrato.Unico)
		fmt.Println("Cdp: ", contrato.Cdp)
		fmt.Println("Rp: ", contrato.Rp)
		fmt.Println("Vigencia: ", contrato.Vigencia)
		fmt.Println("Tipo Nomina: ", contrato.TipoNominaId)
		fmt.Println("Activo: ", contrato.Activo)
		fmt.Println("Resolucion Id ", contrato.ResolucionId)
		fmt.Println("Resolucion ", contrato.Resolucion)
		fmt.Println("Numero Semanas: ", contrato.NumeroSemanas)

		mensaje, codigo, contratoReturn, err := Preliquidacion(contrato)

		if err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": codigo, "Message": mensaje, "Data": contratoReturn}
		} else {
			c.Data["message"] = mensaje + " " + err.Error()
			c.Abort(codigo)
		}
	} else {
		fmt.Println("Error al unmarshal del contrato: ", err)
		c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter"
		c.Abort("404")
	}
	c.ServeJSON()
}

func CargarDatosRetefuente(cedula int) (reglas string, datosRetefuente models.ContratoPreliquidacion, err error) {

	var aux map[string]interface{}
	var tempPersonaNatural models.PersonaNatural
	var alivios models.ContratoPreliquidacion
	reglas = ""

	query := strconv.Itoa(cedula)
	fmt.Println(beego.AppConfig.String("UrlArgoWso2") + "/informacion_persona_natural/" + query)
	if err := request.GetJsonWSO2(beego.AppConfig.String("UrlArgoWso2")+"/informacion_persona_natural/"+query, &aux); err == nil {

		jsonPersonaNatural, errorJSON := json.Marshal(aux["informacion_persona_natural"])
		if errorJSON == nil {
			json.Unmarshal(jsonPersonaNatural, &tempPersonaNatural)
			if tempPersonaNatural.ResponsableIva == "true" {
				reglas = reglas + "reteiva(1)."
				alivios.ResponsableIva = true
			} else {
				reglas = reglas + "reteiva(0)."
				alivios.ResponsableIva = false
			}

			if tempPersonaNatural.PersonasACargo == "true" {
				reglas = reglas + "dependientes(1)."
				alivios.Dependientes = true
			} else {
				reglas = reglas + "dependientes(0)."
				alivios.Dependientes = false
			}
			alivios.MedicinaPrepagadaUvt, _ = strconv.ParseFloat(tempPersonaNatural.ValorUvtPrepagada, 64)
			if alivios.MedicinaPrepagadaUvt >= 0 && tempPersonaNatural.ValorUvtPrepagada != "" {
				reglas = reglas + "medicina_prepagada(" + tempPersonaNatural.ValorUvtPrepagada + ")."
			} else {
				reglas = reglas + "medicina_prepagada(0)."
				alivios.MedicinaPrepagadaUvt = 0
			}

			if tempPersonaNatural.Pensionado == "true" {
				reglas = reglas + "pensionado(1)."
				alivios.Pensionado = true
			} else {
				reglas = reglas + "pensionado(0)."
				alivios.Pensionado = false
			}

			alivios.InteresesVivienda, _ = strconv.ParseFloat(tempPersonaNatural.InteresViviendaAfc, 64)
			if alivios.InteresesVivienda > 0 && tempPersonaNatural.InteresViviendaAfc != "" {
				reglas = reglas + "intereses_vivienda(" + tempPersonaNatural.InteresViviendaAfc + ")."
			} else {
				reglas = reglas + "intereses_vivienda(0)."
				alivios.InteresesVivienda = 0
			}

			alivios.PensionVoluntaria, _ = strconv.ParseFloat(tempPersonaNatural.ValorPensionVoluntaria, 64)
			if alivios.PensionVoluntaria > 0 && tempPersonaNatural.ValorPensionVoluntaria != "" {
				reglas = reglas + "pension_voluntaria(" + tempPersonaNatural.ValorPensionVoluntaria + " )."
			} else {
				alivios.PensionVoluntaria = 0
				reglas = reglas + "pension_voluntaria(0)."
			}
			alivios.Afc, _ = strconv.ParseFloat(tempPersonaNatural.ValorAfc, 64)
			if alivios.Afc > 0 && tempPersonaNatural.ValorAfc != "" {
				reglas = reglas + "afc(" + tempPersonaNatural.ValorAfc + ")."
			} else {
				alivios.Afc = 0
				reglas = reglas + "afc(0)."
			}
		} else {
			fmt.Println("Error al unmarshal del JSON: ", err)
			return "Error al unmarshal del JSON de Ágora: ", alivios, err
		}
	} else {
		fmt.Println("error al consultar en Ágora", err)
		return "error al consultar en Ágora: ", alivios, err
	}
	/*
		reglas = reglas + "dependientes(0)."
		reglas = reglas + "medicina_prepagada(0)."
		reglas = reglas + "pensionado(0)."
		reglas = reglas + "intereses_vivienda(0)."
		reglas = reglas + "reteiva(0)."
		reglas = reglas + "pension_voluntaria(0)."
		reglas = reglas + "afc(0)."
		alivios.PensionVoluntaria = 0
		alivios.Afc = 0
		alivios.ResponsableIva = false
		alivios.Dependientes = false
		alivios.MedicinaPrepagadaUvt = 0
		alivios.Pensionado = false
		alivios.InteresesVivienda = 0
	*/
	return reglas, alivios, nil
}

func CargarDatosRetefuenteOld(cedula int) (reglas string, datosRetefuente models.ContratoPreliquidacionOld, err error) {

	var aux map[string]interface{}
	var tempPersonaNatural models.PersonaNatural
	var alivios models.ContratoPreliquidacionOld
	reglas = ""

	query := strconv.Itoa(cedula)
	fmt.Println(beego.AppConfig.String("UrlArgoWso2") + "/informacion_persona_natural/" + query)
	if err := request.GetJsonWSO2(beego.AppConfig.String("UrlArgoWso2")+"/informacion_persona_natural/"+query, &aux); err == nil {

		jsonPersonaNatural, errorJSON := json.Marshal(aux["informacion_persona_natural"])
		if errorJSON == nil {
			json.Unmarshal(jsonPersonaNatural, &tempPersonaNatural)
			if tempPersonaNatural.ResponsableIva == "true" {
				reglas = reglas + "reteiva(1)."
				alivios.ResponsableIva = true
			} else {
				reglas = reglas + "reteiva(0)."
				alivios.ResponsableIva = false
			}

			if tempPersonaNatural.PersonasACargo == "true" {
				reglas = reglas + "dependientes(1)."
				alivios.Dependientes = true
			} else {
				reglas = reglas + "dependientes(0)."
				alivios.Dependientes = false
			}
			alivios.MedicinaPrepagadaUvt, _ = strconv.ParseFloat(tempPersonaNatural.ValorUvtPrepagada, 64)
			if alivios.MedicinaPrepagadaUvt >= 0 && tempPersonaNatural.ValorUvtPrepagada != "" {
				reglas = reglas + "medicina_prepagada(" + tempPersonaNatural.ValorUvtPrepagada + ")."
			} else {
				reglas = reglas + "medicina_prepagada(0)."
				alivios.MedicinaPrepagadaUvt = 0
			}

			if tempPersonaNatural.Pensionado == "true" {
				reglas = reglas + "pensionado(1)."
				alivios.Pensionado = true
			} else {
				reglas = reglas + "pensionado(0)."
				alivios.Pensionado = false
			}

			alivios.InteresesVivienda, _ = strconv.ParseFloat(tempPersonaNatural.InteresViviendaAfc, 64)
			if alivios.InteresesVivienda > 0 && tempPersonaNatural.InteresViviendaAfc != "" {
				reglas = reglas + "intereses_vivienda(" + tempPersonaNatural.InteresViviendaAfc + ")."
			} else {
				reglas = reglas + "intereses_vivienda(0)."
				alivios.InteresesVivienda = 0
			}

			alivios.PensionVoluntaria, _ = strconv.ParseFloat(tempPersonaNatural.ValorPensionVoluntaria, 64)
			if alivios.PensionVoluntaria > 0 && tempPersonaNatural.ValorPensionVoluntaria != "" {
				reglas = reglas + "pension_voluntaria(" + tempPersonaNatural.ValorPensionVoluntaria + " )."
			} else {
				alivios.PensionVoluntaria = 0
				reglas = reglas + "pension_voluntaria(0)."
			}
			alivios.Afc, _ = strconv.ParseFloat(tempPersonaNatural.ValorAfc, 64)
			if alivios.Afc > 0 && tempPersonaNatural.ValorAfc != "" {
				reglas = reglas + "afc(" + tempPersonaNatural.ValorAfc + ")."
			} else {
				alivios.Afc = 0
				reglas = reglas + "afc(0)."
			}
		} else {
			fmt.Println("Error al unmarshal del JSON: ", err)
			return "Error al unmarshal del JSON de Ágora: ", alivios, err
		}
	} else {
		fmt.Println("error al consultar en Ágora", err)
		return "error al consultar en Ágora: ", alivios, err
	}
	/*
		reglas = reglas + "dependientes(0)."
		reglas = reglas + "medicina_prepagada(0)."
		reglas = reglas + "pensionado(0)."
		reglas = reglas + "intereses_vivienda(0)."
		reglas = reglas + "reteiva(0)."
		reglas = reglas + "pension_voluntaria(0)."
		reglas = reglas + "afc(0)."
		alivios.PensionVoluntaria = 0
		alivios.Afc = 0
		alivios.ResponsableIva = false
		alivios.Dependientes = false
		alivios.MedicinaPrepagadaUvt = 0
		alivios.Pensionado = false
		alivios.InteresesVivienda = 0
	*/
	return reglas, alivios, nil
}

// Resumen de la preliquidacion ...
// @Title Resumen Preliquidacion
// @Description Retorna el total de la preliquidacion
// @Param	ano		path 	string	true		"Año de la preliquidación"
// @Param	mes		path 	string	true		"Mes de la preliquidación"
// @Param	nomina		path 	string	true		"Tipo de nomina"
// @Success 201 {object} models.DetalleMensual
// @Failure 404 No existe el registro
// @router /obtener_resumen_preliquidacion/:mes/:ano/:nomina [get]
func (c *PreliquidacionController) ObtenerResumenPreliquidacion() {
	var aux map[string]interface{}
	var detalles []models.DetallePreliquidacion
	var detalleMesual models.DetalleMensual

	ano := c.Ctx.Input.Param(":ano")
	mes := c.Ctx.Input.Param(":mes")
	nomina := c.Ctx.Input.Param(":nomina")

	//traer todos los detalles de la preliquidacion correspondientes a ese mes
	query := "ContratoPreliquidacionId.PreliquidacionId.Mes:" + mes + ",ContratoPreliquidacionId.PreliquidacionId.Ano:" + ano + ",ContratoPreliquidacionId.PreliquidacionId.NominaId:" + nomina
	fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/detalle_preliquidacion?limit=-1&query=" + query)
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &detalles)
		fmt.Println(len(detalles))
		for i := 0; i < len(detalles); i++ {
			if detalles[i].ConceptoNominaId.Id == 573 {
				detalleMesual.TotalDescuentos = detalleMesual.TotalDescuentos + detalles[i].ValorCalculado
			} else {
				if len(detalles) == 0 && detalles[i].ConceptoNominaId.Id != 574 {
					detalleMesual.Detalle = append(detalleMesual.Detalle, detalles[0])
					if detalles[i].ConceptoNominaId.NaturalezaConceptoNominaId == 423 {
						detalleMesual.TotalDevengado = detalleMesual.TotalDevengado + detalles[i].ValorCalculado
					}
				} else if detalles[i].ConceptoNominaId.Id != 574 {
					res, pos := encontrarConcepto(detalles[i].ConceptoNominaId.Id, detalleMesual.Detalle)
					if res {
						detalleMesual.Detalle[pos].ValorCalculado = detalleMesual.Detalle[pos].ValorCalculado + detalles[i].ValorCalculado
						if detalles[i].ConceptoNominaId.NaturalezaConceptoNominaId == 423 {
							detalleMesual.TotalDevengado = detalleMesual.TotalDevengado + detalles[i].ValorCalculado
						}
					} else {
						detalleMesual.Detalle = append(detalleMesual.Detalle, detalles[i])
						if detalles[i].ConceptoNominaId.NaturalezaConceptoNominaId == 423 {
							detalleMesual.TotalDevengado = detalleMesual.TotalDevengado + detalles[i].ValorCalculado
						}
					}
				}
			}
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": detalleMesual}
	} else {
		fmt.Println("Error al obtener los detalles")
		c.Data["mesaage"] = "Error, la preliquidación no existe"
		c.Abort("404")
	}
	c.ServeJSON()
}

func encontrarConcepto(id int, detalles []models.DetallePreliquidacion) (res bool, pos int) {
	for i := 0; i < len(detalles); i++ {
		if id == detalles[i].ConceptoNominaId.Id {
			return true, i
		}
	}
	return false, 0
}
