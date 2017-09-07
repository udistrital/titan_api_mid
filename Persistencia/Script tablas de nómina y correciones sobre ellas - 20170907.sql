DELETE FROM titan.detalle_liquidacion;
DELETE FROM titan.liquidacion;
DELETE FROM titan.detalle_preliquidacion;
DELETE FROM titan.preliquidacion;
DELETE FROM titan.nomina;
DELETE FROM titan.tipo_nomina;
DELETE FROM titan.tipo_vinculacion;
DELETE FROM titan.concepto_por_persona;
DELETE FROM titan.concepto;

DROP TABLE titan.detalle_liquidacion;
DROP TABLE titan.liquidacion;
DROP TABLE titan.detalle_preliquidacion;
DROP TABLE titan.preliquidacion;
DROP TABLE titan.nomina;
DROP TABLE titan.tipo_nomina;
DROP TABLE titan.tipo_vinculacion;
DROP TABLE titan.concepto_por_persona;
DROP TABLE titan.concepto;
DROP SEQUENCE titan.categoria_parametro_id_seq;
DROP SEQUENCE titan.descuentos_id_seq;
DROP SEQUENCE titan.detalle_novedad_id_seq;
DROP SEQUENCE titan.novedad_aplicada_id_seq;
DROP SEQUENCE titan.novedad_id_seq;
DROP SEQUENCE titan.parametro_liquidacion_id_seq;
DROP SEQUENCE titan.variable_id_seq;
DROP SCHEMA titan CASCADE;

CREATE SEQUENCE administrativa.concepto_nomina_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

CREATE TABLE administrativa.concepto_nomina(
	id integer NOT NULL DEFAULT nextval('administrativa.concepto_nomina_id_seq'::regclass),
	nombre_concepto character varying(80) NOT NULL,
	alias_concepto character varying(80),
	naturaleza_concepto integer NOT NULL,
	tipo_concepto integer NOT NULL,
	CONSTRAINT pk_concepto PRIMARY KEY (id),
	CONSTRAINT uq_concepto_nomina_nombre_concepto UNIQUE (nombre_concepto)

);

COMMENT ON TABLE administrativa.concepto_nomina IS 'Describen a que corresponden los distintos pagos calculados por la nomina';
COMMENT ON COLUMN administrativa.concepto_nomina.nombre_concepto IS 'Nombre de la regla asociada al concepto';
COMMENT ON COLUMN administrativa.concepto_nomina.alias_concepto IS 'Nombre del concepto a mostrar en la interfaz';
COMMENT ON COLUMN administrativa.concepto_nomina.naturaleza_concepto IS 'Llave foranea a naturaleza_concepto.';
COMMENT ON COLUMN administrativa.concepto_nomina.tipo_concepto IS 'Llave foraena a tabla tipo_concepto';

CREATE SEQUENCE administrativa.concepto_nomina_por_persona_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

CREATE TABLE administrativa.concepto_nomina_por_persona(
	id integer NOT NULL DEFAULT nextval('administrativa.concepto_nomina_por_persona_id_seq'::regclass),
	valor_novedad numeric(18,8) NOT NULL,
	num_cuotas numeric(5,0) NOT NULL,
	activo boolean NOT NULL,
	fecha_desde timestamp,
	fecha_hasta timestamp,
	fecha_registro timestamp NOT NULL,
	persona integer NOT NULL,
	nomina integer NOT NULL,
	concepto integer NOT NULL,
	CONSTRAINT pk_concepto_nomina_por_persona PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.concepto_nomina_por_persona IS 'Describe las novedades asociadas a las personas';
COMMENT ON COLUMN administrativa.concepto_nomina_por_persona.activo IS 'Indica si la novedad esta activa o no';
COMMENT ON COLUMN administrativa.concepto_nomina_por_persona.fecha_registro IS 'Fecha de registro de la novedad';
COMMENT ON COLUMN administrativa.concepto_nomina_por_persona.persona IS 'Llave foranea a informacion_proveedor. d de proveedor de la persona a la que le es asociada la novedad';
COMMENT ON COLUMN administrativa.concepto_nomina_por_persona.nomina IS 'Llave foranea a nomina. Nomina sobre la cual se calculara la novedad';
COMMENT ON COLUMN administrativa.concepto_nomina_por_persona.concepto IS 'Llave foreanea a concepto. Concepto asociado a la novedad';


CREATE SEQUENCE administrativa.detalle_preliquidacion_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

CREATE TABLE administrativa.detalle_preliquidacion(
	id integer NOT NULL DEFAULT nextval('administrativa.detalle_preliquidacion_id_seq'::regclass),
	valor_calculado numeric(14,4) NOT NULL,
	numero_contrato character varying,
	vigencia_contrato integer,
	dias_liquidados numeric(2,0),
	tipo_preliquidacion integer NOT NULL,
	preliquidacion integer NOT NULL,
	concepto integer NOT NULL,
	CONSTRAINT pk_detalle_preliquidacion PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.detalle_preliquidacion IS 'Tabla que detalla los pagos realizados a las personas por preliquidacion';
COMMENT ON COLUMN administrativa.detalle_preliquidacion.valor_calculado IS 'Valor pagado a la persona, calculado por reglas de negocio';
COMMENT ON COLUMN administrativa.detalle_preliquidacion.numero_contrato IS 'Numero de contrato de persona a la que se le realiza el pago';
COMMENT ON COLUMN administrativa.detalle_preliquidacion.vigencia_contrato IS 'Vigencia del contrato de persona a la que se le realiza el pago';
COMMENT ON COLUMN administrativa.detalle_preliquidacion.dias_liquidados IS 'Dias bajo los que fueron calculados los conceptos a la persona';
COMMENT ON COLUMN administrativa.detalle_preliquidacion.tipo_preliquidacion IS 'Llave foranea a tipo de preliquidacion. Especifica el tipo de preliquidacion para el que corresponde el pago del concepto';
COMMENT ON COLUMN administrativa.detalle_preliquidacion.preliquidacion IS 'Llave foranea a preliquidacion. Indica a que preliquidacion pertenece cada pago';
COMMENT ON COLUMN administrativa.detalle_preliquidacion.concepto IS 'Llave foranea a concepto. Indica bajo que concepto se realiza el pago';

CREATE SEQUENCE administrativa.nomina_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

CREATE TABLE administrativa.nomina(
	id integer NOT NULL DEFAULT nextval('administrativa.nomina_id_seq'::regclass),
	descripcion character varying(50) NOT NULL,
	activo boolean NOT NULL,
	tipo_nomina integer,
	CONSTRAINT pk_nomina PRIMARY KEY (id),
	CONSTRAINT uq_nomina_tipo_nomina UNIQUE (tipo_nomina)

);

COMMENT ON TABLE administrativa.nomina IS 'Tabla que contiene las diferentes nominas presentes de la Universidad Distrital y sobre las cuales se realizaran calculos de preliquidacion';
COMMENT ON COLUMN administrativa.nomina.descripcion IS 'Nombre de la nomina, formado desde la aplicacion utilizando el tipo de vinculacion y de nomina';
COMMENT ON COLUMN administrativa.nomina.activo IS 'Describe si la nomina se encuentra activa o no';
COMMENT ON COLUMN administrativa.nomina.tipo_nomina IS 'Llave foranea relacionada a la tabla tipo_nomina';

CREATE SEQUENCE administrativa.preliquidacion_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

CREATE TABLE administrativa.preliquidacion(
	id integer NOT NULL DEFAULT nextval('administrativa.preliquidacion_id_seq'::regclass),
	descripcion character varying(100),
	mes integer NOT NULL,
	ano integer NOT NULL,
	fecha_registro timestamp NOT NULL DEFAULT now(),
	estado_preliquidacion integer NOT NULL,
	nomina integer NOT NULL,
	CONSTRAINT pk_preliquidacion PRIMARY KEY (id),
	CONSTRAINT uq_periodo_preliquidacion UNIQUE (mes,ano,nomina)

);

COMMENT ON TABLE administrativa.preliquidacion IS 'Tabla que detalla el mes y el año para el cual se realizaran calculos de pagos a las personas vinculadas contractualmente a la Universidad Distrital';
COMMENT ON COLUMN administrativa.preliquidacion.descripcion IS 'Campo que describe la preliquidacion, creado a partir de aplicacion';
COMMENT ON COLUMN administrativa.preliquidacion.mes IS 'Mes al que corresponde la preliquidacion';
COMMENT ON COLUMN administrativa.preliquidacion.ano IS 'Año al que corresponde la preliquidacion';
COMMENT ON COLUMN administrativa.preliquidacion.fecha_registro IS 'Fecha en la que se realizo la preliquidacion';
COMMENT ON COLUMN administrativa.preliquidacion.estado_preliquidacion IS 'Llave foranea a estado_preliquidacion. ';
COMMENT ON COLUMN administrativa.preliquidacion.nomina IS 'Llave foranea a nomina. Indica bajo que nomina se esta preliquidando';

CREATE SEQUENCE administrativa.tipo_nomina_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

CREATE TABLE administrativa.tipo_nomina(
	id integer NOT NULL DEFAULT nextval('administrativa.tipo_nomina_id_seq'::regclass),
	nombre character varying(50) NOT NULL,
	descripcion character varying(100),
	codigo_abreviacion character varying(20),
	activo boolean NOT NULL,
	numero_orden numeric(5,2),
	CONSTRAINT pk_tipo_nomina PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.tipo_nomina IS 'Tabla parametrica que lista los tipos de nomina dentro de la Universidad Distrital. Ejemplo: Funcionarios de planta, docentes de planta o de HC';

CREATE SEQUENCE administrativa.estado_preliquidacion_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

CREATE TABLE administrativa.estado_preliquidacion(
	id integer NOT NULL DEFAULT nextval('administrativa.estado_preliquidacion_id_seq'::regclass),
	nombre character varying(30) NOT NULL,
	descripcion character varying(100),
	codigo_abreviacion character varying(20),
	activo boolean NOT NULL,
	numero_orden numeric(5,2),
	CONSTRAINT pk_estado_preliquidacion PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.estado_preliquidacion IS 'Tabla que parametriza los diferentes estados que tiene una preliquidacion. Ejemplo: Si está abierta, está cerrada o en solicitud de necesidad';


CREATE SEQUENCE administrativa.tipo_concepto_nomina_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

CREATE TABLE administrativa.tipo_concepto_nomina(
	id integer NOT NULL DEFAULT nextval('administrativa.tipo_concepto_nomina_id_seq'::regclass),
	nombre character varying(50) NOT NULL,
	descripcion character varying(100),
	codigo_abreviacion character varying(20),
	activo boolean NOT NULL,
	numero_orden numeric(5,2),
	CONSTRAINT pk_tipo_concepto PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.tipo_concepto_nomina IS 'Describe si el concepto a la hora de ser calculado corresponde a un valor fijo o porcentual.';


CREATE SEQUENCE administrativa.naturaleza_concepto_nomina_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

CREATE TABLE administrativa.naturaleza_concepto_nomina(
	id integer NOT NULL DEFAULT nextval('administrativa.naturaleza_concepto_nomina_id_seq'::regclass),
	nombre character varying(50) NOT NULL,
	descripcion character varying(100),
	codigo_abreviacion character varying(20),
	activo boolean NOT NULL,
	numero_orden numeric(5,2),
	CONSTRAINT pk_naturaleza_concepto PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.naturaleza_concepto_nomina IS 'Describe si el concepto es un devengo o un descuento, o si hace parte de seguridad social.';

CREATE SEQUENCE administrativa.tipo_preliquidacion_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START WITH 1
	CACHE 1
	NO CYCLE
	OWNED BY NONE;

CREATE TABLE administrativa.tipo_preliquidacion(
	id integer NOT NULL DEFAULT nextval('administrativa.tipo_preliquidacion_id_seq'::regclass),
	nombre character varying(30) NOT NULL,
	descripcion character varying(100),
	codigo_abreviacion character varying(20),
	activo boolean NOT NULL,
	numero_orden numeric(5,2),
	CONSTRAINT pk_tipo_preliquidacion PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.tipo_preliquidacion IS 'Corresponde al periodo a liquidar.
0 es la primera quincena, 1 la segunda quincena, 2 el mes completo, 3 junio y 4 diciembre';

ALTER TABLE administrativa.concepto_nomina ADD CONSTRAINT fk_concepto_tipo_concepto FOREIGN KEY (tipo_concepto)
REFERENCES administrativa.tipo_concepto_nomina (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.concepto_nomina ADD CONSTRAINT fk_concepto_naturaleza_concepto FOREIGN KEY (naturaleza_concepto)
REFERENCES administrativa.naturaleza_concepto_nomina (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.concepto_nomina_por_persona ADD CONSTRAINT fk_concepto_nomina_por_persona_concepto FOREIGN KEY (concepto)
REFERENCES administrativa.concepto_nomina (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.concepto_nomina_por_persona ADD CONSTRAINT fk_concepto_nomina_por_persona_nomina FOREIGN KEY (nomina)
REFERENCES administrativa.nomina (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.detalle_preliquidacion ADD CONSTRAINT fk_detalle_preliquidacion_concepto FOREIGN KEY (concepto)
REFERENCES administrativa.concepto_nomina (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.detalle_preliquidacion ADD CONSTRAINT fk_detalle_preliquidacion_preliquidacion FOREIGN KEY (preliquidacion)
REFERENCES administrativa.preliquidacion (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.detalle_preliquidacion ADD CONSTRAINT fk_detalle_preliquidacion_tipo_preliquidacion FOREIGN KEY (tipo_preliquidacion)
REFERENCES administrativa.tipo_preliquidacion (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.nomina ADD CONSTRAINT fk_nomina_tipo_nomina FOREIGN KEY (tipo_nomina)
REFERENCES administrativa.tipo_nomina (id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.preliquidacion ADD CONSTRAINT fk_preliquidacion_nomina FOREIGN KEY (nomina)
REFERENCES administrativa.nomina (id) MATCH SIMPLE
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.preliquidacion ADD CONSTRAINT fk_preliquidacion_estado_preliquidacion FOREIGN KEY (estado_preliquidacion)
REFERENCES administrativa.estado_preliquidacion (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;


GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA core TO crud_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA core TO crud_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA administrativa TO crud_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA administrativa TO crud_user;
