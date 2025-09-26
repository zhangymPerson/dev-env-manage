package db

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
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

	stmt, err := tx.Prepare(`
        INSERT INTO config_master (
            project_code, env_code, module_code, config_key, config_value, 
            config_alias, auto_alias, config_type, is_encrypted, is_deleted,
            description, sort_order, created_time, updated_time
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`)
	if err != nil {
		tx.Rollback()
		log.Error("准备语句失败: %v", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		config.ProjectCode, config.EnvCode, config.ModuleCode, config.ConfigKey, config.ConfigValue,
		config.ConfigAlias, config.AutoAlias, config.ConfigType, config.IsEncrypted, config.IsDeleted,
		config.Description, config.SortOrder)

	if err != nil {
		tx.Rollback()
		if isUniqueConstraintError(err) {
			log.Error("配置项已存在: 项目[%s] 环境[%s] 模块[%s] 键[%s]", config.ProjectCode, config.EnvCode, config.ModuleCode, config.ConfigKey)
			return err
		}
		log.Error("执行插入失败: %v", err)
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		tx.Rollback()
		return errors.New("没有行被插入，可能违反了约束条件")
	}

	return tx.Commit()
}

// 判断是否是唯一约束错误(SQLite特有)
func isUniqueConstraintError(err error) bool {
	// SQLite的唯一约束错误代码为19
	// 不同数据库驱动可能有不同的错误表示方式
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		return sqliteErr.Code == sqlite3.ErrConstraint &&
			sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	// 其他数据库的判断逻辑...
	return false
}
