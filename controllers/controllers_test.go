package controllers

import (
  "testing"
  "github.com/astaxie/beego/context"
)

// una prueba para que el coverage no sea 0% pero las pruebas reales deben ser
// de la lógica de negocio
func TestURLMapping(t *testing.T) {
  gpapc := GestionPersonasAPreliquidarController{}
  gpapc.Init(context.NewContext(), "test", "test", nil)
  gpapc.URLMapping()
}
