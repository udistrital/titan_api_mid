package golog

import (
	"fmt"
	"strconv"

	. "github.com/udistrital/golog"
	models "github.com/udistrital/titan_api_mid/models"
)



func CargarReglasCT(idProveedor int, reglas string, periodo string) (rest []models.Respuesta) {

	var resultado []models.Respuesta
	var lista_descuentos []models.ConceptosResumen
	var lista_novedades []models.ConceptosResumen
	var lista_retefuente []models.ConceptosResumen
	var tipoPreliquidacion_string = "2";

	reglas = reglas + "cargo(0)."
	reglas = reglas + "periodo("+periodo+")."
	//reglas = reglas + "nomina("+periodo+")."

	lista_descuentos = CalcularConceptosCT(idProveedor,periodo,reglas, tipoPreliquidacion_string)
	lista_novedades = ManejarNovedadesCT(reglas,idProveedor, tipoPreliquidacion_string,periodo)
	lista_retefuente = CalcularReteFuenteSal(tipoPreliquidacion_string,reglas, lista_descuentos);
	total_calculos = append(total_calculos, lista_descuentos...)
	total_calculos = append(total_calculos, lista_novedades...)
	total_calculos = append(total_calculos, lista_retefuente...)
	resultado = GuardarConceptosCT(total_calculos)

	total_calculos = []models.ConceptosResumen{}

	return resultado;

}

func CalcularConceptosCT (idProveedor int, periodo,reglas, tipoPreliquidacion_string string)(rest []models.ConceptosResumen){

	var lista_descuentos []models.ConceptosResumen

	var nombre_archivo string
	var salarioBase float64
	var salarioBase_string string

	nombre_archivo = "reglas" + strconv.Itoa(idProveedor) + ".txt"
	if err := WriteStringToFile(nombre_archivo, reglas); err != nil {
      panic(err)
  }

	m := NewMachine().Consult(reglas)

	valor_pago := m.ProveAll("valor_pago(X,"+periodo+",P).")

	for _, solution := range valor_pago {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("P")), 64)

		salarioBase = Valor

		salarioBase_string = strconv.Itoa(int(salarioBase))
		temp_conceptos := models.ConceptosResumen{Nombre: "salarioBase",
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)
	}


	descuento_reteica := m.ProveAll("calcular_reteica("+salarioBase_string+", "+periodo+", R).")
	for _, solution := range descuento_reteica {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "reteIca",
			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)

	}

	//Estampila UD

	descuento_estampilla := m.ProveAll("calcular_estampilla("+salarioBase_string+", "+periodo+", R).")

	for _, solution := range descuento_estampilla {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "estampillaUD",

			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)


	}

	//proCultura

	descuento_procultura := m.ProveAll("calcular_procultura("+salarioBase_string+", "+periodo+", R).")

	for _, solution := range descuento_procultura {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "proCultura",

			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)


	}

	//proCultura

	descuento_adulto_mayor := m.ProveAll("calcular_adulto_mayor("+salarioBase_string+", "+periodo+", R).")

	for _, solution := range descuento_adulto_mayor {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "adultoMayor",

			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)


	}

	descuento_salud:= m.ProveAll("calcular_salud(si,"+salarioBase_string+", "+periodo+", R).")

	for _, solution := range descuento_salud {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "salud",

			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)


	}

	descuento_pension:= m.ProveAll("calcular_pension(si,"+salarioBase_string+", "+periodo+", R).")

	for _, solution := range descuento_pension {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: "pension",

			Valor: fmt.Sprintf("%.0f", Valor),
		}

		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C, N).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		lista_descuentos = append(lista_descuentos, temp_conceptos)


	}
	return lista_descuentos

}

func GuardarConceptosCT (lista_descuentos []models.ConceptosResumen)(rest []models.Respuesta){
		temp := models.Respuesta{}
		var resultado []models.Respuesta

		temp.Conceptos = &lista_descuentos
		resultado = append(resultado, temp)
		total_devengado_novedad = 0
		total_devengado_no_novedad = 0
		return resultado
}


func ManejarNovedadesCT(reglas string, idProveedor int, tipoPreliquidacion, periodo string) (rest []models.ConceptosResumen){
	fmt.Println("novedades")
	var lista_novedades []models.ConceptosResumen

	f := NewMachine().Consult(reglas)

	idProveedorString := strconv.Itoa(idProveedor)
	novedades := f.ProveAll("info_concepto(" + idProveedorString + ",T,"+periodo+",N,R).")

	for _, solution := range novedades {

		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		temp_conceptos := models.ConceptosResumen{Nombre: fmt.Sprintf("%s", solution.ByName_("N")),
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		codigo := f.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C,N).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.DiasLiquidados = dias_a_liquidar
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacion
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		lista_novedades = append(lista_novedades, temp_conceptos)

	}

	return lista_novedades

}
