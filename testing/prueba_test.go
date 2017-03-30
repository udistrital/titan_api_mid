package testing

import "testing"
import 	"github.com/udistrital/titan_api_mid/models"
import "time"
import 	"github.com/udistrital/titan_api_mid/golog"
import "fmt"

func TestReglas(t *testing.T) {
  	var resultado []models.Respuesta
    var informacion_cargo []models.FuncionarioCargo
    var reglas []string
    var fechaPreliquidacion time.Time
    var conceptos *[]models.ConceptosResumen

    fechaPreliquidacion = time.Date(2017, time.Month(3), 9, 0, 0, 0, 0, time.UTC)
    informacion_cargo = make([]models.FuncionarioCargo, 1)
    reglas = file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/reglas.txt")
    informacion_cargo[0].Id = 35
    informacion_cargo[0].Asignacion_basica =   4114674
    informacion_cargo[0].FechaInicio    = time.Date(1995, time.Month(1), 29, 0, 0, 0, 0, time.UTC)
    reglas_si := processString(reglas)
    resultado = golog.CargarReglasFP(fechaPreliquidacion, reglas_si , 29, informacion_cargo , 7980, "2017", 0, 40, "2")
    valor_correcto_salario := "1920181"
    conceptos = resultado[0].Conceptos
    fmt.Println("Inicio test")
    for i, descuentos := range *conceptos {
        if(i == 0){
          if descuentos.Valor != valor_correcto_salario {
             t.Errorf("Los datos son incorrectos, se obtuvo: ", descuentos.Valor, "y era: ", valor_correcto_salario)
          }
        }

      }

  }
