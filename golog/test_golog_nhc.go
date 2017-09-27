package golog

import (
  "fmt"
  "strconv"
  models "github.com/udistrital/titan_api_mid/models"
  . "github.com/mndrix/golog"
)

func CargarReglasHCS(idProveedor int, reglas string, periodo string) (rest []models.Respuesta) {

  var resultado []models.Respuesta
	var lista_descuentos []models.ConceptosResumen
	var lista_novedades []models.ConceptosResumen
  var lista_retefuente []models.ConceptosResumen
	var tipoPreliquidacion_string = "2";

	reglas = reglas + "cargo(0)."
	reglas = reglas + "periodo("+periodo+")."

  lista_descuentos = CalcularConceptosHCS(idProveedor,periodo,reglas)
	lista_novedades = ManejarNovedadesHCS(reglas,idProveedor, tipoPreliquidacion_string,periodo)
  lista_retefuente = CalcularReteFuenteSal(tipoPreliquidacion_string,reglas, lista_descuentos);
	total_calculos = append(total_calculos, lista_descuentos...)
	total_calculos = append(total_calculos, lista_novedades...)
  total_calculos = append(total_calculos, lista_retefuente...)
	resultado = GuardarConceptosHCS(total_calculos)

	total_calculos = []models.ConceptosResumen{}

	return resultado;
}

func CalcularConceptosHCS(idProveedor int, periodo,reglas string)(rest []models.ConceptosResumen) {

  var lista_descuentos []models.ConceptosResumen

	var nombre_archivo string

	nombre_archivo = "reglas" + strconv.Itoa(idProveedor) + ".txt"
	if err := WriteStringToFile(nombre_archivo, reglas); err != nil {
      panic(err)
  }

	m := NewMachine().Consult(reglas)


    valor_pago := m.ProveAll("valor_pago(X,V,"+periodo+",T).")
    for _, solution := range valor_pago {
      Valor,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
      temp_conceptos := models.ConceptosResumen {Nombre : "salarioBase" ,
                                                 Valor : fmt.Sprintf("%.0f", Valor),
                                                                       }
      codigo := m.ProveAll(`codigo_concepto(`+temp_conceptos.Nombre+`,C).`)

      for _, cod := range codigo{
        temp_conceptos.Id , _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))

       }
      lista_descuentos = append(lista_descuentos,temp_conceptos)

    }



    descuentos := m.ProveAll("concepto_ley(X,Y,"+periodo+",B,N).")
    for _, solution := range descuentos {
      Base,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("B")), 64)
      Valor,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("Y")), 64)


      temp_conceptos := models.ConceptosResumen {Nombre : fmt.Sprintf("%s", solution.ByName_("N")),
                                                 Base : fmt.Sprintf("%.0f", Base),
                                                 Valor : fmt.Sprintf("%.0f", Valor),
                                                                       }
      codigo := m.ProveAll("codigo_concepto("+temp_conceptos.Nombre+",C).")

      for _, cod := range codigo{
        temp_conceptos.Id , _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))

       }

      lista_descuentos = append(lista_descuentos,temp_conceptos)
      }


  	return lista_descuentos

}

func GuardarConceptosHCS (lista_descuentos []models.ConceptosResumen)(rest []models.Respuesta){
		temp := models.Respuesta{}
		var resultado []models.Respuesta

		temp.Conceptos = &lista_descuentos
		resultado = append(resultado, temp)
		total_devengado_novedad = 0
		total_devengado_no_novedad = 0
		return resultado
}


func ManejarNovedadesHCS(reglas string, idProveedor int, tipoPreliquidacion, periodo string) (rest []models.ConceptosResumen){
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
