package golog

import (
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"
	. "github.com/udistrital/golog"

)


var total_devengado_no_novedad float64
var total_devengado_no_novedad_semestral float64
var total_devengado_novedad float64
var diasALiquidar string
var diasNovedadString string
var nombre_archivo string
var ibc float64
var dias_liquidar_prima_semestral  string
var total_calculos []models.ConceptosResumen
var ingresos float64

func CargarReglasFP(dias_a_liq string, MesPreliquidacion int, AnoPreliquidacion int, reglas string, idProveedor int, numero_contrato string, vigencia_contrato int, informacion_cargo []models.FuncionarioCargo, dias_laborados float64, porcentajePT int, tipoPreliquidacion int) (rest []models.Respuesta) {

	//--- Creación de variables
	var resultado []models.Respuesta
	var listaDescuentos []models.ConceptosResumen
	var listaRetefuente []models.ConceptosResumen
	var listaDescuentos_semestral []models.ConceptosResumen
	var listaNovedades []models.ConceptosResumen
	var tipoPreliquidacionString string

	//Conversión de variables
	asignacion_basica_string := strconv.Itoa(informacion_cargo[0].Asignacion_basica)
	id_cargo_string := strconv.Itoa(informacion_cargo[0].Id)
	dias_laborados_string := strconv.Itoa(int(dias_laborados))
	tipoPreliquidacionString = strconv.Itoa(tipoPreliquidacion)
	periodo:= strconv.Itoa(AnoPreliquidacion)
	porcentaje_PT_string := strconv.Itoa(porcentajePT)
	//Asignación de número de días segun tipo de nomina (0 es quince días, 1 es el mes completo)
	if tipoPreliquidacionString  == "0" || tipoPreliquidacionString  == "1" {
		diasALiquidar = "15"

	}

	if tipoPreliquidacionString  == "2" {
		diasALiquidar = "30"

	}	else {
		diasALiquidar = dias_a_liq
	}

	nombre_archivo = "reglas" + strconv.Itoa(idProveedor) + ".txt"
	reglas = reglas + "salario_base(" + asignacion_basica_string + ")."
	reglas = reglas + "tipo_nomina(" + tipoPreliquidacionString + ")."
	reglas = reglas + "porcentaje_pt("+porcentaje_PT_string+")."
	reglas = reglas + "cargo("+id_cargo_string+")."


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
				diasALiquidar = strconv.Itoa(int(30 - dias_novedad))
				diasNovedadString = strconv.Itoa(int(dias_novedad))
				_,total_devengado_novedad = CalcularConceptos(m, reglas, diasNovedadString,asignacion_basica_string,id_cargo_string,dias_laborados_string, tipoPreliquidacionString, porcentajePT, idProveedor,periodo)
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

			doceava_BSPS := CalcularDoceavaBonServPS(reglas,tipoPreliquidacionString, numero_contrato, vigencia_contrato, periodo, MesPreliquidacion , AnoPreliquidacion)
			listaDescuentos_semestral,total_devengado_no_novedad_semestral = CalcularConceptos(m, reglas,dias_liquidar_prima_semestral,asignacion_basica_string,id_cargo_string,dias_laborados_string, "3", porcentajePT, idProveedor,periodo)
			total_calculos = append (total_calculos, listaDescuentos_semestral...)
			total_calculos = append (total_calculos, 	doceava_BSPS...)
			ibc = 0
			}
		//--------------------------------

		//-----NOMINA DE DICIEMBRE ------
		if(MesPreliquidacion  == 12){
			fmt.Println("aqui dic")
			dias_liq_dic := m.ProveAll("dias_liq_dic(FP,TLIQ,D).")
			diasALiquidar = "9"
			for _, solution := range dias_liq_dic{
					tipoLiq := fmt.Sprintf("%s", solution.ByName_("TLIQ"))
					dias_liquidacion_diciembre := fmt.Sprintf("%s", solution.ByName_("D"))
					listaDescuentos_semestral,total_devengado_no_novedad_semestral = CalcularConceptos(m, reglas,dias_liquidacion_diciembre,asignacion_basica_string,id_cargo_string,dias_laborados_string, tipoLiq, porcentajePT, idProveedor,periodo)
					total_calculos = append (total_calculos, listaDescuentos_semestral...)
					doceavas_bsd := CalcularDoceavaBonServDic(reglas,tipoLiq, numero_contrato, vigencia_contrato, periodo)
					doceavas_psd := CalcularDoceavaPSDic(reglas,tipoLiq, numero_contrato, vigencia_contrato, periodo)
					total_calculos = append (total_calculos, doceavas_bsd...)
					total_calculos = append (total_calculos, doceavas_psd...)

					if(tipoLiq == "6"){
						doceavas_pv := CalcularDoceavaPV(reglas, tipoPreliquidacionString, total_calculos)
						total_calculos = append (total_calculos, doceavas_pv...)
					}

					ibc = 0
			}



		}
 		// ----------------------------------

		// ----- Nomina ordinaria ----- Proceso de cálculo, manejo de novedades y guardado de conceptos

		listaDescuentos,total_devengado_no_novedad = CalcularConceptos(m, reglas,diasALiquidar,asignacion_basica_string,id_cargo_string,dias_laborados_string, tipoPreliquidacionString, porcentajePT, idProveedor,periodo)
		ibc = 0
		listaNovedades = ManejarNovedades(reglas,idProveedor, tipoPreliquidacionString, periodo)
		listaRetefuente = CalcularReteFuentePlanta(tipoPreliquidacionString,reglas,periodo, listaDescuentos);
		total_calculos = append(total_calculos, listaDescuentos...)
		total_calculos = append(total_calculos, listaNovedades...)
		total_calculos = append(total_calculos, listaRetefuente...)
		resultado = GuardarConceptos(reglas,total_calculos)
		total_calculos = []models.ConceptosResumen{}
		ingresos = 0

		// ---------------------------
		return resultado;

	}

//Función que, utilizando las reglas, calcula cada uno de los conceptos. Retorna un objeto con los resultados
	func CalcularConceptos(m Machine, reglas,diasALiquidar,asignacion_basica_string,id_cargo_string,dias_laborados_string,tipoPreliquidacionString string, porcentajePT,idProveedor  int,periodo string) (rest []models.ConceptosResumen, total_dev float64){

		var listaDescuentos []models.ConceptosResumen

		valor_devengo := m.ProveAll("liquidacion_planta(CON,"+asignacion_basica_string+","+tipoPreliquidacionString+","+diasALiquidar+","+periodo+","+id_cargo_string+","+dias_laborados_string+",V).")
		for _, solution := range valor_devengo {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
			ingresos = ingresos + Valor
			Nom_Concepto := fmt.Sprintf("%s", solution.ByName_("CON"))
			temp_conceptos := models.ConceptosResumen{Nombre: Nom_Concepto,
			Valor: fmt.Sprintf("%.0f", Valor),
			}

			reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
			codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C, N, D).")

			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				 temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
				temp_conceptos.DiasLiquidados = diasALiquidar
				temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			}

			listaDescuentos = append(listaDescuentos, temp_conceptos)


		}

		//Previo a pagos de salud y pensión se calcula el IBC
		CalcularIBC("whatever",reglas)
		ManejarNovedadesDevengosFP(reglas, tipoPreliquidacionString)
		total_devengado_string := strconv.Itoa(int(ibc))

		valor_descuentos := m.ProveAll("desc_obli_planta(CON, "+total_devengado_string+", "+periodo+","+tipoPreliquidacionString+", V).")
		for _, solution := range valor_descuentos {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
			fmt.Println("valor", Valor)
			Nom_Concepto := fmt.Sprintf("%s", solution.ByName_("CON"))
			fmt.Println(Nom_Concepto)
			temp_conceptos := models.ConceptosResumen{Nombre: Nom_Concepto,
			Valor: fmt.Sprintf("%.0f", Valor),
			}

			reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
			codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C, N,D).")

			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
				temp_conceptos.DiasLiquidados = diasALiquidar
				temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			}

			listaDescuentos = append(listaDescuentos, temp_conceptos)


		}

		return listaDescuentos, ibc
}

//Función que guarda los conceptos que fueron calculados en la anterior
func GuardarConceptos (reglas string,listaDescuentos []models.ConceptosResumen)(rest []models.Respuesta){
		temp := models.Respuesta{}
		var resultado []models.Respuesta

		m := NewMachine().Consult(reglas)

		temp_conceptos := models.ConceptosResumen{Nombre: "ibc_liquidado",
		  	Valor: fmt.Sprintf("%.0f", total_devengado_no_novedad),
		}

		codigo := m.ProveAll(`codigo_concepto(ibc_liquidado,C, N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.DiasLiquidados = diasALiquidar
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		listaDescuentos = append(listaDescuentos, temp_conceptos)

		temp_conceptos_1 := models.ConceptosResumen{Nombre: "ibc_novedad",
			Valor: fmt.Sprintf("%.0f", total_devengado_no_novedad),
		}

		codigo_1 := m.ProveAll(`codigo_concepto(ibc_novedad,C, N,D).`)

		for _, cod := range codigo_1 {
			temp_conceptos_1.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos_1.DiasLiquidados = diasALiquidar
		}


		listaDescuentos = append(listaDescuentos, temp_conceptos_1)

		temp.Conceptos = &listaDescuentos
		resultado = append(resultado, temp)
		total_devengado_novedad = 0
		total_devengado_no_novedad = 0
		return resultado
}



//Función que gestiona las novedades de la persona
func ManejarNovedades(reglas string, idProveedor int, tipoPreliquidacion, periodo string) (rest []models.ConceptosResumen){
	fmt.Println("novedades")
	var listaNovedades []models.ConceptosResumen

	f := NewMachine().Consult(reglas)

	idProveedorString := strconv.Itoa(idProveedor)
	novedades := f.ProveAll("info_concepto(" + idProveedorString + ",T,"+periodo+",N,R).")

	for _, solution := range novedades {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C,N,D).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.DiasLiquidados = diasALiquidar
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		listaNovedades = append(listaNovedades, temp_conceptos)

	}

	return listaNovedades

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
func CalcularDoceavaBonServDic(reglas string,tipoPreliquidacionString string, numero_contrato string, vigencia_contrato int, periodo string) (rest []models.ConceptosResumen){

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

		codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C,N,D).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.DiasLiquidados = diasALiquidar
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		lista_doceavas = append(lista_doceavas, temp_conceptos)

	}

	return lista_doceavas

}


func CalcularDoceavaBonServPS(reglas string,tipoPreliquidacionString string, numero_contrato string, vigencia_contrato int, periodo string, mesPreliq, anoPreliq int) (rest []models.ConceptosResumen){

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

			codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C,N,D).")
			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			  temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
				temp_conceptos.DiasLiquidados = diasALiquidar
				temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			}

			lista_doceavas = append(lista_doceavas, temp_conceptos)
		}
		return lista_doceavas

}

func CalcularDoceavaPSDic(reglas string,tipoPreliquidacionString string, numero_contrato string, vigencia_contrato int, periodo string) (rest []models.ConceptosResumen){

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

		codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C,N,D).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.DiasLiquidados = diasALiquidar
		  temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
		}

		lista_doceavas = append(lista_doceavas, temp_conceptos)

	}

	return lista_doceavas


}

func CalcularDoceavaPV(reglas string, tipoPreliquidacionString string, total_calculado []models.ConceptosResumen) (rest []models.ConceptosResumen){

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

			codigo := e.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C,N,D).")
			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			  temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
				temp_conceptos.DiasLiquidados = diasALiquidar
				temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			}

			lista_doceavas = append(lista_doceavas, temp_conceptos)

}
	return lista_doceavas
}
