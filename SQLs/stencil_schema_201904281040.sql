-- Drop table

-- DROP TABLE public.apps

CREATE TABLE public.apps (
	pk serial NOT NULL,
	app_name varchar(256) NOT NULL,
	"timestamp" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT apps_pkey PRIMARY KEY (pk)
);
CREATE INDEX apps_app_name_idx ON public.apps USING btree (app_name);

-- Drop table

-- DROP TABLE public.app_tables

CREATE TABLE public.app_tables (
	pk serial NOT NULL,
	app_id int4 NOT NULL,
	table_name varchar(256) NOT NULL,
	"timestamp" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT app_tables_pkey PRIMARY KEY (pk),
	CONSTRAINT app_tables_apps_fk FOREIGN KEY (app_id) REFERENCES apps(pk)
);
CREATE INDEX app_tables_app_id_idx ON public.app_tables USING btree (app_id);
CREATE INDEX app_tables_app_id_table_name_idx ON public.app_tables USING btree (app_id, table_name);
CREATE INDEX app_tables_table_name_idx ON public.app_tables USING btree (table_name);

-- Drop table

-- DROP TABLE public.app_schemas

CREATE TABLE public.app_schemas (
	pk serial NOT NULL,
	table_id int4 NOT NULL,
	column_name varchar(256) NOT NULL,
	data_type varchar NULL,
	"constraints" varchar(512) NULL DEFAULT NULL::character varying,
	"timestamp" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT app_schemas_pkey PRIMARY KEY (pk),
	CONSTRAINT app_schemas_app_tables_fk FOREIGN KEY (table_id) REFERENCES app_tables(pk)
);
CREATE INDEX app_schemas_column_name_idx ON public.app_schemas USING btree (column_name);
CREATE INDEX app_schemas_data_type_idx ON public.app_schemas USING btree (data_type);
CREATE INDEX app_schemas_table_id_column_name_idx ON public.app_schemas USING btree (table_id, column_name);
CREATE INDEX app_schemas_table_id_idx ON public.app_schemas USING btree (table_id);

-- Drop table

-- DROP TABLE public.physical_schema

CREATE TABLE public.physical_schema (
	pk serial NOT NULL,
	table_name varchar(256) NOT NULL,
	column_name varchar(256) NOT NULL,
	"timestamp" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT physical_schema_pkey PRIMARY KEY (pk)
);
CREATE INDEX physical_schema_column_name_idx ON public.physical_schema USING btree (column_name);
CREATE INDEX physical_schema_table_name_column_name_idx ON public.physical_schema USING btree (table_name, column_name);
CREATE INDEX physical_schema_table_name_idx ON public.physical_schema USING btree (table_name);

-- Drop table

-- DROP TABLE public.schema_mappings

CREATE TABLE public.schema_mappings (
	pk serial NOT NULL,
	source_attribute int4 NOT NULL,
	dest_attribute int4 NOT NULL,
	rules varchar(512) NULL DEFAULT NULL::character varying,
	"timestamp" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT schema_mappings_pkey PRIMARY KEY (pk),
	CONSTRAINT schema_mappings_app_schemas_destattr_fk FOREIGN KEY (dest_attribute) REFERENCES app_schemas(pk),
	CONSTRAINT schema_mappings_app_schemas_fk FOREIGN KEY (source_attribute) REFERENCES app_schemas(pk)
);
CREATE INDEX schema_mappings_dest_attribute_idx ON public.schema_mappings USING btree (dest_attribute);
CREATE INDEX schema_mappings_source_attr_idx ON public.schema_mappings USING btree (source_attribute);
CREATE INDEX schema_mappings_source_attribute_idx ON public.schema_mappings USING btree (source_attribute, dest_attribute);

-- Drop table

-- DROP TABLE public.physical_mappings

CREATE TABLE public.physical_mappings (
	pk serial NOT NULL,
	logical_attribute int4 NULL,
	physical_attribute int4 NULL,
	"timestamp" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT physical_mappings_pkey PRIMARY KEY (pk),
	CONSTRAINT physical_mappings_app_schemas_fk FOREIGN KEY (logical_attribute) REFERENCES app_schemas(pk),
	CONSTRAINT physical_mappings_physical_schema_fk FOREIGN KEY (physical_attribute) REFERENCES physical_schema(pk)
);
CREATE INDEX physical_mappings_logical_attribute_idx ON public.physical_mappings USING btree (logical_attribute);
CREATE INDEX physical_mappings_logical_attribute_physical_attribute_idx ON public.physical_mappings USING btree (logical_attribute, physical_attribute);
CREATE INDEX physical_mappings_physical_attribute_idx ON public.physical_mappings USING btree (physical_attribute);

-- Drop table

-- DROP TABLE public.supplementary_tables

CREATE TABLE public.supplementary_tables (
	pk serial NOT NULL,
	table_id int4 NOT NULL,
	"timestamp" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT supplementary_tables_pkey PRIMARY KEY (pk),
	CONSTRAINT supplementary_tables_app_tables_fk FOREIGN KEY (table_id) REFERENCES app_tables(pk)
);
CREATE INDEX supplementary_tables_table_id_idx ON public.supplementary_tables USING btree (table_id);

-- Drop table

-- DROP TABLE public.txn_logs

CREATE TABLE public.txn_logs (
	id serial NOT NULL,
	action_id int4 NOT NULL,
	action_type varchar NOT NULL,
	undo_action varchar NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT txn_logs_action_type_check CHECK (((action_type)::text = ANY ((ARRAY['COMMIT'::character varying, 'ABORT'::character varying, 'ABORTED'::character varying, 'CHANGE'::character varying, 'BEGIN_TRANSACTION'::character varying])::text[]))),
	CONSTRAINT txn_logs_pkey PRIMARY KEY (id)
);


-- Drop table

-- DROP TABLE public.display_flags

CREATE TABLE public.display_flags (
	app varchar NOT NULL,
	table_name varchar NOT NULL,
	id int4 NOT NULL,
	display_flag bool NULL DEFAULT true,
	migration_id int4 NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Drop table

-- DROP TABLE public.error_log

CREATE TABLE public.error_log (
	id serial NOT NULL,
	dst_app varchar NULL,
	query varchar NULL,
	args varchar NULL,
	error varchar NULL,
	migration_id varchar NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Drop table

-- DROP TABLE public.evaluation

CREATE TABLE public.evaluation (
	id serial NOT NULL,
	src_app varchar NULL,
	dst_app varchar NULL,
	src_table varchar NULL,
	dst_table varchar NULL,
	src_id varchar NULL,
	dst_id varchar NULL,
	src_cols varchar NULL,
	dst_cols varchar NULL,
	migration_id varchar NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
