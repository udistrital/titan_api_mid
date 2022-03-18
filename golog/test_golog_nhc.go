package golog

import (
	"fmt"
	"strconv"

	. "github.com/udistrital/golog"
	models "github.com/udistrital/titan_api_mid/models"
)

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
		if !EncontrarConcepto(data, detallePreliquidacion.ConceptoNominaId.Id) {
			data = append(data, detallePreliquidacion)
		}
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

func LiquidarMesHCS(reglas string, cedula string, ano int, detallePreliquidacion models.DetallePreliquidacion, mesFinal bool) (data []models.DetallePreliquidacion) {
	var conceptoNomina models.ConceptoNomina

	m := NewMachine().Consult(reglas)
	total := m.ProveAll("liquidar_hcs(" + cedula + "," + strconv.Itoa(ano) + ",N,T).")
	for _, solution := range total {

		detallePreliquidacion.ValorCalculado, _ = strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
		conceptoNomina.NombreConcepto = fmt.Sprintf("%s", solution.ByName_("N"))

		codigo := m.ProveAll(`codigo_concepto(` + conceptoNomina.NombreConcepto + `,C,N).`)
		for _, cod := range codigo {
			conceptoNomina.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			conceptoNomina.NaturalezaConceptoNominaId, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}
		detallePreliquidacion.Id = 0
		detallePreliquidacion.ConceptoNominaId = &models.ConceptoNomina{Id: conceptoNomina.Id}
		if !EncontrarConcepto(data, detallePreliquidacion.ConceptoNominaId.Id) {
			data = append(data, detallePreliquidacion)
		}
	}

	if mesFinal {
		total := m.ProveAll("liquidar_prestacion(" + cedula + "," + strconv.Itoa(ano) + ",N,T).")
		for _, solution := range total {
			detallePreliquidacion.ValorCalculado, _ = strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
			conceptoNomina.NombreConcepto = fmt.Sprintf("%s", solution.ByName_("N"))

			codigo := m.ProveAll(`codigo_concepto(` + conceptoNomina.NombreConcepto + `,C,N).`)
			for _, cod := range codigo {
				conceptoNomina.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				conceptoNomina.NaturalezaConceptoNominaId, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			}

			detallePreliquidacion.Id = 0
			detallePreliquidacion.ConceptoNominaId = &models.ConceptoNomina{Id: conceptoNomina.Id}

			data = append(data, detallePreliquidacion)
		}
	}
	return data
}

func ReliquidarAportes(reglas string, cedula string, ano int, detallePreliquidacion models.DetallePreliquidacion) (data []models.DetallePreliquidacion) {
	var conceptoNomina models.ConceptoNomina

	m := NewMachine().Consult(reglas)
	total := m.ProveAll("reliquidar_aporte(" + cedula + "," + strconv.Itoa(ano) + ",N,T).")
	for _, solution := range total {

		detallePreliquidacion.ValorCalculado, _ = strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
		conceptoNomina.NombreConcepto = fmt.Sprintf("%s", solution.ByName_("N"))

		codigo := m.ProveAll(`codigo_concepto(` + conceptoNomina.NombreConcepto + `,C,N).`)
		for _, cod := range codigo {
			conceptoNomina.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			conceptoNomina.NaturalezaConceptoNominaId, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		detallePreliquidacion.Id = 0
		detallePreliquidacion.ConceptoNominaId = &models.ConceptoNomina{Id: conceptoNomina.Id}
		if !EncontrarConcepto(data, detallePreliquidacion.ConceptoNominaId.Id) {
			data = append(data, detallePreliquidacion)
		}
	}
	return data

}
