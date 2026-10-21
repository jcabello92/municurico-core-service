package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"

	models "core-service/internal/models/postgresql"
)

// RunSeeders ejecuta las inserciones de los datos iniciales en la base de datos de forma idempotente.
func RunSeeders(db *gorm.DB) {
	log.Println("Verificando datos iniciales en la base de datos...")

	if err := seedGeografia(db); err != nil {
		log.Fatalf("ERROR CRÍTICO: Fallo al ejecutar seeders de geografía: %v.", err)
	}

	log.Println("Ejecución de seeders completada exitosamente.")
}

// seedGeografia inserta regiones, provincias y comunas de Chile si no existen previamente.
func seedGeografia(db *gorm.DB) error {
	var totalRegiones int64
	if err := db.Model(&models.Region{}).Count(&totalRegiones).Error; err != nil {
		return fmt.Errorf("error al verificar existencia de regiones: %w", err)
	}

	if totalRegiones > 0 {
		log.Printf("Datos de división geográfica ya existen (%d regiones registradas). Omitiendo seeder de geografía.", totalRegiones)
		return nil
	}

	log.Println("Insertando datos de división geográfica de Chile (16 regiones, 56 provincias, 346 comunas)...")

	systemUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	ctx := context.WithValue(context.Background(), models.UserIDContextKey, systemUserID)

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Inserción de Regiones y mapeo de IDs originales a UUIDs generados
		regionMap := make(map[int]uuid.UUID, len(rawRegiones))
		for _, rawReg := range rawRegiones {
			regID, err := uuid.NewV7()
			if err != nil {
				return fmt.Errorf("error al generar UUID para región %s: %w", rawReg.Nombre, err)
			}

			reg := models.Region{
				Base: models.Base{
					PrimaryKey: models.PrimaryKey{ID: regID},
					CreatedBy:  systemUserID,
					UpdatedBy:  systemUserID,
				},
				Nombre:      rawReg.Nombre,
				Abreviatura: rawReg.Abreviatura,
			}

			if err := tx.Create(&reg).Error; err != nil {
				return fmt.Errorf("error al insertar región %s: %w", rawReg.Nombre, err)
			}
			regionMap[rawReg.ID] = regID
		}

		// 2. Inserción de Provincias y mapeo de IDs originales a UUIDs generados
		provinciaMap := make(map[int]uuid.UUID, len(rawProvincias))
		for _, rawProv := range rawProvincias {
			regID, exists := regionMap[rawProv.RegionID]
			if !exists {
				return fmt.Errorf("región original id %d no encontrada para provincia %s", rawProv.RegionID, rawProv.Nombre)
			}

			provID, err := uuid.NewV7()
			if err != nil {
				return fmt.Errorf("error al generar UUID para provincia %s: %w", rawProv.Nombre, err)
			}

			prov := models.Provincia{
				Base: models.Base{
					PrimaryKey: models.PrimaryKey{ID: provID},
					CreatedBy:  systemUserID,
					UpdatedBy:  systemUserID,
				},
				Nombre:   rawProv.Nombre,
				RegionID: regID,
			}

			if err := tx.Create(&prov).Error; err != nil {
				return fmt.Errorf("error al insertar provincia %s: %w", rawProv.Nombre, err)
			}
			provinciaMap[rawProv.ID] = provID
		}

		// 3. Inserción de Comunas en lotes
		comunas := make([]models.Comuna, 0, len(rawComunas))
		for _, rawCom := range rawComunas {
			provID, exists := provinciaMap[rawCom.ProvinciaID]
			if !exists {
				return fmt.Errorf("provincia original id %d no encontrada para comuna %s", rawCom.ProvinciaID, rawCom.Nombre)
			}

			comID, err := uuid.NewV7()
			if err != nil {
				return fmt.Errorf("error al generar UUID para comuna %s: %w", rawCom.Nombre, err)
			}

			comunas = append(comunas, models.Comuna{
				Base: models.Base{
					PrimaryKey: models.PrimaryKey{ID: comID},
					CreatedBy:  systemUserID,
					UpdatedBy:  systemUserID,
				},
				Nombre:      rawCom.Nombre,
				ProvinciaID: provID,
			})
		}

		if err := tx.CreateInBatches(comunas, 100).Error; err != nil {
			return fmt.Errorf("error al insertar comunas en lotes: %w", err)
		}

		log.Printf("División geográfica inicializada exitosamente: %d regiones, %d provincias y %d comunas creadas.",
			len(regionMap), len(provinciaMap), len(comunas))
		return nil
	})
}
