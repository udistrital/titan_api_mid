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


func ContratosDVE(id_contrato string, vigencia int)(datos models.ObjetoContratoEstado,  err error){

	var temp map[string]interface{}
	var temp_docentes models.ObjetoContratoEstado
	var control_error error

	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/contrato_elaborado_estado/"+id_contrato+"/"+strconv.Itoa(vigencia), &temp); err == nil && temp != nil {
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

func ActaInicioDVE(id_contrato string, vigencia int)(datos models.ObjetoActaInicio,  err error){

	var temp map[string]interface{}
	var temp_docentes models.ObjetoActaInicio
	var control_error error

	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/acta_inicio_elaborado/"+id_contrato+"/"+strconv.Itoa(vigencia), &temp); err == nil && temp != nil {
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

func verificacion_pago(id_proveedor,ano, mes int, num_cont string, vig int, resultado models.Respuesta)(estado int){

	estado_pago := consultar_estado_pago(num_cont, vig, ano, mes);
	disponibilidad := calcular_disponibilidad(id_proveedor,vig,resultado)

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

func consultar_estado_pago(num_cont string, vigencia, ano, mes int)(disponibilidad int){

		//if err := getJson("http://"+beego.AppConfig.String("Urlkronos")+":"+beego.AppConfig.String("Portkronos")+"/"+beego.AppConfig.String("Nskronos")+"/registro_presupuestal/ValorActualRp/"+id_registro_pre, &saldo_rp); err == nil {
		var respuesta_servicio string
		var dispo int

		if err :=getJson("http://"+beego.AppConfig.String("Urlargo")+":"+beego.AppConfig.String("Portargo")+"/"+beego.AppConfig.String("Nsargo")+"/aprobacion_pago/pago_aprobado/DVE10/2017/8/2018", &respuesta_servicio); err == nil {

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

func InformacionContratista(NumeroContrato string, VigenciaContrato int)(Nom, cont, doc string,  err error){


	var temp map[string]interface{}
	var temp_docentes models.ObjetoInformacionContratista
	var nombre_contratista string
	var contrato string
	var documento string

	var control_error error

	if err := getJsonWSO2("http://"+beego.AppConfig.String("Urlwso2argo")+":"+beego.AppConfig.String("Portwso2argo")+"/"+beego.AppConfig.String("Nswso2argo")+"/informacion_contrato_elaborado_contratista/"+NumeroContrato+"/"+strconv.Itoa(VigenciaContrato), &temp); err == nil && temp != nil {
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
