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

func TestFuncionarios(t *testing.T) {

    //pasar aqui JSON proveniente de archivo

    var resultado []models.Respuesta
    var reglas []string
    var conceptos *[]models.ConceptosResumen
    var nombre_archivo string

   var funcionarios_a_probar []string
    var funcionarios string


    funcionarios_a_probar =  file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/json_funcionarios.txt")
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
                  t.Errorf("Los datos son incorrectos, se obtuvo: "+descuentos.Valor+" y era: "+arreglo_funcionarios[x].Valor_correcto_salario)
               }
             }

           }


       }
    }


  }


  func TestContratistas(e *testing.T) {


    var resultado []models.Respuesta
    var reglas []string
    var conceptos *[]models.ConceptosResumen
    var nombre_archivo string

   var contratistas_a_probar []string
    var contratistas string


    contratistas_a_probar =  file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/json_contratistas.txt")
    contratistas = processString(contratistas_a_probar)

    b := []byte(contratistas)

    var arreglo_contratistas []models.ContratistasInfoPruebas
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
         resultado = golog.CargarReglasCT(arreglo_contratistas[x].IdProveedor, arreglo_contratistas[x].Reglas, arreglo_contratistas[x].Periodo)
         conceptos = resultado[0].Conceptos
         for i, descuentos := range *conceptos {
             if(i == 0){
               if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_salario {
                 fmt.Print("Test funcionarios: ")
                  e.Errorf("Los datos son incorrectos, se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_salario)
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
