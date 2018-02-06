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


func ContratosHonorarios(id_contrato string, vigencia int)(datos models.ObjetoContratoEstado,  err error){

	var temp map[string]interface{}
	var temp_docentes models.ObjetoContratoEstado
	var control_error error

	if err := getJsonWSO2("http://jbpm.udistritaloas.edu.co:8280/services/contrato_suscrito_DataService.HTTPEndpoint/contrato_elaborado_estado/"+id_contrato+"/"+strconv.Itoa(vigencia), &temp); err == nil && temp != nil {
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

func ActaInicioHonorarios(id_contrato string, vigencia int)(datos models.ObjetoActaInicio,  err error){

	var temp map[string]interface{}
	var temp_docentes models.ObjetoActaInicio
	var control_error error

	if err := getJsonWSO2("http://jbpm.udistritaloas.edu.co:8280/services/contrato_suscrito_DataService.HTTPEndpoint/acta_inicio_elaborado/"+id_contrato+"/"+strconv.Itoa(vigencia), &temp); err == nil && temp != nil {
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
