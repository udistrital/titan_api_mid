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
import "io"
import "strings"

func TestDocentesPlanta(t *testing.T) {

    //pasar aqui JSON proveniente de archivo

    var resultado []models.Respuesta
    var reglas []string
    var conceptos *[]models.ConceptosResumen
    var nombre_archivo string

    var docentes_planta_a_probar []string
    var docentes_planta string
    var reporte string


    docentes_planta_a_probar =  file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/pruebaDocentesPlanta20174.txt")
    docentes_planta = processString(docentes_planta_a_probar)

    reporte = "Mes de mayo de 2017 - Docentes planta \n"
    b := []byte(docentes_planta)

    var arreglo_docentes_planta []models.PruebaGoDocentes
    err := json.Unmarshal(b, &arreglo_docentes_planta)
    fmt.Println(err)

    if err == nil {

       fmt.Println("Inicio test docentes")
       for x:=0; x < len(arreglo_docentes_planta) ; x++ {
        nombre_archivo = "reglas"
         nombre_archivo = nombre_archivo + strconv.Itoa(arreglo_docentes_planta[x].IdProveedor) +".txt"
         reglas = file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/"+nombre_archivo+"")
         arreglo_docentes_planta[x].Reglas = processString(reglas)
         puntos := strconv.FormatFloat(arreglo_docentes_planta[x].InformacionCargo[0].Puntos, 'f', 6, 64)
 				 regimen := arreglo_docentes_planta[x].InformacionCargo[0].Regimen


         resultado = golog.CargarReglasDP(arreglo_docentes_planta[x].Mes, arreglo_docentes_planta[x].Ano,arreglo_docentes_planta[x].Dias_laborados,arreglo_docentes_planta[x].IdProveedor,"",0,arreglo_docentes_planta[x].Reglas,arreglo_docentes_planta[x].InformacionCargo ,
         puntos, regimen,strconv.Itoa(arreglo_docentes_planta[x].TipoNomina))

         conceptos = resultado[0].Conceptos
         reporte = reporte + "--------------------------------------------------------\n"
         reporte = reporte + strconv.Itoa(arreglo_docentes_planta[x].NumDocumento) + "\n"
         for _, descuentos := range *conceptos {
             if(descuentos.Nombre == "salarioBase"){
               reporte = reporte + "salarioBase \n"
               if descuentos.Valor != arreglo_docentes_planta[x].Valor_correcto_salario {
                 fmt.Print("Test docentes_planta: ")
                  t.Errorf("Los datos son incorrectos para valor salario de funcionario "+strconv.Itoa(arreglo_docentes_planta[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_docentes_planta[x].Valor_correcto_salario)
               }
                 reporte = reporte + "Titan: " + descuentos.Valor + " NOMINAOAS: "+arreglo_docentes_planta[x].Valor_correcto_salario+"\n"

             }


             if(descuentos.Nombre == "salud"){
               reporte = reporte + "Salud \n"
               if descuentos.Valor != arreglo_docentes_planta[x].Valor_correcto_Salud {
                 fmt.Print("Test docentes_planta: ")
                  t.Errorf("Los datos son incorrectos para valor de salud de funcionario "+strconv.Itoa(arreglo_docentes_planta[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_docentes_planta[x].Valor_correcto_Salud)
               }
                 reporte = reporte + "Titan: " + descuentos.Valor + " NOMINAOAS: "+arreglo_docentes_planta[x].Valor_correcto_Salud+"\n"

             }

             if(descuentos.Nombre == "pension"){
               reporte = reporte + "Pension \n"
               if descuentos.Valor != arreglo_docentes_planta[x].Valor_correcto_Pension {
                 fmt.Print("Test docentes_planta: ")
                  t.Errorf("Los datos son incorrectos para valor de pension de funcionario "+strconv.Itoa(arreglo_docentes_planta[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_docentes_planta[x].Valor_correcto_Pension)
               }
                 reporte = reporte + "Titan: " + descuentos.Valor + " NOMINAOAS: "+arreglo_docentes_planta[x].Valor_correcto_Pension+"\n"

             }

           }


       }
    }

    str := fmt.Sprintf("%s", reporte)
    if err := WriteStringToFile("docentes_plantaReporte20175.txt", str); err != nil {
        panic(err)
    }

  }

/*
func TestFuncionarios(t *testing.T) {

    //pasar aqui JSON proveniente de archivo

    var resultado []models.Respuesta
    var reglas []string
    var conceptos *[]models.ConceptosResumen
    var nombre_archivo string

    var funcionarios_a_probar []string
    var funcionarios string
    var reporte string


    funcionarios_a_probar =  file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/pruebaFuncionarios20175.txt")
    funcionarios = processString(funcionarios_a_probar)

    reporte = "Mes de mayo de 2017 - Admin de planta \n"
    b := []byte(funcionarios)

    var arreglo_funcionarios []models.PruebaGo
    err := json.Unmarshal(b, &arreglo_funcionarios)
    fmt.Println(err)

    if err == nil {

       fmt.Println("Inicio test funcionarios")
       for x:=0; x < len(arreglo_funcionarios) ; x++ {
        nombre_archivo = "reglas"
         nombre_archivo = nombre_archivo + strconv.Itoa(arreglo_funcionarios[x].IdProveedor) +".txt"
         reglas = file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/"+nombre_archivo+"")
         arreglo_funcionarios[x].Reglas = processString(reglas)
         resultado = golog.CargarReglasFP(arreglo_funcionarios[x].Mes, arreglo_funcionarios[x].Ano,arreglo_funcionarios[x].Reglas, arreglo_funcionarios[x].IdProveedor,"",0,arreglo_funcionarios[x].InformacionCargo , arreglo_funcionarios[x].Dias_laborados,
         arreglo_funcionarios[x].EsAnual , 	arreglo_funcionarios[x].PorcentajePT ,	arreglo_funcionarios[x].TipoNomina)

         conceptos = resultado[0].Conceptos
         reporte = reporte + "--------------------------------------------------------\n"
         reporte = reporte + strconv.Itoa(arreglo_funcionarios[x].NumDocumento) + "\n"
         for _, descuentos := range *conceptos {
             if(descuentos.Nombre == "salarioBase"){
               reporte = reporte + "salarioBase \n"
               if descuentos.Valor != arreglo_funcionarios[x].Valor_correcto_salario {
                 fmt.Print("Test funcionarios: ")
                  t.Errorf("Los datos son incorrectos para valor salario de funcionario "+strconv.Itoa(arreglo_funcionarios[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_funcionarios[x].Valor_correcto_salario)
               }
                 reporte = reporte + "Titan: " + descuentos.Valor + " NOMINAOAS: "+arreglo_funcionarios[x].Valor_correcto_salario+"\n"

             }

             if(descuentos.Nombre == "primaAnt"){
               reporte = reporte + "Prima antiguedad \n"
               if descuentos.Valor != arreglo_funcionarios[x].Valor_correcto_PrimaAnt {
                 fmt.Print("Test funcionarios: ")
                  t.Errorf("Los datos son incorrectos para valor de prima antiguedad de funcionario "+strconv.Itoa(arreglo_funcionarios[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_funcionarios[x].Valor_correcto_PrimaAnt)
               }
                 reporte = reporte + "Titan: " + descuentos.Valor + " NOMINAOAS: "+arreglo_funcionarios[x].Valor_correcto_PrimaAnt+"\n"

             }

             if(descuentos.Nombre == "priTec"){
               reporte = reporte + "Prima técnica \n"
               if descuentos.Valor != arreglo_funcionarios[x].Valor_correcto_PrimaTecnica {
                 fmt.Print("Test funcionarios: ")
                  t.Errorf("Los datos son incorrectos para valor de prima técnica de funcionario "+strconv.Itoa(arreglo_funcionarios[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_funcionarios[x].Valor_correcto_PrimaTecnica)
               }
                 reporte = reporte + "Titan: " + descuentos.Valor + " NOMINAOAS: "+arreglo_funcionarios[x].Valor_correcto_PrimaTecnica+"\n"

             }

             if(descuentos.Nombre == "salud"){
               reporte = reporte + "Salud \n"
               if descuentos.Valor != arreglo_funcionarios[x].Valor_correcto_Salud {
                 fmt.Print("Test funcionarios: ")
                  t.Errorf("Los datos son incorrectos para valor de salud de funcionario "+strconv.Itoa(arreglo_funcionarios[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_funcionarios[x].Valor_correcto_Salud)
               }
                 reporte = reporte + "Titan: " + descuentos.Valor + " NOMINAOAS: "+arreglo_funcionarios[x].Valor_correcto_Salud+"\n"

             }

             if(descuentos.Nombre == "pension"){
               reporte = reporte + "Pension \n"
               if descuentos.Valor != arreglo_funcionarios[x].Valor_correcto_Pension {
                 fmt.Print("Test funcionarios: ")
                  t.Errorf("Los datos son incorrectos para valor de pension de funcionario "+strconv.Itoa(arreglo_funcionarios[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_funcionarios[x].Valor_correcto_Pension)
               }
                 reporte = reporte + "Titan: " + descuentos.Valor + " NOMINAOAS: "+arreglo_funcionarios[x].Valor_correcto_Pension+"\n"

             }

           }


       }
    }

    str := fmt.Sprintf("%s", reporte)
    if err := WriteStringToFile("FuncionariosReporte20175.txt", str); err != nil {
        panic(err)
    }

  }



  func TestContratistas(e *testing.T) {


    var resultado []models.Respuesta
    var reglas []string
    var conceptos *[]models.ConceptosResumen
    var nombre_archivo string

   var contratistas_a_probar []string
    var contratistas string
    var reporte string


    contratistas_a_probar =  file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/pruebaContratistas201610.txt")
    contratistas = processString(contratistas_a_probar)
    reporte = "Mes de octubre de 2016 - Contratistas \n"
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
         reporte = reporte + "--------------------------------------------------------\n"
         reporte = reporte + strconv.Itoa(arreglo_contratistas[x].NumDocumento) + "\n"
         for _, descuentos := range *conceptos {
             if(descuentos.Nombre == "pagoBruto"){
               reporte = reporte + "pago bruto \n"
               if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_salario {
                 fmt.Print("Test funcionarios: ")
                  e.Errorf("Los datos son incorrectos para valor salario de funcionario "+strconv.Itoa(arreglo_contratistas[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_salario)
               }
                 reporte = reporte + "Titan: " + descuentos.Valor + " Excel: "+arreglo_contratistas[x].Valor_correcto_salario+"\n"


             }

            if(descuentos.Nombre == "reteIca"){
              reporte = reporte + "Reteica \n"
               if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_Reteica {
                 fmt.Print("Test funcionarios: ")
                  e.Errorf("Los datos son incorrectos para descuento reteica de funcionario "+strconv.Itoa(arreglo_contratistas[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_Reteica)
               }
                 reporte = reporte + " Titan: " + descuentos.Valor + " Excel: "+arreglo_contratistas[x].Valor_correcto_Reteica+"\n"

             }

             if(descuentos.Nombre == "estampillaUD"){
               reporte = reporte + "Estampilal UD \n"
                if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_EstampillaUD {
                  fmt.Print("Test funcionarios: ")
                   e.Errorf("Los datos son incorrectos para descuento Estampilla de funcionario "+strconv.Itoa(arreglo_contratistas[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_EstampillaUD)
                }
                  reporte = reporte + " Titan: " + descuentos.Valor + " Excel: "+arreglo_contratistas[x].Valor_correcto_EstampillaUD+"\n"

              }


              if(descuentos.Nombre == "proCultura"){
                reporte = reporte + "proCultura \n"
                 if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_ProCultura {
                   fmt.Print("Test funcionarios: ")
                    e.Errorf("Los datos son incorrectos para descuento ProCultura de funcionario "+strconv.Itoa(arreglo_contratistas[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_ProCultura)
                 }
                   reporte = reporte + "Titan: " + descuentos.Valor + " Excel: "+arreglo_contratistas[x].Valor_correcto_ProCultura+"\n"

               }

               if(descuentos.Nombre == "adultoMayor"){
                  reporte = reporte + "Adulto mayor \n"
                  if descuentos.Valor != arreglo_contratistas[x].Valor_correcto_AdultoMayor {
                    fmt.Print("Test funcionarios: ")
                     e.Errorf("Los datos son incorrectos para descuento AdultoMayor de funcionario "+strconv.Itoa(arreglo_contratistas[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_contratistas[x].Valor_correcto_AdultoMayor)
                  }
                    reporte = reporte + " Titan: " + descuentos.Valor + " Excel: "+arreglo_contratistas[x].Valor_correcto_AdultoMayor+"\n"

                }
           }


       }
    }

    str := fmt.Sprintf("%s", reporte)
  	if err := WriteStringToFile("ContratistasReporte102017.txt", str); err != nil {
  			panic(err)
  	}
  }
*/

func TestHC(e *testing.T) {


  var resultado []models.Respuesta
  var reglas []string
  var conceptos *[]models.ConceptosResumen
  var nombre_archivo string

 var docentes_a_probar []string
  var docentes string
  var reporte string


  docentes_a_probar =  file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/pruebaHC20175.txt")
  docentes = processString(docentes_a_probar)
  reporte = "Mes de mayo de 2017 - Docentes salarios \n"
  b := []byte(docentes)

  var arreglo_docentes []models.PruebaGo
  err := json.Unmarshal(b, &arreglo_docentes)
  fmt.Println(err)

  if err == nil {
    fmt.Println(arreglo_docentes)
     fmt.Println("Inicio test docentes salarios")
     for x:=0; x < len(arreglo_docentes) ; x++ {
       nombre_archivo = "reglas"
       nombre_archivo = nombre_archivo + strconv.Itoa(arreglo_docentes[x].IdProveedor) +".txt"
       reglas = file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/"+nombre_archivo+"")
       arreglo_docentes[x].Reglas = processString(reglas)


       resultado = golog.CargarReglas(arreglo_docentes[x].IdProveedor, arreglo_docentes[x].Reglas, strconv.Itoa(arreglo_docentes[x].Ano))

       conceptos = resultado[0].Conceptos
       reporte = reporte + "--------------------------------------------------------\n"
       reporte = reporte + strconv.Itoa(arreglo_docentes[x].NumDocumento) + "\n"
       for _, descuentos := range *conceptos {
           if(descuentos.Nombre == "pagoBruto"){
             reporte = reporte + "Pago bruto \n"
             if descuentos.Valor != arreglo_docentes[x].Valor_correcto_salario {
               fmt.Print("Test funcionarios: ")
                e.Errorf("Los datos son incorrectos para valor salario de funcionario "+strconv.Itoa(arreglo_docentes[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_docentes[x].Valor_correcto_salario)
             }
               reporte = reporte + "Titan: " + descuentos.Valor + " Excel: "+arreglo_docentes[x].Valor_correcto_salario+"\n"


           }

            if(descuentos.Nombre == "salud"){
              reporte = reporte + "Salud \n"
               if descuentos.Valor != arreglo_docentes[x].Valor_correcto_Salud {
                 fmt.Print("Test funcionarios: ")
                  e.Errorf("Los datos son incorrectos para descuento Salud de funcionario "+strconv.Itoa(arreglo_docentes[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_docentes[x].Valor_correcto_Salud)
               }
                 reporte = reporte + "Titan: " + descuentos.Valor + " Excel: "+arreglo_docentes[x].Valor_correcto_Salud+"\n"

             }

             if(descuentos.Nombre == "pension"){
                reporte = reporte + "Pension \n"
                if descuentos.Valor != arreglo_docentes[x].Valor_correcto_Pension {
                  fmt.Print("Test funcionarios: ")
                   e.Errorf("Los datos son incorrectos para descuento Pension de funcionario "+strconv.Itoa(arreglo_docentes[x].NumDocumento)+", se obtuvo: "+descuentos.Valor+" y era: "+arreglo_docentes[x].Valor_correcto_Pension)
                }
                  reporte = reporte + " Titan: " + descuentos.Valor + " Excel: "+arreglo_docentes[x].Valor_correcto_Pension+"\n"

              }
         }


     }
  }

  str := fmt.Sprintf("%s", reporte)
  if err := WriteStringToFile("DocentesSalariosReporte20175.txt", str); err != nil {
      panic(err)
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

  func WriteStringToFile(filepath, s string) error {
  	fo, err := os.Create(filepath)
  	if err != nil {
  		return err
  	}
  	defer fo.Close()

  	_, err = io.Copy(fo, strings.NewReader(s))
  	if err != nil {
  		return err
  	}

  	return nil
  }


  func processString(reglas []string)(reglas_t string){
    var reglas_temp string = ""
    for i:= 0 ; i < len(reglas) ; i++ {
      reglas_temp = reglas_temp + reglas[i] + "\n"
    }

    return reglas_temp
  }
