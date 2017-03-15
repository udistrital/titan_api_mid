package controllers

import (
	"fmt"
	"strconv"
	"time"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
   "net/http"
	 "net/url"
	 "io/ioutil"
	"github.com/astaxie/beego"
	"encoding/json"
)

// operations for Preliquidaciondp
type PreliquidaciondpController struct {
	beego.Controller
}

func (c *PreliquidaciondpController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	var resumen_preliqu []models.Respuesta
	var idDetaPre interface{}
	var tipoNom string;
	var puntos float64

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {
		var informacion_cargo []models.DocenteCargo
		var cedula int
		var reglasinyectadas string
		var reglas string
		filtrodatos := models.DocenteCargo{Id: datos.PersonasPreLiquidacion[i].IdPersona, Asignacion_basica: 0}
		tipoNom = tipoNomina(datos.Preliquidacion.Tipo)

		//fmt.Println("reglas: ", reglasbase)
		//consulta que envie ID de proveedor en datos y retorne el salario, para que sea enviado a CargarReglas
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/docente_cargo", "POST", &informacion_cargo, &filtrodatos); err == nil {

			num_contrato := datos.PersonasPreLiquidacion[i].NumeroContrato
			if err2 := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/docente_cargo/consultarCedulaDocente", "POST", &cedula, &num_contrato); err == nil {

				fmt.Println(informacion_cargo[0])
				regimen := informacion_cargo[0].Regimen
				puntos = consumir_puntos(cedula)
				tiempo_contrato := CalcularDias(informacion_cargo[0].FechaInicio, time.Now())
				reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos)
				reglas = reglasinyectadas + reglasbase
				temp := golog.CargarReglasDP(datos.PersonasPreLiquidacion[i].IdPersona, reglas, informacion_cargo, tiempo_contrato, datos.Preliquidacion.Nomina.Periodo, puntos, regimen,tipoNom)
				resultado := temp[len(temp)-1]
				resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].IdPersona)
				resumen_preliqu = append(resumen_preliqu, resultado)

				for _, descuentos := range *resultado.Conceptos {
					valor, _ := strconv.ParseInt(descuentos.Valor, 10, 64)

					detallepreliqu := models.DetallePreliquidacion{Concepto: &models.Concepto{Id: descuentos.Id}, Persona: datos.PersonasPreLiquidacion[i].IdPersona, Preliquidacion: datos.Preliquidacion.Id, ValorCalculado: valor, NumeroContrato: &models.ContratoGeneral{Id: datos.PersonasPreLiquidacion[i].NumeroContrato}}
					if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

					} else {
						beego.Debug("error1: ", err)
					}
				}

			} else {
				fmt.Println(err2)
			}

		} else {
			fmt.Println(err)
		}
	}
	return resumen_preliqu
}

func consumir_puntos(cedula int) (res float64) {

	var puntos float64
	resp, err := http.PostForm("http://"+beego.AppConfig.String("UrlKyron")+"/kyron/index.php?pagina=estadoDeCuentaCondor&bloqueNombre=estadoDeCuentaCondor&bloqueGrupo=reportes&procesarAjax=true&action=query&format=jwt", url.Values{"usuario": {beego.AppConfig.String("UsuarioKyron")}, "clave": {beego.AppConfig.String("Clave")}})
	if err != nil {
		fmt.Println("Error al consumir puntos")
		}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	token := string(body[:])
	docente_cedula := strconv.Itoa(cedula)
	//docente_cedula := strconv.Itoa(cedula)
	resp, err = http.Get("http://"+beego.AppConfig.String("UrlKyron")+"/kyron/index.php?data=" + token + "&docente=" + docente_cedula)
	if err != nil {
		fmt.Println("Error2 al consumir Puntos_salariales")
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	puntos_kyron := string(body[:])
	byt := []byte(puntos_kyron)
	var puntos_retorno models.Docente_puntos
	if err := json.Unmarshal(byt, &puntos_retorno); err == nil {
			puntos = puntos_retorno.Puntos_salariales
	}

	return puntos
}
