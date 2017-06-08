package controllers

import (
	"fmt"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"strconv"

	"github.com/astaxie/beego"

)

// PreliquidacionpeController operations for Preliquidacionpe
type PreliquidacionpeController struct {
	beego.Controller
}

func (c *PreliquidacionpeController) Preliquidar(datos *models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	//var predicados []models.Predicado //variable para inyectar reglas
	var resumen_preliqu []models.Respuesta
	var idDetaPre interface{}
	var pensionados []models.InformacionPensionado // arreglo de informacion_pensionado
	var sustitutos []models.Sustituto //arreglo de sustitutos
	var tutor []models.Sustituto
	var beneficiarioF int //Beneficiarios con sub familiar
	var beneficiarioE int //Beneficiarios con aux de estudio
	var beneficiarios []models.Beneficiarios
	var tipoNom string;

	var reglasinyectadas string
	var reglas string


	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {

		filtrodatos := models.InformacionPensionado{Id: datos.PersonasPreLiquidacion[i].IdPersona}
		tipoNom = tipoNomina(datos.Preliquidacion.Tipo)
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/informacion_pensionado", "POST", &pensionados, &filtrodatos); err == nil {

			var idPensionado = pensionados[0].InformacionProveedor
			if err6 := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/beneficiarios/beneficiarioDatos", "POST", &beneficiarios, &idPensionado); err6 == nil {
				fmt.Println("Beneficiarios")
				fmt.Println(beneficiarios)
				}else{
					fmt.Println(err6)
				}

				if pensionados[0].Estado == "R"{
					fmt.Println("Persona retirada")
					if err4 := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/sustituto/sustitutoDatos", "POST", &sustitutos,&idPensionado); err4 == nil {
						fmt.Println("sustitutos")
						fmt.Println(sustitutos)

						if sustitutos[i].Tutor != 0 {
							//fmt.Println("Tutor")
							//fmt.Println(sustitutos[i].Tutor)
							if err5 := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/sustituto/tutorDatos", "POST", &tutor,&idPensionado); err5 == nil {

								}else{
									fmt.Println(err5)
								}
							}
							}else{
								fmt.Println(err4)
							}
						}
						reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(datos.PersonasPreLiquidacion[i].IdPersona, datos)
						reglas = reglasbase + reglasinyectadas


						if len(sustitutos) == 0{
							for i := 0; i < len(beneficiarios); i++ {

								if beneficiarios[i].CategoriaBeneficiario == 1 || beneficiarios[i].CategoriaBeneficiario == 2 || beneficiarios[i].CategoriaBeneficiario == 3{
									if beneficiarios[i].SubFamiliar == "S" && beneficiarios[i].Estado == "A"{
										beneficiarioF = beneficiarioF + 1
									}
								}
								if beneficiarios[i].SubEstudios  == "S" && beneficiarios[i].CategoriaBeneficiario == 3 && beneficiarios[i].Estado == "A"{
									beneficiarioE = beneficiarioE + 1
									}
								}
								fmt.Println("beneficiariooos")
								fmt.Println(beneficiarioE)
								temp := golog.CargarReglasPE(datos.Preliquidacion.Fecha,reglas, datos.Preliquidacion.Nomina.Periodo, pensionados[0],beneficiarioF, beneficiarioE, tipoNom)
								resultado := temp[len(temp)-1]
								resultado.NumDocumento = float64(datos.PersonasPreLiquidacion[i].IdPersona)
								resumen_preliqu = append(resumen_preliqu, resultado)

								fmt.Println("resultado Preliquidacion")
								fmt.Println(resumen_preliqu[0].Conceptos)

								for _, descuentos := range *resultado.Conceptos {
									valor, _ := strconv.ParseInt(descuentos.Valor, 10, 64)
									//fmt.Println("asdfg"+datos.PersonasPreLiquidacion[i].NumeroContrato)
									detallepreliqu := models.DetallePreliquidacion{Concepto: &models.Concepto{Id: descuentos.Id}, Persona: datos.PersonasPreLiquidacion[i].IdPersona, Preliquidacion: datos.Preliquidacion.Id, ValorCalculado: valor, NumeroContrato: &models.ContratoGeneral{Id: datos.PersonasPreLiquidacion[i].NumeroContrato}, DiasLiquidados: descuentos.DiasLiquidados, TipoPreliquidacion: descuentos.TipoPreliquidacion}
									if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

										} else {
											beego.Debug("error1: ", err)
										}
									}
									}else{ //else de sustitos
										for i := 0; i < len(sustitutos); i++{
											//fmt.Println(sustitutos[i].Tutor)

											var cedulaPensionado = strconv.Itoa(pensionados[0].InformacionProveedor)
											var pension = strconv.Itoa(pensionados[0].ValorPensionAsignada)

											temp := golog.CargarReglasSustitutosPE(reglas, sustitutos[i], cedulaPensionado ,pension,datos.Preliquidacion.Nomina.Periodo)
											resultado := temp[len(temp)-1]
											resultado.NumDocumento = float64(sustitutos[0].Proveedor)
											resumen_preliqu = append(resumen_preliqu, resultado)

											fmt.Println("resultado Preliquidacion")
											fmt.Println(resumen_preliqu[0].Conceptos)

											for _, descuentos := range *resultado.Conceptos {
												valor, _ := strconv.ParseInt(descuentos.Valor, 10, 64)
												if sustitutos[i].Tutor == 0 {
													detallepreliqu := models.DetallePreliquidacion{Concepto: &models.Concepto{Id: descuentos.Id}, Persona: sustitutos[i].Proveedor, Preliquidacion: datos.Preliquidacion.Id, ValorCalculado: valor, NumeroContrato: &models.ContratoGeneral{Id: sustitutos[i].NumeroContrato}}

													if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

														} else {
															beego.Debug("error1: ", err)
														}
														}else{
															detallepreliqu := models.DetallePreliquidacion{Concepto: &models.Concepto{Id: descuentos.Id}, Persona: tutor[i].Proveedor, Preliquidacion: datos.Preliquidacion.Id, ValorCalculado: valor, NumeroContrato: &models.ContratoGeneral{Id: tutor[i].NumeroContrato}}
															fmt.Println("Id Sustituto")
															fmt.Println(sustitutos[0].NumeroContrato)
															//fmt.Println(sustitutos[i].NumeroContrato)
															if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

																} else {
																	beego.Debug("error1: ", err)
																}
															}
														}
													}
												}
												}else {
													fmt.Println(err)
												}
													reglasinyectadas = "";
											}
											return resumen_preliqu
										}
