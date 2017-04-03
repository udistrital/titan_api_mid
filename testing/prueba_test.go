package testing

import "testing"
import 	"github.com/udistrital/titan_api_mid/models"
import "time"
import 	"github.com/udistrital/titan_api_mid/golog"
import "fmt"

func TestReglas(t *testing.T) {

    var arreglo_funcionarios []models.FuncionarioInfoPruebas
    arreglo_funcionarios = make([]models.FuncionarioInfoPruebas, 1)
  //  var   InformacionCargo []models.FuncionarioCargo
  //  InformacionCargo = make([]models.FuncionarioCargo, 1)
    var resultado []models.Respuesta
    var reglas []string
    var conceptos *[]models.ConceptosResumen

    arreglo_funcionarios[0].InformacionCargo = make([]models.FuncionarioCargo, 1)
    arreglo_funcionarios[0].InformacionCargo[0].Id = 35
    arreglo_funcionarios[0].InformacionCargo[0].Asignacion_basica =   4114674
    arreglo_funcionarios[0].InformacionCargo[0].FechaInicio    = time.Date(1995, time.Month(1), 29, 0, 0, 0, 0, time.UTC)

    reglas = file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/reglas.txt")

    arreglo_funcionarios[0].Reglas = processString(reglas)
    arreglo_funcionarios[0].FechaPreliquidacion = time.Date(2017, time.Month(3), 9, 0, 0, 0, 0, time.UTC)

    arreglo_funcionarios[0].Valor_correcto_salario = "1920181"
    arreglo_funcionarios[0].IdProveedor = 29
    arreglo_funcionarios[0].Dias_laborados = 7980
    arreglo_funcionarios[0].Periodo = "2017"
  	arreglo_funcionarios[0].EsAnual  = 0
  	arreglo_funcionarios[0].PorcentajePT = 40
  	arreglo_funcionarios[0].TipoNomina = "2"

    fmt.Println("Inicio test")
    for x:=0; x < len(arreglo_funcionarios) ; x++ {


      resultado = golog.CargarReglasFP(arreglo_funcionarios[x].FechaPreliquidacion, arreglo_funcionarios[x].Reglas, arreglo_funcionarios[x].IdProveedor,arreglo_funcionarios[x].InformacionCargo , arreglo_funcionarios[x].Dias_laborados,
      arreglo_funcionarios[x].Periodo, 	arreglo_funcionarios[x].EsAnual , 	arreglo_funcionarios[x].PorcentajePT ,	arreglo_funcionarios[x].TipoNomina)
      conceptos = resultado[0].Conceptos
      for i, descuentos := range *conceptos {
          if(i == 0){
            if descuentos.Valor != arreglo_funcionarios[x].Valor_correcto_salario {
               t.Errorf("Los datos son incorrectos, se obtuvo: ", descuentos.Valor, "y era: ", arreglo_funcionarios[x].Valor_correcto_salario)
            }
          }

        }


    }








  }
