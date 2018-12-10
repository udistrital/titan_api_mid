package golog

import (
	"fmt"
	"strconv"

	. "github.com/udistrital/golog"
	models "github.com/udistrital/titan_api_mid/models"
)



func CargarReglasCT(idProveedor int, reglas string,preliquidacion models.Preliquidacion, periodo string, objeto_datos_acta models.ObjetoActaInicio) (rest []models.Respuesta) {

	var resultado []models.Respuesta
	var lista_descuentos []models.ConceptosResumen
	var lista_novedades []models.ConceptosResumen
	var lista_retefuente []models.ConceptosResumen
	var tipoPreliquidacion_string = "2";
	var dias_novedad_string = "0"
	var ano,_ =  strconv.Atoi(periodo)

	reglas = reglas + "cargo(0)."
	reglas = reglas + "periodo("+periodo+")."

	dias_a_liquidar, _:= CalcularPeriodoLiquidacion(preliquidacion,objeto_datos_acta)
	fmt.Println("dias liquidados", dias_a_liquidar)


	m := NewMachine().Consult(reglas)
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
				dias_novedad := CalcularDiasNovedades(preliquidacion.Mes,ano,AnoDesde, MesDesde, DiaDesde, AnoHasta, MesHasta, DiaHasta)
				dias_a_liquidar = strconv.Itoa(int(30 - dias_novedad))
				dias_novedad_string = strconv.Itoa(int(dias_novedad))
      	_,total_devengado_novedad =  CalcularConceptosCT(idProveedor,periodo,reglas, tipoPreliquidacion_string, dias_novedad_string)
				ibc = 0;
		}



		}

	lista_descuentos,total_devengado_no_novedad = CalcularConceptosCT(idProveedor,periodo,reglas, tipoPreliquidacion_string, dias_a_liquidar)
	lista_novedades = ManejarNovedadesCT(reglas,idProveedor, tipoPreliquidacion_string,periodo)
	lista_retefuente = CalcularReteFuenteSal(tipoPreliquidacion_string,reglas, lista_descuentos,dias_a_liquidar);
	total_calculos = append(total_calculos, lista_descuentos...)
	total_calculos = append(total_calculos, lista_novedades...)
	total_calculos = append(total_calculos, lista_retefuente...)
	resultado = GuardarConceptosCT(reglas,total_calculos,dias_a_liquidar, dias_novedad_string)

	total_calculos = []models.ConceptosResumen{}
	ibc = 0;

	return resultado;

}

func CalcularConceptosCT (idProveedor int, periodo,reglas, tipoPreliquidacion_string, dias_liq string)(rest []models.ConceptosResumen, total_dev float64){

	var lista_descuentos []models.ConceptosResumen

	var salarioBase float64
	var salarioBase_string string

	reglas = reglas + "dias_liquidados("+strconv.Itoa(idProveedor)+","+dias_liq+")."

	m := NewMachine().Consult(reglas)

	valor_pago := m.ProveAll("valor_pago(X,"+periodo+",P).")

	for _, solution := range valor_pago {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("P")), 64)
		Nom_Concepto := "salarioBase"
		salarioBase = Valor
		salarioBase_string = strconv.Itoa(int(salarioBase))
		temp_conceptos := models.ConceptosResumen{Nombre: "salarioBase",
			Valor: fmt.Sprintf("%.0f", Valor),
		}

		reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			temp_conceptos.DiasLiquidados = dias_liq
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)
	}

	fmt.Println("salariobase", salarioBase_string)
	descuento_reteica := m.ProveAll("calcular_reteica("+salarioBase_string+", "+periodo+", R).")
	for _, solution := range descuento_reteica {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		Nom_Concepto := "reteIca"
		temp_conceptos := models.ConceptosResumen{Nombre: "reteIca",
			Valor: fmt.Sprintf("%.0f", Valor),
		}

		reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C,N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			temp_conceptos.DiasLiquidados = dias_liq
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)

	}

	/*Estampila UD - Se elimina porque a partir de 2018

	descuento_estampilla := m.ProveAll("calcular_estampilla("+salarioBase_string+", "+periodo+", R).")

	for _, solution := range descuento_estampilla {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		Nom_Concepto := "estampillaUD"
		temp_conceptos := models.ConceptosResumen{Nombre: "estampillaUD",

			Valor: fmt.Sprintf("%.0f", Valor),
		}

		reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			temp_conceptos.DiasLiquidados = dias_liq
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)


	}
	*/
	//proCultura

	descuento_procultura := m.ProveAll("calcular_procultura("+salarioBase_string+", "+periodo+", R).")

	for _, solution := range descuento_procultura {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		Nom_Concepto := "proCultura"
		temp_conceptos := models.ConceptosResumen{Nombre: "proCultura",

			Valor: fmt.Sprintf("%.0f", Valor),
		}

		reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			temp_conceptos.DiasLiquidados = dias_liq
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)


	}

	//proCultura

	descuento_adulto_mayor := m.ProveAll("calcular_adulto_mayor("+salarioBase_string+", "+periodo+", R).")

	for _, solution := range descuento_adulto_mayor {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		Nom_Concepto := "adultoMayor"
		temp_conceptos := models.ConceptosResumen{Nombre: "adultoMayor",

			Valor: fmt.Sprintf("%.0f", Valor),
		}

		reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
		  temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			temp_conceptos.DiasLiquidados = dias_liq
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)


	}

	descuento_salud:= m.ProveAll("calcular_salud(si,"+salarioBase_string+", "+periodo+", R).")

	for _, solution := range descuento_salud {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		Nom_Concepto := "salud"
		temp_conceptos := models.ConceptosResumen{Nombre: "salud",

			Valor: fmt.Sprintf("%.0f", Valor),
		}

		reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			temp_conceptos.DiasLiquidados = dias_liq
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)


	}

	descuento_pension:= m.ProveAll("calcular_pension(si,"+salarioBase_string+", "+periodo+", R).")

	for _, solution := range descuento_pension {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		Nom_Concepto := "pension"
		temp_conceptos := models.ConceptosResumen{Nombre: "pension",

			Valor: fmt.Sprintf("%.0f", Valor),
		}

		reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			temp_conceptos.DiasLiquidados = dias_liq
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)


	}

	CalcularIBC(reglas)
	return lista_descuentos,ibc

}

func GuardarConceptosCT (reglas string,lista_descuentos []models.ConceptosResumen, dias_a_liq_no_nov, dias_a_liq_nov string)(rest []models.Respuesta){
		temp := models.Respuesta{}
		var resultado []models.Respuesta
		m := NewMachine().Consult(reglas)

		temp_conceptos := models.ConceptosResumen{Nombre: "ibc_liquidado",
			Valor: fmt.Sprintf("%.0f", total_devengado_no_novedad),
		}

		codigo := m.ProveAll(`codigo_concepto(ibc_liquidado,C, N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
      temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.DiasLiquidados = dias_a_liq_no_nov
		}


		lista_descuentos = append(lista_descuentos, temp_conceptos)

		temp_conceptos_1 := models.ConceptosResumen{Nombre: "ibc_novedad",
			Valor: fmt.Sprintf("%.0f", total_devengado_novedad),
		}

		codigo_1 := m.ProveAll(`codigo_concepto(ibc_novedad,C, N,D).`)

		for _, cod := range codigo_1 {
			temp_conceptos_1.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos_1.DiasLiquidados = dias_a_liq_nov
			temp_conceptos_1.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos_1.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))

		}


		lista_descuentos = append(lista_descuentos, temp_conceptos_1)

		temp.Conceptos = &lista_descuentos
		resultado = append(resultado, temp)

		total_devengado_novedad = 0
		total_devengado_no_novedad = 0
		return resultado
}


func ManejarNovedadesCT(reglas string, idProveedor int, tipoPreliquidacion, periodo string) (rest []models.ConceptosResumen){

	var lista_novedades []models.ConceptosResumen

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
			temp_conceptos.DiasLiquidados = dias_a_liquidar
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		lista_novedades = append(lista_novedades, temp_conceptos)

	}

	return lista_novedades

}
