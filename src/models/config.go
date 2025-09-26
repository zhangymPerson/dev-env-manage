package models

import (
	"time"
)

// ConfigMaster represents the config_master table structure
type ConfigMaster struct {
	ID          int       `db:"id" json:"id"`
	ProjectName string    `db:"project_name" json:"project_name"`
	ProjectCode string    `db:"project_code" json:"project_code"`
	EnvName     string    `db:"env_name" json:"env_name"`
	EnvCode     string    `db:"env_code" json:"env_code"`
	ModuleName  string    `db:"module_name" json:"module_name"`
	ModuleCode  string    `db:"module_code" json:"module_code"`
	ConfigKey   string    `db:"config_key" json:"config_key"`
	AutoAlias   string    `db:"auto_alias" json:"auto_alias"`
	ConfigAlias string    `db:"config_alias" json:"config_alias"`
	ConfigValue string    `db:"config_value" json:"config_value"`
	ConfigType  string    `db:"config_type" json:"config_type"`
	Description string    `db:"description" json:"description"`
	IsEncrypted int       `db:"is_encrypted" json:"is_encrypted"`
	SortOrder   int       `db:"sort_order" json:"sort_order"`
	CreatedTime time.Time `db:"created_time" json:"created_time"`
	UpdatedTime time.Time `db:"updated_time" json:"updated_time"`
	IsDeleted   int       `db:"is_deleted" json:"is_deleted"`
}
