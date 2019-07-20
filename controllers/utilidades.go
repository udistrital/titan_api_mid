package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
	"github.com/udistrital/titan_api_mid/models"
	"strconv"
	 "fmt"
	 "github.com/udistrital/titan_api_mid/golog"
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

// ActaInicioDVE ...
// @Title ActaInicioDVE
// @Description Trae el acta de inicio segun contrato y vigencia
func ActaInicioDVE(id_contrato, vigencia string)(datos models.ObjetoActaInicio,  err error){

	var temp map[string]interface{}
	var tempDocentes models.ObjetoActaInicio
	var controlError error

	if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/acta_inicio_elaborado/"+id_contrato+"/"+vigencia, &temp); err == nil && temp != nil {
		jsonDocentes, errorJSON := json.Marshal(temp)

		if errorJSON == nil {

			json.Unmarshal(jsonDocentes, &tempDocentes)

		} else {
			controlError = errorJSON
			fmt.Println("error al traer contratos docentes DVE")
		}
	} else {
		controlError = err
		fmt.Println("Error al unmarshal datos de nómina",err)


	}

		return tempDocentes, controlError;
}

func verificacionPago(id_proveedor,ano, mes int, num_cont, vig string)(estado int){

	fmt.Println("verificación de cumplido 1: ")
	estadoPago := consultarEstadoPago(num_cont, vig, ano, mes);
	//disponibilidad := calcular_disponibilidad(id_proveedor,vig,resultado)
	disponibilidad := 2;

	if(estadoPago == 2 && disponibilidad == 2){
		return 2
	}else{
		return 1
	}

}

func consultarRp(id_proveedor, vigencia int) (saldo float64){
		var registroPresupuestal []models.RegistroPresupuestal
		var saldoRP float64
		var IDProveedorString = strconv.Itoa(id_proveedor)
		var vigenciaString = strconv.Itoa(vigencia)
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlkronos")+":"+beego.AppConfig.String("Portkronos")+"/"+beego.AppConfig.String("Nskronos")+"/registroPresupuestal?limit=-1&query=Beneficiario:"+IDProveedorString+",Vigencia:"+vigenciaString, &registroPresupuestal); err == nil && registroPresupuestal != nil {
			var id_registro_pre = strconv.Itoa(registroPresupuestal[0].Id)
			if err := request.GetJson("http://"+beego.AppConfig.String("Urlkronos")+":"+beego.AppConfig.String("Portkronos")+"/"+beego.AppConfig.String("Nskronos")+"/registroPresupuestal/ValorActualRp/"+id_registro_pre, &saldoRP); err == nil {
				fmt.Println("saldo rp")
				fmt.Println(saldoRP)
			}else{
				fmt.Println("error al consultar saldo de rp")
				fmt.Println(err)
				saldoRP = 0;
			}



		}else{
			fmt.Println("error en consulta de rp")
			fmt.Println(err)
			saldoRP = 0;
		}

		return saldoRP
}


func totalAPagar(respuesta models.Respuesta)(total float64){
	var total_dev float64
	for _, descuentos := range *respuesta.Conceptos {
		if(descuentos.NaturalezaConcepto == 1){
			valor, _ := strconv.ParseFloat(descuentos.Valor,64)
			total_dev = total_dev + valor
		}


}
 return total_dev
}

func calcularDisponibilidad(id_proveedor, vigencia int,respuesta models.Respuesta)(disp int){
	var valorAPagar float64
	var saldoRP float64
	var disponibilidad int
	saldoRP = consultarRp(id_proveedor, vigencia)
	valorAPagar = totalAPagar(respuesta)
	if(valorAPagar > saldoRP){
		disponibilidad = 1;
		fmt.Println("no hay dinero")
	}else{
		disponibilidad = 2;
		fmt.Println("si hay dinero ")
	}

	return disponibilidad
}

func consultarEstadoPago(num_cont, vigencia string,  ano, mes int)(disponibilidad int){


		var respuesta_servicio string
		var dispo int
		fmt.Println("pago:","http://"+beego.AppConfig.String("Urlargomid")+":"+beego.AppConfig.String("Portargomid")+"/"+beego.AppConfig.String("Nsargomid")+"/aprobacion_pago/pago_aprobado/"+num_cont+"/"+vigencia+"/"+strconv.Itoa(mes)+"/"+strconv.Itoa(ano)+"")
		if err :=request.GetJson("http://"+beego.AppConfig.String("Urlargomid")+":"+beego.AppConfig.String("Portargomid")+"/"+beego.AppConfig.String("Nsargomid")+"/aprobacion_pago/pago_aprobado/"+num_cont+"/"+vigencia+"/"+strconv.Itoa(mes)+"/"+strconv.Itoa(ano)+"", &respuesta_servicio); err == nil {

			if(respuesta_servicio == "True"){
				dispo = 2;
			}else{
				dispo = 1;
			}

			fmt.Println("consulta exitosa de aprobación de pago")
		}else{
			fmt.Println("error en consulta de aprobación de pago")
			dispo = 1;
		}

		fmt.Println("verificación de cumplido:",dispo)
		return dispo

}

// GetIDProveedor ...
func GetIDProveedor(Documento string)(IDProveedor int){


		var idProveedor int

		var respuesta_servicio []models.InformacionProveedor
		if controlError :=request.GetJson("http://"+beego.AppConfig.String("Urlargoamazon")+":"+beego.AppConfig.String("Portargoamazon")+"/"+beego.AppConfig.String("Nsargoamazon")+"/informacion_proveedor?query=NumDocumento:"+Documento, &respuesta_servicio); controlError == nil {
			idProveedor = respuesta_servicio[0].Id;
		}else{
			idProveedor= 0
			fmt.Println("error en consulta id de persona", controlError)

		}

		return idProveedor;



}

// InformacionPersonaProveedor ...
func InformacionPersonaProveedor(idPersona int)(Nom string, doc int,  err error){

		var nombre_persona string
		var documento int
		var respuesta_servicio []models.InformacionProveedor
		var controlError error
		fmt.Println("URL ARGO","http://"+beego.AppConfig.String("Urlargoamazon")+":"+beego.AppConfig.String("Portargoamazon")+"/"+beego.AppConfig.String("Nsargoamazon")+"/informacion_proveedor?query=Id:"+strconv.Itoa(idPersona))
		if controlError :=request.GetJson("http://"+beego.AppConfig.String("Urlargoamazon")+":"+beego.AppConfig.String("Portargoamazon")+"/"+beego.AppConfig.String("Nsargoamazon")+"/informacion_proveedor?query=Id:"+strconv.Itoa(idPersona), &respuesta_servicio); controlError == nil {

			nombre_persona = respuesta_servicio[0].NomProveedor;
			documento,_ = strconv.Atoi(respuesta_servicio[0].NumDocumento);

		}else{
			nombre_persona = "No encontrado"
			nombre_persona = "0"
			fmt.Println("error en consulta de información de persona", controlError)

		}

		return nombre_persona, documento,controlError;



}

func InformacionPersona(tipoNomina string, NumeroContrato string, VigenciaContrato int)(Nom, cont, doc string,  err error){


	var temp map[string]interface{}
	var tempDocentes models.ObjetoInformacionContratista
	var nombre_contratista string
	var contrato string
	var documento string
	var endpoint string

	var controlError error


	if(tipoNomina == "CT" || tipoNomina == "HCS" || tipoNomina == "HCH"){

			if(tipoNomina == "CT"){
				endpoint = "informacion_contrato_contratista"
			}

			if(tipoNomina == "HCS" || tipoNomina == "HCH"){
				endpoint = "informacion_contrato_elaborado_contratista"
			}

			if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/"+endpoint+"/"+NumeroContrato+"/"+strconv.Itoa(VigenciaContrato), &temp); err == nil && temp != nil {

				jsonDocentes, errorJSON := json.Marshal(temp)

				if errorJSON == nil {

					json.Unmarshal(jsonDocentes, &tempDocentes)
					nombre_contratista = tempDocentes.InformacionContratista.NombreCompleto
					documento = tempDocentes.InformacionContratista.Documento.Numero
					contrato = tempDocentes.InformacionContratista.Contrato.Numero


				} else {
					controlError = errorJSON
					fmt.Println("error al traer contratos docentes DVE")
				}
			} else {
				controlError = err
				fmt.Println("Error al unmarshal datos de nómina",err)


			}
		}

		if(tipoNomina == "FP"){
			fmt.Println("asdafadada1")
			var datosPlanta []models.Funcionario_x_Proveedor
			if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/informacion_proveedor/get_informacion_personas_planta?numero_contrato="+NumeroContrato+"&vigencia="+strconv.Itoa(VigenciaContrato), &datosPlanta); err == nil {
				fmt.Println("asdafadada", datosPlanta)
				nombre_contratista = datosPlanta[0].NombreProveedor ;
				contrato = datosPlanta[0].NumeroContrato;
				documento = strconv.Itoa(datosPlanta[0].NumDocumento);
				controlError = err;
			}else{
				fmt.Println(err)
			}


			}

		return nombre_contratista, contrato, documento,controlError;



}


func CalcularTotalesPorPersona(conceptos  []models.ConceptosResumen)(total_dev, total_des, total_pag int){

	var totalDevengos float64;
	var totalDescuentos float64
	var total_a_pagar float64;

	for _, descuentos := range conceptos{
		valor, _ := strconv.ParseFloat(descuentos.Valor,64)
		if descuentos.NaturalezaConcepto == 1 {
			totalDevengos = totalDevengos + valor;
		}

		if descuentos.NaturalezaConcepto == 2 {
			totalDescuentos = totalDescuentos + valor;
		}
	}

	total_a_pagar = totalDevengos - totalDescuentos
	return int(totalDevengos), int(totalDescuentos), int(total_a_pagar)
}


func CalcularDescuentosTotales(reglas string, preliquidacion models.Preliquidacion, resumen  []models.Respuesta)(concepto []models.ConceptosResumen){

	info_total_persona := make(map[string]string)
	info_total_personas := make(map[string]interface{})

	for _,dato_resumen := range resumen {


	for _,dato_conceptos := range *dato_resumen.Conceptos {

			if(dato_conceptos.NaturalezaConcepto == 1) {

					_, ok := info_total_persona[strconv.Itoa(dato_resumen.Id)]
					if ok {

									info_total_persona_temp := make(map[string]string)
									tempValor_actual,_ := strconv.Atoi(info_total_persona[strconv.Itoa(dato_resumen.Id)])
									tempValor_a_sumar,_ := strconv.Atoi(dato_conceptos.Valor)
									tempValor := tempValor_actual + tempValor_a_sumar
									info_total_persona[strconv.Itoa(dato_resumen.Id)] =  strconv.Itoa(tempValor)
									info_total_persona_temp["Total"] =  info_total_persona[strconv.Itoa(dato_resumen.Id)]
									info_total_personas[strconv.Itoa(dato_resumen.Id)] = info_total_persona_temp

					} else {

									info_total_persona_temp := make(map[string]string)
									tempValor,_ := strconv.Atoi(dato_conceptos.Valor)
									info_total_persona[strconv.Itoa(dato_resumen.Id)] =  strconv.Itoa(tempValor)
									info_total_persona_temp["Total"] =  info_total_persona[strconv.Itoa(dato_resumen.Id)]
									info_total_personas[strconv.Itoa(dato_resumen.Id)] = info_total_persona_temp

					}

					}
				}
			}


			var temp  []models.ConceptosResumen
			for key,_ := range info_total_personas {
				aux := models.TotalPersona{}
			 if err := formatdata.FillStruct(info_total_personas [key], &aux); err == nil{
				 temp = append(temp,golog.CalcularDescuentosTotalesHCS(key, aux.Total ,aux.Id,reglas,preliquidacion, strconv.Itoa(preliquidacion.Ano))...)
				fmt.Println("fondo solidaridad total",temp)
			 }else{
				 fmt.Println("error al guardar información agrupada",err)
			 }
			}

			return temp;
		}


		// ContratosContratistas ...
		// @Title ContratosContratistas
		// @Description Trae de Argo la informacion del contrato por su número y su vigencia
		func ContratosContratistas(id_contrato string, vigencia int)(datos models.ObjetoContratoEstado,  err error){

			var temp map[string]interface{}
			var tempDocentes models.ObjetoContratoEstado
			var controlError error
			if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contrato_estado/"+id_contrato+"/"+strconv.Itoa(vigencia), &temp); err == nil && temp != nil {
				jsonDocentes, errorJSON := json.Marshal(temp)

				if errorJSON == nil {

					json.Unmarshal(jsonDocentes, &tempDocentes)

				} else {
					controlError = errorJSON
					fmt.Println("error al traer contratos docentes DVE")
				}
			} else {
				controlError = err
				fmt.Println("Error al unmarshal datos de nómina",err)


			}

				return tempDocentes, controlError;
		}

		// ActaInicioContratistas ...
		// @Title ActaInicioContratistas
		// @Description Trae el acta de inicio por contrato y vigencia
		func ActaInicioContratistas(id_contrato string, vigencia int)(datos models.ObjetoActaInicio,  err error){

			var temp map[string]interface{}
			var tempDocentes models.ObjetoActaInicio
			var controlError error

			if err := request.GetJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/acta_inicio/"+id_contrato+"/"+strconv.Itoa(vigencia), &temp); err == nil && temp != nil {
				jsonDocentes, errorJSON := json.Marshal(temp)

				if errorJSON == nil {

					json.Unmarshal(jsonDocentes, &tempDocentes)

				} else {
					controlError = errorJSON
					fmt.Println("error al traer contratos docentes DVE")
				}
			} else {
				controlError = err
				fmt.Println("Error al unmarshal datos de nómina",err)


			}

				return tempDocentes, controlError;
		}
