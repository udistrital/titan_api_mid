package golog

import (
	"fmt"
	"strconv"
	"time"
	models "github.com/udistrital/titan_api_mid/models"

	. "github.com/mndrix/golog"
)



var salario string


func CargarReglasDP(fechaPreliquidacion time.Time, dias_laborados float64, idProveedor int, reglas string, informacion_cargo []models.DocenteCargo, dias_trabajados float64, periodo string, puntos string, regimen string,tipoPreliquidacion string) (rest []models.Respuesta) {
	//Definición de variables

	var resultado []models.Respuesta
	var lista_descuentos []models.ConceptosResumen
	var lista_novedades []models.ConceptosResumen
	var lista_descuentos_semestral []models.ConceptosResumen
	var tipoPreliquidacion_string string
	var regimen_numero string
	var cargo string


	fechaInicio := informacion_cargo[0].FechaInicio
	fechaActual := time.Now().Local()
	asignacion_basica_string := strconv.Itoa(informacion_cargo[0].Asignacion_basica)
	tipoPreliquidacion_string = tipoPreliquidacion
	dias_laborados_string := strconv.Itoa(int(dias_laborados))

	if informacion_cargo[0].Cargo == "DC" {
		cargo = "1"
	} else {
		cargo = "2"
	}

	if regimen == "N" {
		fmt.Println("Nuevo")
		regimen_numero = "1"
	} else {
		fmt.Println("antiguo")
		regimen_numero = "2"
	}

	if tipoPreliquidacion_string  == "0" || tipoPreliquidacion_string  == "1" {
		dias_a_liquidar = "15"

	} else {
		dias_a_liquidar = "30"
	}

		m := NewMachine().Consult(reglas)

	fmt.Println(regimen_numero + " " + " " + puntos + " " + asignacion_basica_string + " " + cargo)

	if(int(fechaPreliquidacion.Month()) == 6){
		dias_liq_ps := m.ProveAll("dias_liq_ps("+dias_laborados_string+","+regimen_numero+",V).")
		for _, solution := range dias_liq_ps{
				dias_liquidar_prima_semestral = fmt.Sprintf("%s", solution.ByName_("V"))
		}

		doceava_BSPS := CalcularDoceavaBonServPSDP(reglas,"3", idProveedor, periodo, fechaPreliquidacion)
		lista_descuentos_semestral = CalcularConceptosDP(m, reglas,dias_liquidar_prima_semestral,asignacion_basica_string, "3",regimen_numero, puntos, cargo, fechaInicio, fechaActual)
		total_calculos = append (total_calculos, lista_descuentos_semestral...)
		total_calculos = append (total_calculos, 	doceava_BSPS...)
		ibc = 0
		}

		//-----NOMINA DE DICIEMBRE ------
		if(int(fechaPreliquidacion.Month()) == 12){
			dias_liq_dic := m.ProveAll("dias_liq_dic("+regimen_numero+",TLIQ,D).")
			dias_a_liquidar = "9"
			for _, solution := range dias_liq_dic{

					tipoLiq := fmt.Sprintf("%s", solution.ByName_("TLIQ"))
					fmt.Println("hola hola")
					fmt.Println(tipoLiq)
					dias_liquidacion_diciembre := fmt.Sprintf("%s", solution.ByName_("D"))
					lista_descuentos_semestral = CalcularConceptosDP(m, reglas,dias_liquidacion_diciembre,asignacion_basica_string, tipoLiq ,regimen_numero, puntos, cargo, fechaInicio, fechaActual)
					total_calculos = append (total_calculos, lista_descuentos_semestral...)
					doceavas_bsd := CalcularDoceavaBonServDicDP(reglas,tipoLiq, idProveedor, periodo, fechaPreliquidacion)
					doceavas_psd := CalcularDoceavaPSDicDP(reglas,tipoLiq, idProveedor, periodo, fechaPreliquidacion)
					total_calculos = append (total_calculos, doceavas_bsd...)
					total_calculos = append (total_calculos, doceavas_psd...)

					if(tipoLiq == "6"){
						fmt.Println("hola 13")
						doceavas_pv := CalcularDoceavaPVDP(reglas, tipoLiq, total_calculos)
						fmt.Println(doceavas_pv)
						total_calculos = append (total_calculos, doceavas_pv...)
					}

					ibc = 0
			}
		}
		
	// ----- Nomina ordinaria ----- Proceso de cálculo, manejo de novedades y guardado de conceptos
	lista_descuentos = CalcularConceptosDP(m, reglas,dias_a_liquidar,asignacion_basica_string, tipoPreliquidacion_string,regimen_numero, puntos, cargo, fechaInicio, fechaActual)
	ibc = 0
	lista_novedades = ManejarNovedadesDP(reglas,idProveedor, tipoPreliquidacion_string)
	total_calculos = append(total_calculos, lista_descuentos...)
	total_calculos = append(total_calculos, lista_novedades...)
	resultado = GuardarConceptosDP(total_calculos)
	total_calculos = []models.ConceptosResumen{}

	// ---------------------------
	return resultado;

	//falta arreglar el periodo para que sea congruente con los valores provenientes de la bd liquidar(R,P,V,T,C,L)
}

	func CalcularConceptosDP(m Machine, reglas, dias_a_liquidar, asignacion_basica_string, tipoPreliquidacion_string, regimen_numero, puntos, cargo string, fechaInicio, fechaActual time.Time) (rest []models.ConceptosResumen){
		var lista_descuentos []models.ConceptosResumen

		novedades_devengo := m.ProveAll("novedades_devengos(X).")
		for _, solution := range novedades_devengo {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("X")), 64)
			ibc = ibc + Valor

		}

		valor_salario := m.ProveAll("liquidar(" + regimen_numero + "," + puntos + "," + asignacion_basica_string + ", "+dias_a_liquidar+"," + cargo + ",L ).")
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
				temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			}

			lista_descuentos = append(lista_descuentos, temp_conceptos)

		}

		if fechaInicio.Month() == fechaActual.Month() && regimen_numero == "2" {
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
					temp_conceptos.DiasLiquidados = dias_a_liquidar
					temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
				}

				lista_descuentos = append(lista_descuentos, temp_conceptos)
		  }
		}

		//Previo a pagos de salud y pensión se calcula el IBC
		CalcularIBC(reglas)
		total_devengado_string := strconv.Itoa(int(ibc))
		fmt.Println("total devengado")


		salud_empleado := m.ProveAll("salud(" + total_devengado_string + ",S).")
		for _, solution := range salud_empleado {
		  Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("S")), 64)
		  temp_conceptos := models.ConceptosResumen{Nombre: "salud",
		    Valor: fmt.Sprintf("%.0f", Valor),
		  }

		codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

 		 for _, cod := range codigo {
 			 temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
 			 temp_conceptos.DiasLiquidados = dias_a_liquidar
 			 temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
 		 }
 		 lista_descuentos = append(lista_descuentos, temp_conceptos)

		}

		pension_empleado := m.ProveAll("pension(" + total_devengado_string + ",S).")
		for _, solution := range pension_empleado {
		  Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("S")), 64)
		  temp_conceptos := models.ConceptosResumen{Nombre: "pension",
		    Valor: fmt.Sprintf("%.0f", Valor),
		  }

		codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

	 		 for _, cod := range codigo {
	 			 temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
	 			 temp_conceptos.DiasLiquidados = dias_a_liquidar
	 			 temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
	 		 }
	 		 lista_descuentos = append(lista_descuentos, temp_conceptos)

			}

		return lista_descuentos
	}

	func ManejarNovedadesDP(reglas string, idProveedor int, tipoPreliquidacion string) (rest []models.ConceptosResumen){
		var lista_novedades []models.ConceptosResumen

		f := NewMachine().Consult(reglas)

		idProveedorString := strconv.Itoa(idProveedor)
		novedades := f.ProveAll("info_concepto(" + idProveedorString + ",T,2017,N,R).")

		for _, solution := range novedades {

			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
			temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
				Valor: fmt.Sprintf("%.0f", Valor),
			}
			codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion
			}

			lista_novedades = append(lista_novedades, temp_conceptos)

		}

		return lista_novedades

	}


func GuardarConceptosDP(lista_descuentos []models.ConceptosResumen)(rest []models.Respuesta){
			temp := models.Respuesta{}
			var resultado []models.Respuesta

			/*temp_conceptos := models.ConceptosResumen{Nombre: "ibc_liquidado",
				Valor: fmt.Sprintf("%.0f", total_devengado_no_novedad),
			}
			temp_conceptos.Id = 2322
			temp_conceptos.DiasLiquidados = dias_a_liquidar

			lista_descuentos = append(lista_descuentos, temp_conceptos)

			temp_conceptos = models.ConceptosResumen{Nombre: "ibc_novedad",
				Valor: fmt.Sprintf("%.0f", total_devengado_novedad),
			}
			temp_conceptos.Id = 2327
			temp_conceptos.DiasLiquidados = dias_novedad_string


			lista_descuentos = append(lista_descuentos, temp_conceptos)
*/
			temp.Conceptos = &lista_descuentos
			resultado = append(resultado, temp)
			total_devengado_novedad = 0
			total_devengado_no_novedad = 0
			return resultado
	}

func CalcularDoceavaBonServPSDP(reglas string,tipoPreliquidacion_string string, idProveedor int, periodo string, fechaPreliquidacion time.Time) (rest []models.ConceptosResumen){

		var lista_doceavas []models.ConceptosResumen
		var total_sumado int64

			tipoPreliquidacion_string = "3"
			f := NewMachine().Consult(reglas)
			consultar_valores_bonificacion := f.ProveAll("concepto_bon_serv_ps(X).")
			 for _, solution := range consultar_valores_bonificacion {

				codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
				total_sumado = total_sumado + ConsultarValoresBonServPS(fechaPreliquidacion, idProveedor,codigo_concepto, periodo)

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
					temp_conceptos.DiasLiquidados = dias_a_liquidar
					temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
				}

				lista_doceavas = append(lista_doceavas, temp_conceptos)
			}
			return lista_doceavas

	}

	//Función que calcula doceavas de bonificacion por servicios, doceava de prima semestral, doceava de prima vacaciones y adhiere resutltado a lo anterior
	func CalcularDoceavaBonServDicDP(reglas string,tipoPreliquidacion_string string, idProveedor int, periodo string, fechaPreliquidacion time.Time) (rest []models.ConceptosResumen){

		var lista_doceavas []models.ConceptosResumen
		var total_sumado int64
		f := NewMachine().Consult(reglas)

	 	consultar_valores_bonificacion := f.ProveAll("concepto_bon_serv_dic(X).")
		 for _, solution := range consultar_valores_bonificacion {
			codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
			total_sumado = total_sumado + ConsultarValoresBonServDic(fechaPreliquidacion, idProveedor,codigo_concepto, periodo)

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
				temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			}

			lista_doceavas = append(lista_doceavas, temp_conceptos)

		}

		return lista_doceavas

	}

	func CalcularDoceavaPSDicDP(reglas string,tipoPreliquidacion_string string, idProveedor int, periodo string, fechaPreliquidacion time.Time) (rest []models.ConceptosResumen){

		var lista_doceavas []models.ConceptosResumen
		var total_sumado int64
		f := NewMachine().Consult(reglas)

		total_sumado = ConsultarValoresPriServDic(fechaPreliquidacion, idProveedor, periodo)



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
				temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			}

			lista_doceavas = append(lista_doceavas, temp_conceptos)

		}

		return lista_doceavas


	}

	func CalcularDoceavaPVDP(reglas string, tipoPreliquidacion_string string, total_calculado []models.ConceptosResumen) (rest []models.ConceptosResumen){

		var lista_doceavas []models.ConceptosResumen
		var total_sumado int64

		for _, solution := range total_calculado {
				if(solution.TipoPreliquidacion == "4"){
					i, _ := strconv.ParseInt(solution.Valor, 10, 64)
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
					temp_conceptos.DiasLiquidados = dias_a_liquidar
					temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
				}

				lista_doceavas = append(lista_doceavas, temp_conceptos)

	}
		return lista_doceavas
	}
