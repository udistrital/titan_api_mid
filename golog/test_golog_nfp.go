package golog

import (
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"
	. "github.com/mndrix/golog"

)


var total_devengado_no_novedad float64
var total_devengado_no_novedad_semestral float64
var total_devengado_novedad float64
var dias_a_liquidar string
var dias_novedad_string string
var nombre_archivo string
var ibc float64
var dias_liquidar_prima_semestral  string
var total_calculos []models.ConceptosResumen

func CargarReglasFP(MesPreliquidacion int, AnoPreliquidacion int, reglas string, idProveedor int, numero_contrato string, vigencia_contrato int, informacion_cargo []models.FuncionarioCargo, dias_laborados float64, esAnual int, porcentajePT int, tipoPreliquidacion int) (rest []models.Respuesta) {

	//--- Creación de variables
	var resultado []models.Respuesta
	var lista_descuentos []models.ConceptosResumen
	var lista_descuentos_semestral []models.ConceptosResumen
	var lista_novedades []models.ConceptosResumen
	var tipoPreliquidacion_string string

	//Conversión de variables
	asignacion_basica_string := strconv.Itoa(informacion_cargo[0].Asignacion_basica)
	id_cargo_string := strconv.Itoa(informacion_cargo[0].Id)
	dias_laborados_string := strconv.Itoa(int(dias_laborados))
	tipoPreliquidacion_string = strconv.Itoa(tipoPreliquidacion)
	periodo:= strconv.Itoa(AnoPreliquidacion)
	//Asignación de número de días segun tipo de nomina (0 es quince días, 1 es el mes completo)
	if tipoPreliquidacion_string  == "0" || tipoPreliquidacion_string  == "1" {
		dias_a_liquidar = "15"

	} else {
		dias_a_liquidar = "30"
	}


  nombre_archivo = "reglas" + strconv.Itoa(idProveedor) + ".txt"
	reglas = reglas + "salario_base(" + asignacion_basica_string + ")."
	reglas = reglas + "tipo_nomina(" + tipoPreliquidacion_string + ")."

	if err := WriteStringToFile(nombre_archivo, reglas); err != nil {
      panic(err)
  }

	m := NewMachine().Consult(reglas)

	//-- NOVEDADES DE SEGURIDAD SOCIAL --
	novedades_seg_social := m.ProveAll("seg_social(N,A,M,D,AA,MM,DD).")

	for _, solution := range novedades_seg_social {

		fmt.Println("aqui nov")
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
				dias_novedad := CalcularDiasNovedades(MesPreliquidacion, AnoPreliquidacion,AnoDesde, MesDesde, DiaDesde, AnoHasta, MesHasta, DiaHasta)
				dias_a_liquidar = strconv.Itoa(int(30 - dias_novedad))
				fmt.Println("dias a liquidar")

				dias_novedad_string = strconv.Itoa(int(dias_novedad))
				_,total_devengado_novedad = CalcularConceptos(m, reglas, dias_novedad_string,asignacion_basica_string,id_cargo_string,dias_laborados_string, tipoPreliquidacion_string,esAnual, porcentajePT, idProveedor,periodo)
				ibc = 0;
		}

		}

		//-------------------

		// -- PRIMA SEMESTRAL --

		if(MesPreliquidacion == 6){
			fmt.Println("aqui semestral")
			dias_liq_ps := m.ProveAll("dias_liq_ps("+dias_laborados_string+",V).")
			for _, solution := range dias_liq_ps{
					dias_liquidar_prima_semestral = fmt.Sprintf("%s", solution.ByName_("V"))
			}

			doceava_BSPS := CalcularDoceavaBonServPS(reglas,tipoPreliquidacion_string, numero_contrato, vigencia_contrato, periodo, MesPreliquidacion , AnoPreliquidacion)
			lista_descuentos_semestral,total_devengado_no_novedad_semestral = CalcularConceptos(m, reglas,dias_liquidar_prima_semestral,asignacion_basica_string,id_cargo_string,dias_laborados_string, "3",esAnual, porcentajePT, idProveedor,periodo)
			total_calculos = append (total_calculos, lista_descuentos_semestral...)
			total_calculos = append (total_calculos, 	doceava_BSPS...)
			ibc = 0
			}
		//--------------------------------

		//-----NOMINA DE DICIEMBRE ------
		if(MesPreliquidacion  == 12){
			fmt.Println("aqui dic")
			dias_liq_dic := m.ProveAll("dias_liq_dic(FP,TLIQ,D).")
			dias_a_liquidar = "9"
			for _, solution := range dias_liq_dic{
					tipoLiq := fmt.Sprintf("%s", solution.ByName_("TLIQ"))
					dias_liquidacion_diciembre := fmt.Sprintf("%s", solution.ByName_("D"))
					lista_descuentos_semestral,total_devengado_no_novedad_semestral = CalcularConceptos(m, reglas,dias_liquidacion_diciembre,asignacion_basica_string,id_cargo_string,dias_laborados_string, tipoLiq,esAnual, porcentajePT, idProveedor,periodo)
					total_calculos = append (total_calculos, lista_descuentos_semestral...)
					doceavas_bsd := CalcularDoceavaBonServDic(reglas,tipoLiq, numero_contrato, vigencia_contrato, periodo)
					doceavas_psd := CalcularDoceavaPSDic(reglas,tipoLiq, numero_contrato, vigencia_contrato, periodo)
					total_calculos = append (total_calculos, doceavas_bsd...)
					total_calculos = append (total_calculos, doceavas_psd...)

					if(tipoLiq == "6"){
						doceavas_pv := CalcularDoceavaPV(reglas, tipoPreliquidacion_string, total_calculos)
						total_calculos = append (total_calculos, doceavas_pv...)
					}

					ibc = 0
			}



		}
 		// ----------------------------------

		// ----- Nomina ordinaria ----- Proceso de cálculo, manejo de novedades y guardado de conceptos

		lista_descuentos,total_devengado_no_novedad = CalcularConceptos(m, reglas,dias_a_liquidar,asignacion_basica_string,id_cargo_string,dias_laborados_string, tipoPreliquidacion_string,esAnual, porcentajePT, idProveedor,periodo)
		ibc = 0
		lista_novedades = ManejarNovedades(reglas,idProveedor, tipoPreliquidacion_string, periodo)
		total_calculos = append(total_calculos, lista_descuentos...)
		total_calculos = append(total_calculos, lista_novedades...)
		resultado = GuardarConceptos(total_calculos)
		total_calculos = []models.ConceptosResumen{}

		// ---------------------------
		return resultado;

	}

//Función que, utilizando las reglas, calcula cada uno de los conceptos. Retorna un objeto con los resultados
	func CalcularConceptos(m Machine, reglas,dias_a_liquidar,asignacion_basica_string,id_cargo_string,dias_laborados_string,tipoPreliquidacion_string string,  esAnual,  porcentajePT,idProveedor  int,periodo string) (rest []models.ConceptosResumen, total_dev float64){

		var lista_descuentos []models.ConceptosResumen
		porcentaje_PT_string := strconv.Itoa(porcentajePT)

			//Cálculo de cada uno de los devengos, si aplican a la persona
		valor_salario := m.ProveAll("sb(" + asignacion_basica_string + "," + tipoPreliquidacion_string + "," + dias_a_liquidar + ",V).")
		for _, solution := range valor_salario {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
			temp_conceptos := models.ConceptosResumen{Nombre: "salarioBase",
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

		valor_gastos_representacion := m.ProveAll("gr(" + asignacion_basica_string + "," + dias_a_liquidar + "," + tipoPreliquidacion_string + ","+periodo+"," + id_cargo_string + ",V).")
		for _, solution := range valor_gastos_representacion {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
			temp_conceptos := models.ConceptosResumen{Nombre: "gastosRep",
				Valor: fmt.Sprintf("%.0f", Valor),
			}

			reglas = reglas + "sumar_ibc(gastosRep,"+strconv.Itoa(int(Valor))+")."
			codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string

			}
			lista_descuentos = append(lista_descuentos, temp_conceptos)


		}

		valor_prima_antiguedad := m.ProveAll("prima_ant(" + asignacion_basica_string + "," + dias_a_liquidar + "," + tipoPreliquidacion_string + ","+periodo+"," + dias_laborados_string + ",V).")
		for _, solution := range valor_prima_antiguedad {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
			temp_conceptos := models.ConceptosResumen{Nombre: "primaAnt",
				Valor: fmt.Sprintf("%.0f", Valor),
			}

			reglas = reglas + "sumar_ibc(primaAnt,"+strconv.Itoa(int(Valor))+")."
			codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string

			}

			lista_descuentos = append(lista_descuentos, temp_conceptos)


		}

		if esAnual == 1 {
			valor_bonificacion_servicios := m.ProveAll("bon_ser(" + asignacion_basica_string + "," + dias_a_liquidar + "," + tipoPreliquidacion_string + ","+periodo+"," + dias_laborados_string + "," + id_cargo_string + ",V).")
			for _, solution := range valor_bonificacion_servicios {
				Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
				temp_conceptos := models.ConceptosResumen{Nombre: "bonServ",
					Valor: fmt.Sprintf("%.0f", Valor),
				}

				reglas = reglas + "sumar_ibc(bonServ,"+strconv.Itoa(int(Valor))+")."
				codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

				for _, cod := range codigo {
					temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
					temp_conceptos.DiasLiquidados = dias_a_liquidar
					temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
				}

				lista_descuentos = append(lista_descuentos, temp_conceptos)


			}
		}

		if porcentajePT != 0 {

			valor_prima_tecnica := m.ProveAll("prima_tecnica(" + asignacion_basica_string + "," + dias_a_liquidar + "," + tipoPreliquidacion_string + "," + porcentaje_PT_string + ",V).")
			for _, solution := range valor_prima_tecnica {
				Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
				temp_conceptos := models.ConceptosResumen{Nombre: "priTec",
					Valor: fmt.Sprintf("%.0f", Valor),
				}
				reglas = reglas + "sumar_ibc(priTec,"+strconv.Itoa(int(Valor))+")."
				codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

				for _, cod := range codigo {
					temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
					temp_conceptos.DiasLiquidados = dias_a_liquidar
				  temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
				}

				lista_descuentos = append(lista_descuentos, temp_conceptos)


			}
		}

		valor_prima_secretarial := m.ProveAll("prima_secretarial(" + asignacion_basica_string + ","+periodo+"," + id_cargo_string + "," + tipoPreliquidacion_string + "," + dias_laborados_string + ",V).")
		for _, solution := range valor_prima_secretarial {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
			temp_conceptos := models.ConceptosResumen{Nombre: "primaSecr",
				Valor: fmt.Sprintf("%.0f", Valor),
			}

			reglas = reglas + "sumar_ibc(primaSecr,"+strconv.Itoa(int(Valor))+")."
			codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string

			}
			lista_descuentos = append(lista_descuentos, temp_conceptos)


		}

		//Previo a pagos de salud y pensión se calcula el IBC
		CalcularIBC(reglas)
		ManejarNovedadesDevengosFP(reglas, tipoPreliquidacion_string)
		total_devengado_string := strconv.Itoa(int(ibc))

		valor_salud := m.ProveAll("salud_fun(" + total_devengado_string + ","+periodo+","+tipoPreliquidacion_string+",V).")
		for _, solution := range valor_salud {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
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

		valor_pension := m.ProveAll("pension_fun(" + total_devengado_string + ","+periodo+","+tipoPreliquidacion_string+",V).")
		for _, solution := range valor_pension {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
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

		if tipoPreliquidacion_string == "1" || tipoPreliquidacion_string == "2" {

			valor_fondo_solidaridad := m.ProveAll("fondo_solidaridad_fun(" + total_devengado_string + ","+periodo+",V).")
			for _, solution := range valor_fondo_solidaridad {
				Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
				temp_conceptos := models.ConceptosResumen{Nombre: "fondoSolidaridad",
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
		}



		return lista_descuentos, ibc
}

//Función que guarda los conceptos que fueron calculados en la anterior
func GuardarConceptos (lista_descuentos []models.ConceptosResumen)(rest []models.Respuesta){
		temp := models.Respuesta{}
		var resultado []models.Respuesta

		temp_conceptos := models.ConceptosResumen{Nombre: "ibc_liquidado",
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

		temp.Conceptos = &lista_descuentos
		resultado = append(resultado, temp)
		total_devengado_novedad = 0
		total_devengado_no_novedad = 0
		return resultado
}

//Función que calcula IBC, basado en hechos de golog
func CalcularIBC(reglas string){

	e := NewMachine().Consult(reglas)

	valor_ibc := e.ProveAll("calcular_ibc(V).")
	for _, solution := range valor_ibc {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
		ibc = ibc + Valor

		}
}

//Función que gestiona las novedades de la persona
func ManejarNovedades(reglas string, idProveedor int, tipoPreliquidacion, periodo string) (rest []models.ConceptosResumen){
	fmt.Println("novedades")
	var lista_novedades []models.ConceptosResumen

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
			temp_conceptos.DiasLiquidados = dias_a_liquidar
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion
		}

		lista_novedades = append(lista_novedades, temp_conceptos)

	}

	return lista_novedades

}

func ManejarNovedadesDevengosFP(reglas string, tipoPreliquidacion string){

	f := NewMachine().Consult(reglas)
	novedades_devengo := f.ProveAll("novedades_devengos(X).")
		for _, solution := range novedades_devengo {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("X")), 64)
			ibc = ibc + Valor
		}
}
//Función que calcula doceavas de bonificacion por servicios, doceava de prima semestral, doceava de prima vacaciones y adhiere resutltado a lo anterior
func CalcularDoceavaBonServDic(reglas string,tipoPreliquidacion_string string, numero_contrato string, vigencia_contrato int, periodo string) (rest []models.ConceptosResumen){

	var lista_doceavas []models.ConceptosResumen
	var total_sumado float64
	f := NewMachine().Consult(reglas)

 	consultar_valores_bonificacion := f.ProveAll("concepto_bon_serv_dic(X).")
	 for _, solution := range consultar_valores_bonificacion {
		codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
		total_sumado = total_sumado + ConsultarValoresBonServDic(numero_contrato, vigencia_contrato, codigo_concepto, periodo)

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


func CalcularDoceavaBonServPS(reglas string,tipoPreliquidacion_string string, numero_contrato string, vigencia_contrato int, periodo string, mesPreliq, anoPreliq int) (rest []models.ConceptosResumen){

	var lista_doceavas []models.ConceptosResumen
	var total_sumado float64

		tipoPreliquidacion_string = "3"
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
				temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			}

			lista_doceavas = append(lista_doceavas, temp_conceptos)
		}
		return lista_doceavas

}

func CalcularDoceavaPSDic(reglas string,tipoPreliquidacion_string string, numero_contrato string, vigencia_contrato int, periodo string) (rest []models.ConceptosResumen){

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
			temp_conceptos.DiasLiquidados = dias_a_liquidar
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
		}

		lista_doceavas = append(lista_doceavas, temp_conceptos)

	}

	return lista_doceavas


}

func CalcularDoceavaPV(reglas string, tipoPreliquidacion_string string, total_calculado []models.ConceptosResumen) (rest []models.ConceptosResumen){

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
