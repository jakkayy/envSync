package database

import (
	"fmt"

	"gorm.io/gorm"
)

// AutoMigrate runs GORM automatic migrations for all data models
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Project{},
		&Environment{},
		&EnvVersion{},
		&AuditLog{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate database schemas: %w", err)
	}
	return nil
}
