TRUNCATE supplementary_1;
TRUNCATE supplementary_2;
TRUNCATE base_1;

INSERT INTO `physical_mappings` (`logical_attribute`, `physical_attribute`) VALUES (, );

SELECT app_schemas.column_name, base_table_attributes.column_name FROM app_schemas JOIN physical_mappings ON physical_mappings.logical_attribute = app_schemas.PK JOIN base_table_attributes ON base_table_attributes.pk = physical_mappings.physical_attribute WHERE app_schemas.table_id IN (3,4)

TRUNCATE tweet;
TRUNCATE user;


SELECT a1.app_name AS app1, as1.table_name AS table1, as1.column_name AS col1, 
       a2.app_name AS app2, as2.table_name AS table2, as2.column_name AS col2
FROM app_schemas as1 
JOIN schema_mappings sm ON as1.row_id = sm.source_attribute
JOIN app_schemas as2 ON as2.row_id = sm.dest_attribute
JOIN apps a1 ON as1.app_id = a1.row_id
JOIN apps a2 ON as2.app_id = a2.row_id
LIMIT 1000

SELECT a1.app_name AS app1, as1.table_name AS table1, as1.column_name AS col1, sm.source_attribute,
       a2.app_name AS app2, as2.table_name AS table2, as2.column_name AS col2, sm.dest_attribute
FROM app_schemas as1 
JOIN schema_mappings sm ON as1.row_id = sm.source_attribute
JOIN app_schemas as2 ON as2.row_id = sm.dest_attribute
JOIN apps a1 ON as1.app_id = a1.row_id
JOIN apps a2 ON as2.app_id = a2.row_id
WHERE a1.app_name in ('app1', 'app2') 
  -- LIMIT 1000

CREATE VIEW oneway_schema_mappings as
SELECT sm1.row_id, sm1.source_attribute, sm1.dest_attribute, sm1.rules, sm1.created_at
FROM schema_mappings sm1 
JOIN schema_mappings sm2 
  ON sm1.source_attribute = sm2.dest_attribute AND sm1.dest_attribute = sm2.source_attribute
WHERE sm1.row_id < sm2.row_id



SELECT a1.app_name AS app1, as1.table_name AS table1, as1.column_name AS col1, sm.source_attribute,
       a2.app_name AS app2, as2.table_name AS table2, as2.column_name AS col2, sm.dest_attribute
FROM app_schemas as1 
JOIN schema_mappings sm ON as1.row_id = sm.source_attribute
JOIN schema_mappings sm2 ON sm.source_attribute = sm2.dest_attribute AND sm.dest_attribute = sm2.source_attribute
JOIN app_schemas as2 ON as2.row_id = sm.dest_attribute
JOIN apps a1 ON as1.app_id = a1.row_id
JOIN apps a2 ON as2.app_id = a2.row_id
WHERE sm.row_id < sm2.row_id AND a1.app_name in ('app1', 'app2') 
ORDER BY a1.app_name, a2.app_name ASC




SELECT a1.app_name AS app1, as1.table_name AS table1, as1.column_name AS col1, sm.source_attribute,
       a2.app_name AS app2, as2.table_name AS table2, as2.column_name AS col2, sm.dest_attribute
FROM app_schemas as1 
JOIN schema_mappings sm ON as1.row_id = sm.source_attribute
JOIN schema_mappings sm2 ON sm.source_attribute = sm2.dest_attribute AND sm.dest_attribute = sm2.source_attribute
JOIN app_schemas as2 ON as2.row_id = sm.dest_attribute
JOIN apps a1 ON as1.app_id = a1.row_id
JOIN apps a2 ON as2.app_id = a2.row_id
WHERE sm.row_id < sm2.row_id 
-- AND a1.app_name in ('app1', 'app2')
ORDER BY a1.app_name, a2.app_name ASC



