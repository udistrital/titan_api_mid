package test

import (
	"testing"

	"github.com/astaxie/beego/context"
	"github.com/udistrital/titan_api_mid/controllers"
)

func TestURLMapping(t *testing.T) {
	gpapc := controllers.GestionPersonasAPreliquidarController{}
	gpapc.Init(context.NewContext(), "test", "test", nil)
	gpapc.URLMapping()
}
