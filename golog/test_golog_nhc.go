package golog

import (
	"fmt"
	"strconv"

	. "github.com/udistrital/golog"
	models "github.com/udistrital/titan_api_mid/models"
)

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

func LiquidarMesHCH(reglas string, cedula string, ano int, detallePreliquidacion models.DetallePreliquidacion) (data []models.DetallePreliquidacion) {
	var conceptoNomina models.ConceptoNomina
	var totalDevengado float64
	var totalDescuentos float64
	var totalAPagar float64

	m := NewMachine().Consult(reglas)
	total := m.ProveAll("liquidar_hch(" + cedula + "," + strconv.Itoa(ano) + ",N,T).")
	totalDescuentos = 0
	totalDevengado = 0
	totalAPagar = 0
	for _, solution := range total {

		detallePreliquidacion.ValorCalculado, _ = strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
		conceptoNomina.NombreConcepto = fmt.Sprintf("%s", solution.ByName_("N"))

		codigo := m.ProveAll(`codigo_concepto(` + conceptoNomina.NombreConcepto + `,C,N).`)
		for _, cod := range codigo {
			conceptoNomina.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			conceptoNomina.NaturalezaConceptoNominaId, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		if conceptoNomina.NaturalezaConceptoNominaId == 423 {
			totalDevengado = totalDevengado + detallePreliquidacion.ValorCalculado
		}

		if conceptoNomina.NaturalezaConceptoNominaId == 424 {
			totalDescuentos = totalDescuentos + detallePreliquidacion.ValorCalculado
		}

		detallePreliquidacion.Id = 0
		detallePreliquidacion.ConceptoNominaId = &models.ConceptoNomina{Id: conceptoNomina.Id}
		data = append(data, detallePreliquidacion)
	}

	totalAPagar = totalDevengado - totalDescuentos
	//se agrega el detalle del total a pagar
	detallePreliquidacion.Id = 0
	detallePreliquidacion.ValorCalculado = totalAPagar
	detallePreliquidacion.ConceptoNominaId = &models.ConceptoNomina{Id: 574}
	detallePreliquidacion.Activo = true
	data = append(data, detallePreliquidacion)

	//se agrega el detalle del total de los descuentos
	detallePreliquidacion.Id = 0
	detallePreliquidacion.ValorCalculado = totalDescuentos
	detallePreliquidacion.ConceptoNominaId = &models.ConceptoNomina{Id: 573}
	detallePreliquidacion.Activo = true
	data = append(data, detallePreliquidacion)

	return data
}
