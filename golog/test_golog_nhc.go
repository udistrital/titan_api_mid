package golog

import (
	"fmt"
	"strconv"
	"strings"
	. "github.com/udistrital/golog"
	models "github.com/udistrital/titan_api_mid/models"
)

func CargarReglasHCS(idProveedor int, reglas string, preliquidacion models.Preliquidacion, periodo string, objeto_datos_acta models.ObjetoActaInicio) (rest []models.Respuesta) {

	var resultado []models.Respuesta
	var listaDescuentos []models.ConceptosResumen
	var listaNovedades []models.ConceptosResumen
	var listaRetefuente []models.ConceptosResumen
	var tipoPreliquidacionString = "2"
	var diasNovedadString = "0"
	var ano, _ = strconv.Atoi(periodo)
	reglas = reglas + "cargo(0)."
	reglas = reglas + "periodo(" + periodo + ")."

	//llamar funcion que calculaDias
	diasALiquidar, meses := CalcularPeriodoLiquidacion(preliquidacion, objeto_datos_acta)
	fmt.Println("dias liquidados", diasALiquidar)
	fmt.Println("meses", meses)

	reglas = reglas + "duracion_contrato(" + strconv.Itoa(idProveedor) + "," + meses + "," + periodo + ")."

	m := NewMachine().Consult(reglas)

	novedades_seg_social := m.ProveAll("seg_social(N,A,M,D,AA,MM,DD).")

	for _, solution := range novedades_seg_social {

		fmt.Println("existe novedad de SS")

		novedad := fmt.Sprintf("%s", solution.ByName_("N"))
		AnoDesde, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("A")), 64)
		MesDesde, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("M")), 64)
		DiaDesde, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("D")), 64)
		AnoHasta, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("AA")), 64)
		MesHasta, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("MM")), 64)
		DiaHasta, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("DD")), 64)

		afectacion_seg_social := m.ProveAll("afectacion_seguridad(" + novedad + ").")
		for _, _ = range afectacion_seg_social {

			dias_novedad := CalcularDiasNovedades(preliquidacion.Mes, ano, AnoDesde, MesDesde, DiaDesde, AnoHasta, MesHasta, DiaHasta)
			diasALiquidar = strconv.Itoa(int(30 - dias_novedad))
			diasNovedadString = strconv.Itoa(int(dias_novedad))
			_, total_devengado_novedad = CalcularConceptosHCS(idProveedor, periodo, reglas, tipoPreliquidacionString, diasNovedadString)
			fmt.Println("- dias_a_liq", diasALiquidar)
			fmt.Println("- diasNovedadString", dias_novedad)
			fmt.Println("- total novedad", total_devengado_novedad)
			ibc = 0
		}

	}

	fmt.Println("dias_a_liq", diasALiquidar)
	fmt.Println("diasNovedadString", diasNovedadString)
	fmt.Println("total novedad", total_devengado_novedad)

	listaDescuentos, total_devengado_no_novedad = CalcularConceptosHCS(idProveedor, periodo, reglas, tipoPreliquidacionString, diasALiquidar)
	listaNovedades = ManejarNovedadesHCS(reglas, idProveedor, tipoPreliquidacionString, periodo, diasALiquidar)
	//listaRetefuente = CalcularReteFuenteSal(tipoPreliquidacionString, reglas, listaDescuentos, diasALiquidar)
	total_calculos = append(total_calculos, listaDescuentos...)
	total_calculos = append(total_calculos, listaNovedades...)
	total_calculos = append(total_calculos, listaRetefuente...)
	resultado = GuardarConceptosHCS(reglas, total_calculos, diasALiquidar, diasNovedadString)

	total_calculos = []models.ConceptosResumen{}
	ibc = 0
	return resultado
}

func CalcularRetefuenteHCS(reglas string, listaConceptos []models.ConceptosResumen, datos models.DatosPreliquidacion) (rest []models.ConceptosResumen) {

	reglas = reglas + "periodo(" + strconv.Itoa(datos.Preliquidacion.Ano) + ")."
	reglas = reglas + "intereses_vivienda(0)."

	return CalcularReteFuenteSal("2", reglas, listaConceptos, datos.DiasALiquidar)

}

func CalcularConceptosHCS(idProveedor int, periodo, reglas, tipoPreliquidacionString, dias_liq string) (rest []models.ConceptosResumen, total_dev float64) {

	var listaDescuentos []models.ConceptosResumen
	reglas = reglas + "dias_liquidados(" + strconv.Itoa(idProveedor) + "," + dias_liq + ")."

	m := NewMachine().Consult(reglas)

	valor_pago := m.ProveAll("valor_pago(X," + periodo + ",T).")
	for _, solution := range valor_pago {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
		Nom_Concepto := "salarioBase"
		temp_conceptos := models.ConceptosResumen{Nombre: "salarioBase",
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		fmt.Println("salario ->", Valor)
		reglas = reglas + "sumar_ibc(" + Nom_Concepto + "," + strconv.Itoa(int(Valor)) + ")."
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C,N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.DiasLiquidados = dias_liq
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
		}
		listaDescuentos = append(listaDescuentos, temp_conceptos)

	}

	descuentos := m.ProveAll("concepto_ley(X,Y," + periodo + ",B,N).")
	for _, solution := range descuentos {

		Base, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("B")), 64)
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("Y")), 64)
		Nom_Concepto := fmt.Sprintf("%s", solution.ByName_("N"))

		temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
			Base:  fmt.Sprintf("%.0f", Base),
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		reglas = reglas + "sumar_ibc(" + Nom_Concepto + "," + strconv.Itoa(int(Valor)) + ")."
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C,N,D).`)

		for _, cod := range codigo {

			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.DiasLiquidados = dias_liq
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
		}

		listaDescuentos = append(listaDescuentos, temp_conceptos)
	}

	CalcularIBC(strconv.Itoa(idProveedor), reglas)

	return listaDescuentos, ibc

}

func GuardarConceptosHCS(reglas string, listaDescuentos []models.ConceptosResumen, dias_a_liq_no_nov, dias_a_liq_nov string) (rest []models.Respuesta) {
	temp := models.Respuesta{}
	var resultado []models.Respuesta
	m := NewMachine().Consult(reglas)

	temp_conceptos := models.ConceptosResumen{Nombre: "ibc_liquidado",
		Valor: fmt.Sprintf("%.0f", total_devengado_no_novedad),
	}

	codigo := m.ProveAll(`codigo_concepto(ibc_liquidado,C,N,D).`)

	for _, cod := range codigo {
		temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
		temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
		temp_conceptos.DiasLiquidados = dias_a_liq_no_nov
	}

	listaDescuentos = append(listaDescuentos, temp_conceptos)

	temp_conceptos_1 := models.ConceptosResumen{Nombre: "ibc_novedad",
		Valor: fmt.Sprintf("%.0f", total_devengado_novedad),
	}

	codigo_1 := m.ProveAll(`codigo_concepto(ibc_novedad,C, N,D).`)

	for _, cod := range codigo_1 {
		temp_conceptos_1.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
		temp_conceptos_1.DiasLiquidados = dias_a_liq_nov //DIAS NOVEDAD
		temp_conceptos_1.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		temp_conceptos_1.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))

	}

	listaDescuentos = append(listaDescuentos, temp_conceptos_1)

	temp.Conceptos = &listaDescuentos
	resultado = append(resultado, temp)
	total_devengado_novedad = 0
	total_devengado_no_novedad = 0
	return resultado
}

func ManejarNovedadesHCS(reglas string, idProveedor int, tipoPreliquidacion, periodo, dias_a_liq string) (rest []models.ConceptosResumen) {

	var listaNovedades []models.ConceptosResumen
	reglas = reglas + "dias_liquidados(" + strconv.Itoa(idProveedor) + "," + dias_a_liq + ")."

	f := NewMachine().Consult(reglas)

	idProveedorString := strconv.Itoa(idProveedor)
	novedades := f.ProveAll("info_concepto(" + idProveedorString + ",T," + periodo + ",N,R).")
	for _, solution := range novedades {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C,N,D).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.DiasLiquidados = dias_a_liq
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		listaNovedades = append(listaNovedades, temp_conceptos)

	}

	return listaNovedades

}

func CalcularTotalesContratoHCS(NumDocumento, MesesContrato, VigenciaContrato, Categoria,TotalHoras ,reglas string) (rest []models.ConceptosResumen) {

	var listaDescuentos []models.ConceptosResumen

	cat := strings.ToLower(Categoria) 
	reglas = reglas + "valor_contrato(" +NumDocumento + ",3677929)."
	reglas = reglas + "fin_contrato(" + NumDocumento + ",si)."

	reglas = reglas + "duracion_contrato(12,0,V)."

	reglas = reglas + "pensionado(no)."
	reglas = reglas + "valor_punto(2019,14210)."
	reglas = reglas + "factor(auxiliar,1.8,2019)."
	reglas = reglas + "factor(asistente,2.3,2019)."
	reglas = reglas + "factor(asociado,2.7,2019)."
	reglas = reglas + "factor(titular,3,2019)."
	reglas = reglas + "porcentaje_devengo_v2(2019,prima_servicios,0.0740740929887414)."
	reglas = reglas + "porcentaje_devengo_v2(2019,prima_vacaciones,0.0414814800000000)."
	reglas = reglas + "porcentaje_devengo_v2(2019,intereses_cesantias,0.0100000000000000)."
	reglas = reglas + "porcentaje_devengo_v2(2019,cesantias,0.0833333333333333)."
	reglas = reglas + "porcentaje_devengo_v2(2019,vacaciones,0.0414814800000000)."
	reglas = reglas + "porcentaje_devengo_v2(2019,prima_navidad,0.0740740929887414)."
	reglas = reglas + "valor_pago_total_v2(X,C,V,NHT,T):- valor_punto(V,VP),factor(C,F,P), VH is VP * F, T is VH * NHT."
	reglas = reglas + "conceptos_total_contrato_v2(X,C,V,NHT,MC,T,prima_servicios):- (MC @> 6 -> valor_pago_total_v2(X,C,V,NHT,VP), porcentaje_devengo_v2(V,prima_servicios,PPV), R is VP *PPV, T is (R rnd 0); T is 0).	"
	reglas = reglas + "conceptos_total_contrato_v2(X,C,V,NHT,MC,T,prima_vacaciones):- valor_pago_total_v2(X,C,V,NHT,VP), porcentaje_devengo_v2(V,prima_vacaciones,PPV), R is VP *PPV, T is (R rnd 0)."
	reglas = reglas + "conceptos_total_contrato_v2(X,C,V,NHT,MC,T,intereses_cesantias):- valor_pago_total_v2(X,C,V,NHT,VP), porcentaje_devengo_v2(V,intereses_cesantias,PPV), R is VP *PPV, T is (R rnd 0)."
	reglas = reglas + "conceptos_total_contrato_v2(X,C,V,NHT,MC,T,cesantias):- valor_pago_total_v2(X,C,V,NHT,VP), porcentaje_devengo_v2(V,cesantias,PPV), R is VP *PPV, T is (R rnd 0)."
	reglas = reglas + "conceptos_total_contrato_v2(X,C,V,NHT,MC,T,vacaciones):- valor_pago_total_v2(X,C,V,NHT,VP), porcentaje_devengo_v2(V,vacaciones,PPV), R is VP *PPV, T is (R rnd 0)."
	reglas = reglas + "conceptos_total_contrato_v2(X,C,V,NHT,MC,T,prima_navidad):- valor_pago_total_v2(X,C,V,NHT,VP), porcentaje_devengo_v2(V,prima_navidad,PPV), R is VP *PPV, T is (R rnd 0)."
	m := NewMachine().Consult(reglas)

	valor_pago_1 := m.ProveAll("valor_pago_total(" +NumDocumento + "," +VigenciaContrato + ",T).")
	for _, solution := range valor_pago_1 {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "salarioBase1",
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C,N,D).`)

		for _, cod := range codigo {
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))

		}
		listaDescuentos = append(listaDescuentos, temp_conceptos)

	}

	valor_pago := m.ProveAll("valor_pago_total_v2(" + NumDocumento + "," + cat+ "," + VigenciaContrato+","+TotalHoras+",T).")
	for _, solution := range valor_pago {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "salarioBase",
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C,N,D).`)

		for _, cod := range codigo {
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))

		}
		listaDescuentos = append(listaDescuentos, temp_conceptos)

	}
	
	descuentos := m.ProveAll("conceptos_total_contrato_v2("+NumDocumento + "," +cat+ "," + VigenciaContrato+","+TotalHoras+","+MesesContrato+",T,N).")
	for _, solution := range descuentos {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)

		temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C,N,D).`)

		for _, cod := range codigo {

			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))

		}

		listaDescuentos = append(listaDescuentos, temp_conceptos)
	}
	
	return listaDescuentos

}


func CalcularDescuentosTotalesHCS(IdPersona, valor_total string, idProveedor int, reglas string, preliquidacion models.Preliquidacion, periodo string) (rest []models.ConceptosResumen) {

	var listaDescuentos []models.ConceptosResumen

	m := NewMachine().Consult(reglas)

	fondo_sol := m.ProveAll("calcular_fondo_sol(X," + valor_total + "," + periodo + ",V).")
	for _, solution := range fondo_sol {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "fondoSolidaridad",
			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C,N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.IdPersona, _ = strconv.Atoi(IdPersona)
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.DiasLiquidados = "0"
			temp_conceptos.TipoPreliquidacion = "2"
		}
		listaDescuentos = append(listaDescuentos, temp_conceptos)

	}

	fondo_sub := m.ProveAll("calcular_fondo_sub(X," + valor_total + "," + periodo + ",V).")
	for _, solution := range fondo_sub {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "fondoSubsistencia",
			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C,N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.IdPersona, _ = strconv.Atoi(IdPersona)
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.DiasLiquidados = "0"
			temp_conceptos.TipoPreliquidacion = "2"
		}
		listaDescuentos = append(listaDescuentos, temp_conceptos)

	}

	return listaDescuentos
}
