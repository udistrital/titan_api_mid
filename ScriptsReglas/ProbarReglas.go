package main

import (
	"fmt"
	//"strconv"
  "os"
  "bufio"
  "io"
  "strings"
	. "github.com/mndrix/golog"

)

func main() {
  var reglas_a_probar []string
  var reglas string
  var archivo_reglas_a_cargar = "HCSReglas"

  reglas_a_probar =  file2lines("/home/mariaalejandra9404/Documentos/ProyectosGo/src/github.com/udistrital/titan_api_mid/ScriptsReglas/"+archivo_reglas_a_cargar+".txt")
  reglas = processString(reglas_a_probar)
  m := NewMachine().Consult(reglas)
  predicado_a_probar:= "evaluar_uvt(31859,99,X)."

  resultado_prueba := m.ProveAll(predicado_a_probar)

  for _, solution := range resultado_prueba {
   retorno := fmt.Sprintf("%s", solution.ByName_("X"))
   fmt.Println(retorno)
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
