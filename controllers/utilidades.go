package controllers
import (
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
	"io"
	"os"
	"strings"
	"github.com/udistrital/titan_api_mid/models"
	"strconv"
	 "fmt"
	 "github.com/udistrital/titan_api_mid/golog"
	 	"github.com/udistrital/utils_oas/formatdata"
)



func sendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil{
			json.NewEncoder(b).Encode(datajson)
	}
	client := &http.Client{}
	req, err := http.NewRequest(trequest, url, b)
	r, err:= client.Do(req)
  //r, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}


func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getJsonWSO2(urlp string, target interface{}) error {
	b := new(bytes.Buffer)
	//proxyUrl, err := url.Parse("http://10.20.4.15:3128")
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlp, b)
	req.Header.Set("Accept", "application/json")
	r, err := client.Do(req)
	//r, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		beego.Error("error", err)
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

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



    year = int(y2 - y1)
    month = int(M2 - M1)
    day = int(d2 - d1) + 1


    // Normalize negative values
		/*if day < 0{
			day = 0
		}
		if month < 0 {
        month = 0
    }*/
    if day < 0 {
        // days in month:
        t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
        day += 32 - t.Day()
        month--
    }
    if month < 0 {
        month += 12
        year--
    }

    return
}

func tipoNomina(tipoNomina string)(tipo string){
	if tipoNomina == "151"  {
		tipo = "0"
	}

	if tipoNomina == "152" {
		tipo = "1"
	}

	if tipoNomina == "30" {
		tipo = "2"
	}
	 return tipo
}

func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
}


func ContratosDVE(id_contrato , vigencia string)(datos models.ObjetoContratoEstado,  err error){

	var temp map[string]interface{}
	var temp_docentes models.ObjetoContratoEstado
	var control_error error

	fmt.Println("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contrato_elaborado_estado/"+id_contrato+"/"+vigencia)
	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contrato_elaborado_estado/"+id_contrato+"/"+vigencia, &temp); err == nil && temp != nil {
		jsonDocentes, error_json := json.Marshal(temp)

		if error_json == nil {

			json.Unmarshal(jsonDocentes, &temp_docentes)
			fmt.Println(temp_docentes)
		} else {
			control_error = error_json
			fmt.Println("error al traer contratos docentes DVE")
		}
	} else {
		control_error = err
		fmt.Println("Error al unmarshal datos de nómina",err)


	}

		return temp_docentes, control_error;
}

func ActaInicioDVE(id_contrato, vigencia string)(datos models.ObjetoActaInicio,  err error){

	var temp map[string]interface{}
	var temp_docentes models.ObjetoActaInicio
	var control_error error

	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/acta_inicio_elaborado/"+id_contrato+"/"+vigencia, &temp); err == nil && temp != nil {
		jsonDocentes, error_json := json.Marshal(temp)

		if error_json == nil {

			json.Unmarshal(jsonDocentes, &temp_docentes)

		} else {
			control_error = error_json
			fmt.Println("error al traer contratos docentes DVE")
		}
	} else {
		control_error = err
		fmt.Println("Error al unmarshal datos de nómina",err)


	}

		return temp_docentes, control_error;
}

func verificacion_pago(id_proveedor,ano, mes int, num_cont, vig string,  resultado models.Respuesta)(estado int){

	estado_pago := consultar_estado_pago(num_cont, vig, ano, mes);
	//disponibilidad := calcular_disponibilidad(id_proveedor,vig,resultado)
	disponibilidad := 2;

	if(estado_pago == 2 && disponibilidad == 2){
		return 2
	}else{
		return 1
	}

}
func consultar_rp (id_proveedor, vigencia int) (saldo float64){
		var registro_presupuestal []models.RegistroPresupuestal
		var saldo_rp float64
		var id_proveedor_string = strconv.Itoa(id_proveedor)
		var vigencia_string = strconv.Itoa(vigencia)
		if err := getJson("http://"+beego.AppConfig.String("Urlkronos")+":"+beego.AppConfig.String("Portkronos")+"/"+beego.AppConfig.String("Nskronos")+"/registro_presupuestal?limit=-1&query=Beneficiario:"+id_proveedor_string+",Vigencia:"+vigencia_string, &registro_presupuestal); err == nil && registro_presupuestal != nil {
			var id_registro_pre = strconv.Itoa(registro_presupuestal[0].Id)
			if err := getJson("http://"+beego.AppConfig.String("Urlkronos")+":"+beego.AppConfig.String("Portkronos")+"/"+beego.AppConfig.String("Nskronos")+"/registro_presupuestal/ValorActualRp/"+id_registro_pre, &saldo_rp); err == nil {
				fmt.Println("saldo rp")
				fmt.Println(saldo_rp)
			}else{
				fmt.Println("error al consultar saldo de rp")
				fmt.Println(err)
				saldo_rp = 0;
			}



		}else{
			fmt.Println("error en consulta de rp")
			fmt.Println(err)
			saldo_rp = 0;
		}

		return saldo_rp
}


func total_a_pagar(respuesta models.Respuesta)(total float64){
	var total_dev float64
	for _, descuentos := range *respuesta.Conceptos {
		if(descuentos.NaturalezaConcepto == 1){
			valor, _ := strconv.ParseFloat(descuentos.Valor,64)
			total_dev = total_dev + valor
		}


}
 return total_dev
}

func calcular_disponibilidad(id_proveedor, vigencia int,respuesta models.Respuesta)(disp int){
	var valor_a_pagar float64
	var saldo_rp float64
	var disponibilidad int
	saldo_rp = consultar_rp(id_proveedor, vigencia)
	valor_a_pagar = total_a_pagar(respuesta)
	if(valor_a_pagar > saldo_rp){
		disponibilidad = 1;
		fmt.Println("no hay dinero")
	}else{
		disponibilidad = 2;
		fmt.Println("si hay dinero ")
	}

	return disponibilidad
}

func consultar_estado_pago(num_cont, vigencia string,  ano, mes int)(disponibilidad int){

		//if err := getJson("http://"+beego.AppConfig.String("Urlkronos")+":"+beego.AppConfig.String("Portkronos")+"/"+beego.AppConfig.String("Nskronos")+"/registro_presupuestal/ValorActualRp/"+id_registro_pre, &saldo_rp); err == nil {
		var respuesta_servicio string
		var dispo int
		if err :=getJson("http://"+beego.AppConfig.String("Urlargomid")+":"+beego.AppConfig.String("Portargomid")+"/"+beego.AppConfig.String("Nsargomid")+"/aprobacion_pago/pago_aprobado/"+num_cont+"/"+vigencia+"/"+strconv.Itoa(mes)+"/"+strconv.Itoa(ano)+"", &respuesta_servicio); err == nil {

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

		return dispo

}

func GetIdProveedor(Documento string)(IdProveedor int){


		var idProveedor int

		var respuesta_servicio []models.InformacionProveedor
		if control_error :=getJson("http://"+beego.AppConfig.String("Urlargoamazon")+"/"+beego.AppConfig.String("Nsargoamazon")+"/informacion_proveedor?query=NumDocumento:"+Documento, &respuesta_servicio); control_error == nil {
			idProveedor = respuesta_servicio[0].Id;
		}else{
			idProveedor= 0
			fmt.Println("error en consulta id de persona", control_error)

		}

		return idProveedor;



}

func InformacionPersonaProveedor(idPersona int)(Nom string, doc int,  err error){

		var nombre_persona string
		var documento int
		var respuesta_servicio []models.InformacionProveedor
		var control_error error
		if control_error :=getJson("http://"+beego.AppConfig.String("Urlargoamazon")+"/"+beego.AppConfig.String("Nsargoamazon")+"/informacion_proveedor?query=Id:"+strconv.Itoa(idPersona), &respuesta_servicio); control_error == nil {

			nombre_persona = respuesta_servicio[0].NomProveedor;
			documento,_ = strconv.Atoi(respuesta_servicio[0].NumDocumento);

		}else{
			nombre_persona = "No encontrado"
			nombre_persona = "0"
			fmt.Println("error en consulta de información de persona", control_error)

		}

		return nombre_persona, documento,control_error;



}

func InformacionPersona(tipoNomina string, NumeroContrato string, VigenciaContrato int)(Nom, cont, doc string,  err error){


	var temp map[string]interface{}
	var temp_docentes models.ObjetoInformacionContratista
	var nombre_contratista string
	var contrato string
	var documento string
	var endpoint string

	var control_error error


	if(tipoNomina == "CT" || tipoNomina == "HCS" || tipoNomina == "HCH"){

			if(tipoNomina == "CT"){
				endpoint = "informacion_contrato_contratista"
			}

			if(tipoNomina == "HCS" || tipoNomina == "HCH"){
				endpoint = "informacion_contrato_elaborado_contratista"
			}

			if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/"+endpoint+"/"+NumeroContrato+"/"+strconv.Itoa(VigenciaContrato), &temp); err == nil && temp != nil {

				jsonDocentes, error_json := json.Marshal(temp)

				if error_json == nil {

					json.Unmarshal(jsonDocentes, &temp_docentes)
					nombre_contratista = temp_docentes.InformacionContratista.NombreCompleto
					documento = temp_docentes.InformacionContratista.Documento.Numero
					contrato = temp_docentes.InformacionContratista.Contrato.Numero


				} else {
					control_error = error_json
					fmt.Println("error al traer contratos docentes DVE")
				}
			} else {
				control_error = err
				fmt.Println("Error al unmarshal datos de nómina",err)


			}
		}

		if(tipoNomina == "FP"){
			fmt.Println("asdafadada1")
			var datos_planta []models.Funcionario_x_Proveedor
			if err = getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/informacion_proveedor/get_informacion_personas_planta?numero_contrato="+NumeroContrato+"&vigencia="+strconv.Itoa(VigenciaContrato), &datos_planta); err == nil {
				fmt.Println("asdafadada", datos_planta)
				nombre_contratista = datos_planta[0].NombreProveedor ;
				contrato = datos_planta[0].NumeroContrato;
				documento = strconv.Itoa(datos_planta[0].NumDocumento);
				control_error = err;
			}else{
				fmt.Println(err)
			}


			}

		return nombre_contratista, contrato, documento,control_error;



}

func CrearResultado(detalles_a_totalizar []models.DetallePreliquidacion)(respuesta models.Respuesta){
	var res models.Respuesta

	conceptos := make([]models.ConceptosResumen, len(detalles_a_totalizar))

	for x,pos := range detalles_a_totalizar{
		conceptos[x].Valor = strconv.FormatFloat(pos.ValorCalculado, 'E', -1, 64)
		conceptos[x].NaturalezaConcepto = pos.Concepto.NaturalezaConcepto.Id;

	}

	res.Conceptos = &conceptos;
	return res

}

func CalcularTotalesPorPersona(conceptos  []models.ConceptosResumen)(total_dev, total_des, total_pag int){

	var total_devengos float64;
	var total_descuentos float64
	var total_a_pagar float64;

	for _, descuentos := range conceptos{
		valor, _ := strconv.ParseFloat(descuentos.Valor,64)
		if descuentos.NaturalezaConcepto == 1 {
			total_devengos = total_devengos + valor;
		}

		if descuentos.NaturalezaConcepto == 2 {
			total_descuentos = total_descuentos + valor;
		}
	}

	total_a_pagar = total_devengos - total_descuentos
	return int(total_devengos), int(total_descuentos), int(total_a_pagar)
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
									temp_valor_actual,_ := strconv.Atoi(info_total_persona[strconv.Itoa(dato_resumen.Id)])
									temp_valor_a_sumar,_ := strconv.Atoi(dato_conceptos.Valor)
									temp_valor := temp_valor_actual + temp_valor_a_sumar
									info_total_persona[strconv.Itoa(dato_resumen.Id)] =  strconv.Itoa(temp_valor)
									info_total_persona_temp["Total"] =  info_total_persona[strconv.Itoa(dato_resumen.Id)]
									info_total_personas[strconv.Itoa(dato_resumen.Id)] = info_total_persona_temp

					} else {

									info_total_persona_temp := make(map[string]string)
									temp_valor,_ := strconv.Atoi(dato_conceptos.Valor)
									info_total_persona[strconv.Itoa(dato_resumen.Id)] =  strconv.Itoa(temp_valor)
									info_total_persona_temp["Total"] =  info_total_persona[strconv.Itoa(dato_resumen.Id)]
									info_total_personas[strconv.Itoa(dato_resumen.Id)] = info_total_persona_temp

					}

					}
				}
			}

			fmt.Println("info_total_personas",info_total_personas, len(info_total_personas))
			var temp  []models.ConceptosResumen
			for key,_ := range info_total_personas {
				aux := models.TotalPersona{}
			 if err := formatdata.FillStruct(info_total_personas [key], &aux); err == nil{
				 temp = append(temp,golog.CalcularDescuentosTotalesHCS(key, aux.Total ,aux.Id,reglas,preliquidacion, strconv.Itoa(preliquidacion.Ano))...)
				fmt.Println("fondo soliwis",temp)
			 }else{
				 fmt.Println("error al guardar información agrupada",err)
			 }
			}

			return temp;
		}
