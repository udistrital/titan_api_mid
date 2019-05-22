package golog

import (
	"fmt"
	"strconv"
	"time"
	models "github.com/udistrital/titan_api_mid/models"

	. "github.com/udistrital/golog"
)



var salario string

func CargarReglasDP(MesPreliquidacion int, AnoPreliquidacion int, dias_laborados float64, idProveedor int, numero_contrato string, vigencia_contrato int,reglas string, informacion_cargo []models.DocenteCargo, puntos string, regimen string,tipoPreliquidacion string) (rest []models.Respuesta) {
	//Definición de variables

	var resultado []models.Respuesta
	var listaDescuentos []models.ConceptosResumen
	var listaNovedades []models.ConceptosResumen
	var listaRetefuente []models.ConceptosResumen
	var listaDescuentos_semestral []models.ConceptosResumen
	var tipoPreliquidacionString string
	var regimen_numero string
	var cargo string
	var periodo string

	fechaInicio := informacion_cargo[0].FechaInicio
	fechaActual := time.Now().Local()
	asignacion_basica_string := strconv.Itoa(informacion_cargo[0].Asignacion_basica)
	tipoPreliquidacionString = tipoPreliquidacion
	dias_laborados_string := strconv.Itoa(int(dias_laborados))
	periodo = strconv.Itoa(AnoPreliquidacion)

	if informacion_cargo[0].Cargo == "DC" {
		cargo = "1"
	} else {
		cargo = "2"
	}

	if regimen == "N" {
		fmt.Println("Nuevo")
		regimen_numero = "1"
		reglas = reglas + "regimen(1)."
	} else {
		fmt.Println("antiguo")
		regimen_numero = "2"
		reglas = reglas + "regimen(2)."
	}

	if tipoPreliquidacionString  == "0" || tipoPreliquidacionString  == "1" {
		diasALiquidar = "15"

	} else {
		diasALiquidar = "30"
	}

	nombre_archivo = "reglas" + strconv.Itoa(idProveedor) + ".txt"
	reglas = reglas + "salario_base(" + asignacion_basica_string + ")."
	reglas = reglas + "tipo_nomina(" + tipoPreliquidacionString + ")."
	reglas = reglas + "cargo("+cargo+")."

			m := NewMachine().Consult(reglas)



		//-- NOVEDADES DE SEGURIDAD SOCIAL --
		novedades_seg_social := m.ProveAll("seg_social(N,A,M,D,AA,MM,DD).")

		for _, solution := range novedades_seg_social {


			novedad := fmt.Sprintf("%s", solution.ByName_("N"))
			AnoDesde,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("A")), 64)
			MesDesde,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("M")), 64)
			DiaDesde,_:= strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("D")), 64)
			AnoHasta,_:= strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("AA")), 64)
			MesHasta,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("MM")), 64)
			DiaHasta,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("DD")), 64)

			afectacion_seg_social := m.ProveAll("afectacion_seguridad("+novedad+").")
			for _, solution := range afectacion_seg_social {

					fmt.Println(solution)
					dias_novedad := CalcularDiasNovedades(MesPreliquidacion, AnoPreliquidacion, AnoDesde, MesDesde, DiaDesde, AnoHasta, MesHasta, DiaHasta)
					diasALiquidar = strconv.Itoa(int(30 - dias_novedad))
					fmt.Println("dias a liquidar")
					fmt.Println(diasALiquidar)
					diasNovedadString = strconv.Itoa(int(dias_novedad))
					_,total_devengado_novedad = CalcularConceptosDP(m, reglas,diasNovedadString,asignacion_basica_string, tipoPreliquidacionString,regimen_numero, puntos, cargo, fechaInicio, fechaActual)
					ibc = 0;
			}

			}

			//-------------------
	fmt.Println(regimen_numero + " " + " " + puntos + " " + asignacion_basica_string + " " + cargo)

	if(MesPreliquidacion == 6){
		dias_liq_ps := m.ProveAll("dias_liq_ps("+dias_laborados_string+","+regimen_numero+",V).")
		for _, solution := range dias_liq_ps{
				dias_liquidar_prima_semestral = fmt.Sprintf("%s", solution.ByName_("V"))
		}

		doceava_BSPS := CalcularDoceavaBonServPSDP(reglas,"3", numero_contrato, vigencia_contrato, periodo, MesPreliquidacion , AnoPreliquidacion)
		listaDescuentos_semestral,total_devengado_no_novedad_semestral = CalcularConceptosDP(m, reglas,dias_liquidar_prima_semestral,asignacion_basica_string, "3",regimen_numero, puntos, cargo, fechaInicio, fechaActual)
		fmt.Println(listaDescuentos_semestral)
		total_calculos = append (total_calculos, listaDescuentos_semestral...)
		total_calculos = append (total_calculos, 	doceava_BSPS...)
		ibc = 0
		}

		//-----NOMINA DE DICIEMBRE ------
		if(MesPreliquidacion == 12){
			dias_liq_dic := m.ProveAll("dias_liq_dic("+regimen_numero+",TLIQ,D).")
			diasALiquidar = "9" //Preguntar num dias
			for _, solution := range dias_liq_dic{

					tipoLiq := fmt.Sprintf("%s", solution.ByName_("TLIQ"))
					dias_liquidacion_diciembre := fmt.Sprintf("%s", solution.ByName_("D"))
					listaDescuentos_semestral,total_devengado_no_novedad_semestral = CalcularConceptosDP(m, reglas,dias_liquidacion_diciembre,asignacion_basica_string, tipoLiq ,regimen_numero, puntos, cargo, fechaInicio, fechaActual)
					total_calculos = append (total_calculos, listaDescuentos_semestral...)
					doceavas_bsd := CalcularDoceavaBonServDicDP(reglas,tipoLiq, numero_contrato, vigencia_contrato, periodo)
					doceavas_psd := CalcularDoceavaPSDicDP(reglas,tipoLiq, numero_contrato, vigencia_contrato, periodo)
					total_calculos = append (total_calculos, doceavas_bsd...)
					total_calculos = append (total_calculos, doceavas_psd...)
					ibc = 0
			}

			doceavas_pv := CalcularDoceavaPVDP(reglas, "6", total_calculos)
			total_calculos = append (total_calculos, doceavas_pv...)
		}

	// ----- Nomina ordinaria ----- Proceso de cálculo, manejo de novedades y guardado de conceptos
	listaDescuentos,total_devengado_no_novedad = CalcularConceptosDP(m, reglas,diasALiquidar,asignacion_basica_string, tipoPreliquidacionString,regimen_numero, puntos, cargo, fechaInicio, fechaActual)
	ibc = 0
	listaNovedades = ManejarNovedadesDP(reglas,idProveedor, tipoPreliquidacionString,periodo)
	listaRetefuente = CalcularReteFuentePlanta(tipoPreliquidacionString,reglas,periodo, listaDescuentos);
	total_calculos = append(total_calculos, listaDescuentos...)
	total_calculos = append(total_calculos, listaNovedades...)
	total_calculos = append(total_calculos, listaRetefuente...)
	resultado = GuardarConceptosDP(reglas,total_calculos)

	total_calculos = []models.ConceptosResumen{}

	// ---------------------------
	return resultado;

	//falta arreglar el periodo para que sea congruente con los valores provenientes de la bd liquidar(R,P,V,T,C,L)
}

	func CalcularConceptosDP(m Machine, reglas, diasALiquidar, asignacion_basica_string, tipoPreliquidacionString, regimen_numero, puntos, cargo string, fechaInicio, fechaActual time.Time) (rest []models.ConceptosResumen,  total_dev float64){

		var listaDescuentos []models.ConceptosResumen

		valor_salario := m.ProveAll("liquidar(" + regimen_numero + "," + puntos + "," + asignacion_basica_string + ", "+diasALiquidar+"," + cargo + ",L ).")
		for _, solution := range valor_salario {
		  Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("L")), 64)
		  temp_conceptos := models.ConceptosResumen{Nombre: "salarioBase",
		    Valor: fmt.Sprintf("%.0f", Valor),
		  }
			salario = strconv.FormatFloat(Valor, 'f', 6, 64)

			reglas = reglas + "sumar_ibc(salarioBase,"+strconv.Itoa(int(Valor))+")."
			codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = diasALiquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			}

			listaDescuentos = append(listaDescuentos, temp_conceptos)

		}


		  bonificacion := m.ProveAll("bonificacionServicios(" + salario + ",S).")
		  for _, solution := range bonificacion {

		    Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("S")), 64)
		    temp_conceptos := models.ConceptosResumen{Nombre: "bonServ",
		      Valor: fmt.Sprintf("%.0f", Valor),
		    }

				reglas = reglas + "sumar_ibc(salarioBase,"+strconv.Itoa(int(Valor))+")."
				codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

				for _, cod := range codigo {
					temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
					temp_conceptos.DiasLiquidados = diasALiquidar
					temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
				}

				listaDescuentos = append(listaDescuentos, temp_conceptos)
		  }


		//Previo a pagos de salud y pensión se calcula el IBC
		CalcularIBC("whatever",reglas)
		ManejarNovedadesDevengosDP(reglas, tipoPreliquidacionString)
		total_devengado_string := strconv.Itoa(int(ibc))

		salud_empleado := m.ProveAll("salud("+tipoPreliquidacionString+"," + total_devengado_string + ",S).")
		for _, solution := range salud_empleado {
		  Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("S")), 64)
		  temp_conceptos := models.ConceptosResumen{Nombre: "salud",
		    Valor: fmt.Sprintf("%.0f", Valor),
		  }

		codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

 		 for _, cod := range codigo {
 			 temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
 			 temp_conceptos.DiasLiquidados = diasALiquidar
 			 temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
 		 }
 		 listaDescuentos = append(listaDescuentos, temp_conceptos)

		}

		pension_empleado := m.ProveAll("pension("+tipoPreliquidacionString+"," + total_devengado_string + ",S).")
		for _, solution := range pension_empleado {
		  Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("S")), 64)
		  temp_conceptos := models.ConceptosResumen{Nombre: "pension",
		    Valor: fmt.Sprintf("%.0f", Valor),
		  }

		codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

	 		 for _, cod := range codigo {
	 			 temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
	 			 temp_conceptos.DiasLiquidados = diasALiquidar
	 			 temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
	 		 }
	 		 listaDescuentos = append(listaDescuentos, temp_conceptos)

			}

		return listaDescuentos,ibc
	}

	func ManejarNovedadesDP(reglas string, idProveedor int, tipoPreliquidacion,periodo string) (rest []models.ConceptosResumen){
		var listaNovedades []models.ConceptosResumen

		f := NewMachine().Consult(reglas)

		idProveedorString := strconv.Itoa(idProveedor)
		novedades := f.ProveAll("info_concepto(" + idProveedorString + ",T,"+periodo+",N,R).")

		for _, solution := range novedades {

			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
			temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
				Valor: fmt.Sprintf("%.0f", Valor),
			}
			codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = diasALiquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion
			}

			listaNovedades = append(listaNovedades, temp_conceptos)

		}

		return listaNovedades

	}

	func ManejarNovedadesDevengosDP(reglas string, tipoPreliquidacion string){

		f := NewMachine().Consult(reglas)
 		novedades_devengo := f.ProveAll("novedades_devengos(X).")
			for _, solution := range novedades_devengo {
				Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("X")), 64)
				ibc = ibc + Valor
			}


	}

func GuardarConceptosDP(reglas string,listaDescuentos []models.ConceptosResumen)(rest []models.Respuesta){
			temp := models.Respuesta{}
			var resultado []models.Respuesta

			m := NewMachine().Consult(reglas)

	    temp_conceptos := models.ConceptosResumen{Nombre: "ibc_liquidado",
	      Valor: fmt.Sprintf("%.0f", total_devengado_no_novedad),
	    }

	    codigo := m.ProveAll(`codigo_concepto(ibc_liquidado,C, N).`)

	    for _, cod := range codigo {
	      temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
	      temp_conceptos.DiasLiquidados = diasALiquidar
	    }

			listaDescuentos = append(listaDescuentos, temp_conceptos)

			temp_conceptos_1 := models.ConceptosResumen{Nombre: "ibc_novedad",
				Valor: fmt.Sprintf("%.0f", total_devengado_no_novedad),
			}

			codigo_1 := m.ProveAll(`codigo_concepto(ibc_novedad,C, N).`)

			for _, cod := range codigo_1 {
				temp_conceptos_1.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos_1.DiasLiquidados = diasALiquidar
			}


			listaDescuentos = append(listaDescuentos, temp_conceptos_1)


			temp.Conceptos = &listaDescuentos
			resultado = append(resultado, temp)
			total_devengado_novedad = 0
			total_devengado_no_novedad = 0
			return resultado
	}

func CalcularDoceavaBonServPSDP(reglas string,tipoPreliquidacionString string, numero_contrato string, vigencia_contrato int, periodo string, mesPreliq, anoPreliq int) (rest []models.ConceptosResumen){

		var lista_doceavas []models.ConceptosResumen
		var total_sumado float64

			tipoPreliquidacionString = "3"
			f := NewMachine().Consult(reglas)
			consultar_valores_bonificacion := f.ProveAll("concepto_bon_serv_ps(X).")
			 for _, solution := range consultar_valores_bonificacion {

				codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
				total_sumado = total_sumado + ConsultarValoresBonServPS(mesPreliq, anoPreliq, numero_contrato, vigencia_contrato,codigo_concepto, periodo)

			}

			reglas = reglas + "bonificacion_servicio_ps(bonServ,"+strconv.Itoa(int(total_sumado))+")."

			e := NewMachine().Consult(reglas)
		 	doc_bonServ := e.ProveAll("doceava_bs(N,"+dias_liquidar_prima_semestral+",V).")
			for _, solution := range doc_bonServ {

					Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
					temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
					Valor: fmt.Sprintf("%.0f", Valor),
				}

				codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
				for _, cod := range codigo {
					temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
					temp_conceptos.DiasLiquidados = diasALiquidar
					temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
				}

				lista_doceavas = append(lista_doceavas, temp_conceptos)
			}
			return lista_doceavas

	}

	//Función que calcula doceavas de bonificacion por servicios, doceava de prima semestral, doceava de prima vacaciones y adhiere resutltado a lo anterior
	func CalcularDoceavaBonServDicDP(reglas string,tipoPreliquidacionString string, numero_contrato string, vigencia_contrato int, periodo string) (rest []models.ConceptosResumen){

		var lista_doceavas []models.ConceptosResumen
		var total_sumado float64
		f := NewMachine().Consult(reglas)

	 	consultar_valores_bonificacion := f.ProveAll("concepto_bon_serv_dic(X).")
		 for _, solution := range consultar_valores_bonificacion {
			codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
			total_sumado = total_sumado + ConsultarValoresBonServDic(numero_contrato, vigencia_contrato,codigo_concepto, periodo)

		}

		reglas = reglas + "bonificacion_servicio(bonServ,"+strconv.Itoa(int(total_sumado))+")."

		e := NewMachine().Consult(reglas)
	 	doc_bonServ := e.ProveAll("doceava(N,V).")
		for _, solution := range doc_bonServ {
				Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
				temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
				Valor: fmt.Sprintf("%.0f", Valor),
			}

			codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = diasALiquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			}

			lista_doceavas = append(lista_doceavas, temp_conceptos)

		}

		return lista_doceavas

	}

	func CalcularDoceavaPSDicDP(reglas string,tipoPreliquidacionString string, numero_contrato string,vigencia_contrato int, periodo string) (rest []models.ConceptosResumen){

		var lista_doceavas []models.ConceptosResumen
		var total_sumado float64
		f := NewMachine().Consult(reglas)

		total_sumado = ConsultarValoresPriServDic(numero_contrato, vigencia_contrato, periodo)



		reglas = reglas + "prima_servicios(priServ,"+strconv.Itoa(int(total_sumado))+")."

		e := NewMachine().Consult(reglas)
		doc_bonServ := e.ProveAll("doceava_ps(N,V).")
		for _, solution := range doc_bonServ {
				Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
				temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
				Valor: fmt.Sprintf("%.0f", Valor),
			}

			codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = diasALiquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			}

			lista_doceavas = append(lista_doceavas, temp_conceptos)

		}

		return lista_doceavas


	}

	func CalcularDoceavaPVDP(reglas string, tipoPreliquidacionString string, total_calculado []models.ConceptosResumen) (rest []models.ConceptosResumen){

		var lista_doceavas []models.ConceptosResumen
		var total_sumado float64

		for _, solution := range total_calculado {
				if(solution.TipoPreliquidacion == "4"){
					i, _ := strconv.ParseFloat(solution.Valor, 64)
					total_sumado = total_sumado + i
				}
			}

			reglas = reglas + "prima_vacaciones(priVac,"+strconv.Itoa(int(total_sumado))+")."

			e := NewMachine().Consult(reglas)
			doc_bonServ := e.ProveAll("doceava_pv(N,V).")
			for _, solution := range doc_bonServ {
					Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
					temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
					Valor: fmt.Sprintf("%.0f", Valor),
				}

				codigo := e.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
				for _, cod := range codigo {
					temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
					temp_conceptos.DiasLiquidados = diasALiquidar
					temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
				}

				lista_doceavas = append(lista_doceavas, temp_conceptos)

	}
		return lista_doceavas
	}
