CREATE TABLE administrativa.pension(
	id serial NOT NULL,
	valor_pension numeric(8,0) NOT NULL,
	tipo_pension integer NOT NULL,
	fecha_pension timestamp,
	CONSTRAINT pk_pension PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.pension IS 'Entidad que consigna el valor de la pension y a que tipo corresponde';
COMMENT ON COLUMN administrativa.pension.id IS 'Llave primaria de tabla pension';
COMMENT ON COLUMN administrativa.pension.valor_pension IS 'Valor de la pension que fue asignada';
COMMENT ON COLUMN administrativa.pension.tipo_pension IS 'Llave foranea a tabla parametica tipo_pension';
COMMENT ON COLUMN administrativa.pension.fecha_pension IS 'Fecha en la que se asigno la pension';

CREATE TABLE administrativa.tipo_pension(
	id integer NOT NULL DEFAULT nextval('titan.tipo_nomina_id_seq'::regclass),
	nombre character varying(50) NOT NULL,
	descripcion character varying(100),
	codigo_abreviacion character varying(20),
	activo boolean NOT NULL,
	numero_orden numeric(5,2),
	CONSTRAINT pk_tipo_pension PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.tipo_pension IS 'Tabla parametrica que indica el motivo de la pension asignada';
ALTER TABLE administrativa.tipo_pension OWNER TO postgres;

CREATE TABLE administrativa.asignacion_pension_user_role(
	id smallint NOT NULL,
	pension integer NOT NULL,
	persona_asignada integer NOT NULL,
	estado_rol_pension integer NOT NULL,
	fecha_inicio_beneficio_pension timestamp,
	fecha_fin_beneficio_pension timestamp,
	tutor integer,
	porcentaje numeric(3,2),
	CONSTRAINT pk_asignacion_pension_user_rol PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.asignacion_pension_user_role IS 'Relaciona una pension a una persona con un rol de pensionado, beneficiario o sustituto';
COMMENT ON COLUMN administrativa.asignacion_pension_user_role.id IS 'Llave primaria de la tabla asignacion_pension_user_rol';
COMMENT ON COLUMN administrativa.asignacion_pension_user_role.pension IS 'Llave foranea a pension';
COMMENT ON COLUMN administrativa.asignacion_pension_user_role.persona_asignada IS 'Relacion a tabla um_user_rol, indica que esta persona tiene a ella relacionada cierta pension';
COMMENT ON COLUMN administrativa.asignacion_pension_user_role.estado_rol_pension IS 'Llave foranea a tabla estado_rol_pension';
COMMENT ON COLUMN administrativa.asignacion_pension_user_role.fecha_inicio_beneficio_pension IS 'Indica la fecha en la que la persona asociada comenzo a recibir dicha pension';
COMMENT ON COLUMN administrativa.asignacion_pension_user_role.fecha_fin_beneficio_pension IS 'Indica la fecha en la que la persona asignada dejo de recibir la pension';
COMMENT ON COLUMN administrativa.asignacion_pension_user_role.tutor IS 'Indica quien en tabla um_user_rol es tutor de la persona asignada. Solo aplica para persona asignada menor de edad';
COMMENT ON COLUMN administrativa.asignacion_pension_user_role.porcentaje IS 'Indica el porcentaje de la pension que le corresponde como sustituto';

CREATE TABLE administrativa.estado_rol_pension(
	id integer NOT NULL DEFAULT nextval('titan.tipo_nomina_id_seq'::regclass),
	nombre character varying(50) NOT NULL,
	descripcion character varying(100),
	codigo_abreviacion character varying(20),
	activo boolean NOT NULL,
	numero_orden numeric(5,2),
	CONSTRAINT pk_estado_rol_pension PRIMARY KEY (id)

);

COMMENT ON TABLE administrativa.estado_rol_pension IS 'Tabla parametrica que indica el estado de la pension recibida';
ALTER TABLE administrativa.estado_rol_pension OWNER TO postgres;

CREATE TABLE core.relacion(
	id serial NOT NULL,
	persona_principal integer NOT NULL,
	persona_relacionada integer NOT NULL,
	tipo_relacion integer NOT NULL,
	CONSTRAINT pk_relacion PRIMARY KEY (id)

);

COMMENT ON COLUMN core.relacion.id IS 'llave primaria de tabla relacion';
COMMENT ON COLUMN core.relacion.persona_principal IS 'Persona principal, la cual cuenta con el parentesco a describir. Se relaciona con tabla um_user_rol ';
COMMENT ON COLUMN core.relacion.persona_relacionada IS 'Persona que tiene parentesco con persona_principal. Se relaciona con tabla um_user_rol';

CREATE TABLE core.tipo_relacion(
	id integer NOT NULL DEFAULT nextval('titan.tipo_nomina_id_seq'::regclass),
	nombre character varying(50) NOT NULL,
	descripcion character varying(100),
	codigo_abreviacion character varying(20),
	activo boolean NOT NULL,
	numero_orden numeric(5,2),
	CONSTRAINT pk_tipo_relacion PRIMARY KEY (id)

);

COMMENT ON TABLE core.tipo_relacion IS 'Tabla parametrica que indica el tipo de relacion que pueden tener dos user';
ALTER TABLE core.tipo_relacion OWNER TO postgres;

CREATE TABLE core.informacion_asignacion_rol(
	id serial NOT NULL,
	um_user_role integer NOT NULL,
	fecha_inicio timestamp NOT NULL,
	fecha_fin timestamp,
	documento_soporte integer NOT NULL,
	CONSTRAINT pk_informacion_asignacion_rol PRIMARY KEY (id),
	CONSTRAINT uk_um_user_rol UNIQUE (um_user_role)

);
-- ddl-end --
COMMENT ON TABLE core.informacion_asignacion_rol IS 'Esta tabla agrega información adicional a la tabla um_user_role, cuando un rol se le ha asignado a una persona. ';
COMMENT ON COLUMN core.informacion_asignacion_rol.id IS 'llave primaria de tabla informacion_asignacion_rol';
COMMENT ON COLUMN core.informacion_asignacion_rol.um_user_role IS 'llave foranea a registro de usuario y rol en WSO2';
COMMENT ON COLUMN core.informacion_asignacion_rol.fecha_inicio IS 'Indica desde cuándo la persona cuenta con un vinculo laboral con la organización';
COMMENT ON COLUMN core.informacion_asignacion_rol.fecha_fin IS 'Indica cuándo terminó el vínculo laboral de la persona con la organización';
COMMENT ON COLUMN core.informacion_asignacion_rol.documento_soporte IS 'referencia a soporte de rol asignado a la persona';

ALTER TABLE administrativa.pension ADD CONSTRAINT fk_pension_tipo_pension FOREIGN KEY (tipo_pension)
REFERENCES administrativa.tipo_pension (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.asignacion_pension_user_role ADD CONSTRAINT fk_asignacion_pension_user_rol_pension FOREIGN KEY (pension)
REFERENCES administrativa.pension (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE administrativa.asignacion_pension_user_role ADD CONSTRAINT fk_asignacion_pension_user_rol_estado_pension FOREIGN KEY (estado_rol_pension)
REFERENCES administrativa.estado_rol_pension (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE core.relacion ADD CONSTRAINT fk_relacion_tipo_relacion FOREIGN KEY (tipo_relacion)
REFERENCES core.tipo_relacion (id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --



