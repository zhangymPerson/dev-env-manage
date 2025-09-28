package models

import "time"

// ConfigMaster 映射数据库表 config_master
// 所有字段均使用指针类型以兼容数据库中的 NULL 值
type ConfigMaster struct {
	ID        int64   `json:"id"`
	Project   *string `json:"project,omitempty"`
	Env       *string `json:"env,omitempty"`
	Module    *string `json:"module,omitempty"`
	ConfigKey *string `json:"config_key,omitempty"`

	AutoAlias   *string `json:"auto_alias,omitempty"`
	ConfigAlias *string `json:"config_alias,omitempty"`
	ConfigValue *string `json:"config_value,omitempty"`
	ConfigType  *string `json:"config_type,omitempty"`
	Description *string `json:"description,omitempty"`
	IsEncrypted *int    `json:"is_encrypted,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`

	CreatedTime *time.Time `json:"created_time,omitempty"`
	UpdatedTime *time.Time `json:"updated_time,omitempty"`
}
