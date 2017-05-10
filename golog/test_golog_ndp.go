package golog

import (
	"fmt"
	"strconv"
	"time"
	models "github.com/udistrital/titan_api_mid/models"

	. "github.com/mndrix/golog"
)



var salario string


func CargarReglasDP(idProveedor int, reglas string, informacion_cargo []models.DocenteCargo, dias_trabajados float64, periodo string, puntos string, regimen string,tipoPreliquidacion string) (rest []models.Respuesta) {
	//Definición de variables

	var resultado []models.Respuesta
	var lista_descuentos []models.ConceptosResumen
	var lista_novedades []models.ConceptosResumen
	var tipoPreliquidacion_string string
	var regimen_numero string
	var cargo string


	fechaInicio := informacion_cargo[0].FechaInicio
	fechaActual := time.Now().Local()
	asignacion_basica_string := strconv.Itoa(informacion_cargo[0].Asignacion_basica)
	tipoPreliquidacion_string = tipoPreliquidacion

	if informacion_cargo[0].Cargo == "DC" {
		cargo = "1"
	} else {
		cargo = "2"
	}

	if regimen == "N" {
		fmt.Println("Nuevo")
		regimen_numero = "1"
	} else {
		fmt.Println("antiguo")
		regimen_numero = "2"
	}

	if tipoPreliquidacion_string  == "0" || tipoPreliquidacion_string  == "1" {
		dias_a_liquidar = "15"

	} else {
		dias_a_liquidar = "26"
	}

		m := NewMachine().Consult(reglas)

	fmt.Println(regimen_numero + " " + " " + puntos + " " + asignacion_basica_string + " " + cargo)

	// ----- Nomina ordinaria ----- Proceso de cálculo, manejo de novedades y guardado de conceptos
	lista_descuentos = CalcularConceptosDP(m, reglas,dias_a_liquidar,asignacion_basica_string, tipoPreliquidacion_string,regimen_numero, puntos, cargo, fechaInicio, fechaActual)
	ibc = 0
	lista_novedades = ManejarNovedadesDP(reglas,idProveedor, tipoPreliquidacion_string)
	total_calculos = append(total_calculos, lista_descuentos...)
	total_calculos = append(total_calculos, lista_novedades...)
	resultado = GuardarConceptosDP(total_calculos)
	total_calculos = []models.ConceptosResumen{}

	// ---------------------------
	return resultado;

	//falta arreglar el periodo para que sea congruente con los valores provenientes de la bd liquidar(R,P,V,T,C,L)
}

	func CalcularConceptosDP(m Machine, reglas, dias_a_liquidar, asignacion_basica_string, tipoPreliquidacion_string, regimen_numero, puntos, cargo string, fechaInicio, fechaActual time.Time) (rest []models.ConceptosResumen){
		var lista_descuentos []models.ConceptosResumen

		novedades_devengo := m.ProveAll("novedades_devengos(X).")
		for _, solution := range novedades_devengo {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("X")), 64)
			ibc = ibc + Valor

		}
		
		valor_salario := m.ProveAll("liquidar(" + regimen_numero + "," + puntos + "," + asignacion_basica_string + ", "+dias_a_liquidar+"," + cargo + ",L ).")
		for _, solution := range valor_salario {
		  Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("L")), 64)
		  temp_conceptos := models.ConceptosResumen{Nombre: "salarioBase",
		    Valor: fmt.Sprintf("%.0f", Valor),
		  }
			salario = strconv.FormatFloat(Valor, 'f', 6, 64)
			reglas = reglas + "sumar_ibc(salarioBase,"+strconv.Itoa(int(Valor))+")."
			codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
			}

			lista_descuentos = append(lista_descuentos, temp_conceptos)

		}

		if fechaInicio.Month() == fechaActual.Month() && regimen_numero == "2" {
		  bonificacion := m.ProveAll("bonificacionServicios(" + salario + ",S).")
		  for _, solution := range bonificacion {
		    Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("S")), 64)
		    temp_conceptos := models.ConceptosResumen{Nombre: "bonServ",
		      Valor: fmt.Sprintf("%.0f", Valor),
		    }

				reglas = reglas + "sumar_ibc(salarioBase,"+strconv.Itoa(int(Valor))+")."
				codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

				for _, cod := range codigo {
					temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
					temp_conceptos.DiasLiquidados = dias_a_liquidar
					temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
				}

				lista_descuentos = append(lista_descuentos, temp_conceptos)
		  }
		}

		//Previo a pagos de salud y pensión se calcula el IBC
		CalcularIBC(reglas)
		total_devengado_string := strconv.Itoa(int(ibc))
		fmt.Println("total devengado")


		salud_empleado := m.ProveAll("salud(" + total_devengado_string + ",S).")
		for _, solution := range salud_empleado {
		  Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("S")), 64)
		  temp_conceptos := models.ConceptosResumen{Nombre: "salud",
		    Valor: fmt.Sprintf("%.0f", Valor),
		  }

			codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

 		 for _, cod := range codigo {
 			 temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
 			 temp_conceptos.DiasLiquidados = dias_a_liquidar
 			 temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
 		 }
 		 lista_descuentos = append(lista_descuentos, temp_conceptos)

		}

		return lista_descuentos
	}

	func ManejarNovedadesDP(reglas string, idProveedor int, tipoPreliquidacion string) (rest []models.ConceptosResumen){
		var lista_novedades []models.ConceptosResumen

		f := NewMachine().Consult(reglas)

		idProveedorString := strconv.Itoa(idProveedor)
		novedades := f.ProveAll("info_concepto(" + idProveedorString + ",T,2017,N,R).")

		for _, solution := range novedades {

			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
			temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
				Valor: fmt.Sprintf("%.0f", Valor),
			}
			codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")
			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion
			}

			lista_novedades = append(lista_novedades, temp_conceptos)

		}

		return lista_novedades

	}


func GuardarConceptosDP(lista_descuentos []models.ConceptosResumen)(rest []models.Respuesta){
			temp := models.Respuesta{}
			var resultado []models.Respuesta

			/*temp_conceptos := models.ConceptosResumen{Nombre: "ibc_liquidado",
				Valor: fmt.Sprintf("%.0f", total_devengado_no_novedad),
			}
			temp_conceptos.Id = 2322
			temp_conceptos.DiasLiquidados = dias_a_liquidar

			lista_descuentos = append(lista_descuentos, temp_conceptos)

			temp_conceptos = models.ConceptosResumen{Nombre: "ibc_novedad",
				Valor: fmt.Sprintf("%.0f", total_devengado_novedad),
			}
			temp_conceptos.Id = 2327
			temp_conceptos.DiasLiquidados = dias_novedad_string


			lista_descuentos = append(lista_descuentos, temp_conceptos)
*/
			temp.Conceptos = &lista_descuentos
			resultado = append(resultado, temp)
			total_devengado_novedad = 0
			total_devengado_no_novedad = 0
			return resultado
	}
