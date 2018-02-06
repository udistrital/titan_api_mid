package golog

import (
  "testing"
  "fmt"
)

func Test1(t *testing.T) {
  defer func() {
    if fail := recover() ; fail != nil {
      if fmt.Sprintf("%s", fail) != "Undefined predicate: valor_pago/3" {
        t.Fatal("Esperaba cierto error y no lo recibí") // I have no idea what I'm testing here
      }
    }
  }()
  CargarReglasCT(0, "test", "test")
}
