TRUNCATE supplementary_1;
TRUNCATE supplementary_2;
TRUNCATE base_1;

INSERT INTO `physical_mappings` (`logical_attribute`, `physical_attribute`) VALUES (, );

SELECT app_schemas.column_name, base_table_attributes.column_name FROM app_schemas JOIN physical_mappings ON physical_mappings.logical_attribute = app_schemas.PK JOIN base_table_attributes ON base_table_attributes.pk = physical_mappings.physical_attribute WHERE app_schemas.table_id IN (3,4)

TRUNCATE tweet;
TRUNCATE user;

