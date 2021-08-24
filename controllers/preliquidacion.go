package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
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

// PersonasPorPreliquidacion ...
// @Title create PersonasPorPreliquidacion
// @Description create PersonasPorPreliquidacion
// @Param	body		body 	models.Preliquidacion	true		"body for Preliquidacion content"
// @Success 201 {object}  []models.PersonasPreliquidacion
// @Failure 403 body is empty
// @router /personas_x_preliquidacion [post]
func (c *PreliquidacionController) PersonasPorPreliquidacion() {
	var v models.Preliquidacion
	var personasPreliquidacion []models.PersonasPreliquidacion
	var errorConsultaInformacionAgora error
	var aux map[string]interface{}

	//fmt.Println(c.Ctx.Input.RequestBody)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		query := "?query=PreliquidacionId:" + strconv.Itoa(v.Id) + ",EstadoDisponibilidadId:2&fields=PersonaId,NumeroContrato,VigenciaContrato,EstadoDisponibilidadId"
		if err := request.GetJson(beego.AppConfig.String("UrlCrudTitan")+"/detalle_preliquidacion/"+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &personasPreliquidacion)
			if len(personasPreliquidacion) != 0 {
				for x, dato := range personasPreliquidacion {
					personasPreliquidacion[x].NombreCompleto, personasPreliquidacion[x].NumDocumento, errorConsultaInformacionAgora = InformacionPersonaProveedor(dato.IdPersona)
				}

				if errorConsultaInformacionAgora == nil {
					c.Data["json"] = personasPreliquidacion
				} else {
					c.Data["json"] = errorConsultaInformacionAgora
					fmt.Println("error al consultar información en Agora")
				}
			} else {
				fmt.Printf("Sin personas por preliquidar")
			}

		} else {
			fmt.Println("error al traer resumen de preliquidacion")
			c.Data["json"] = err
		}

	} else {

		c.Data["json"] = err

		fmt.Println("error al leer datos de preliquidación a listar")
	}
	//fmt.Println(c.Data)
	c.ServeJSON()

}

// ResumenConceptos ...
// @Title create ResumenConceptos
// @Description create ResumenConceptos
// @Param	body		body 	 models.Preliquidacion	true		"body for Preliquidacion content"
// @Success 201 {object}  []models.ResumentCompleto
// @Failure 403 body is empty
// @router /resumen_conceptos [post]
func (c *PreliquidacionController) ResumenConceptos() {
	var v models.Preliquidacion
	infoDetalle := make(map[string]string)
	infoDetalles := make(map[string]interface{})
	var totalDevengos float64
	var totalDescuentos float64
	var detallePreliquidacion []models.DetallePreliquidacion
	var aux map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		query := "?limit=-1&query=PreliquidacionId:" + strconv.Itoa(v.Id) + ",EstadoDisponibilidadId:2"
		if err := request.GetJson(beego.AppConfig.String("UrlCrudTitan")+"/detalle_preliquidacion"+query, &aux); err == nil && detallePreliquidacion != nil {
			LimpiezaRespuestaRefactor(aux, &detallePreliquidacion)

			if len(detallePreliquidacion) != 0 {
				for _, dato := range detallePreliquidacion {

					_, ok := infoDetalle[strconv.Itoa(dato.ConceptoNominaId.Id)]
					if ok {

						infoDetalleTemp := make(map[string]string)
						tempValor, _ := strconv.Atoi(infoDetalle[strconv.Itoa(dato.ConceptoNominaId.Id)])
						tempValor = tempValor + int(dato.ValorCalculado)
						infoDetalle[strconv.Itoa(dato.ConceptoNominaId.Id)] = strconv.Itoa(tempValor)
						infoDetalleTemp["NombreConcepto"] = dato.ConceptoNominaId.AliasConcepto
						infoDetalleTemp["NaturalezaConcepto"] = dato.ConceptoNominaId.NaturalezaConceptoNominaId.Nombre
						infoDetalleTemp["NaturalezaConceptoId"] = strconv.Itoa(dato.ConceptoNominaId.EstadoConceptoNominaId.Id)
						infoDetalleTemp["Total"] = infoDetalle[strconv.Itoa(dato.ConceptoNominaId.Id)]
						infoDetalles[strconv.Itoa(dato.ConceptoNominaId.Id)] = infoDetalleTemp

					} else {

						infoDetalleTemp := make(map[string]string)
						tempValor := int(dato.ValorCalculado)
						infoDetalle[strconv.Itoa(dato.ConceptoNominaId.Id)] = strconv.Itoa(tempValor)
						infoDetalleTemp["NombreConcepto"] = dato.ConceptoNominaId.AliasConcepto
						infoDetalleTemp["NaturalezaConcepto"] = dato.ConceptoNominaId.EstadoConceptoNominaId.Nombre
						infoDetalleTemp["NaturalezaConceptoId"] = strconv.Itoa(dato.ConceptoNominaId.NaturalezaConceptoNominaId.Id)
						infoDetalleTemp["Total"] = infoDetalle[strconv.Itoa(dato.ConceptoNominaId.Id)]
						infoDetalles[strconv.Itoa(dato.ConceptoNominaId.Id)] = infoDetalleTemp

					}

					if dato.ConceptoNominaId.NaturalezaConceptoNominaId.Nombre == "devengo" {
						totalDevengos = totalDevengos + dato.ValorCalculado
					}

					if dato.ConceptoNominaId.NaturalezaConceptoNominaId.Nombre == "descuento" {
						totalDescuentos = totalDescuentos + dato.ValorCalculado
					}
				}
			}

			var resumenConceptos []models.Resumen
			for key := range infoDetalles {

				aux := models.Resumen{}
				if err := formatdata.FillStruct(infoDetalles[key], &aux); err == nil {
					resumenConceptos = append(resumenConceptos, aux)
				} else {
					fmt.Println("error al guardar información agrupada", err)
				}
			}

			var resumenTotal models.ResumentCompleto
			resumenTotal.TotalDevengos = int(totalDevengos)
			resumenTotal.TotalDescuentos = int(totalDescuentos)
			resumenTotal.ResumenTotalConceptos = resumenConceptos
			c.Data["json"] = resumenTotal

		} else {

			fmt.Println("error al traer resumen de preliquidacion")
			c.Data["json"] = err
		}

	} else {
		c.Data["json"] = err

		fmt.Println("error al leer datos de preliquidación a listar")
	}

	c.ServeJSON()

}

// GetIBCPorNovedad ...
// @Title create GetIBCPorNovedad
// @Description create GetIBCPorNovedad
// @Param	body		body 	models.IBCPorNovedad	true		"body for models.IBCPorNovedad content"
// @Success 201
// @Failure 403 body is empty
// @router /get_ibc_novedad [post]
func (c *PreliquidacionController) GetIBCPorNovedad() {
	var v models.IBCPorNovedad
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//carga de reglas desde el ruler
		reglasbase := cargarReglasBase(v.NombreNomina) //funcion general para dar formato a reglas cargadas desde el ruler
		reglasbase = reglasbase + cargarReglasSS()
		//-----------------------------

		if v.NombreNomina == "HCS" {
			var n *PreliquidacionHcSController
			resumen := n.GetIBCPorNovedad(v.Ano, v.Mes, v.NumDocumento, v.IdPersona, reglasbase, v.Novedad)
			c.Data["json"] = resumen
			c.ServeJSON()

		}

		if v.NombreNomina == "CT" {
			var n *PreliquidacionctController
			resumen := n.GetIBCPorNovedad(v.Ano, v.Mes, v.NumDocumento, v.IdPersona, reglasbase, v.Novedad)
			c.Data["json"] = resumen
			c.ServeJSON()

		}

		/*

					if v.Preliquidacion.Nomina.TipoNomina.Nombre == "CT" {
			  		var n *PreliquidacionctController //aca se esta creando un objeto del controlador especico
						resumen := n.Preliquidar(&v, reglasbase)
						c.Data["json"] = resumen
						c.ServeJSON()
					}

			    if v.Preliquidacion.Nomina.TipoNomina.Nombre == "HCH" {
			  		var n *PreliquidacionhchController //aca se esta creando un objeto del controlador especico
						resumen := n.Preliquidar(v, reglasbase)
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

	} else {
		fmt.Println("error al leer datos para calcular ibc", err)
	}

}

// Preliquidar ...
// @Title create Preliquidar
// @Description create Preliquidacion
// @Param	body		body 	models.DatosPreliquidacion	true		"body for DatosPreliquidacion content"
// @Success 201 {object} models.Preliquidacion
// @Failure 403 body is empty
// @router / [post]
func (c *PreliquidacionController) Preliquidar() {
	var v models.DatosPreliquidacion

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		//carga de reglas desde el ruler
		reglasbase := cargarReglasBase(v.Preliquidacion.NominaId.TipoNominaId.Nombre) //funcion general para dar formato a reglas cargadas desde el ruler
		reglasbase = reglasbase + cargarReglasSS()
		//-----------------------------

		if v.Preliquidacion.NominaId.TipoNominaId.Nombre == "HCS" {
			var n *PreliquidacionHcSController
			resumen := n.Preliquidar(v, reglasbase)

			c.Data["json"] = resumen
			c.ServeJSON()

		}

		if v.Preliquidacion.NominaId.TipoNominaId.Nombre == "CT" {
			fmt.Println("Datos preliquidacion: ", v.PersonasPreLiquidacion)
			var n *PreliquidacionctController //aca se esta creando un objeto del controlador especico
			resumen := n.Preliquidar(v, reglasbase)
			c.Data["json"] = resumen
			c.ServeJSON()
		}

		if v.Preliquidacion.NominaId.TipoNominaId.Nombre == "HCH" {
			var n *PreliquidacionhchController //aca se esta creando un objeto del controlador especico
			resumen := n.Preliquidar(v, reglasbase)
			c.Data["json"] = resumen
			c.ServeJSON()
		}

		if v.Preliquidacion.NominaId.TipoNominaId.Nombre == "FP" {

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

	} else {
		fmt.Println("error al leer datos de preliquidacion ", err)
	}

}

func cargarReglasBase(dominio string) (reglas string) {
	//carga de reglas desde el ruler
	var reglasbase string = ``

	var v []models.Predicado
	var aux map[string]interface{}
	var datosConceptos []models.ConceptoNomina

	if err := request.GetJson(beego.AppConfig.String("UrlCrudTitan")+"/concepto_nomina?limit=-1", &aux); err == nil {
		//fmt.Println("entre a cargar reglas")
		LimpiezaRespuestaRefactor(aux, &datosConceptos)
		//fmt.Println("concepto, nomina", datosConceptos)
		for _, datos := range datosConceptos {
			reglasbase = reglasbase + "codigo_concepto(" + datos.NombreConcepto + "," + strconv.Itoa(datos.Id) + "," + strconv.Itoa(datos.NaturalezaConceptoNominaId.Id) + ",'" + datos.AliasConcepto + "')." + "\n"
		}
	} else {
		fmt.Println("error al cargar conceptos como reglas", err)
	}

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

// FormatoReglas ...
// @Title FormatoReglas
// @Description Convierte las reglas en objeto Predicado a string
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

// CargarNovedadesPersona ...
// @Title CargarNovedadesPersona
// @Description Carga las novedades para una persona y las valida para esa preliquidacion
func CargarNovedadesPersona(id_persona int, numero_contrato, vigencia string, datos_preliqu models.Preliquidacion) (reglas string) {

	//consulta de la(s) novedades que pueda tener la persona para la pre-liquidacion
	var v []models.ConceptoNominaPorPersona
	var aux map[string]interface{}
	reglasNovDev = ""
	reglas = "" //inicializacion de la variable donde se inyectaran las novedades como reglas
	query := "Activo:true,NumeroContrato:" + numero_contrato + ",VigenciaContrato:" + vigencia
	if err := request.GetJson(beego.AppConfig.String("UrlCrudTitan")+"/concepto_nomina_por_persona?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &v)
		if v != nil {
			for i := 0; i < len(v); i++ {
				esActiva := validarNovedadesSegSocial(datos_preliqu.Mes, datos_preliqu.Ano, v[i].FechaDesde, v[i].FechaHasta)
				if esActiva == 1 {
					fmt.Println("Novedad activa")
					//todo debe estar dentro de este if para que se verifiquen fechas siempre
					if v[i].ConceptoNominaId.NaturalezaConceptoNominaId.Nombre == "seguridad_social" {

						year, month, day := v[i].FechaDesde.Date()
						year2, month2, day2 := v[i].FechaHasta.Date()

						reglas = reglas + "seg_social(" + v[i].ConceptoNominaId.NombreConcepto + "," + strconv.Itoa(year) + "," + strconv.Itoa(int(month)) + "," + strconv.Itoa(day+1) + "," + strconv.Itoa(year2) + "," + strconv.Itoa(int(month2)) + "," + strconv.Itoa(day2+1) + ")." + "\n"
						reglas = reglas + "concepto(" + strconv.Itoa(id_persona) + "," + v[i].ConceptoNominaId.NaturalezaConceptoNominaId.Nombre + ", " + v[i].ConceptoNominaId.TipoConceptoNominaId.Nombre + ", " + v[i].ConceptoNominaId.NombreConcepto + ", " + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + ", " + strconv.Itoa(datos_preliqu.Ano) + "). " + "\n"

					}
				}

				if v[i].NumCuotas != 999 && v[i].NumCuotas != 0 {

					numCuotas := cuotasPagas(numero_contrato, vigencia, v[i].ConceptoNominaId.Id)
					if numCuotas == v[i].NumCuotas {

						v[i].Activo = false
						desactivarNovedad(v[i].ConceptoNominaId.Id, v[i])
						//inhabilitar cuotas
					} else {

						reglas = reglas + "concepto(" + strconv.Itoa(id_persona) + "," + v[i].ConceptoNominaId.NaturalezaConceptoNominaId.Nombre + ", " + v[i].ConceptoNominaId.TipoConceptoNominaId.Nombre + ", " + v[i].ConceptoNominaId.NombreConcepto + ", " + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + ", " + strconv.Itoa(datos_preliqu.Ano) + "). " + "\n"
						if v[i].ConceptoNominaId.NaturalezaConceptoNominaId.Nombre == "devengo" {
							reglas = reglas + "devengo(" + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + "," + v[i].ConceptoNominaId.NombreConcepto + ")." + "\n"

						}
					}
				}

				if v[i].NumCuotas == 999 && v[i].ConceptoNominaId.NaturalezaConceptoNominaId.Nombre != "seguridad_social" {

					reglas = reglas + "concepto(" + strconv.Itoa(id_persona) + "," + v[i].ConceptoNominaId.NaturalezaConceptoNominaId.Nombre + ", " + v[i].ConceptoNominaId.TipoConceptoNominaId.Nombre + ", " + v[i].ConceptoNominaId.NombreConcepto + ", " + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + ", " + strconv.Itoa(datos_preliqu.Ano) + "). " + "\n"
					if v[i].ConceptoNominaId.NaturalezaConceptoNominaId.Nombre == "devengo" {
						reglas = reglas + "devengo(" + strconv.FormatFloat(v[i].ValorNovedad, 'f', -1, 64) + "," + v[i].ConceptoNominaId.NombreConcepto + ")." + "\n"

					}

				}
			}
		}
	} else {
		fmt.Println("Error al traer novedades", err)
	}

	//------------------------------------------------------------------------------
	return reglas

}

func validarNovedadesSegSocial(Mes, Ano int, FechaDesde, FechaHasta time.Time) (flag int) {

	/*
	  Se verifica si el año de la novedad es mayor (mas no igual) a la fecha final de la novedad
	  ya que de ser así, ya es valida
	  Si no, se verifica que las fechas desde y hasta cubren el mes de Liquidacion
	*/

	if FechaHasta.Year() > Ano {
		return 1
	} else if FechaDesde.Year() <= Ano && FechaHasta.Year() >= Ano && int(FechaDesde.Month()) <= Mes && int(FechaHasta.Month()) >= Mes {
		return 1
	} else {

		return 0
	}

}

func cuotasPagas(numero_contrato, vigencia string, idConcepto int) (cuotas_pagas int) {

	var idConceptoString string

	idConceptoString = strconv.Itoa(idConcepto)

	var numeroCuotasPagas int
	var aux map[string]interface{}

	var detallePreliquidacion []models.DetallePreliquidacion
	if err := request.GetJson(beego.AppConfig.String("UrlCrudTitan")+"/detalle_preliquidacion?limit=-1&query=PreliquidacionId.EstadoPreliquidacionId.Nombre:Cerrada,VigenciaContrato:"+vigencia+",NumeroContrato:"+numero_contrato+",ConceptoNominaId:"+idConceptoString, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &detallePreliquidacion)
		if len(detallePreliquidacion) != 0 {
			if detallePreliquidacion[0].VigenciaContrato != 0 {
				numeroCuotasPagas = len(detallePreliquidacion)
			} else {
				numeroCuotasPagas = 0
			}
		}
	}

	return numeroCuotasPagas
}

func desactivarNovedad(idNovedad int, v models.ConceptoNominaPorPersona) {
	var idNovedadString string
	var idCPP interface{}
	idNovedadString = strconv.Itoa(idNovedad)
	if err2 := request.SendJson(beego.AppConfig.String("UrlCrudTitan")+"/concepto_nomina_por_persona/"+idNovedadString, "PUT", &idCPP, &v); err2 == nil {

	}
}

// CargarDatosRetefuente ...
// @Title CargarDatosRetefuente
// @Description Carga lo referente al calculo de la retefuente según cédula de la persona
func CargarDatosRetefuente(cedula int) (reglas string) {

	var v []models.InformacionPersonaNatural
	reglas = ""
	query := "Id:" + strconv.Itoa(cedula)
	if err := request.GetJson(beego.AppConfig.String("UrlAdministrativaAmazon")+"/informacion_persona_natural?limit=-1&query="+query, &v); err == nil {
		fmt.Println("reglas SS:", query)
		if len(v) != 0 {

			if v[0].PersonasACargo == true {
				reglas = reglas + "dependiente(si)."
			} else {
				reglas = reglas + "dependiente(no)."
			}

			if v[0].DeclaranteRenta == true {
				reglas = reglas + "declarante(si)."
			} else {
				reglas = reglas + "declarante(no)."
			}

			if v[0].MedicinaPrepagada == true {
				reglas = reglas + "medicina_prepagada(si)."
			} else {
				reglas = reglas + "medicina_prepagada(no)."
			}

			if v[0].IdFondoPension == 119 {
				reglas = reglas + "pensionado(si)."
			} else {
				reglas = reglas + "pensionado(no)."
			}

			reglas = reglas + "intereses_vivienda(" + strconv.Itoa(int(v[0].InteresViviendaAfc)) + ")."

		} else {
			fmt.Println("No existen datos sobre esa persona")
			reglas = reglas + "dependiente(no)."
			reglas = reglas + "declarante(no)."
			reglas = reglas + "medicina_prepagada(no)."
			reglas = reglas + "pensionado(no)."
			reglas = reglas + "intereses_vivienda(0)."
		}
	} else {
		fmt.Println("error al consultar en Ágora", err)
		reglas = reglas + "dependiente(no)."
		reglas = reglas + "declarante(no)."
		reglas = reglas + "medicina_prepagada(no)."
		reglas = reglas + "pensionado(no)."
		reglas = reglas + "intereses_vivienda(0)."
	}

	//fmt.Println("reglas RETEARGO:", v)
	//fmt.Println("reglas RETEFUENTE:", reglas)
	return reglas
}
