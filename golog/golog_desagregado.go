package golog

import (
	"fmt"
	"strconv"

	. "github.com/udistrital/golog"
	models "github.com/udistrital/titan_api_mid/models"
)

func DesagregarContrato(reglas string, categoria string, cedula string, dedicacion string, ano string) (desagregado models.DesagregadoContratoHCS) {

	var nombreConcepto string
	var valorCalculado float64
	var idConcepto int
	m := NewMachine().Consult(reglas)
	fmt.Println("desagregado(" + categoria + "," + cedula + "," + dedicacion + "," + ano + ",N,V).")
	total := m.ProveAll("desagregado(" + categoria + "," + cedula + "," + dedicacion + "," + ano + ",N,V).")
	for _, solution := range total {
		valorCalculado, _ = strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
		nombreConcepto = fmt.Sprintf("%s", solution.ByName_("N"))

		codigo := m.ProveAll(`codigo_concepto(` + nombreConcepto + `,C,N).`)
		for _, cod := range codigo {
			idConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
		}

		switch idConcepto {
		case 152:
			desagregado.SueldoBasico = valorCalculado
		case 540:
			desagregado.Cesantias = valorCalculado
		case 163:
			desagregado.PrimaServicios = valorCalculado
		case 539:
			desagregado.PrimaNavidad = valorCalculado
		case 549:
			desagregado.PrimaVacaciones = valorCalculado
		case 550:
			desagregado.Vacaciones = valorCalculado
		case 541:
			desagregado.InteresesCesantias = valorCalculado
		case 179:
			desagregado.BonificacionServicios = valorCalculado
		default:
			fmt.Println("No evalué nada")
		}
	}

	return desagregado
}
