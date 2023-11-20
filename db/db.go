package db

import (
	"os"
	"log"
	"fmt"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func ConnectDb() (*Database, error) {
	var err error
	if err := godotenv.Load(); err != nil {
        log.Fatalf("Error cargando variables de entorno: %v", err)
    }
	
	dbUsed := os.Getenv("DB_USED")
	dbPostgres := os.Getenv("DB_POSTGRES")
	dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")
    dbSSLMode := os.Getenv("DB_SSLMODE")

	dbUri := fmt.Sprintf("%s://%s:%s@%s:%s/%ssslmode=%s", dbPostgres, dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	db, err := sql.Open(dbUsed,dbUri)
	
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (database *Database) Close() error {
	return database.db.Close()
}

func (database *Database) GetDB() *sql.DB {
	return database.db
}