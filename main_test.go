package main

import (
  "testing"
  "github.com/astaxie/beego"
  "time"
)

// test something real here
func TestDev(t *testing.T) {
  beego.BConfig.RunMode = "dev"
  go func() {
    main()
  }()
  time.Sleep(5 * time.Second)
}
