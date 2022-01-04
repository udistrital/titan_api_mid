package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

func diff(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	oneDay := time.Hour * 5
	a = a.Add(oneDay)
	b = b.Add(oneDay)
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year = y2 - y1
	month = int(M2 - M1)
	day = d2 - d1

	if day < 0 {

		day = (30 - d1) + d2
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// GetIDProveedor ...
func GetIDProveedor(Documento string) (IDProveedor int) {

	var idProveedor int

	var respuesta_servicio []models.InformacionProveedor
	if controlError := request.GetJson("http://"+beego.AppConfig.String("Urlargoamazon")+":"+beego.AppConfig.String("Portargoamazon")+"/"+beego.AppConfig.String("Nsargoamazon")+"/informacion_proveedor?query=NumDocumento:"+Documento, &respuesta_servicio); controlError == nil {
		idProveedor = respuesta_servicio[0].Id
	} else {
		idProveedor = 0
		fmt.Println("error en consulta id de persona", controlError)

	}

	return idProveedor

}

// InformacionPersonaProveedor ...
func InformacionPersonaProveedor(idPersona int) (Nom string, doc int, err error) {

	var nombre_persona string
	var documento int
	var respuesta_servicio []models.InformacionProveedor
	var controlError error
	fmt.Println("URL ARGO", "http://"+beego.AppConfig.String("Urlargoamazon")+":"+beego.AppConfig.String("Portargoamazon")+"/"+beego.AppConfig.String("Nsargoamazon")+"/informacion_proveedor?query=Id:"+strconv.Itoa(idPersona))
	if controlError := request.GetJson("http://"+beego.AppConfig.String("Urlargoamazon")+":"+beego.AppConfig.String("Portargoamazon")+"/"+beego.AppConfig.String("Nsargoamazon")+"/informacion_proveedor?query=Id:"+strconv.Itoa(idPersona), &respuesta_servicio); controlError == nil {

		nombre_persona = respuesta_servicio[0].NomProveedor
		documento, _ = strconv.Atoi(respuesta_servicio[0].NumDocumento)
	} else {
		nombre_persona = "No encontrado"
		nombre_persona = "0"
		fmt.Println("error en consulta de información de persona", controlError)

	}
	return nombre_persona, documento, controlError
}

func InformacionPersona(tipoNomina string, NumeroContrato string, VigenciaContrato int) (Nom, cont, doc string, err error) {

	var temp map[string]interface{}
	var tempDocentes models.ObjetoInformacionContratista
	var nombre_contratista string
	var contrato string
	var documento string
	var endpoint string

	var controlError error

	if tipoNomina == "CT" || tipoNomina == "HCS" || tipoNomina == "HCH" {

		if tipoNomina == "CT" {
			endpoint = "informacion_contrato_contratista"
		}

		if tipoNomina == "HCS" || tipoNomina == "HCH" {
			endpoint = "informacion_contrato_elaborado_contratista"
		}

		if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/"+endpoint+"/"+NumeroContrato+"/"+strconv.Itoa(VigenciaContrato), &temp); err == nil && temp != nil {

			jsonDocentes, errorJSON := json.Marshal(temp)

			if errorJSON == nil {

				json.Unmarshal(jsonDocentes, &tempDocentes)
				nombre_contratista = tempDocentes.InformacionContratista.NombreCompleto
				documento = tempDocentes.InformacionContratista.Documento.Numero
				contrato = tempDocentes.InformacionContratista.Contrato.NumeroContrato

			} else {
				controlError = errorJSON
				fmt.Println("error al traer contratos docentes DVE")
			}
		} else {
			controlError = err
			fmt.Println("Error al unmarshal datos de nómina", err)

		}
	}

	if tipoNomina == "FP" {
		fmt.Println("asdafadada1")
		var datosPlanta []models.Funcionario_x_Proveedor
		if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/informacion_proveedor/get_informacion_personas_planta?numero_contrato="+NumeroContrato+"&vigencia="+strconv.Itoa(VigenciaContrato), &datosPlanta); err == nil {
			fmt.Println("asdafadada", datosPlanta)
			nombre_contratista = datosPlanta[0].NombreProveedor
			contrato = datosPlanta[0].NumeroContrato
			documento = strconv.Itoa(datosPlanta[0].NumDocumento)
			controlError = err
		} else {
			fmt.Println(err)
		}

	}

	return nombre_contratista, contrato, documento, controlError

}

func CalcularTotalesPorPersona(conceptos []models.ConceptosResumen) (total_dev, total_des, total_pag int) {

	var totalDevengos float64
	var totalDescuentos float64
	var total_a_pagar float64

	for _, descuentos := range conceptos {
		valor, _ := strconv.ParseFloat(descuentos.Valor, 64)
		if descuentos.NaturalezaConcepto == 1 {
			totalDevengos = totalDevengos + valor
		}

		if descuentos.NaturalezaConcepto == 2 {
			totalDescuentos = totalDescuentos + valor
		}
	}

	total_a_pagar = totalDevengos - totalDescuentos
	return int(totalDevengos), int(totalDescuentos), int(total_a_pagar)
}

func CalcularDescuentosTotales(reglas string, preliquidacion models.Preliquidacion, resumen []models.Respuesta) (concepto []models.ConceptosResumen) {

	info_total_persona := make(map[string]string)
	info_total_personas := make(map[string]interface{})

	for _, dato_resumen := range resumen {

		for _, dato_conceptos := range *dato_resumen.Conceptos {

			if dato_conceptos.NaturalezaConcepto == 1 {

				_, ok := info_total_persona[strconv.Itoa(dato_resumen.Id)]
				if ok {

					info_total_persona_temp := make(map[string]string)
					tempValor_actual, _ := strconv.Atoi(info_total_persona[strconv.Itoa(dato_resumen.Id)])
					tempValor_a_sumar, _ := strconv.Atoi(dato_conceptos.Valor)
					tempValor := tempValor_actual + tempValor_a_sumar
					info_total_persona[strconv.Itoa(dato_resumen.Id)] = strconv.Itoa(tempValor)
					info_total_persona_temp["Total"] = info_total_persona[strconv.Itoa(dato_resumen.Id)]
					info_total_personas[strconv.Itoa(dato_resumen.Id)] = info_total_persona_temp

				} else {

					info_total_persona_temp := make(map[string]string)
					tempValor, _ := strconv.Atoi(dato_conceptos.Valor)
					info_total_persona[strconv.Itoa(dato_resumen.Id)] = strconv.Itoa(tempValor)
					info_total_persona_temp["Total"] = info_total_persona[strconv.Itoa(dato_resumen.Id)]
					info_total_personas[strconv.Itoa(dato_resumen.Id)] = info_total_persona_temp

				}

			}
		}
	}

	var temp []models.ConceptosResumen
	for key := range info_total_personas {
		aux := models.TotalPersona{}
		if err := formatdata.FillStruct(info_total_personas[key], &aux); err == nil {

			if preliquidacion.NominaId == 415 {
				auxhch, _ := strconv.Atoi(aux.Total)

				auxhch2 := float64(auxhch) * 0.4

				aux.Total = fmt.Sprintf("%.0f", auxhch2)
			}
			temp = append(temp, golog.CalcularDescuentosTotalesHCS(key, aux.Total, aux.Id, reglas, preliquidacion, strconv.Itoa(preliquidacion.Ano))...)

			vrefFondoSol, _ := strconv.ParseFloat(temp[0].Valor, 64)
			vrefFondoSub, _ := strconv.ParseFloat(temp[1].Valor, 64)

			vTotal, _ := strconv.ParseFloat(aux.Total, 64)

			auxTemp := temp

			temp = nil
			for _, rPreliq := range resumen {
				//var detallePreliq []models.DetallePreliquidacion
				//if vrefFondoSol != 0 {

				auxConceptos := rPreliq.Conceptos

				//fmt.Println("AUXCONCEPTOS: ", *auxConceptos)
				var valorIBC float64

				for _, auxConcepto := range *auxConceptos {

					if auxConcepto.Nombre == "ibc_liquidado" {

						valorIBC, _ = strconv.ParseFloat(auxConcepto.Valor, 64)
					}

				}
				//Concepto 36 es el IBC
				//if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query=Preliquidacion:"+strconv.Itoa(preliquidacion.Id)+",VigenciaContrato:"+rPreliq.VigenciaContrato+",NumeroContrato:"+rPreliq.NumeroContrato+",Concepto.Id:36", &detallePreliq); err == nil {

				//valorIBC := detallePreliq[0].ValorCalculado

				valorFondoSol := (valorIBC * vrefFondoSol) / vTotal
				valorFondoSub := (valorIBC * vrefFondoSub) / vTotal

				conceptoFondoSol := auxTemp[0]
				conceptoFondoSub := auxTemp[1]

				conceptoFondoSol.Valor = fmt.Sprintf("%.0f", valorFondoSol)
				conceptoFondoSub.Valor = fmt.Sprintf("%.0f", valorFondoSub)

				temp = append(temp, conceptoFondoSol)
				temp = append(temp, conceptoFondoSub)

				//} else {
				//	fmt.Println("error al guardar información agrupada", err)
				//}

				//	}

			}
			fmt.Println("fondo solidaridad total", temp)
		} else {
			fmt.Println("error al guardar información agrupada", err)
		}
	}

	return temp
}

func LimpiezaRespuestaRefactor(respuesta map[string]interface{}, v interface{}) {
	b, err := json.Marshal(respuesta["Data"])
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &v)
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

func CalcularDias(FechaInicio time.Time, FechaFin time.Time) (diasLaborados float64, meses float64) {

	var a, m, d int
	var mesesContrato float64
	var diasContrato float64
	if FechaFin.IsZero() {
		FechaFin2 := time.Now()
		a, m, d = diff(FechaInicio, FechaFin2)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
		diasContrato = mesesContrato * 30
	} else {
		a, m, d = diff(FechaInicio, FechaFin)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
		diasContrato = mesesContrato * 30
	}
	return diasContrato, mesesContrato

}

func calcularSemanasContratoHCH(FechaInicio time.Time, FechaFin time.Time) (semanas float64) {
	var a, m, d int
	var mesesContrato float64
	if FechaFin.IsZero() {
		FechaFin2 := time.Now()
		a, m, d = diff(FechaInicio, FechaFin2)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)

	} else {
		a, m, d = diff(FechaInicio, FechaFin)
		mesesContrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
	}
	if mesesContrato/float64(int(mesesContrato)) != 1 {
		return (mesesContrato * 4) + 1
	} else {
		return (mesesContrato * 4)
	}
}

func registrarPreliquidacion(año, mes, estadoPreliquidacion, nomina int) (preliquidacion models.Preliquidacion) {
	var aux map[string]interface{}
	var nombre string

	if nomina == 412 {
		nombre = "Funcionarios Administrativos-Planta-"
	} else if nomina == 413 {
		nombre = "Docentes de Planta-"
	} else if nomina == 414 {
		nombre = "Contratistas-"
	} else if nomina == 415 {
		nombre = "Hora cátedra honorarios-"
	} else {
		nombre = "Hora cátedra salarios-"
	}
	preliquidacion.Descripcion = nombre + strconv.Itoa(año) + strconv.Itoa(mes)
	preliquidacion.Ano = año
	preliquidacion.Mes = mes
	preliquidacion.NominaId = nomina
	preliquidacion.EstadoPreliquidacionId = estadoPreliquidacion
	preliquidacion.Activo = true
	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/preliquidacion", "POST", &aux, preliquidacion); err == nil {
		LimpiezaRespuestaRefactor(aux, &preliquidacion)
		fmt.Println("Preliquidación creada con éxito")
	} else {
		fmt.Println("Error al guardar preliquidacion: ", err)
	}

	return preliquidacion
}

func registrarContratoPreliquidacion(preliquidacionId int, contratoId int, contratoPreliq models.ContratoPreliquidacion) (contratoPreliquidacion models.ContratoPreliquidacion) {
	var aux map[string]interface{}
	var respuesta models.ContratoPreliquidacion
	contratoPreliq.Id = 0
	contratoPreliquidacion = contratoPreliq
	contratoPreliquidacion.ContratoId = &models.Contrato{Id: contratoId}
	contratoPreliquidacion.PreliquidacionId = &models.Preliquidacion{Id: preliquidacionId}
	contratoPreliquidacion.Cumplido = false
	contratoPreliquidacion.Preliquidado = false
	contratoPreliquidacion.Activo = true

	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion", "POST", &aux, contratoPreliquidacion); err == nil {
		LimpiezaRespuestaRefactor(aux, &respuesta)
		fmt.Println("Contrato_preliquidacion guardado con exito")
	} else {
		fmt.Println("Error al guardar contrato_preliquidacion: ", err)
	}
	return respuesta
}

func registrarDetallePreliquidacion(detallePreliquidacion models.DetallePreliquidacion) {
	var response models.DetallePreliquidacion
	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion", "POST", &response, detallePreliquidacion); err == nil {
		fmt.Println("Detalle guardado con éxito")
	} else {
		fmt.Println("Error al guardar detalle", err)
	}
}

func registratContrato(contrato models.Contrato) (respuesta models.Contrato) {
	var response map[string]interface{}
	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato", "POST", &response, contrato); err == nil {
		fmt.Println("Contrato guardado con éxito")
		LimpiezaRespuestaRefactor(response, &respuesta)
	} else {
		fmt.Println("Error al guardar contrato", err)
	}
	return respuesta
}

func CalcularPeriodoLiquidacion(anoPreliquidacion, mesPreliquidacion int, fechaInicio, fechaFin time.Time) (periodoLiquidacion, periodoEspecifico string) {

	var FechaControl time.Time
	var periodo_liquidacion float64

	//En caso de que sea el mes de inicio y el mes final el mismo
	if fechaInicio.Month() == fechaFin.Month() && fechaInicio.Year() == fechaFin.Year() {
		periodo_liquidacion, _ = CalcularDias(fechaInicio, fechaFin)
		periodo_liquidacion = periodo_liquidacion + 1
		periodoEspecifico = "Del " + strconv.Itoa(fechaInicio.Day()) + " al " + strconv.Itoa(fechaFin.Day()) + " del mes " + strconv.Itoa(mesPreliquidacion) + " del año " + strconv.Itoa(anoPreliquidacion)
		//Para el mes de inicio
	} else if int(fechaInicio.Month()) == mesPreliquidacion && int(fechaInicio.Year()) == anoPreliquidacion {
		//En caso de que sea febrero, si no se coloca 28 tomará días de más
		if mesPreliquidacion == 2 {
			FechaControl = time.Date(anoPreliquidacion, time.Month(mesPreliquidacion), 28, 0, 0, 0, 0, time.UTC)
			periodo_liquidacion, _ = CalcularDias(fechaInicio, FechaControl)
			periodo_liquidacion = periodo_liquidacion + 3 //Dia inclusive y 2 días de desface de febrero
			periodoEspecifico = "Del " + strconv.Itoa(fechaInicio.Day()) + " al " + strconv.Itoa(FechaControl.Day()) + " del mes " + strconv.Itoa(mesPreliquidacion) + " del año " + strconv.Itoa(anoPreliquidacion)
		} else {
			FechaControl = time.Date(anoPreliquidacion, time.Month(mesPreliquidacion), 30, 0, 0, 0, 0, time.UTC)
			periodo_liquidacion, _ = CalcularDias(fechaInicio, FechaControl)
			periodo_liquidacion = periodo_liquidacion + 1 //Dia inclusive
			periodoEspecifico = "Del " + strconv.Itoa(fechaInicio.Day()) + " al " + strconv.Itoa(FechaControl.Day()) + " del mes " + strconv.Itoa(mesPreliquidacion) + " del año " + strconv.Itoa(anoPreliquidacion)
		}
		//Para el mes final
	} else if int(fechaFin.Month()) == mesPreliquidacion && int(fechaFin.Year()) == anoPreliquidacion {
		FechaControl = time.Date(anoPreliquidacion, time.Month(mesPreliquidacion), 1, 0, 0, 0, 0, time.UTC)
		periodo_liquidacion, _ = CalcularDias(FechaControl, fechaFin)
		periodo_liquidacion = periodo_liquidacion + 1 //Dia Inclusivo
		periodoEspecifico = "Del " + strconv.Itoa(FechaControl.Day()) + " al " + strconv.Itoa(fechaFin.Day()) + " del mes " + strconv.Itoa(mesPreliquidacion) + " del año " + strconv.Itoa(anoPreliquidacion)
	} else {
		periodo_liquidacion = 30
		periodoEspecifico = "Del 1 al 30 del mes " + strconv.Itoa(mesPreliquidacion) + " del año " + strconv.Itoa(anoPreliquidacion)
	}

	periodo := strconv.Itoa(int(periodo_liquidacion))

	return periodo, periodoEspecifico
}

func CalcularSemanas(diasLiquidados float64) (semanas int) {
	aux := diasLiquidados / 7

	if aux <= 1 {
		return 1
	} else if aux <= 2 {
		return 2
	} else if aux <= 3 {
		return 3
	} else {
		return 4
	}
}

func Remove(s []models.Contrato, i int) []models.Contrato {
	s = append(s[:i], s[i+1:]...)
	return s
}

func Roundf(x float64) float64 {
	t := math.Trunc(x)
	if math.Abs(x-t) >= 0.5 {
		return t + math.Copysign(1, x)
	}
	return t
}
