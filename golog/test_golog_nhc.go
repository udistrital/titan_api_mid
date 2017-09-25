package golog

import (
  "fmt"
  "strconv"
  models "github.com/udistrital/titan_api_mid/models"
  . "github.com/mndrix/golog"
)

func CargarReglas(idProveedor int, reglas string, periodo string) (rest []models.Respuesta) {
//******QUITAR ARREGLO, DEJAR UNA SOLA VARIABLE PARA LAS REGLAS ******
  var nombre_archivo string
  nombre_archivo = "reglas" + strconv.Itoa(idProveedor) + ".txt"
  if err := WriteStringToFile(nombre_archivo, reglas); err != nil {
      panic(err)
  }
  m := NewMachine().Consult(reglas)
  var resultado []models.Respuesta


    temp := models.Respuesta{}
    var lista_descuentos []models.ConceptosResumen

    valor_pago := m.ProveAll("valor_pago(X,V,"+periodo+",T).")
    for _, solution := range valor_pago {
      Valor,_ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("T")), 64)
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
    fmt.Println(periodo)


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
