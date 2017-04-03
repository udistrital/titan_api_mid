package golog
import (

	"time"

	"io"
	"os"
	"strings"

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
