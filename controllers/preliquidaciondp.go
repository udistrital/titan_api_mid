package controllers

import (
	"fmt"
	"strconv"

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
	var datos_pruebas []models.DatosPruebas
	var arreglo_pruebas []models.PruebaGoDocentes
	arreglo_pruebas = make([]models.PruebaGoDocentes, len(datos.PersonasPreLiquidacion))


	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {
		var informacion_cargo []models.DocenteCargo
		var cedula int
		var reglasinyectadas string
		var reglas string
		filtrodatos := models.DocenteCargo{Id: datos.PersonasPreLiquidacion[i].IdPersona, Asignacion_basica: 0}
		tipoNom = "2"
		tipoNom_int, _ :=  strconv.Atoi(tipoNom)
		//fmt.Println("reglas: ", reglasbase)
		//consulta que envie ID de proveedor en datos y retorne el salario, para que sea enviado a CargarReglas
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/docente_cargo", "POST", &informacion_cargo, &filtrodatos); err == nil {
			fmt.Println("infocargo")
			fmt.Println(informacion_cargo)
			num_contrato := datos.PersonasPreLiquidacion[i].NumeroContrato
			if err2 := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/docente_cargo/consultarCedulaDocente", "POST", &cedula, &num_contrato); err == nil {
				dias_laborados := CalcularDias(informacion_cargo[0].FechaInicio, informacion_cargo[0].FechaFin)
				puntos := strconv.FormatFloat(informacion_cargo[0].Puntos, 'f', 6, 64)
				regimen := informacion_cargo[0].Regimen
				//puntos = consumir_puntos(cedula)

				reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.Preliquidacion)
				reglas = reglasinyectadas + reglasbase

				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/datos_pruebas?limit=-1&query=MesPreliq:"+strconv.Itoa(datos.Preliquidacion.Mes)+",AnoPreliq:"+strconv.Itoa(datos.Preliquidacion.Ano)+",NumDocumento:"+strconv.Itoa(datos.PersonasPreLiquidacion[i].NumDocumento), &datos_pruebas); err == nil && datos_pruebas != nil{
					arreglo_pruebas[i] = models.PruebaGoDocentes{informacion_cargo, "",datos.Preliquidacion.FechaRegistro, datos_pruebas[0].ValorSalario,"","","","",datos_pruebas[0].ValorPrimaTecnica,datos_pruebas[0].ValorPrimaAnt,datos_pruebas[0].ValorSalud,datos_pruebas[0].ValorPension,datos.PersonasPreLiquidacion[i].IdPersona,datos.PersonasPreLiquidacion[i].NumDocumento,dias_laborados,datos.Preliquidacion.Mes,datos.Preliquidacion.Ano, 0,0, tipoNom_int}

				}else{
					fmt.Println(err)
				}

				temp := golog.CargarReglasDP(datos.Preliquidacion.Mes, datos.Preliquidacion.Ano,dias_laborados, datos.PersonasPreLiquidacion[i].IdPersona, datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato,reglas, informacion_cargo, puntos, regimen,tipoNom)
				resultado := temp[len(temp)-1]
				resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].IdPersona)
				resumen_preliqu = append(resumen_preliqu, resultado)

				for _, descuentos := range *resultado.Conceptos {
					valor, _ := strconv.ParseFloat(descuentos.Valor,64)
					dias_liquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados,64)
					tipo_preliquidacion,_ := strconv.Atoi(descuentos.TipoPreliquidacion)
					detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: datos.Preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: datos.PersonasPreLiquidacion[i].NumeroContrato,VigenciaContrato: datos.PersonasPreLiquidacion[i].VigenciaContrato, DiasLiquidados: dias_liquidados, TipoPreliquidacion: &models.TipoPreliquidacion {Id: tipo_preliquidacion}}

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
			reglasinyectadas = "";
	}

	data, err := json.Marshal(arreglo_pruebas)
	if err != nil {
			fmt.Println("error en json")
		}
	str := fmt.Sprintf("%s", data)
	mes := strconv.Itoa(datos.Preliquidacion.Mes)
	ano := strconv.Itoa(datos.Preliquidacion.Ano)
	if err := WriteStringToFile("pruebaDocentesPlanta"+ano+mes+".txt", str); err != nil {
			panic(err)
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
	fmt.Println("puntos")
	fmt.Println(puntos)
	return puntos
}
