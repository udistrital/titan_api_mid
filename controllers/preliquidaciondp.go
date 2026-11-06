package controllers

import (
	"github.com/astaxie/beego"
)

// PreliquidaciondpController operations for Preliquidaciondp
type PreliquidaciondpController struct {
	beego.Controller
}

// Preliquidar ...
// @Title Preliquidar
// @Description Preliquidacion para docentes de planta
//func (c *PreliquidaciondpController) Preliquidar(ctx context.Context, datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
/*
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

		//consulta que envie ID de proveedor en datos y retorne el salario, para que sea enviado a CargarReglas
		if status, err := request.PostWithContext(ctx, "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/docente_cargo", &filtrodatos, &informacionCargo); err == nil && status >= 200 && status <= 226 {
			numContrato := datos.PersonasPreLiquidacion[i].NumeroContrato
			if status, err2 := request.PostWithContext(ctx, "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/docente_cargo/consultarCedulaDocente", &numContrato, &cedula); err2 == nil && status >= 200 && status <= 226 {
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

					if status, err := request.PostWithContext(ctx, "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", &detallepreliqu, &idDetaPre); err == nil && status >= 200 && status <= 226 {

					} else {
					}
				}

			} else {
			}

		} else {
		}
			reglasinyectadas = "";
	}


	return resumenPreliqu
}
*/
/*
func consumirPuntos(cedula int) (res float64) {

	var puntos float64
	resp, err := http.PostForm("http://"+beego.AppConfig.String("UrlKyron")+"/kyron/index.php?pagina=estadoDeCuentaCondor&bloqueNombre=estadoDeCuentaCondor&bloqueGrupo=reportes&procesarAjax=true&action=query&format=jwt", url.Values{"usuario": {beego.AppConfig.String("UsuarioKyron")}, "clave": {beego.AppConfig.String("Clave")}})
	if err != nil {
		}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	token := string(body[:])
	docenteCedula := strconv.Itoa(cedula)
	//docenteCedula := strconv.Itoa(cedula)
	resp, err = http.Get("http://"+beego.AppConfig.String("UrlKyron")+"/kyron/index.php?data=" + token + "&docente=" + docenteCedula)
	if err != nil {
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	puntosKyron := string(body[:])
	byt := []byte(puntosKyron)
	var puntosRetorno models.Docente_puntos
	if err := json.Unmarshal(byt, &puntosRetorno); err == nil {
			puntos = puntosRetorno.Puntos_salariales
	}
	return puntos
}
*/
