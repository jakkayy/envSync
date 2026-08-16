package database

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID            string         `gorm:"primaryKey;size:64" json:"id"`
	Name          string         `gorm:"size:128;not null;unique" json:"name"`
	SecretKeyHash string         `gorm:"size:256;not null" json:"-"`
	Environments  []Environment  `gorm:"foreignKey:ProjectID" json:"environments,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Environment struct {
	ID            uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID     string       `gorm:"size:64;not null;index" json:"project_id"`
	Name          string       `gorm:"size:32;not null" json:"name"`
	LatestVersion int          `gorm:"default:0" json:"latest_version"`
	Versions      []EnvVersion `gorm:"foreignKey:EnvironmentID" json:"versions,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type EnvVersion struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EnvironmentID    uint      `gorm:"not null;index" json:"environment_id"`
	Version          int       `gorm:"not null" json:"version"`
	EncryptedPayload []byte    `gorm:"type:blob" json:"-"`
	CreatedBy        string    `gorm:"size:64" json:"created_by"`
	Message          string    `gorm:"size:256" json:"message"`
	CreatedAt        time.Time `json:"created_at"`
}

type AuditLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID string    `gorm:"size:64;not null;index" json:"project_id"`
	Action    string    `gorm:"size:32;not null" json:"action"`
	EnvName   string    `gorm:"size:32" json:"env_name"`
	User      string    `gorm:"size:64" json:"user"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	UserAgent string    `gorm:"size:256" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}
