package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/zhangymPerson/dev-env-manage/src/db"
)

func HandleDeleteCommand(project, env, module string, verbose bool, key string) {

	// 首先检查配置项是否存在
	var configID int
	var configKey string
	err := db.DB.QueryRow(`
		SELECT id, config_key FROM config_master 
		WHERE project=? AND env=? AND module=? AND config_key=?`,
		project, env, module, key,
	).Scan(&configID, &configKey)

	if err == sql.ErrNoRows {
		// 如果通过config_key找不到，尝试通过config_alias查找
		err = db.DB.QueryRow(`
			SELECT id, config_key FROM config_master 
			WHERE project=? AND env=? AND module=? AND config_alias=?`,
			project, env, module, key,
		).Scan(&configID, &configKey)

		if err == sql.ErrNoRows {
			// 如果通过config_alias找不到，尝试通过auto_alias查找
			err = db.DB.QueryRow(`
				SELECT id, config_key FROM config_master 
				WHERE project=? AND env=? AND module=? AND auto_alias=?`,
				project, env, module, key,
			).Scan(&configID, &configKey)

			if err == sql.ErrNoRows {
				log.Fatalf("Config not found for key: %s", key)
			} else if err != nil {
				log.Fatalf("Failed to query config by auto_alias: %v", err)
			}
		} else if err != nil {
			log.Fatalf("Failed to query config by config_alias: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Failed to query config by config_key: %v", err)
	}

	// Confirm deletion with user
	fmt.Printf("Are you sure you want to delete configuration item '%s'? (Y/N): ", configKey)
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "Y" && confirm != "y" {
		fmt.Println("Deletion cancelled.")
		return
	}

	// 执行物理删除
	tx, err := db.DB.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(`
		DELETE FROM config_master 
		WHERE id = ?`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare delete statement: %v", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(configID)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to execute delete: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		tx.Rollback()
		log.Fatalf("No configuration item was deleted")
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	if verbose {
		fmt.Printf("Configuration item deleted successfully:\n")
		fmt.Printf("  Project: %s\n", project)
		fmt.Printf("  Environment: %s\n", env)
		fmt.Printf("  Module: %s\n", module)
		fmt.Printf("  Key: %s\n", configKey)
		fmt.Printf("  Deleted using identifier: %s\n", key)
	} else {
		fmt.Printf("Deleted: %s\n", configKey)
	}
}
