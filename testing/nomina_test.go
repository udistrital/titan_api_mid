package testing

import "testing"
import 	"github.com/udistrital/titan_api_mid/models"
//import "time"
import 	"github.com/udistrital/titan_api_mid/golog"
import "fmt"
import 	"encoding/json"
import 	"strconv"

func TestFuncionarios(t *testing.T) {

    //pasar aqui JSON proveniente de archivo

    var resultado []models.Respuesta
    var reglas []string
    var conceptos *[]models.ConceptosResumen
    var nombre_archivo string

    var ejemplo_json string
    ejemplo_json = `[
  {

    "InformacionCargo":[{
        "Id": 35,
        "Asignacion_basica": 4114674,
        "FechaInicio": "1995-01-24T19:00:00-05:00"
     }],
    "Reglas":"reglitas",
    "FechaPreliquidacion": "2017-03-09T19:00:00-05:00",
    "Valor_correcto_salario": "1920180",
    "IdProveedor": 29,
    "Dias_laborados": 7980,
    "Periodo": "2017",
    "EsAnual": 0,
    "PorcentajePT" : 40,
    "TipoNomina": "2"

 },
 {

    "InformacionCargo":[{
        "Id": 95,
        "Asignacion_basica": 2193636,
        "FechaInicio": "1995-04-20T19:00:00-05:00"
     }],
    "Reglas":"reglitas",
    "FechaPreliquidacion": "2017-03-09T19:00:00-05:00",
    "Valor_correcto_salario": "1023698",
    "IdProveedor": 2,
    "Dias_laborados": 6462,
    "Periodo": "2017",
    "EsAnual": 0,
    "PorcentajePT" : 0,
    "TipoNomina": "2"

 }
]
`
    b := []byte(ejemplo_json)

    var arreglo_funcionarios []models.FuncionarioInfoPruebas
    err := json.Unmarshal(b, &arreglo_funcionarios)

    if err == nil {

       fmt.Println("Inicio test funcionarios")
       for x:=0; x < len(arreglo_funcionarios) ; x++ {
        nombre_archivo = "reglas"
         nombre_archivo = nombre_archivo + strconv.Itoa(arreglo_funcionarios[x].IdProveedor) +".txt"
         reglas = file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/"+nombre_archivo+"")
         arreglo_funcionarios[x].Reglas = processString(reglas)
         resultado = golog.CargarReglasFP(arreglo_funcionarios[x].FechaPreliquidacion, arreglo_funcionarios[x].Reglas, arreglo_funcionarios[x].IdProveedor,arreglo_funcionarios[x].InformacionCargo , arreglo_funcionarios[x].Dias_laborados,
         arreglo_funcionarios[x].Periodo, 	arreglo_funcionarios[x].EsAnual , 	arreglo_funcionarios[x].PorcentajePT ,	arreglo_funcionarios[x].TipoNomina)
         conceptos = resultado[0].Conceptos
         for i, descuentos := range *conceptos {
             if(i == 0){
               if descuentos.Valor != arreglo_funcionarios[x].Valor_correcto_salario {
                 fmt.Print("Test funcionarios: ")
                  t.Errorf("Los datos son incorrectos, se obtuvo: ", descuentos.Valor, "y era: ", arreglo_funcionarios[x].Valor_correcto_salario)
               }
             }

           }


       }
    }


  }

/*
  func TestContratistas(e *testing.T) {
    var resultado []models.Respuesta
    var reglas_arreglo []string
    var reglas string
    var conceptos *[]models.ConceptosResumen

    reglas_arreglo = file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/reglascontratistas.txt")
    reglas = processString(reglas_arreglo)
    fmt.Println("Inicio test contratistas")

      resultado = golog.CargarReglasCT(reglas,"2016")
      conceptos = resultado[0].Conceptos
      for i, descuentos := range *conceptos {
          if(i == 0){
            if descuentos.Valor != "689455" {
               e.Errorf("El valor es incorrecto, se obtuvo:"+descuentos.Valor+" se deseaba:689455")
            }
          }

        }

  }
*/
