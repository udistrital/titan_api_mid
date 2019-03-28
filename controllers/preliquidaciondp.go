package controllers

import (
	"fmt"
	"strconv"

	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
   "net/http"
	 "net/url"
	 "io/ioutil"
	"github.com/astaxie/beego"
	"encoding/json"
)

// PreliquidaciondpController operations for Preliquidaciondp
type PreliquidaciondpController struct {
	beego.Controller
}

// Preliquidar ...
// @Title Preliquidar
// @Description Preliquidacion para docentes de planta
func (c *PreliquidaciondpController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	var resumenPreliqu []models.Respuesta
	var idDetaPre interface{}
	var tipoNom string;



	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {
		var informacionCargo []models.DocenteCargo
		var cedula int
		var reglasinyectadas string
		var reglas string
		filtrodatos := models.DocenteCargo{Id: datos.PersonasPreLiquidacion[i].IdPersona, Asignacion_basica: 0}
		tipoNom = "2"

		//fmt.Println("reglas: ", reglasbase)
		//consulta que envie ID de proveedor en datos y retorne el salario, para que sea enviado a CargarReglas
		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/docente_cargo", "POST", &informacionCargo, &filtrodatos); err == nil {
			fmt.Println("infocargo")
			fmt.Println(informacionCargo)
			numContrato := datos.PersonasPreLiquidacion[i].NumeroContrato
			if err2 := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/docente_cargo/consultarCedulaDocente", "POST", &cedula, &numContrato); err == nil {
				diasLaborados := CalcularDias(informacionCargo[0].FechaInicio, informacionCargo[0].FechaFin)
				puntos := strconv.FormatFloat(informacionCargo[0].Puntos, 'f', 6, 64)
				regimen := informacionCargo[0].Regimen
				esAnual := esAnual(datos.Preliquidacion.Mes, informacionCargo[0].FechaInicio)
				//puntos = consumirPuntos(cedula)


				reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato,  strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato), datos.Preliquidacion)
				reglas = reglasinyectadas + reglasbase + esAnual

				temp := golog.CargarReglasDP(datos.Preliquidacion.Mes, datos.Preliquidacion.Ano,diasLaborados, datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato,reglas, informacionCargo, puntos, regimen,tipoNom)
				resultado := temp[len(temp)-1]
				resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].IdPersona)
				resumenPreliqu = append(resumenPreliqu, resultado)

				for _, descuentos := range *resultado.Conceptos {
					valor, _ := strconv.ParseFloat(descuentos.Valor,64)
					diasLiquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados,64)
					tipoPreliquidacion,_ := strconv.Atoi(descuentos.TipoPreliquidacion)
					detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: datos.Preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: datos.PersonasPreLiquidacion[i].NumeroContrato,VigenciaContrato: datos.PersonasPreLiquidacion[i].VigenciaContrato, DiasLiquidados: diasLiquidados, TipoPreliquidacion: &models.TipoPreliquidacion {Id: tipoPreliquidacion}}

					if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

					} else {
						fmt.Println("error1: ", err)
					}
				}

			} else {
				fmt.Println(err2)
			}

		} else {
			fmt.Println(err)
		}
			reglasinyectadas = "";
	}


	return resumenPreliqu
}

func consumirPuntos(cedula int) (res float64) {

	var puntos float64
	resp, err := http.PostForm("http://"+beego.AppConfig.String("UrlKyron")+"/kyron/index.php?pagina=estadoDeCuentaCondor&bloqueNombre=estadoDeCuentaCondor&bloqueGrupo=reportes&procesarAjax=true&action=query&format=jwt", url.Values{"usuario": {beego.AppConfig.String("UsuarioKyron")}, "clave": {beego.AppConfig.String("Clave")}})
	if err != nil {
		fmt.Println("Error al consumir puntos")
		}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	token := string(body[:])
	docenteCedula := strconv.Itoa(cedula)
	//docenteCedula := strconv.Itoa(cedula)
	resp, err = http.Get("http://"+beego.AppConfig.String("UrlKyron")+"/kyron/index.php?data=" + token + "&docente=" + docenteCedula)
	if err != nil {
		fmt.Println("Error2 al consumir Puntos_salariales")
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	puntosKyron := string(body[:])
	byt := []byte(puntosKyron)
	var puntosRetorno models.Docente_puntos
	if err := json.Unmarshal(byt, &puntosRetorno); err == nil {
			puntos = puntosRetorno.Puntos_salariales
	}
	fmt.Println("puntos")
	fmt.Println(puntos)
	return puntos
}
