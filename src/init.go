package src

import (
	"embed"
	"database/sql"
	"os"
	"path/filepath"

	"github.com/zhangymPerson/dev-env-manage/src/log"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed sql/*
var sqlFiles embed.FS

const db_file = "dem_config.db"

// Init  初始化函数
func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	demDir := filepath.Join(homeDir, ".dem")
	if _, err := os.Stat(demDir); os.IsNotExist(err) {
		err := os.Mkdir(demDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// 检查 sql 目录是否存在文件
	entries, err := sqlFiles.ReadDir("sql")
	if err != nil {
		log.Debug("Warning: No SQL files found or sql directory is empty")
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		content, err := sqlFiles.ReadFile(filepath.Join("sql", entry.Name()))
		if err != nil {
			log.Debug("Error reading SQL file %s: %v\\n", entry.Name(), err)
			continue
		}
		log.Debug("File: %s\\nContent: %s\\n", entry.Name(), string(content))
	}

	// 检查 demDir 目录下是否有 db_file 文件
	dbPath := filepath.Join(demDir, db_file)
	if _, err := os.Stat(dbPath); err == nil {
		log.Debug("Database file already exists")
		return
	} else {
		log.Debug("Database file does not exist, creating...")
		// 执行 SQL 文件创建数据库
		if err := executeSQLFiles(dbPath); err != nil {
			log.Error("Failed to execute SQL files: %v", err)
			panic(err)
		}
		log.Info("Database created successfully at %s", dbPath)
	}
}

// executeSQLFiles 执行 SQL 文件内容到数据库
func executeSQLFiles(dbPath string) error {
	// 打开数据库连接
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// 检查 sql 目录是否存在文件
	entries, err := sqlFiles.ReadDir("sql")
	if err != nil {
		return err
	}

	// 执行所有 SQL 文件
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		content, err := sqlFiles.ReadFile(filepath.Join("sql", entry.Name()))
		if err != nil {
			log.Debug("Error reading SQL file %s: %v\\n", entry.Name(), err)
			continue
		}

		sqlContent := string(content)
		log.Debug("Executing SQL file: %s\\nContent: %s\\n", entry.Name(), sqlContent)

		// 执行 SQL 语句
		if _, err := db.Exec(sqlContent); err != nil {
			log.Warning("Error executing SQL from file %s: %v", entry.Name(), err)
			return err
		}
	}

	return nil
}