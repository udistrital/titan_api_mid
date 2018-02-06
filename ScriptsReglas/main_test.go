// no sé muy bien que hay que probar aquí
package main

import "testing"

// este test está básicamente para que no esté en 0% el coverage, se necesita
// este paquete en este repo?
func Test1(t *testing.T) {
  defer func() {
    if fail := recover(); fail != nil {
      t.Fatal("main hizo panic, pero no lo debío haber hecho")
    }
  }()
  main()
}
