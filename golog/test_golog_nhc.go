package golog

import (
  "fmt"
  "strconv"
  models "github.com/udistrital/titan_api_mid/models"
  . "github.com/udistrital/golog"
)

func CargarReglasHCS(idProveedor int, reglas string, periodo string) (rest []models.Respuesta) {

  var resultado []models.Respuesta
	var lista_descuentos []models.ConceptosResumen
	var lista_novedades []models.ConceptosResumen
  var lista_retefuente []models.ConceptosResumen
	var tipoPreliquidacion_string = "2";

	reglas = reglas + "cargo(0)."
	reglas = reglas + "periodo("+periodo+")."
  fmt.Println("reglas", reglas)
  lista_descuentos,total_devengado_no_novedad = CalcularConceptosHCS(idProveedor,periodo,reglas, tipoPreliquidacion_string)
	lista_novedades = ManejarNovedadesHCS(reglas,idProveedor, tipoPreliquidacion_string,periodo)
  lista_retefuente = CalcularReteFuenteSal(tipoPreliquidacion_string,reglas, lista_descuentos);
	total_calculos = append(total_calculos, lista_descuentos...)
	total_calculos = append(total_calculos, lista_novedades...)
  total_calculos = append(total_calculos, lista_retefuente...)
	resultado = GuardarConceptosHCS(reglas,total_calculos)

	total_calculos = []models.ConceptosResumen{}
  ibc = 0;
	return resultado;
}

func CalcularConceptosHCS(idProveedor int, periodo,reglas,tipoPreliquidacion_string string)(rest []models.ConceptosResumen, total_dev float64) {

  var lista_descuentos []models.ConceptosResumen


	m := NewMachine().Consult(reglas)


    valor_pago := m.ProveAll("valor_pago(X,"+periodo+",T).")
    for _, solution := range valor_pago {

      Valor,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
      Nom_Concepto := "salarioBase"
      temp_conceptos := models.ConceptosResumen {Nombre : "salarioBase" ,
                                                 Valor : fmt.Sprintf("%.0f", Valor),
                                                                       }
    	reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
      codigo := m.ProveAll(`codigo_concepto(`+temp_conceptos.Nombre+`,C,N).`)

      for _, cod := range codigo{
        temp_conceptos.Id , _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
        temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
        temp_conceptos.DiasLiquidados = dias_a_liquidar
				temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
       }
      lista_descuentos = append(lista_descuentos,temp_conceptos)

    }



    descuentos := m.ProveAll("concepto_ley(X,Y,"+periodo+",B,N).")
    for _, solution := range descuentos {

      Base,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("B")), 64)
      Valor,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("Y")), 64)
      Nom_Concepto := fmt.Sprintf("%s", solution.ByName_("N"))

      temp_conceptos := models.ConceptosResumen {Nombre : fmt.Sprintf("%s", solution.ByName_("N")),
                                                 Base : fmt.Sprintf("%.0f", Base),
                                                 Valor : fmt.Sprintf("%.0f", Valor),
                                                                       }
      reglas = reglas + "sumar_ibc("+Nom_Concepto+","+strconv.Itoa(int(Valor))+")."
      codigo := m.ProveAll("codigo_concepto("+temp_conceptos.Nombre+",C,N).")

      for _, cod := range codigo{

        temp_conceptos.Id , _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
        temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
        temp_conceptos.TipoPreliquidacion = tipoPreliquidacion_string
       }

      lista_descuentos = append(lista_descuentos,temp_conceptos)
      }

    CalcularIBC(reglas)

  	return lista_descuentos,ibc

}

func GuardarConceptosHCS (reglas string,lista_descuentos []models.ConceptosResumen)(rest []models.Respuesta){
		temp := models.Respuesta{}
		var resultado []models.Respuesta
    m := NewMachine().Consult(reglas)

    temp_conceptos := models.ConceptosResumen{Nombre: "ibc_liquidado",
      Valor: fmt.Sprintf("%.0f", total_devengado_no_novedad),
    }

    codigo := m.ProveAll(`codigo_concepto(ibc_liquidado,C, N).`)

    for _, cod := range codigo {
      temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
      temp_conceptos.DiasLiquidados = dias_a_liquidar
    }

		lista_descuentos = append(lista_descuentos, temp_conceptos)

		temp.Conceptos = &lista_descuentos
		resultado = append(resultado, temp)
		total_devengado_novedad = 0
		total_devengado_no_novedad = 0
		return resultado
}


func ManejarNovedadesHCS(reglas string, idProveedor int, tipoPreliquidacion, periodo string) (rest []models.ConceptosResumen){

	var lista_novedades []models.ConceptosResumen

	f := NewMachine().Consult(reglas)

	idProveedorString := strconv.Itoa(idProveedor)
	novedades := f.ProveAll("info_concepto(" + idProveedorString + ",T,"+periodo+",N,R).")

	for _, solution := range novedades {
    fmt.Println("novedades",solution)
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
