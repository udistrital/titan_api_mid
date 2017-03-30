package golog

import (
	"fmt"
	"strconv"
	"github.com/udistrital/titan_api_mid/models"

	. "github.com/mndrix/golog"

	"time"
)


var total_devengado_no_novedad float64
var total_devengado_novedad float64
var dias_a_liquidar string
var dias_novedad_string string

func CargarReglasFP(fechaPreliquidacion time.Time, reglas string, idProveedor int, informacion_cargo []models.FuncionarioCargo, dias_laborados float64, periodo string, esAnual int, porcentajePT int, tipoNomina string) (rest []models.Respuesta) {

	fmt.Println("prueba test")
	var resultado []models.Respuesta
	var lista_descuentos []models.ConceptosResumen

	asignacion_basica_string := strconv.Itoa(informacion_cargo[0].Asignacion_basica)
	id_cargo_string := strconv.Itoa(informacion_cargo[0].Id)
	dias_laborados_string := strconv.Itoa(int(dias_laborados))

	var tipoNomina_string string

	tipoNomina_string = tipoNomina

	if tipoNomina_string  == "0" || tipoNomina_string  == "1" {
		dias_a_liquidar = "15"

	} else {
		dias_a_liquidar = "30"
	}


	reglas = reglas + "salario_base(" + asignacion_basica_string + ")."
	reglas = reglas + "tipo_nomina(" + tipoNomina_string + ")."
	if err := WriteStringToFile("reglas.txt", reglas); err != nil {
      panic(err)
  }
	m := NewMachine().Consult(reglas)

	novedades_seg_social := m.ProveAll("seg_social(N,A,M,D,AA,MM,DD).")

	for _, solution := range novedades_seg_social {

		novedad := fmt.Sprintf("%s", solution.ByName_("N"))
		AnoDesde,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("A")), 64)
		MesDesde,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("M")), 64)
		DiaDesde,_:= strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("D")), 64)
		AnoHasta,_:= strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("AA")), 64)
		MesHasta,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("MM")), 64)
		DiaHasta,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("DD")), 64)

		afectacion_seg_social := m.ProveAll("afectacion_seguridad("+novedad+").")
		for _, solution := range afectacion_seg_social {

				fmt.Println(solution)
				dias_novedad := CalcularDiasNovedades(fechaPreliquidacion, AnoDesde, MesDesde, DiaDesde, AnoHasta, MesHasta, DiaHasta)
				dias_a_liquidar = strconv.Itoa(int(30 - dias_novedad))
				dias_novedad_string = strconv.Itoa(int(dias_novedad))
				_,total_devengado_novedad = CalcularConceptos(m, dias_novedad_string,asignacion_basica_string,id_cargo_string,dias_laborados_string, tipoNomina_string,esAnual, porcentajePT, idProveedor)

		}

		}


		lista_descuentos,total_devengado_no_novedad = CalcularConceptos(m, dias_a_liquidar,asignacion_basica_string,id_cargo_string,dias_laborados_string, tipoNomina_string,esAnual, porcentajePT, idProveedor)
		resultado = GuardarConceptos(lista_descuentos)
		fmt.Println("resultado")
		fmt.Println(resultado)
		return resultado;


	}

	func CalcularConceptos(m Machine, dias_a_liquidar,asignacion_basica_string,id_cargo_string,dias_laborados_string,tipoNomina_string string,  esAnual,  porcentajePT,idProveedor  int) (rest []models.ConceptosResumen, total_dev float64){

		var lista_descuentos []models.ConceptosResumen
		porcentaje_PT_string := strconv.Itoa(porcentajePT)
		var total_devengado float64

		novedades_devengo := m.ProveAll("novedades_devengos(X).")
		for _, solution := range novedades_devengo {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("X")), 64)
			total_devengado = total_devengado + Valor

			}

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
				temp_conceptos.DiasLiquidados = dias_a_liquidar
			}

			lista_descuentos = append(lista_descuentos, temp_conceptos)


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
				temp_conceptos.DiasLiquidados = dias_a_liquidar

			}
			lista_descuentos = append(lista_descuentos, temp_conceptos)


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
				temp_conceptos.DiasLiquidados = dias_a_liquidar

			}

			lista_descuentos = append(lista_descuentos, temp_conceptos)


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
					temp_conceptos.DiasLiquidados = dias_a_liquidar
				}

				lista_descuentos = append(lista_descuentos, temp_conceptos)


			}
		}

		if porcentajePT != 0 {

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
					temp_conceptos.DiasLiquidados = dias_a_liquidar
				}

				lista_descuentos = append(lista_descuentos, temp_conceptos)


			}
		}

		valor_prima_secretarial := m.ProveAll("prima_secretarial(" + asignacion_basica_string + ",2016," + id_cargo_string + "," + tipoNomina_string + "," + dias_laborados_string + ",V).")
		for _, solution := range valor_prima_secretarial {
			Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
			temp_conceptos := models.ConceptosResumen{Nombre: "primaSecr",
				Valor: fmt.Sprintf("%.0f", Valor),
			}
			codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

			for _, cod := range codigo {
				temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
				temp_conceptos.DiasLiquidados = dias_a_liquidar

			}
			lista_descuentos = append(lista_descuentos, temp_conceptos)


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
				temp_conceptos.DiasLiquidados = dias_a_liquidar
			}
			lista_descuentos = append(lista_descuentos, temp_conceptos)


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
				temp_conceptos.DiasLiquidados = dias_a_liquidar
			}
			lista_descuentos = append(lista_descuentos, temp_conceptos)


		}

		if tipoNomina_string == "1" || tipoNomina_string == "2" {

			valor_fondo_solidaridad := m.ProveAll("fondo_solidaridad_fun(" + asignacion_basica_string + ",2017,V).")
			for _, solution := range valor_fondo_solidaridad {
				Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
				temp_conceptos := models.ConceptosResumen{Nombre: "fondoSolidaridad",
					Valor: fmt.Sprintf("%.0f", Valor),
				}
				codigo := m.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C).")

				for _, cod := range codigo {
					temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
					temp_conceptos.DiasLiquidados = dias_a_liquidar
				}
				lista_descuentos = append(lista_descuentos, temp_conceptos)


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
				temp_conceptos.DiasLiquidados = dias_a_liquidar
			}

			lista_descuentos = append(lista_descuentos, temp_conceptos)

		}

		return lista_descuentos, total_devengado
}

func GuardarConceptos (lista_descuentos []models.ConceptosResumen)(rest []models.Respuesta){
		temp := models.Respuesta{}
		var resultado []models.Respuesta

		temp_conceptos := models.ConceptosResumen{Nombre: "ibc_liquidado",
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

		temp.Conceptos = &lista_descuentos
		resultado = append(resultado, temp)
		total_devengado_novedad = 0
		total_devengado_no_novedad = 0
		return resultado
}
