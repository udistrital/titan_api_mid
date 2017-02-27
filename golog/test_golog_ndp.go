package golog

import (
	"fmt"
	"strconv"
	models "titan_api_mid/models"

	. "github.com/mndrix/golog"
)

func CargarReglasDP(idProveedor int, reglas string, informacion_cargo []models.DocenteCargo, dias_trabajados float64, periodo string, puntos string, regimen string) (rest []models.Respuesta) {
	var resultado []models.Respuesta
	temp := models.Respuesta{}
	var lista_descuentos []models.ConceptosResumen
	var regimen_numero string
	var total_devengado string
	asignacion_basica_string := strconv.Itoa(informacion_cargo[0].Asignacion_basica)
	m := NewMachine().Consult(reglas)
	//liquidar(R,P,V,T,L).
	if regimen == "N" {
		regimen_numero = "1"
	} else {
		regimen_numero = "2"
	}
	fmt.Println(asignacion_basica_string)
	//falta arreglar el periodo para que sea congruente con los valores provenientes de la bd
	valor_salario := m.ProveAll("liquidar(" + regimen_numero + "," + puntos + "," + asignacion_basica_string + "," + periodo + ",L ).")
	for _, solution := range valor_salario {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("L")), 64)
		total_devengado = strconv.FormatFloat(Valor, 'f', 6, 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "pagoBruto",
			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)
		temp.Conceptos = &lista_descuentos
		resultado = append(resultado, temp)
	}

	salud_empleado := m.ProveAll("salud(" + total_devengado + ",S).")
	for _, solution := range salud_empleado {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("S")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "salud",
			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)
		temp.Conceptos = &lista_descuentos
		resultado = append(resultado, temp)
	}

	pension_empleado := m.ProveAll("salud(" + total_devengado + ",S).")
	for _, solution := range pension_empleado {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("S")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "pension",
			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)
		temp.Conceptos = &lista_descuentos
		resultado = append(resultado, temp)
	}
	idProveedorString := strconv.Itoa(idProveedor)
	novedades := m.ProveAll("info_concepto(" + idProveedorString + ",T,2017,N,R).")

	for _, solution := range novedades {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)
		temp.Conceptos = &lista_descuentos
		resultado = append(resultado, temp)
	}

	return resultado
}
