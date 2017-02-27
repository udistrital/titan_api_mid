package golog

import (
  "fmt"
  "strconv"
  models "titan_api_mid/models"
  . "github.com/mndrix/golog"
)

func CargarReglas(reglas string, periodo string) (rest []models.Respuesta) {
//******QUITAR ARREGLO, DEJAR UNA SOLA VARIABLE PARA LAS REGLAS ******
  m := NewMachine().Consult(reglas)

  var resultado []models.Respuesta

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
      Valor,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("P")), 64)
      temp.Nombre_Cont = fmt.Sprintf("%s", solution.ByName_("X"))

      temp_conceptos := models.ConceptosResumen {Nombre : "pagoBruto" ,
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
      novedades := m.ProveAll("info_concepto("+temp.Nombre_Cont+",T,"+periodo+",N,R).")

      for _, solution := range novedades {

        Valor,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("R")), 64)


        temp_conceptos := models.ConceptosResumen {Nombre : fmt.Sprintf("%s", solution.ByName_("N")),
                                                   Valor : fmt.Sprintf("%.0f", Valor),
                                                                         }
        codigo := m.ProveAll("codigo_concepto("+temp_conceptos.Nombre+",C).")
        for _, cod := range codigo{
          temp_conceptos.Id , _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
         }
        lista_descuentos = append(lista_descuentos,temp_conceptos)
        }
      temp.Conceptos = &lista_descuentos



    resultado = append(resultado,temp)
//  }



  return resultado

}
