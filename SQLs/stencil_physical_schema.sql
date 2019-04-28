-- Stencil Storage Schema

CREATE TABLE apps (
  PK  SERIAL,
  app_name varchar(256)  NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

CREATE TABLE app_tables (
  PK  SERIAL,
  app_id  int NOT NULL,
  table_name varchar(256)  NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

CREATE TABLE app_schemas (
  PK  SERIAL,
  table_id  int NOT NULL,
  column_name varchar(256)  NOT NULL,
  data_type  varchar,
  constraints varchar(512)  DEFAULT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

CREATE TABLE physical_schema (
  PK  SERIAL,
  table_name varchar(256)  NOT NULL,
  column_name varchar(256)  NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

CREATE TABLE physical_mappings (
  PK  SERIAL,
  logical_attribute  int DEFAULT NULL,
  physical_attribute  int DEFAULT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

CREATE TABLE schema_mappings (
  PK  SERIAL,
  source_attribute  int NOT NULL,
  dest_attribute  int NOT NULL,
  rules varchar(512)  DEFAULT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

CREATE TABLE supplementary_tables (
  PK  SERIAL,
  table_id  int NOT NULL,
  supplementary_table varchar(256)  NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)  ;	