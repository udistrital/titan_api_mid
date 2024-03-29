package golog

import (
	"fmt"
	"strconv"

	. "github.com/udistrital/golog"
	models "github.com/udistrital/titan_api_mid/models"
)

func LiquidarMesHCH(reglas string, cedula string, ano int, detallePreliquidacion models.DetallePreliquidacion) (data []models.DetallePreliquidacion) {
	var conceptoNomina models.ConceptoNomina
	m := NewMachine().Consult(reglas)
	total := m.ProveAll("liquidar_hch(" + cedula + "," + strconv.Itoa(ano) + ",N,T).")
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

func LiquidarMesHCHOld(reglas string, cedula string, ano int, detallePreliquidacion models.DetallePreliquidacionOld) (data []models.DetallePreliquidacionOld) {
	var conceptoNomina models.ConceptoNomina
	m := NewMachine().Consult(reglas)
	total := m.ProveAll("liquidar_hch(" + cedula + "," + strconv.Itoa(ano) + ",N,T).")
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

func LiquidarMesHCS(reglas string, contrato models.Contrato, detallePreliquidacion models.DetallePreliquidacion, mesFinal bool) (data []models.DetallePreliquidacion) {
	var conceptoNomina models.ConceptoNomina
	//var desagregado models.Desagregado
	cedula := contrato.Documento
	ano := contrato.Vigencia
	m := NewMachine().Consult(reglas)
	//fmt.Println(reglas)
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
		detallePreliquidacion.ConceptoNominaId = &models.ConceptoNomina{
			Id:                         conceptoNomina.Id,
			NombreConcepto:             conceptoNomina.NombreConcepto,
			NaturalezaConceptoNominaId: conceptoNomina.NaturalezaConceptoNominaId,
		}
		data = append(data, detallePreliquidacion)
	}
	if mesFinal {
		total := m.ProveAll("liquidar_prestacion(" + cedula + "," + strconv.Itoa(ano) + ",N,T).")
		for _, solution := range total {
			//detallePreliquidacion.ValorCalculado, _ = strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
			conceptoNomina.NombreConcepto = fmt.Sprintf("%s", solution.ByName_("N"))
			codigo := m.ProveAll(`codigo_concepto(` + conceptoNomina.NombreConcepto + `,C,N).`)
			for _, cod := range codigo {
				conceptoNomina.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				switch {
				case conceptoNomina.NombreConcepto == "primaNavidad":
					detallePreliquidacion.ValorCalculado = contrato.Desagregado.PrimaNavidad
				case conceptoNomina.NombreConcepto == "cesantias":
					detallePreliquidacion.ValorCalculado = contrato.Desagregado.Cesantias
				case conceptoNomina.NombreConcepto == "priServ":
					detallePreliquidacion.ValorCalculado = contrato.Desagregado.PrimaServicios
				case conceptoNomina.NombreConcepto == "primaVacaciones":
					detallePreliquidacion.ValorCalculado = contrato.Desagregado.PrimaVacaciones
				case conceptoNomina.NombreConcepto == "vacaciones":
					detallePreliquidacion.ValorCalculado = contrato.Desagregado.Vacaciones
				case conceptoNomina.NombreConcepto == "interesCesantias":
					detallePreliquidacion.ValorCalculado = contrato.Desagregado.InteresesCesantias
				case conceptoNomina.NombreConcepto == "bonServ":
					detallePreliquidacion.ValorCalculado = contrato.Desagregado.BonificacionServicios
				}
				conceptoNomina.NaturalezaConceptoNominaId, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
			}

			detallePreliquidacion.Id = 0
			detallePreliquidacion.ConceptoNominaId = &models.ConceptoNomina{
				Id:                         conceptoNomina.Id,
				NombreConcepto:             conceptoNomina.NombreConcepto,
				NaturalezaConceptoNominaId: conceptoNomina.NaturalezaConceptoNominaId,
			}
			data = append(data, detallePreliquidacion)
		}
	}
	return data
}

func LiquidarMesHCSOld(reglas string, cedula string, ano int, detallePreliquidacion models.DetallePreliquidacionOld, mesFinal bool) (data []models.DetallePreliquidacionOld) {
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
		data = append(data, detallePreliquidacion)
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
