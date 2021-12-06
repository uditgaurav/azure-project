package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB

var server = Getenv("SERVER", "")
var port, _ = strconv.Atoi(Getenv("PORT", "1433"))
var user = Getenv("USERNAME", "")
var password = Getenv("PASSWORD", "")
var database = Getenv("DB_NAME", "test123")
var tableName = Getenv("TABLE_NAME", "load")

func main() {

	log.InfoWithValues("The mssql information is as follows", logrus.Fields{
		"Server":     server,
		"Port":       port,
		"Database":   database,
		"Table Name": tableName,
	})

	log.Info("[CONNECTION]: Trying to establish connection with the sql server")

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Error creating connection pool: err: %v", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("fail to connect the db, err: %v", err.Error())
	}
	log.Info("[CONNECTION]: Connection established successfully...")

	go AbortWatcher(db, tableName)

	// Create table to generate load
	err = CreateTable(db, tableName)
	if err != nil {
		log.Errorf("Error in creating table: %v", err.Error())
		_ = DeleteTable(db, tableName)
		os.Exit(1)
	}
	log.Info("[Info]: Table created successfully...")
	log.Info("[Load]: Starting load generator ...")

	// start generating load
	if err = GenerateLoad(db, tableName); err != nil {
		log.Errorf("fail to generate load, err: %v", err.Error())
		_ = DeleteTable(db, tableName)
		os.Exit(1)
	}

	// Delete table
	err = DeleteTable(db, tableName)
	if err != nil {
		log.Fatalf("Error deleting Table: %v", err.Error())
	}
	fmt.Printf("load table deleted successfully.\n")

}

// CreateTable inserts an employee record
func CreateTable(db *sql.DB, tableName string) error {

	ctx := context.Background()
	var err error

	if db == nil {
		err = errors.New("CreateTable: db is null")
		return err
	}

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE " + tableName + " ( recipe_id int NOT NULL, ingredient_id INT NOT NULL, amount INT NOT NULL);")
	if err != nil {
		return errors.Errorf("error in creating a table, err: %v", err)
	}
	return nil
}

// GenerateLoad will contineously generate load on the mssql server

func GenerateLoad(db *sql.DB, tableName string) error {

	ctx := context.Background()
	var err error

	if db == nil {
		err = errors.New("DeleteTable: db is null")
		return err
	}

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
        DECLARE @Counter INT
        SET @Counter=1
        WHILE ( @Counter <= 10)
        BEGIN
            INSERT INTO ` + tableName + `
            (recipe_id, ingredient_id, amount)
        VALUES      
            (1,1,1), 
            (1,2,2), 
            (1,3,2),
            (1,4,3),
            (1,5,1),
            (2,3,2),                             
            (2,6,1),
            (3,5,1),
            (1,5,1),
            (2,3,2),
            (2,6,1),
            (3,5,1),
            (1,5,1),
            (2,3,2),
            (2,6,1),
            (1,5,1),
            (2,3,2),
            (2,6,1),
            (3,5,1),
            (3,5,1),
            (1,5,1),
            (2,3,2),
            (2,6,1),
            (3,5,1),             
            (3,7,2);
        END         
    `)
	if err != nil {
		return errors.Errorf("error in generating load, err: %v", tableName, err)
	}
	return nil
}

// DeleteTable inserts an employee record
func DeleteTable(db *sql.DB, tableName string) error {
	ctx := context.Background()
	var err error

	if db == nil {
		err = errors.New("DeleteTable: db is null")
		return err
	}

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE " + tableName + ";")
	if err != nil {
		return errors.Errorf("error in deleting the table %v, err: %v", tableName, err)
	}
	return nil
}

// AbortWatcher continuously watch for the abort signals
func AbortWatcher(db *sql.DB, tableName string) {

	// signChan channel is used to transmit signal notifications.
	signChan := make(chan os.Signal, 1)
	// Catch and relay certain signal(s) to signChan channel.
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)

	// waiting until the abort signal received
	<-signChan

	// Delete table
	err := DeleteTable(db, tableName)
	if err != nil {
		log.Errorf("Error deleting Table: %v", err.Error())
	}
	log.Info("Load has been successfully removed...")

}

// Getenv fetch the env and set the default value, if any
func Getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
