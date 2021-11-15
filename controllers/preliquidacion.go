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

var reglasNovDev string = ""

// URLMapping ...
func (c *PreliquidacionController) URLMapping() {
	c.Mapping("Preliquidar", c.Preliquidar)
}

// Preliquidar ...
// @Title Preliquidar Contrato
// @Description Preliquida todos los meses del contrato que se le pase como parámetro
// @Param	body		body 	models.Contrato		true		"body for DatosPreliquidacion content"
// @Success 201 body is empty
// @Failure 403 body is empty
// @router / [post]
func (c *PreliquidacionController) Preliquidar() {
	var contrato models.Contrato
	var aux map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &contrato); err == nil {
		if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato", "POST", &aux, contrato); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato)
			if contrato.TipoNominaId == 411 {
				liquidarCPS(contrato, 0)
			} else if contrato.TipoNominaId == 409 {
				liquidarHCH(contrato)
			} else if contrato.TipoNominaId == 410 {

			}
		} else {
			fmt.Println("No se pudo guardar el contrato", err)
		}
	} else {
		fmt.Println("Error al obtener contratos", err)
	}
}

func CargarDatosRetefuente(cedula int) (reglas string, datosRetefuente models.ContratoPreliquidacion) {

	var tempPersonaNatural []models.PersonaNatural
	var alivios models.ContratoPreliquidacion
	reglas = ""
	query := "Id:" + strconv.Itoa(cedula)
	if err := request.GetJson(beego.AppConfig.String("UrlAdministrativaAmazon")+"/informacion_persona_natural?limit=-1&query="+query, &tempPersonaNatural); err == nil {

		if cedula == 80771924 || cedula == 14970179 || cedula == 51716336 || cedula == 79307484 || cedula == 73184070 || cedula == 71594083 || cedula == 11377590 || cedula == 79634185 {
			reglas = reglas + "reteiva(1)."
			alivios.ResponsableIva = true
		} else {
			reglas = reglas + "reteiva(0)."
			alivios.ResponsableIva = false
		}
		if tempPersonaNatural[0].Dependientes == true {
			reglas = reglas + "dependientes(1)."
			alivios.Dependientes = true
		} else {
			reglas = reglas + "dependientes(0)."
			alivios.Dependientes = false
		}

		if tempPersonaNatural[0].ValorUvtPrepagada > 0 {
			reglas = reglas + "medicina_prepagada(" + fmt.Sprintf("%f", tempPersonaNatural[0].ValorUvtPrepagada) + ")."
			alivios.MedicinaPrepagadaUvt = tempPersonaNatural[0].ValorUvtPrepagada
		} else {
			reglas = reglas + "medicina_prepagada(0)."
			alivios.MedicinaPrepagadaUvt = 0
		}

		if tempPersonaNatural[0].Pensionado == "true" {
			reglas = reglas + "pensionado(1)."
			alivios.Pensionado = true
		} else {
			reglas = reglas + "pensionado(0)."
			alivios.Pensionado = false
		}

		reglas = reglas + "intereses_vivienda(" + fmt.Sprintf("%f", tempPersonaNatural[0].InteresViviendaAfc) + ")."
		alivios.InteresesVivienda = tempPersonaNatural[0].InteresViviendaAfc

		alivios.PensionVoluntaria = 0 //Verificar campo cuando se tenga el endpoint de ágora
		reglas = reglas + "pension_voluntaria(0)."
		alivios.Afc = 0 //Verificar campo cuando se tenga el endpoint de ágora
		reglas = reglas + "afc(0)."
	} else {
		fmt.Println("error al consultar en Ágora", err)
		reglas = reglas + "dependientes(0)."
		reglas = reglas + "medicina_prepagada(0)."
		reglas = reglas + "pensionado(no)."
		reglas = reglas + "intereses_vivienda(0)."
		reglas = reglas + "reteiva(0)."
		reglas = reglas + "pension_voluntaria(0)."
		reglas = reglas + "afc(0)."
		alivios.PensionVoluntaria = 0 //Verificar campo cuando se tenga el endpoint de ágora
		alivios.Afc = 0               //Verificar campo cuando se tenga el endpoint de ágora
		alivios.ResponsableIva = false
		alivios.Dependientes = false
		alivios.MedicinaPrepagadaUvt = 0
		alivios.Pensionado = false
		alivios.InteresesVivienda = 0
	}
	return reglas, alivios
}
