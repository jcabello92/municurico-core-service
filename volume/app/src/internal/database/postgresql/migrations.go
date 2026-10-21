package postgresql

import (
	"log"

	"gorm.io/gorm"

	models "core-service/internal/models/postgresql"
)

// RunMigrations ejecuta las migraciones para crear o actualizar las tablas en la base de datos.
func RunMigrations(db *gorm.DB) {
	log.Println("Ejecutando migraciones en la base de datos...")

	err := db.AutoMigrate(
		// --- MODELOS DEL SISTEMA (No borrar) ---
		&models.Auditoria{}, // Esta tabla debe existir obligatoriamente en la base de datos

		// --- MODELOS DEL NEGOCIO ---
		&models.Region{},
		&models.Provincia{},
		&models.Comuna{},
		&models.Persona{},
		&models.Usuario{},
		&models.Rol{},
		&models.Permiso{},
		&models.UsuarioRol{},
		&models.RolPermiso{},
		&models.Sesion{},
		&models.Miembro{},
		&models.Cargo{},
		&models.Area{},
		&models.MiembroAsignacion{},
	)

	if err != nil {
		log.Fatalf("ERROR CRÍTICO: Fallo al ejecutar migraciones: %v.", err)
	}

	log.Println("Ejecución completada exitosamente.")
}
