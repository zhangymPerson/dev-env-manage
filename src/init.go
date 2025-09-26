package src

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zhangymPerson/dev-env-manage/src/constant"
	"github.com/zhangymPerson/dev-env-manage/src/log"
)

//go:embed sql/*
var sqlFiles embed.FS

// Init  初始化函数
func Init() {
	// check sql directory
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
	dbPath := constant.GetDBFilePath()
	if _, err := os.Stat(dbPath); err == nil {
		log.Debug("Database file already exists")
		return
	} else {
		log.Debug("Database file does not exist, creating...")
		// execute sql files and create database and tables
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
