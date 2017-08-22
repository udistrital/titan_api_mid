package testing

import "testing"
import 	"github.com/udistrital/titan_api_mid/models"
//import "time"
import 	"github.com/udistrital/titan_api_mid/golog"
import "fmt"
import 	"encoding/json"
import 	"strconv"
import "os"
import "bufio"

/*
func TestFuncionarios(t *testing.T) {

    //pasar aqui JSON proveniente de archivo

    var resultado []models.Respuesta
    var reglas []string
    var conceptos *[]models.ConceptosResumen
    var nombre_archivo string

   var funcionarios_a_probar []string
    var funcionarios string


    funcionarios_a_probar =  file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/prueba.txt")
    funcionarios = processString(funcionarios_a_probar)

    b := []byte(funcionarios)

    var arreglo_funcionarios []models.FuncionarioInfoPruebas
    err := json.Unmarshal(b, &arreglo_funcionarios)
    fmt.Println(err)

    if err == nil {
      fmt.Println(arreglo_funcionarios)
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
                  t.Errorf("Los datos son incorrectos para salario, se obtuvo: "+descuentos.Valor+" y era: "+arreglo_funcionarios[x].Valor_correcto_salario)
               }
             }

           }


       }
    }


  }
*/

  func TestContratistas(e *testing.T) {


    var resultado []models.Respuesta
    var reglas []string
    var conceptos *[]models.ConceptosResumen
    var nombre_archivo string

   var contratistas_a_probar []string
    var contratistas string


    contratistas_a_probar =  file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/pruebaContratistas10.txt")
    contratistas = processString(contratistas_a_probar)

    b := []byte(contratistas)

    var arreglo_contratistas []models.PruebaGo
    err := json.Unmarshal(b, &arreglo_contratistas)
    fmt.Println(err)

    if err == nil {
      fmt.Println(arreglo_contratistas)
       fmt.Println("Inicio test contratistas")
       for x:=0; x < len(arreglo_contratistas) ; x++ {
         nombre_archivo = "reglas"
         nombre_archivo = nombre_archivo + strconv.Itoa(arreglo_contratistas[x].IdProveedor) +".txt"
         reglas = file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/"+nombre_archivo+"")
         arreglo_contratistas[x].Reglas = processString(reglas)
         resultado = golog.CargarReglasCT(arreglo_contratistas[x].IdProveedor, arreglo_contratistas[x].Reglas, strconv.Itoa(arreglo_contratistas[x].Ano))
         conceptos = resultado[0].Conceptos
         fmt.Println(conceptos)
         for _, descuentos := range *conceptos {
             if(descuentos.Nombre == "pagoBruto"){
               if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_salario {
                 fmt.Print("Test funcionarios: ")
                  e.Errorf("Los datos son incorrectos para valor salario de funcionario "+strconv.Itoa(arreglo_contratistas[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_salario)
               }
             }

            if(descuentos.Nombre == "reteIca"){
               if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_Reteica {
                 fmt.Print("Test funcionarios: ")
                  e.Errorf("Los datos son incorrectos para descuento reteica de funcionario "+strconv.Itoa(arreglo_contratistas[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_Reteica)
               }
             }

             if(descuentos.Nombre == "estampillaUD"){
                if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_EstampillaUD {
                  fmt.Print("Test funcionarios: ")
                   e.Errorf("Los datos son incorrectos para descuento Estampilla de funcionario "+strconv.Itoa(arreglo_contratistas[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_EstampillaUD)
                }
              }


              if(descuentos.Nombre == "proCultura"){
                 if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_ProCultura {
                   fmt.Print("Test funcionarios: ")
                    e.Errorf("Los datos son incorrectos para descuento ProCultura de funcionario "+strconv.Itoa(arreglo_contratistas[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_ProCultura)
                 }
               }

               if(descuentos.Nombre == "adultoMayor"){
                  if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_AdultoMayor {
                    fmt.Print("Test funcionarios: ")
                     e.Errorf("Los datos son incorrectos para descuento AdultoMayor de funcionario "+strconv.Itoa(arreglo_contratistas[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_AdultoMayor)
                  }
                }
           }


       }
    }

  }

  func file2lines(filePath string) []string {
        f, err := os.Open(filePath)
        if err != nil {
                panic(err)
        }
        defer f.Close()

        var lines []string
        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
                lines = append(lines, scanner.Text())
        }
        if err := scanner.Err(); err != nil {
                fmt.Fprintln(os.Stderr, err)
        }

        return lines
  }

  func processString(reglas []string)(reglas_t string){
    var reglas_temp string = ""
    for i:= 0 ; i < len(reglas) ; i++ {
      reglas_temp = reglas_temp + reglas[i] + "\n"
    }

    return reglas_temp
  }
