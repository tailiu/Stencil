CREATE TABLE IF NOT EXISTS warehouse (
	row_id SERIAL,
	w_id INTEGER NOT NULL,
	w_name STRING(10) NULL,
	w_street_1 STRING(20) NULL,
	w_street_2 STRING(20) NULL,
	w_city STRING(20) NULL,
	w_state STRING(2) NULL,
	w_zip STRING(9) NULL,
	w_tax REAL NULL,
	w_ytd DOUBLE PRECISION NULL,
	mark_delete BOOL NULL DEFAULT 'false',
	CONSTRAINT wareh1 PRIMARY KEY (w_id ASC),
	FAMILY "primary" (w_id, w_name, w_street_1, w_street_2, w_city, w_state, w_zip, w_tax, w_ytd, mark_delete)
);

CREATE TABLE IF NOT EXISTS district (
	row_id SERIAL,
	d_id INTEGER NOT NULL,
	d_w_id INTEGER NOT NULL,
	d_name STRING(10) NULL,
	d_street_1 STRING(20) NULL,
	d_street_2 STRING(20) NULL,
	d_city STRING(20) NULL,
	d_state STRING(2) NULL,
	d_zip STRING(9) NULL,
	d_tax REAL NULL,
	d_ytd DOUBLE PRECISION NULL,
	d_next_o_id INTEGER NULL,
	mark_delete BOOL NULL DEFAULT 'false',
	CONSTRAINT dist1 PRIMARY KEY (d_w_id ASC, d_id ASC),
	INDEX district_d_w_id_idx (d_w_id ASC),
	-- CONSTRAINT district_fk1 FOREIGN KEY (d_w_id) REFERENCES warehouse (w_id) ON DELETE CASCADE,
	FAMILY "primary" (d_id, d_w_id, d_name, d_street_1, d_street_2, d_city, d_state, d_zip, d_tax, d_ytd, d_next_o_id, mark_delete)
);

CREATE TABLE IF NOT EXISTS customer (
	row_id SERIAL,
	c_id INTEGER NOT NULL,
	c_d_id INTEGER NOT NULL,
	c_w_id INTEGER NOT NULL,
	c_first STRING(16) NULL,
	c_middle STRING(2) NULL,
	c_last STRING(16) NULL,
	c_street_1 STRING(20) NULL,
	c_street_2 STRING(20) NULL,
	c_city STRING(20) NULL,
	c_state STRING(2) NULL,
	c_zip STRING(9) NULL,
	c_phone STRING(16) NULL,
	c_since TIMESTAMP NULL DEFAULT '1970-01-01 00:00:00+00:00':::TIMESTAMP,
	c_credit STRING(2) NULL,
	c_credit_lim DOUBLE PRECISION NULL,
	c_discount REAL NULL,
	c_balance DOUBLE PRECISION NULL,
	c_ytd_payment DOUBLE PRECISION NULL,
	c_payment_cnt SMALLINT NULL,
	c_delivery_cnt SMALLINT NULL,
	c_data STRING(500) NULL,
	mark_delete BOOL NULL DEFAULT 'false',
	CONSTRAINT custom1 PRIMARY KEY (c_w_id ASC, c_d_id ASC, c_id ASC),
	-- CONSTRAINT customer_fk1 FOREIGN KEY (c_w_id, c_d_id) REFERENCES district (d_w_id, d_id) ON DELETE CASCADE,
	FAMILY "primary" (c_id, c_d_id, c_w_id, c_first, c_middle, c_last, c_street_1, c_street_2, c_city, c_state, c_zip, c_phone, c_since, c_credit, c_credit_lim, c_discount, c_balance, c_ytd_payment, c_payment_cnt, c_delivery_cnt, c_data, mark_delete)
);

CREATE TABLE IF NOT EXISTS history (
	row_id SERIAL,
	h_c_id INTEGER NULL,
	h_c_d_id INTEGER NULL,
	h_c_w_id INTEGER NULL,
	h_d_id INTEGER NULL,
	h_w_id INTEGER NULL,
	h_date TIMESTAMP NULL DEFAULT '1970-01-01 00:00:00+00:00':::TIMESTAMP,
	h_amount REAL NULL,
	h_data STRING(24) NULL,
	mark_delete BOOL NULL DEFAULT 'false',
	-- CONSTRAINT history_fk2 FOREIGN KEY (h_c_w_id, h_c_d_id, h_c_id) REFERENCES customer (c_w_id, c_d_id, c_id) ON DELETE CASCADE,
	INDEX history_h_c_w_id_h_c_d_id_h_c_id_idx (h_c_w_id ASC, h_c_d_id ASC, h_c_id ASC),
	-- CONSTRAINT history_fk1 FOREIGN KEY (h_w_id, h_d_id) REFERENCES district (d_w_id, d_id) ON DELETE CASCADE,
	INDEX history_h_w_id_h_d_id_idx (h_w_id ASC, h_d_id ASC),
	FAMILY "primary" (h_c_id, h_c_d_id, h_c_w_id, h_d_id, h_w_id, h_date, h_amount, h_data, rowid, mark_delete)
);

CREATE TABLE IF NOT EXISTS item (
	row_id SERIAL,
	i_id INTEGER NOT NULL,
	i_im_id INTEGER NULL,
	i_name STRING(24) NULL,
	i_price DOUBLE PRECISION NULL,
	i_data STRING(50) NULL,
	mark_delete BOOL NULL DEFAULT 'false',
	CONSTRAINT item1 PRIMARY KEY (i_id ASC),
	FAMILY "primary" (i_id, i_im_id, i_name, i_price, i_data, mark_delete)
);

CREATE TABLE IF NOT EXISTS new_order (
	row_id SERIAL,
	no_o_id INTEGER NOT NULL,
	no_d_id INTEGER NOT NULL,
	no_w_id INTEGER NOT NULL,
	mark_delete BOOL NULL DEFAULT 'false',
	CONSTRAINT no1 PRIMARY KEY (no_w_id ASC, no_d_id ASC, no_o_id ASC),
	FAMILY "primary" (no_o_id, no_d_id, no_w_id, mark_delete)
);

CREATE TABLE IF NOT EXISTS orderr (
	row_id SERIAL,
	o_id INTEGER NOT NULL,
	o_w_id INTEGER NOT NULL,
	o_d_id INTEGER NOT NULL,
	o_c_id INTEGER NULL,
	o_entry_d TIMESTAMP NULL DEFAULT '1970-01-01 00:00:00+00:00':::TIMESTAMP,
	o_carrier_id SMALLINT NULL DEFAULT 0:::INT,
	o_ol_cnt SMALLINT NULL,
	o_all_local SMALLINT NULL,
	mark_delete BOOL NULL DEFAULT 'false',
	CONSTRAINT orderr1 PRIMARY KEY (o_w_id ASC, o_d_id ASC, o_id ASC),
	-- CONSTRAINT orderr_fk1 FOREIGN KEY (o_w_id, o_d_id, o_c_id) REFERENCES customer (c_w_id, c_d_id, c_id) ON DELETE CASCADE,
	INDEX orderr_o_w_id_o_d_id_o_c_id_idx (o_w_id ASC, o_d_id ASC, o_c_id ASC),
	FAMILY "primary" (o_id, o_w_id, o_d_id, o_c_id, o_entry_d, o_carrier_id, o_ol_cnt, o_all_local, mark_delete)
);

CREATE TABLE IF NOT EXISTS stock (
	row_id SERIAL,
	s_i_id INTEGER NOT NULL,
	s_w_id INTEGER NOT NULL,
	s_quantity SMALLINT NULL,
	s_dist_01 STRING(24) NULL,
	s_dist_02 STRING(24) NULL,
	s_dist_03 STRING(24) NULL,
	s_dist_04 STRING(24) NULL,
	s_dist_05 STRING(24) NULL,
	s_dist_06 STRING(24) NULL,
	s_dist_07 STRING(24) NULL,
	s_dist_08 STRING(24) NULL,
	s_dist_09 STRING(24) NULL,
	s_dist_10 STRING(24) NULL,
	s_ytd DECIMAL(8,2) NULL,
	s_order_cnt SMALLINT NULL,
	s_remote_cnt SMALLINT NULL,
	s_data STRING(50) NULL,
	mark_delete BOOL NULL DEFAULT 'false',
	CONSTRAINT stock1 PRIMARY KEY (s_w_id ASC, s_i_id ASC),
	INDEX stock_s_w_id_idx (s_w_id ASC),
	-- CONSTRAINT stock_fk2 FOREIGN KEY (s_i_id) REFERENCES item (i_id) ON DELETE CASCADE,
	INDEX stock_s_i_id_idx (s_i_id ASC),
	-- CONSTRAINT stock_fk1 FOREIGN KEY (s_w_id) REFERENCES warehouse (w_id) ON DELETE CASCADE,
	FAMILY "primary" (s_i_id, s_w_id, s_quantity, s_dist_01, s_dist_02, s_dist_03, s_dist_04, s_dist_05, s_dist_06, s_dist_07, s_dist_08, s_dist_09, s_dist_10, s_ytd, s_order_cnt, s_remote_cnt, s_data, mark_delete)
);

CREATE TABLE IF NOT EXISTS order_line (
	row_id SERIAL,
	ol_o_id INTEGER NOT NULL,
	ol_d_id INTEGER NOT NULL,
	ol_w_id INTEGER NOT NULL,
	ol_number SMALLINT NOT NULL,
	ol_i_id INTEGER NULL,
	ol_supply_w_id INTEGER NULL,
	ol_delivery_d TIMESTAMP NULL DEFAULT '1970-01-01 00:00:00+00:00':::TIMESTAMP,
	ol_quantity SMALLINT NULL,
	ol_amount DECIMAL(6,2) NULL,
	ol_dist_info STRING(24) NULL,
	mark_delete BOOL NULL DEFAULT 'false',
	CONSTRAINT ol1 PRIMARY KEY (ol_w_id ASC, ol_d_id ASC, ol_o_id ASC, ol_number ASC),
	-- CONSTRAINT order_line_fk2 FOREIGN KEY (ol_supply_w_id, ol_i_id) REFERENCES stock (s_w_id, s_i_id) ON DELETE CASCADE,
	INDEX order_line_ol_supply_w_id_ol_i_id_idx (ol_supply_w_id ASC, ol_i_id ASC),
	-- CONSTRAINT order_line_fk1 FOREIGN KEY (ol_w_id, ol_d_id, ol_o_id) REFERENCES orderr (o_w_id, o_d_id, o_id) ON DELETE CASCADE,
	FAMILY "primary" (ol_o_id, ol_d_id, ol_w_id, ol_number, ol_i_id, ol_supply_w_id, ol_delivery_d, ol_quantity, ol_amount, ol_dist_info, mark_delete)
);
