package golog
import (

	"time"
	"github.com/udistrital/titan_api_mid/models"
	. "github.com/udistrital/golog"
	"github.com/udistrital/utils_oas/request"
	"strconv"
	"github.com/astaxie/beego"
	"fmt"
)

func CalcularDiasNovedades(MesPreliq, AnoPreliq int,  AnoDesde float64, MesDesde float64, DiaDesde float64, AnoHasta float64, MesHasta float64, DiaHasta float64) (dias_liquidar float64) {
	var FechaDesde time.Time
	var FechaHasta time.Time
	var FechaControl time.Time
	var periodo_liquidacion float64

	FechaDesde = time.Date(int(AnoDesde), time.Month(int(MesDesde)), int(DiaDesde), 0, 0, 0, 0, time.UTC)
	FechaHasta = time.Date(int(AnoHasta), time.Month(int(MesHasta)), int(DiaHasta), 0, 0, 0, 0, time.UTC)

	esActiva := validarNovedades_segSocial(MesPreliq,AnoPreliq , FechaDesde, FechaHasta)
	if esActiva == 1 {
		fmt.Println("novedad activa")
		//Si la fechas de las novedades son del mismo año y del mismo mes se debe calcular los dias entre ambas fechas
		if FechaHasta.Month() == FechaDesde.Month() && FechaHasta.Year() == FechaDesde.Year() {
			periodo_liquidacion = CalcularDias(FechaDesde, FechaHasta) + 1
		}else{
			/*En caso de que la novedad no empiece y termine el mismo mes,
			la novedad cubre varios meses, por lo que se deben calcular los dias de ese mes en que afectó la novedad.
			Si la novedad comienza en el mes y el año de la preliquidacion a realizar, se calculan los dias a partir de ella hasta el día 30 del mes
			*/
			if (int(FechaDesde.Month()) == MesPreliq && int(FechaDesde.Year()) == AnoPreliq ){
				FechaControl = time.Date(AnoPreliq, time.Month(MesPreliq), 30, 0, 0, 0, 0, time.UTC)
				fmt.Println("fecha cont",FechaControl, FechaDesde)
				periodo_liquidacion = CalcularDias(FechaDesde, FechaControl) + 1
			}else if int(FechaHasta.Month()) == MesPreliq && int(FechaHasta.Year()) == AnoPreliq {
				/*
				Si la novedad termina en el mes y el año de la preliquidacion a realizar, se calculan los dias a
				partir del primero hasta el dia que comienza la novedad
				*/
				FechaControl = time.Date(AnoPreliq, time.Month(MesPreliq), 1, 0, 0, 0, 0, time.UTC)
				periodo_liquidacion = CalcularDias(FechaControl, FechaHasta) + 1
			} else{
				/*
				Si no se cumple que las novedades comiencen o terminen el mes de la preliquidacion
				significa que afectó a todo el mes completo, por lo que los dias de la novedad son 30
				*/
				periodo_liquidacion = 30
			}

		}
	}else{
		fmt.Println("soy inactiva")

		periodo_liquidacion = 0;
	}

	return periodo_liquidacion

}

//Función que calcula IBC, basado en hechos de golog
func CalcularIBC(reglas string){

	e := NewMachine().Consult(reglas)

	valor_ibc := e.ProveAll("calcular_ibc(V).")
	for _, solution := range valor_ibc {
		Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
		ibc = ibc + Valor

		}
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

    if day < 0 {
        // days in month:
				day = (30 - d1) + d2
        month--
    }
    if month < 0 {
        month += 12
        year--
    }

    return
}



func ConsultarValoresBonServPS(mesPreliq, anoPreliq int, numero_contrato string, vigencia_contrato int, codigo_concepto string, periodo string ) (valor_con float64){

	var valor float64

	ano_preliquidacion := anoPreliq
	ano_preliquidacion_string := strconv.Itoa(ano_preliquidacion)
	ano_busqueda := ano_preliquidacion - 1
	ano_busqueda_string := strconv.Itoa(ano_busqueda)

	var valor_concepto []models.DetallePreliquidacion
	var vigencia_contrato_string string = strconv.Itoa(vigencia_contrato)

	for i:=1; i< mesPreliq; i++{
		mesPreliq_string := strconv.Itoa(i)
			//http://localhost:8082/v1/detalle_liquidacion?limit=-1&query=Liquidacion.FechaLiquidacion__gte:2016-05-30,Liquidacion.FechaLiquidacion__lte:2017-06-30,Concepto.Id:1195,Persona:29
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query=Preliquidacion.Ano:"+ano_preliquidacion_string+",Preliquidacion.Mes:"+mesPreliq_string+",Concepto.Id:"+codigo_concepto+",NumeroContrato:"+numero_contrato+",VigenciaContrato:"+vigencia_contrato_string+",TipoPreliquidacion.Id:2,Preliquidacion.EstadoPreliquidacion.Nombre:Cerrada", &valor_concepto); err == nil {
			for _, solution := range valor_concepto {
		 	valor = valor + solution.ValorCalculado
		 }
		}else{
			fmt.Println(err)
		}

		}

		for i:=5; i<=12 ; i++{
			mesPreliq_string := strconv.Itoa(i)
				//http://localhost:8082/v1/detalle_liquidacion?limit=-1&query=Liquidacion.FechaLiquidacion__gte:2016-05-30,Liquidacion.FechaLiquidacion__lte:2017-06-30,Concepto.Id:1195,Persona:29
			if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query=Preliquidacion.Ano:"+ano_busqueda_string+",Preliquidacion.Mes:"+mesPreliq_string+",Concepto.Id:"+codigo_concepto+",NumeroContrato:"+numero_contrato+",VigenciaContrato:"+vigencia_contrato_string+",TipoPreliquidacion.Id:2,Preliquidacion.EstadoPreliquidacion.Nombre:Cerrada", &valor_concepto); err == nil {
				for _, solution := range valor_concepto {
			 	valor = valor + solution.ValorCalculado
			 }
			}else{
				fmt.Println(err)
			}

			}
		return valor

}

func ConsultarValoresBonServDic(numero_contrato string, vigencia_contrato int, codigo_concepto string, periodo string ) (valor_con float64){

	var valor float64
	periodo_nomina := periodo
	var valor_concepto []models.DetallePreliquidacion

	var vigencia_contrato_string string = strconv.Itoa(vigencia_contrato)
	//AGREGAR TIPO DE LIQUIDACION!! porque habran varios 29, y se necesita el pagado en nomina 2
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query=Preliquidacion.Ano:"+periodo_nomina+",Concepto.Id:"+codigo_concepto+",NumeroContrato:"+numero_contrato+",VigenciaContrato:"+vigencia_contrato_string+",TipoLiquidacion:2,Preliquidacion.EstadoPreliquidacion.Nombre:Cerrada", &valor_concepto); err == nil {
			for _, solution := range valor_concepto {
		 	valor = valor + solution.ValorCalculado
		 }
	}
		//http://localhost:8082/v1/detalle_liquidacion?limit=-1&query=Liquidacion.Nomina.Periodo:2017,Persona:29,TipoLiquidacion:3 <-- CONSULTA DOCEAVA PRIMA SEMESTRAL
		//nuevaRegla = "bonificacion_servicio(bonServ,1540945)."
		//hacer consulta de conceptos con codigo 129,139,1195 que se le hayan pagado a la persona en el presente año y se crea este hecho

		return valor

}

func ConsultarValoresPriServDic(numero_contrato string, vigencia_contrato int, periodo string ) (valor_con float64){

	var valor float64
	periodo_nomina := periodo
	var valor_concepto []models.DetallePreliquidacion

	var vigencia_contrato_string string = strconv.Itoa(vigencia_contrato)
	//AGREGAR TIPO DE LIQUIDACION!! porque habran varios 29, y se necesita el pagado en nomina 2

		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query=Preliquidacion.Ano:"+periodo_nomina+",NumeroContrato:"+numero_contrato+",VigenciaContrato:"+vigencia_contrato_string+",TipoLiquidacion:3,Preliquidacion.EstadoPreliquidacion.Nombre:Cerrada", &valor_concepto); err == nil {
			for _, solution := range valor_concepto {
		 	valor = valor + solution.ValorCalculado
		 }
	}
		//http://localhost:8082/v1/detalle_liquidacion?limit=-1&query=Liquidacion.Nomina.Periodo:2017,Persona:29,TipoLiquidacion:3 <-- CONSULTA DOCEAVA PRIMA SEMESTRAL
		//nuevaRegla = "bonificacion_servicio(bonServ,1540945)."
		//hacer consulta de conceptos con codigo 129,139,1195 que se le hayan pagado a la persona en el presente año y se crea este hecho

		return valor

}


func CalcularReteFuentePlanta(tipoPreliquidacionString, reglas,periodo string, listaDescuentos []models.ConceptosResumen)(rest []models.ConceptosResumen){
	fmt.Println("retefuente",periodo)
	var listaRetefuente []models.ConceptosResumen

	var ingresos int
	var deduccion_salud int
	var deduccion_pen_vol int
	var valor_gastos_rep int
	var Valor_alivio_beneficiario float64
	var Valor_alivio_vivienda float64
	var Valor_alivio_salud_prepagada float64
	var definitivo_deduccion int
	fmt.Println(listaDescuentos)
	temp_reglas := reglas
	temp_reglas = temp_reglas + "beneficiario(no)."
	temp_reglas = temp_reglas + "intereses_vivienda(0)."
	temp_reglas = temp_reglas + "salud_prepagada(0)."
	temp_reglas = temp_reglas + "declarante(si)."
	temp_reglas = temp_reglas + "porcentaje_diciembre(0.48)."

	m := NewMachine().Consult(temp_reglas)

	consultar_conceptos_ingresos_retencion := m.ProveAll("aplica_ingreso_retencion(X).")
	 for _, solution := range consultar_conceptos_ingresos_retencion {
		codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
		ingresos = ingresos + BuscarValorConcepto(listaDescuentos, codigo_concepto)
	}

	consultar_conceptos_deduccion_retencion := m.ProveAll("aplica_deduccion_retencion(X).")
	 for _, solution := range consultar_conceptos_deduccion_retencion {
		codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
		deduccion_salud = deduccion_salud + BuscarValorConcepto(listaDescuentos, codigo_concepto)
	}

	consultar_conceptos_deduccionpenvol_retencion := m.ProveAll("aplica_deduccion_penvol_retencion(X).")
	 for _, solution := range consultar_conceptos_deduccionpenvol_retencion {
		codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
		deduccion_pen_vol = deduccion_pen_vol + BuscarValorConcepto(listaDescuentos, codigo_concepto)
	}

	consultar_gastos_rep := m.ProveAll("aplica_gastos_rep(X).")
 	for _, solution := range consultar_gastos_rep {
 	 codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
 	 valor_gastos_rep = valor_gastos_rep + BuscarValorConcepto(listaDescuentos, codigo_concepto)
  }

	temp_reglas = temp_reglas + "ingreso_retencion("+strconv.Itoa(ingresos)+")."
	temp_reglas = temp_reglas + "deduccion_salud("+strconv.Itoa(deduccion_salud)+")."
	temp_reglas = temp_reglas + "deduccion_pen_vol("+strconv.Itoa(deduccion_pen_vol)+")."
	temp_reglas = temp_reglas + "valor_gastos_rep("+strconv.Itoa(valor_gastos_rep)+")."

	n := NewMachine().Consult(temp_reglas)

 deduccion_gastos_rep_rector := n.ProveAll("deduccion_gastos_rep_rector(DGR).")
 for _, solution := range deduccion_gastos_rep_rector {
	Valor, _ := strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("DGR")), 64)
	ingresos = int(Valor)
	}


	alivios := n.ProveAll("calcular_alivios(B,V,SP,D,"+periodo+").")
	 for _, solution := range alivios {
		Valor_alivio_beneficiario, _ = strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("B")), 64)
		Valor_alivio_vivienda, _ = strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("V")), 64)
		Valor_alivio_salud_prepagada, _ = strconv.ParseFloat(fmt.Sprintf("%s", solution.ByName_("SP")), 64)

	}

	ajuste_deduccion := n.ProveAll("ajustar_deducciones(AD,"+periodo+").")
	 for _, solution := range ajuste_deduccion {
		deduccion_pen_vol,_ = strconv.Atoi(fmt.Sprintf("%s", solution.ByName_("AD")))

	}

/*
	fmt.Println("ingresos")
	fmt.Println(ingresos)
	fmt.Println(deduccion_salud)
	fmt.Println(deduccion_pen_vol)
	fmt.Println(Valor_alivio_beneficiario)
	fmt.Println(Valor_alivio_vivienda)
	fmt.Println(Valor_alivio_salud_prepagada)
 */
	temp_reglas = temp_reglas + "alivio_ben("+strconv.Itoa(int(Valor_alivio_beneficiario))+")."
	temp_reglas = temp_reglas + "alivio_vin("+strconv.Itoa(int(Valor_alivio_vivienda))+")."
	temp_reglas = temp_reglas + "alivio_salud("+strconv.Itoa(int(Valor_alivio_salud_prepagada))+")."

	x := NewMachine().Consult(temp_reglas)


		deduccion_def := x.ProveAll("definitivo_deduccion(DD).")
		 for _, solution := range deduccion_def {
			 fmt.Println("deduccion")
			definitivo_deduccion,_ = strconv.Atoi(fmt.Sprintf("%s", solution.ByName_("DD")))
		}


		temp_reglas = temp_reglas + "resultado_deduccion("+strconv.Itoa(definitivo_deduccion)+")."

		o := NewMachine().Consult(temp_reglas)

	valor_retencion := o.ProveAll("valor_retencion(VR).")
	 for _, solution := range valor_retencion {
		 fmt.Println("asdf")
		val_reten:= fmt.Sprintf("%s", solution.ByName_("VR"))
		temp_conceptos := models.ConceptosResumen{Nombre: "reteFuente",
		Valor: val_reten,
		}


		codigo := o.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C,N,D).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			temp_conceptos.DiasLiquidados = diasALiquidar
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}

		listaRetefuente = append(listaRetefuente, temp_conceptos)

	}


	return listaRetefuente
}

func BuscarValorConcepto(listaDescuentos []models.ConceptosResumen,codigo_concepto string)(valor int){
	var temp int
	 for _, solution := range listaDescuentos {
				if(strconv.Itoa(solution.Id) == codigo_concepto ){
							temp,_ = strconv.Atoi(solution.Valor)

				}
		}

		return temp
}

func CalcularReteFuenteSal(tipoPreliquidacionString, reglas string, listaDescuentos []models.ConceptosResumen, dias_a_liq string)(rest []models.ConceptosResumen){

	var listaRetefuente []models.ConceptosResumen
	var ingresos int
	var deduccion_salud int

	temp_reglas := reglas
	m := NewMachine().Consult(reglas)

	consultar_conceptos_ingresos_retencion := m.ProveAll("aplica_ingreso_retencion(X).")
	 for _, solution := range consultar_conceptos_ingresos_retencion {
		codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
		ingresos = ingresos + BuscarValorConcepto(listaDescuentos, codigo_concepto)
	}

	temp_reglas = temp_reglas + "ingresos("+strconv.Itoa(ingresos)+")."

	consultar_conceptos_deduccion_retencion := m.ProveAll("aplica_deduccion_retencion(X).")
	 for _, solution := range consultar_conceptos_deduccion_retencion {
		codigo_concepto := fmt.Sprintf("%s", solution.ByName_("X"))
		deduccion_salud = deduccion_salud + BuscarValorConcepto(listaDescuentos, codigo_concepto)
	}


	temp_reglas = temp_reglas + "deducciones("+strconv.Itoa(deduccion_salud)+")."

	o := NewMachine().Consult(temp_reglas)

	valor_retencion := o.ProveAll("valor_retencion(VR).")
	 for _, solution := range valor_retencion {
		val_reten:= fmt.Sprintf("%s", solution.ByName_("VR"))
		temp_conceptos := models.ConceptosResumen{Nombre: "reteFuente",
		Valor: val_reten,
		}


		codigo := o.ProveAll("codigo_concepto(" + temp_conceptos.Nombre + ",C,N,D).")
		for _, cod := range codigo {
			temp_conceptos.Id, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("C")))
			 temp_conceptos.AliasConcepto = fmt.Sprintf("%s", cod.ByName_("D"))
			temp_conceptos.DiasLiquidados = dias_a_liq
			temp_conceptos.TipoPreliquidacion = tipoPreliquidacionString
			temp_conceptos.NaturalezaConcepto, _ = strconv.Atoi(fmt.Sprintf("%s", cod.ByName_("N")))
		}
		fmt.Println("dias a liq rete", dias_a_liq)
		listaRetefuente = append(listaRetefuente, temp_conceptos)

	}

	return listaRetefuente
}


func CalcularPeriodoLiquidacion(preliquidacion models.Preliquidacion, objeto_datos_acta models.ObjetoActaInicio) (periodoLiquidacion, mesesContrato string) {

	var FechaInicioContrato time.Time
	var FechaFinContrato time.Time
	var FechaControl time.Time
	var FechaInicio time.Time
	var FechaFin time.Time
	var periodo_liquidacion float64
	var meses_contrato float64

	layout := "2006-01-02"
	//FechaInicio, _ = time.Parse(layout , "2018-08-01")
	//FechaFin, _ = time.Parse(layout , "2018-12-15")

	datos_acta := objeto_datos_acta.ActaInicio

	FechaInicio, _ = time.Parse(layout , datos_acta.FechaInicioTemp)
	FechaFin, _ = time.Parse(layout , datos_acta.FechaFinTemp)
	a,m,d := diff(FechaInicio,FechaFin)

	fmt.Println("Fecha inicio",FechaInicio)
	fmt.Println("Fecha fin",FechaFin)
FechaInicioContrato = time.Date(FechaInicio.Year(), FechaInicio.Month(), FechaInicio.Day(), 0, 0, 0, 0, time.UTC)
FechaFinContrato = time.Date(FechaFin.Year(), FechaFin.Month(), FechaFin.Day(), 0, 0, 0, 0, time.UTC)

if int(FechaInicioContrato.Month()) == preliquidacion.Mes && int(FechaInicioContrato.Year()) == preliquidacion.Ano {
	FechaControl = time.Date(preliquidacion.Ano, time.Month(preliquidacion.Mes), 30, 0, 0, 0, 0, time.UTC)
	periodo_liquidacion = CalcularDias(FechaInicioContrato, FechaControl) + 1


} else if int(FechaFinContrato.Month()) == preliquidacion.Mes && int(FechaFinContrato.Year()) == preliquidacion.Ano {

	FechaControl = time.Date(preliquidacion.Ano, time.Month(preliquidacion.Mes), 1, 0, 0, 0, 0, time.UTC)
	periodo_liquidacion = CalcularDias(FechaControl, FechaFinContrato) + 1

} else {
	periodo_liquidacion = 30

}

	fmt.Println("aaaaaaaa", a, "m",m , "d",d)
	if (FechaFin.Day() != FechaInicio.Day()){
		d = d+1;
	}
	//meses_contrato = 4.5;
	meses_contrato = (float64(a*12))+float64(m)+(float64(d)/30)
	periodo := strconv.Itoa(int(periodo_liquidacion))
	meses := strconv.FormatFloat(meses_contrato, 'E', -1, 64)

	return periodo,meses
}

func validarNovedades_segSocial(Mes, Ano int, FechaDesde, FechaHasta time.Time) (flag int) {

/*
Se verifica si el año de la novedad es mayor (mas no igual) a la fecha final de la novedad
ya que de ser así, ya es valida
Si no, se verifica que las fechas desde y hasta cubren el mes de Liquidacion
*/

	if(FechaHasta.Year() > Ano){
		return 1
	}else if (FechaDesde.Year() <= Ano && FechaHasta.Year() >= Ano && int(FechaDesde.Month()) <= Mes && int(FechaHasta.Month()) >= Mes ){
		return 1
	}else {

		return 0
	}
}
