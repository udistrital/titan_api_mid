package testing

import (

	"fmt"
	"os"
	"bufio"

)

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
    reglas_temp = reglas_temp + reglas[i]
  }

  return reglas_temp
}
