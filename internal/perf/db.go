package perf

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rsevilla87/perfapp/pkg/utils"
)

// DBInfo Database connection information
type DBInfo struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	RetryInt   int
	conn       *sql.DB
}

// DB DBInfo instance
var DB DBInfo = DBInfo{
	DBUser:     os.Getenv("POSTGRESQL_USER"),
	DBPassword: os.Getenv("POSTGRESQL_PASSWORD"),
	DBHost:     os.Getenv("POSTGRESQL_HOSTNAME"),
	DBPort:     os.Getenv("POSTGRESQL_PORT"),
	DBName:     os.Getenv("POSTGRESQL_DATABASE"),
	RetryInt:   5,
}

const dbTImeout = 10

// Connect2Db Connects to a Postgres database using DBInfo
func Connect2Db() {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable connect_timeout=%d", DB.DBUser, DB.DBPassword, DB.DBHost, DB.DBPort, DB.DBName, dbTImeout)
	for {
		log.Infof("Connecting with database %s:%s", DB.DBHost, DB.DBPort)
		DB.conn, _ = sql.Open("postgres", connStr)
		if err := DB.conn.Ping(); err != nil {
			log.Warnln(err)
			log.Warnf("Retrying connection with %s:%s in %d seconds", DB.DBHost, DB.DBPort, DB.RetryInt)
			time.Sleep(time.Duration(DB.RetryInt) * time.Second)
			continue
		}
		break
	}
	log.Println("Database connection successfully established")
}

// QueryDB Performs a query on the database
func QueryDB(query string) error {
	// Verify database connection by pinging database
	if err := DB.conn.Ping(); err != nil {
		return err
	}
	if _, err := DB.conn.Exec(query); err != nil {
		utils.ErrorHandler(err)
	}
	return nil
}

// CreateTables Creates all tables at tableList
func CreateTables(tableList []map[string]string) error {
	for k := range tableList {
		for t, q := range tableList[k] {
			log.Infof("Creating table %s: %s", t, q)
			if err := QueryDB(q); err != nil {
				return err
			}
		}
	}
	return nil
}

// DropTables Drops all tables at tableList
func DropTables(tableList []map[string]string) error {
	for k := range tableList {
		for t := range tableList[k] {
			log.Infof("Dropping %s table", t)
			if err := QueryDB(fmt.Sprintf("DROP TABLE %s", t)); err != nil {
				return err
			}
		}
	}
	return nil
}
