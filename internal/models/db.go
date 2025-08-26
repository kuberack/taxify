package models

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
)

var taxifyDb *sql.DB
var taxifyMock sqlmock.Sqlmock

func getDbConnection() (*sql.DB, error) {
	db, _, err := getDbMockConnection()
	return db, err
}

func getDbMockConnection() (*sql.DB, sqlmock.Sqlmock, error) {

	// singleton
	if taxifyDb != nil {
		return taxifyDb, taxifyMock, nil
	}
	env, exists := os.LookupEnv("TAXIFY_DEPLOY_TYPE")
	if exists && env == "UNIT_TEST" {
		// sql mock
		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("error creating mock database")
			return nil, nil, errors.New("error creating mock database")
		}
		taxifyDb, taxifyMock = db, mock
		return db, mock, nil
	}

	// Make a connection to the sql instance
	username, exists := os.LookupEnv("TAXIFY_DB_USERNAME")
	if !exists {
		return nil, nil, errors.New("db username not configured")
	}
	password, exists := os.LookupEnv("TAXIFY_DB_PASSWORD")
	if !exists {
		return nil, nil, errors.New("db password not configured")
	}
	name, exists := os.LookupEnv("TAXIFY_DB_NAME")
	if !exists {
		return nil, nil, errors.New("db name not configured")
	}
	ip, exists := os.LookupEnv("TAXIFY_DB_IP_ADDRESS")
	if !exists {
		return nil, nil, errors.New("db ip address not configured")
	}
	port, exists := os.LookupEnv("TAXIFY_DB_IP_PORT")
	if !exists {
		return nil, nil, errors.New("db ip port not configured")
	}
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, ip, port, name)
	db, err := sql.Open("mysql", connstring)
	if err != nil {
		return nil, nil, err
	}
	taxifyDb, taxifyMock = db, nil
	return db, nil, nil
}
