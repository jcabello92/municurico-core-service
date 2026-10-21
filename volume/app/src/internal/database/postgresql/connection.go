package postgresql

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// BuildPostgreSQLDSN genera una cadena de conexión para PostgreSQL.
func BuildPostgreSQLDSN(host, port, user, password, dbName string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
}

// ConnectPostgreSQL conecta a la base de datos PostgreSQL.
func ConnectPostgreSQL(dbName, dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("ERROR CRÍTICO: No se pudo conectar a PostgreSQL (%s): %v.", dbName, err)
	}

	log.Printf("Conexión a PostgreSQL (%s) establecida exitosamente.", dbName)

	return db
}