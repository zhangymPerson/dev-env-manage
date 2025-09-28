package db

import (
	"database/sql"
	"errors"

	"github.com/zhangymPerson/dev-env-manage/src/log"
	"github.com/zhangymPerson/dev-env-manage/src/models"
)

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Error("failed to open database: %v", err)
		return err
	}
	return nil
}

func AddConfig(config models.ConfigMaster) error {
	db := DB
	tx, err := db.Begin()
	if err != nil {
		log.Error("开始事务失败: %v", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // 重新抛出panic
		}
	}()

	// 首先尝试查询是否已存在相同配置
	var existingID int
	err = tx.QueryRow(`
		SELECT id FROM config_master 
		WHERE project_code = ? AND env_code = ? AND module_code = ? AND config_key = ? AND is_deleted = 0`,
		config.ProjectCode, config.EnvCode, config.ModuleCode, config.ConfigKey).Scan(&existingID)

	if err == nil {
		// 配置已存在，执行更新操作
		stmt, err := tx.Prepare(`
			UPDATE config_master SET 
				config_value = ?, config_alias = ?, auto_alias = ?, config_type = ?, 
				is_encrypted = ?, description = ?, sort_order = ?, updated_time = CURRENT_TIMESTAMP
			WHERE id = ?`)
		if err != nil {
			tx.Rollback()
			log.Error("准备更新语句失败: %v", err)
			return err
		}
		defer stmt.Close()

		res, err := stmt.Exec(
			config.ConfigValue, config.ConfigAlias, config.AutoAlias, config.ConfigType,
			config.IsEncrypted, config.Description, config.SortOrder, existingID)

		if err != nil {
			tx.Rollback()
			log.Error("执行更新失败: %v", err)
			return err
		}

		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			tx.Rollback()
			return errors.New("没有行被更新")
		}

		log.Info("配置项已更新: 项目[%s] 环境[%s] 模块[%s] 键[%s]", config.ProjectCode, config.EnvCode, config.ModuleCode, config.ConfigKey)
	} else if err == sql.ErrNoRows {
		// 配置不存在，执行插入操作
		stmt, err := tx.Prepare(`
			INSERT INTO config_master (
				project_code, env_code, module_code, config_key, config_value, 
				config_alias, auto_alias, config_type, is_encrypted, is_deleted,
				description, sort_order, created_time, updated_time
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`)
		if err != nil {
			tx.Rollback()
			log.Error("准备插入语句失败: %v", err)
			return err
		}
		defer stmt.Close()

		res, err := stmt.Exec(
			config.ProjectCode, config.EnvCode, config.ModuleCode, config.ConfigKey, config.ConfigValue,
			config.ConfigAlias, config.AutoAlias, config.ConfigType, config.IsEncrypted, config.IsDeleted,
			config.Description, config.SortOrder)

		if err != nil {
			tx.Rollback()
			log.Error("执行插入失败: %v", err)
			return err
		}

		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			tx.Rollback()
			return errors.New("没有行被插入")
		}

		log.Info("配置项已新增: 项目[%s] 环境[%s] 模块[%s] 键[%s]", config.ProjectCode, config.EnvCode, config.ModuleCode, config.ConfigKey)
	} else {
		// 查询过程中出现其他错误
		tx.Rollback()
		log.Error("查询配置项失败: %v", err)
		return err
	}

	return tx.Commit()
}
