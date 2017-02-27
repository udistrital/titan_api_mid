package controllers

import (
	"fmt"
	"strconv"
	"time"
	"titan_api_mid/golog"
	"titan_api_mid/models"

	"github.com/astaxie/beego"
)

// operations for Preliquidaciondp
type PreliquidaciondpController struct {
	beego.Controller
}

func (c *PreliquidaciondpController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	var resumen_preliqu []models.Respuesta
	var idDetaPre interface{}
	//var puntos []models.Docente_puntos

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {
		var informacion_cargo []models.DocenteCargo
		var cedula int
		var reglasinyectadas string
		var reglas string
		filtrodatos := models.DocenteCargo{Id: datos.PersonasPreLiquidacion[i].IdPersona, Asignacion_basica: 0}
		//fmt.Println("reglas: ", reglasbase)
		//consulta que envie ID de proveedor en datos y retorne el salario, para que sea enviado a CargarReglas
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/docente_cargo", "POST", &informacion_cargo, &filtrodatos); err == nil {

			num_contrato := datos.PersonasPreLiquidacion[i].NumeroContrato
			if err2 := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/docente_cargo/consultarCedulaDocente", "POST", &cedula, &num_contrato); err == nil {
				puntos := strconv.FormatFloat(informacion_cargo[0].Puntos, 'f', 6, 64)
				fmt.Println(informacion_cargo[0])
				regimen := informacion_cargo[0].Regimen
				//puntos := consumir_puntos(cedula)
				tiempo_contrato := CalcularDias(informacion_cargo[0].FechaInicio, time.Now())
				reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos)
				reglas = reglasinyectadas + reglasbase
				temp := golog.CargarReglasDP(datos.PersonasPreLiquidacion[i].IdPersona, reglas, informacion_cargo, tiempo_contrato, datos.Preliquidacion.Nomina.Periodo, puntos, regimen)
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

func consumir_puntos(cedula int) (res string) {
	var docente_puntos models.Docente_puntos
	var puntos models.Puntos
	docente_puntos.Id = 408
	//docente_puntos.Documento = "79708124"
	docente_puntos.Puntos = 10
	ruta := ""
	fmt.Println(ruta)
	if err := getJson("http://10.20.0.127/kyron/index.php?data="+ruta, &puntos); err == nil {
		fmt.Println(puntos.Puntos_salariales)
		fmt.Println(puntos.Puntos_bonificacion)
	} else {
		fmt.Println(err)
	}
	puntos_retorno := strconv.FormatFloat(puntos.Puntos_salariales, 'f', 6, 64)
	return puntos_retorno
}

/*
func aes_256(cedula string) (res string) {
	key := []byte("")
	data := "pagina=estadoDeCuentaCondor&bloqueNombre=estadoDeCuentaCondor&bloqueGrupo=reportes&docente=" + cedula + "&expiracion=1483946906&procesarAjax=true&action=query&format=json"
	plaintext := []byte(data)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return b64.URLEncoding.EncodeToString(ciphertext)
}

*/
