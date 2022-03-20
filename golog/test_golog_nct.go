package golog

import (
	"fmt"
	"strconv"

	. "github.com/udistrital/golog"
	models "github.com/udistrital/titan_api_mid/models"
)

func LiquidarMesCPS(reglas string, cedula string, ano int, detallePreliquidacion models.DetallePreliquidacion) (data []models.DetallePreliquidacion) {
	var conceptoNomina models.ConceptoNomina
	m := NewMachine().Consult(reglas)
	total := m.ProveAll("liquidar_ct(" + cedula + "," + strconv.Itoa(ano) + ",N,T).")
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

	return data
}
