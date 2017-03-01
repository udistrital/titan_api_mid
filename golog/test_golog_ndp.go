package golog

import (
	"fmt"
	"strconv"
	"time"
	models "titan_api_mid/models"

	. "github.com/mndrix/golog"
)

func CargarReglasDP(idProveedor int, reglas string, informacion_cargo []models.DocenteCargo, dias_trabajados float64, periodo string, puntos string, regimen string) (rest []models.Respuesta) {
	var resultado []models.Respuesta
	temp := models.Respuesta{}
	var lista_descuentos []models.ConceptosResumen
	var regimen_numero string
	var total_devengado float64
	var devengo string
	var cargo string
	var salario string

	fechaInicio := informacion_cargo[0].FechaInicio
	fechaActual := time.Now().Local()
	asignacion_basica_string := strconv.Itoa(informacion_cargo[0].Asignacion_basica)
	m := NewMachine().Consult(reglas)
	//liquidar(R,P,V,T,L).
	if informacion_cargo[0].Cargo == "DC" {
		cargo = "1"
	} else {
		cargo = "2"
	}
	if regimen == "N" {
		regimen_numero = "1"
	} else {
		regimen_numero = "2"
	}
	fmt.Println("ASFDGASDSASSSSFS0FSSDF")
	fmt.Println(informacion_cargo[0])
	fmt.Println(regimen_numero + " " + " " + puntos + " " + asignacion_basica_string + " " + cargo)
	//falta arreglar el periodo para que sea congruente con los valores provenientes de la bd liquidar(R,P,V,T,C,L)
	valor_salario := m.ProveAll("liquidar(" + regimen_numero + "," + puntos + "," + asignacion_basica_string + ", 1," + cargo + ",L ).")
	for _, solution := range valor_salario {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("L")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "pagoBruto",
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		salario = strconv.FormatFloat(Valor, 'f', 6, 64)
		total_devengado = total_devengado + Valor
		fmt.Println(total_devengado)
		codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)
		resultado = append(resultado, temp)
		temp.Conceptos = &lista_descuentos
	}

	//condicional para saber si debe aplicarse la bonificacion por servicios
	if fechaInicio.Month() == fechaActual.Month() {
		bonificacion := m.ProveAll("bonificacionServicios(" + salario + ",S).")
		for _, solution := range bonificacion {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("S")), 64)
			temp_conceptos := models.ConceptosResumen{Nombre: "bonServ",
				Valor: fmt.Sprintf("%.0f", Valor),
			}
			total_devengado = total_devengado + Valor
			fmt.Println(Valor)
			codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			}

			lista_descuentos = append(lista_descuentos, temp_conceptos)
			temp.Conceptos = &lista_descuentos
			resultado = append(resultado, temp)
		}
	}

	devengo = strconv.FormatFloat(total_devengado, 'f', 6, 64)
	salud_empleado := m.ProveAll("salud(" + devengo + ",S).")
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

	pension_empleado := m.ProveAll("pension(" + devengo + ",S).")
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
