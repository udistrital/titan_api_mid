package golog
import (

	"time"
	"github.com/udistrital/titan_api_mid/models"
	"io"
	"os"
	"strings"
	"strconv"
	"github.com/astaxie/beego"
	"encoding/json"
	"net/http"
	"fmt"
)

func CalcularDiasNovedades(FechaPreliq time.Time, AnoDesde float64, MesDesde float64, DiaDesde float64, AnoHasta float64, MesHasta float64, DiaHasta float64) (dias_liquidar float64) {
	var FechaDesde time.Time
	var FechaHasta time.Time
	var FechaControl time.Time
	var periodo_liquidacion float64

	FechaDesde = time.Date(int(AnoDesde), time.Month(int(MesDesde)), int(DiaDesde), 0, 0, 0, 0, time.UTC)
	FechaHasta = time.Date(int(AnoHasta), time.Month(int(MesHasta)), int(DiaHasta), 0, 0, 0, 0, time.UTC)

	if FechaDesde.Month() == FechaPreliq.Month() && FechaDesde.Year() == FechaPreliq.Year() {
		FechaControl = time.Date(FechaPreliq.Year(), FechaPreliq.Month(), 30, 0, 0, 0, 0, time.UTC)
		periodo_liquidacion = CalcularDias(FechaDesde, FechaControl) + 1
	} else if FechaHasta.Month() == FechaPreliq.Month() && FechaHasta.Year() == FechaPreliq.Year() {
		FechaControl = time.Date(FechaPreliq.Year(), FechaPreliq.Month(), 1, 0, 0, 0, 0, time.UTC)
		periodo_liquidacion = CalcularDias(FechaControl, FechaHasta) + 1
	} else if FechaHasta.Month() == FechaDesde.Month() && FechaHasta.Year() == FechaDesde.Year() {
		periodo_liquidacion = CalcularDias(FechaDesde, FechaHasta) + 1
	}

	return periodo_liquidacion

}

func CalcularDias(FechaInicio time.Time, FechaFin time.Time) (dias_laborados float64) {
	var a, m, d int
	var meses_contrato float64
	var dias_contrato float64
	if FechaFin.IsZero() {
		var FechaFin2 time.Time
		FechaFin2 = time.Now()
		a, m, d = diff(FechaInicio, FechaFin2)
		meses_contrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
		dias_contrato = meses_contrato * 30

	} else {
		a, m, d = diff(FechaInicio, FechaFin)
		meses_contrato = (float64(a * 12)) + float64(m) + (float64(d) / 30)
		dias_contrato = meses_contrato * 30

	}

	return dias_contrato

}

func diff(a, b time.Time) (year, month, day int) {
    if a.Location() != b.Location() {
        b = b.In(a.Location())
    }
    if a.After(b) {
        a, b = b, a
    }
		 oneDay := time.Hour * 5
		 a = a.Add(oneDay)
		 b = b.Add(oneDay)
    y1, M1, d1 := a.Date()
    y2, M2, d2 := b.Date()



    year = int(y2 - y1)
    month = int(M2 - M1)
    day = int(d2 - d1)


    // Normalize negative values
		/*if day < 0{
			day = 0
		}
		if month < 0 {
        month = 0
    }*/
    if day < 0 {
        // days in month:
        t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
        day += 32 - t.Day()
        month--
    }
    if month < 0 {
        month += 12
        year--
    }

    return
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

func ConsultarValoresBonServPS(fechaPreliquidacion time.Time, idPersona int, codigo_concepto string, periodo string ) (valor_con int64){


	mes_preliquidacion := int(fechaPreliquidacion.Month())
	ano_preliquidacion := int(fechaPreliquidacion.Year())
	ano_preliquidacion_string := strconv.Itoa(ano_preliquidacion)
	dia_preliquidacion_string := strconv.Itoa(int(fechaPreliquidacion.Day()))
	mes_preliquidacion_string := strconv.Itoa(mes_preliquidacion)
	ano_busqueda := ano_preliquidacion - 1
	ano_busqueda_string := strconv.Itoa(ano_busqueda)
	var valor_concepto []models.DetalleLiquidacion
	var valor int64
	var id_persona_string string = strconv.Itoa(idPersona)
			//http://localhost:8082/v1/detalle_liquidacion?limit=-1&query=Liquidacion.FechaLiquidacion__gte:2016-05-30,Liquidacion.FechaLiquidacion__lte:2017-06-30,Concepto.Id:1195,Persona:29
		if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_liquidacion?limit=-1&query=Liquidacion.FechaLiquidacion__gte:"+ano_busqueda_string+"-05-30,Liquidacion.FechaLiquidacion__lte:"+ano_preliquidacion_string+"-"+mes_preliquidacion_string+"-"+dia_preliquidacion_string+",Concepto.Id:"+codigo_concepto+",Persona:"+id_persona_string+",TipoLiquidacion:2", &valor_concepto); err == nil {
			for _, solution := range valor_concepto {
		 	valor = valor + solution.ValorCalculado
		 }
		}else{
			fmt.Println(err)
		}
		return valor

}

func ConsultarValoresBonServDic(fechaPreliquidacion time.Time, idPersona int, codigo_concepto string, periodo string ) (valor_con int64){

	periodo_nomina := periodo
	var valor_concepto []models.DetalleLiquidacion
	var valor int64
	var id_persona_string string = strconv.Itoa(idPersona)
	//AGREGAR TIPO DE LIQUIDACION!! porque habran varios 29, y se necesita el pagado en nomina 2
		if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_liquidacion?limit=-1&query=Liquidacion.Nomina.Periodo:"+periodo_nomina+",Concepto.Id:"+codigo_concepto+",Persona:"+id_persona_string+",TipoLiquidacion:2", &valor_concepto); err == nil {
			for _, solution := range valor_concepto {
		 	valor = valor + solution.ValorCalculado
		 }
	}
		//http://localhost:8082/v1/detalle_liquidacion?limit=-1&query=Liquidacion.Nomina.Periodo:2017,Persona:29,TipoLiquidacion:3 <-- CONSULTA DOCEAVA PRIMA SEMESTRAL
		//nuevaRegla = "bonificacion_servicio(bonServ,1540945)."
		//hacer consulta de conceptos con codigo 129,139,1195 que se le hayan pagado a la persona en el presente año y se crea este hecho
		return valor

}

func ConsultarValoresPriServDic(fechaPreliquidacion time.Time, idPersona int, periodo string ) (valor_con int64){

	periodo_nomina := periodo
	var valor_concepto []models.DetalleLiquidacion
	var valor int64
	var id_persona_string string = strconv.Itoa(idPersona)
	//AGREGAR TIPO DE LIQUIDACION!! porque habran varios 29, y se necesita el pagado en nomina 2

		if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_liquidacion?limit=-1&query=Liquidacion.Nomina.Periodo:"+periodo_nomina+",Persona:"+id_persona_string+",TipoLiquidacion:3", &valor_concepto); err == nil {
			for _, solution := range valor_concepto {
		 	valor = valor + solution.ValorCalculado
		 }
	}
		//http://localhost:8082/v1/detalle_liquidacion?limit=-1&query=Liquidacion.Nomina.Periodo:2017,Persona:29,TipoLiquidacion:3 <-- CONSULTA DOCEAVA PRIMA SEMESTRAL
		//nuevaRegla = "bonificacion_servicio(bonServ,1540945)."
		//hacer consulta de conceptos con codigo 129,139,1195 que se le hayan pagado a la persona en el presente año y se crea este hecho
		return valor

}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
