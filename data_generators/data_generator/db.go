package data_generator

import (
	"database/sql"
	"log"
	"fmt"
)

func GetDBConn(dbname string) *sql.DB {

	// dbConnAddr := "postgresql://%s@%s:%s/%s?sslmode=disable"
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", DB_ADDR, DB_PORT, DB_USER, DB_PASSWORD, dbname)

	dbConn, err := sql.Open("postgres", psqlInfo)
	// sql.Open("postgres",fmt.Sprintf(dbConnAddr, config.DB_USER, config.DB_ADDR, config.DB_PORT, dbname))
	if err != nil {
		log.Println("Can't connect to DB:", dbname)
		log.Fatal(err)
	} else {
		log.Println("Connected to DB:", dbname)
	}
	return dbConn
}