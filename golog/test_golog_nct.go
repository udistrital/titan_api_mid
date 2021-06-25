package golog

import (
	"fmt"
	"strconv"

	. "github.com/udistrital/golog"
	models "github.com/udistrital/titan_api_mid/models"
)

func CargarReglasCT(idProveedor int, reglas string, preliquidacion models.Preliquidacion, periodo string, objeto_datos_acta models.ObjetoActaInicio, pensionado bool, dependientes bool) (rest []models.Respuesta) {

	var resultado []models.Respuesta
	var listaDescuentos []models.ConceptosResumen
	var listaNovedades []models.ConceptosResumen
	var listaRetefuente []models.ConceptosResumen
	var tipoPreliquidacionString = "2"
	var diasNovedadString = "0"
	var ano, _ = strconv.Atoi(periodo)

	reglas = reglas + "cargo(0)."
	reglas = reglas + "periodo(" + periodo + ")."

	//debe verificar si tiene una novedad tipo cesion y calcular el periodo de acuerdo a ello
	diasALiquidar, _ := CalcularPeriodoLiquidacion(preliquidacion, objeto_datos_acta)
	fmt.Println("dias liquidados", diasALiquidar)

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
		fmt.Println("existe novedad de SS", novedad)
		fmt.Println("AnoDesde", AnoDesde)
		fmt.Println("MesDesde", MesDesde)
		fmt.Println("DiaDesde", DiaDesde)
		fmt.Println("AnoHasta", AnoHasta)
		fmt.Println("MesHasta", MesHasta)
		fmt.Println("DiaHasta", DiaHasta)

		afectacion_seg_social := m.ProveAll("afectacion_seguridad(" + novedad + ").")
		for _, solution := range afectacion_seg_social {

			fmt.Println(solution)
			dias_novedad := CalcularDiasNovedades(preliquidacion.Mes, ano, AnoDesde, MesDesde, DiaDesde, AnoHasta, MesHasta, DiaHasta)
			diasALiquidar = strconv.Itoa(int(30 - dias_novedad))
			diasNovedadString = strconv.Itoa(int(dias_novedad))
			_, total_devengado_novedad = CalcularConceptosCT(idProveedor, periodo, reglas, tipoPreliquidacionString, diasNovedadString)
			ibc = 0
		}

	}

	listaDescuentos, total_devengado_no_novedad = CalcularConceptosCT(idProveedor, periodo, reglas, tipoPreliquidacionString, diasALiquidar, pensionado)
	listaNovedades = ManejarNovedadesCT(reglas, idProveedor, tipoPreliquidacionString, periodo, diasALiquidar)
	listaRetefuente = CalcularReteFuenteSal(tipoPreliquidacionString, reglas, listaDescuentos, diasALiquidar, dependientes)
	total_calculos = append(total_calculos, listaDescuentos...)
	total_calculos = append(total_calculos, listaNovedades...)
	total_calculos = append(total_calculos, listaRetefuente...)
	resultado = GuardarConceptosCT(reglas, total_calculos, diasALiquidar, diasNovedadString)

	total_calculos = []models.ConceptosResumen{}
	ibc = 0

	return resultado

}

func CalcularConceptosCT(idProveedor int, periodo, reglas, tipoPreliquidacionString, dias_liq string) (rest []models.ConceptosResumen, total_dev float64) {

	var listaDescuentos []models.ConceptosResumen

	reglas = reglas + "dias_liquidados(" + strconv.Itoa(idProveedor) + "," + dias_liq + ")."

	m := NewMachine().Consult(reglas)

	total := m.ProveAll("liquidar_ct(" + strconv.Itoa(idProveedor) + "," + periodo + ", N,T).")

	for _, solution := range total {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
		Nom_Concepto := fmt.Sprintf("%s", solution.ByName_("N"))
		temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
			Valor: fmt.Sprintf("%.0f", Valor),
		}

		reglas = reglas + "sumar_ibc(" + Nom_Concepto + "," + strconv.Itoa(int(Valor)) + ")."
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N,D).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			temp_conceptos.DiasLiquidados = dias_liq
		}

		listaDescuentos = append(listaDescuentos, temp_conceptos)
	}

	CalcularIBC(strconv.Itoa(idProveedor), reglas)
	return listaDescuentos, ibc

}

func GuardarConceptosCT(reglas string, listaDescuentos []models.ConceptosResumen, dias_a_liq_no_nov, dias_a_liq_nov string) (rest []models.Respuesta) {
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

	listaDescuentos = append(listaDescuentos, temp_conceptos)

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

	listaDescuentos = append(listaDescuentos, temp_conceptos_1)

	temp.Conceptos = &listaDescuentos
	resultado = append(resultado, temp)

	total_devengado_novedad = 0
	total_devengado_no_novedad = 0
	return resultado
}

func ManejarNovedadesCT(reglas string, idProveedor int, tipoPreliquidacion, periodo, dias_a_liq string) (rest []models.ConceptosResumen) {

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
