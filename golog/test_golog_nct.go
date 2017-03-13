package golog

import (
	"fmt"
	"strconv"

	. "github.com/mndrix/golog"
	models "github.com/udistrital/titan_api_mid/models"
)

func CargarReglasCT(reglas string, periodo string) (rest []models.Respuesta) {
	//******QUITAR ARREGLO, DEJAR UNA SOLA VARIABLE PARA LAS REGLAS ******
	fmt.Println(reglas)
	m := NewMachine().Consult(reglas)

	var resultado []models.Respuesta
	var salarioBase float64
	/*preliqu := m.ProveAll("valor_pago_neto(X,Y,"+periodo+",V,L,L2).")
	  for _, solution := range preliqu {
	    Neto,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("Y")), 64)
	    Bruto,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
	    temp := models.Respuesta{Valor_neto:fmt.Sprintf("%.0f", Neto),
	                            Nombre_Cont : fmt.Sprintf("%s", solution.ByName_("X")),
	                            Valor_bruto  : fmt.Sprintf("%.0f", Bruto),}*/
	temp := models.Respuesta{}
	valor_pago := m.ProveAll("valor_pago(X,V,P).")
	var lista_descuentos []models.ConceptosResumen
	for _, solution := range valor_pago {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("P")), 64)
		temp.Nombre_Cont = fmt.Sprintf("%s", solution.ByName_("X"))
		salarioBase = Valor
		fmt.Println(salarioBase)
		temp_conceptos := models.ConceptosResumen{Nombre: "pagoBruto",
			Valor: fmt.Sprintf("%.0f", Valor),
		}
		codigo := m.ProveAll(`codigo_concepto(` + temp_conceptos.Nombre + `,C).`)

		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))

		}
		lista_descuentos = append(lista_descuentos, temp_conceptos)
	}

	//DESCUENTOS
	//Reteica

	descuento_reteica := m.ProveAll("calcular_reteica(896292, 2016, R).")
	fmt.Println("valor1")
	for _, solution := range descuento_reteica {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		fmt.Println("valor")
		fmt.Println(Valor)
		temp_conceptos := models.ConceptosResumen{Nombre: "reteIca",

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

	//Estampila UD

	descuento_estampilla := m.ProveAll("calcular_estampilla(896292, 2016, R).")
	fmt.Println("valor1")
	for _, solution := range descuento_estampilla {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		fmt.Println("valor")
		fmt.Println(Valor)
		temp_conceptos := models.ConceptosResumen{Nombre: "estampillaUD",

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

	//proCultura

	descuento_procultura := m.ProveAll("calcular_procultura(827346, 2016, R).")
	fmt.Println("valor1")
	for _, solution := range descuento_procultura {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		fmt.Println("valor")
		fmt.Println(Valor)
		temp_conceptos := models.ConceptosResumen{Nombre: "proCultura",

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

	//proCultura

	descuento_adulto_mayor := m.ProveAll("calcular_adulto_mayor(827346, 2016, R).")
	fmt.Println("valor1")
	for _, solution := range descuento_adulto_mayor {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)
		fmt.Println("valor")
		fmt.Println(Valor)
		temp_conceptos := models.ConceptosResumen{Nombre: "adultoMayor",

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

	//NO TOCAR
	novedades := m.ProveAll("info_concepto(" + temp.Nombre_Cont + ",T," + periodo + ",N,R).")

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
	}
	temp.Conceptos = &lista_descuentos

	resultado = append(resultado, temp)
	//  }

	return resultado

}
