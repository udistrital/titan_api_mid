package golog

import (
	"fmt"
	"strconv"
	"titan_api_mid/models"

	. "github.com/mndrix/golog"
)

func CargarReglasFP(reglas string, idProveedor int, informacion_cargo []models.FuncionarioCargo, dias_laborados float64, periodo string, esAnual int, porcentajePT int) (rest []models.Respuesta) {

	var resultado []models.Respuesta
	temp := models.Respuesta{}
	var lista_descuentos []models.ConceptosResumen

	asignacion_basica_string := strconv.Itoa(informacion_cargo[0].Asignacion_basica)
	id_cargo_string := strconv.Itoa(informacion_cargo[0].Id)
	dias_laborados_string := strconv.Itoa(int(dias_laborados))
	porcentaje_PT_string := strconv.Itoa(porcentajePT)

	var total_devengado float64
	var dias_a_liquidar string
	var tipoNomina_string string = "1"
	var tipoNomina int = 1

	if tipoNomina == 0 || tipoNomina == 1 {
		dias_a_liquidar = "15"

	} else {
		dias_a_liquidar = "30"

	}
	reglas = reglas + "salario_base(" + asignacion_basica_string + ")."
	reglas = reglas + "tipo_nomina(" + tipoNomina_string + ")."
	m := NewMachine().Consult(reglas)
	valor_salario := m.ProveAll("sb(" + asignacion_basica_string + "," + tipoNomina_string + "," + dias_a_liquidar + ",V).")
	for _, solution := range valor_salario {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
		total_devengado = total_devengado + Valor
		temp_conceptos := models.ConceptosResumen{Nombre: "salarioBase",
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

	valor_gastos_representacion := m.ProveAll("gr(" + asignacion_basica_string + "," + dias_a_liquidar + "," + tipoNomina_string + ",2016," + id_cargo_string + ",V).")
	for _, solution := range valor_gastos_representacion {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
		total_devengado = total_devengado + Valor
		temp_conceptos := models.ConceptosResumen{Nombre: "gastosRep",
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

	valor_prima_antiguedad := m.ProveAll("prima_ant(" + asignacion_basica_string + "," + dias_a_liquidar + "," + tipoNomina_string + ",2016," + dias_laborados_string + ",V).")
	for _, solution := range valor_prima_antiguedad {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
		total_devengado = total_devengado + Valor
		temp_conceptos := models.ConceptosResumen{Nombre: "primaAnt",
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

	if esAnual == 1 {
		valor_bonificacion_servicios := m.ProveAll("bon_ser(" + asignacion_basica_string + "," + dias_a_liquidar + "," + tipoNomina_string + ",2016," + dias_laborados_string + "," + id_cargo_string + ",V).")
		for _, solution := range valor_bonificacion_servicios {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
			total_devengado = total_devengado + Valor
			temp_conceptos := models.ConceptosResumen{Nombre: "bonServ",
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
	}

	if porcentajePT != 0 {
		fmt.Println("asdasdasdsad")

		valor_prima_tecnica := m.ProveAll("prima_tecnica(" + asignacion_basica_string + "," + dias_a_liquidar + "," + tipoNomina_string + "," + porcentaje_PT_string + ",V).")
		for _, solution := range valor_prima_tecnica {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
			total_devengado = total_devengado + Valor
			temp_conceptos := models.ConceptosResumen{Nombre: "priTec",
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
	}

	valor_prima_secretarial := m.ProveAll("prima_secretarial(" + asignacion_basica_string + ",2016," + id_cargo_string + "," + tipoNomina_string + "," + dias_laborados_string + ",V).")
	for _, solution := range valor_prima_secretarial {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
		total_devengado = total_devengado + Valor
		temp_conceptos := models.ConceptosResumen{Nombre: "primaSecr",
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

	total_devengado_string := strconv.Itoa(int(total_devengado))

	valor_salud := m.ProveAll("salud_fun(" + total_devengado_string + ",2016,V).")
	for _, solution := range valor_salud {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
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

	valor_pension := m.ProveAll("pension_fun(" + total_devengado_string + ",2016,V).")
	for _, solution := range valor_pension {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
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

	if tipoNomina == 1 || tipoNomina == 2 {

		valor_fondo_solidaridad := m.ProveAll("fondo_solidaridad_fun(" + asignacion_basica_string + ",2017,V).")
		for _, solution := range valor_fondo_solidaridad {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
			temp_conceptos := models.ConceptosResumen{Nombre: "fondoSolidaridad",
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
