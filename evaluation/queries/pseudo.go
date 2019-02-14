dbConn, err := sql.Open("postgres", "postgresql://root@10.230.12.75:26257/mastodon?sslmode=disable")

tx, err := dbConn.Begin()

tx.Query(sql_SELECT_Query)

tx.Exec(sql_UPDATE_DELETE_INSERT_Query)

tx.rollback()

tx.Commit()

///////////////////////////////////////////////////////////////////////////////

dbConn, err := sql.Open("stencil", "stencil://root@10.230.12.75:26257/mastodon?sslmode=disable")

tx, err := dbConn.Begin()

tx.Stencil_query(sql_SELECT_Query)

tx.Stencil_exec(sql_UPDATE_DELETE_INSERT_Query)

tx.Stencil_rollback()

tx.Stencil_commit()

tx.Stencil_migration()