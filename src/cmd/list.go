package cmd

import (
	"fmt"
	"log"

	"github.com/zhangymPerson/dev-env-manage/src/db"
)

type ConfigItem struct {
	ProjectCode string
	EnvCode     string
	ModuleCode  string
	ConfigKey   string
	ConfigValue string
	ConfigAlias string
	AutoAlias   string
}

func HandleListCommand(project, env, module string, verbose bool) {
	// 查询所有配置项
	rows, err := db.DB.Query(`
		SELECT project_code, env_code, module_code, config_key, config_value, config_alias, auto_alias
		FROM config_master 
		WHERE is_deleted = 0
		ORDER BY project_code, env_code, module_code, config_key`)
	if err != nil {
		log.Fatalf("Failed to query config items: %v", err)
	}
	defer rows.Close()

	var configs []ConfigItem
	for rows.Next() {
		var config ConfigItem
		err := rows.Scan(&config.ProjectCode, &config.EnvCode, &config.ModuleCode,
			&config.ConfigKey, &config.ConfigValue, &config.ConfigAlias, &config.AutoAlias)
		if err != nil {
			log.Fatalf("Failed to scan config item: %v", err)
		}
		configs = append(configs, config)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating config items: %v", err)
	}

	if len(configs) == 0 {
		fmt.Println("No configuration items found.")
		return
	}

	// 根据verbose参数决定输出格式
	if verbose {
		printVerboseList(configs)
	} else {
		printSimpleList(configs)
	}
}

func printSimpleList(configs []ConfigItem) {
	// 简单输出：只显示config_key
	for _, config := range configs {
		fmt.Println(config.ConfigKey)
	}
}

func printVerboseList(configs []ConfigItem) {
	// 详细输出：显示所有信息
	currentProject := ""
	currentEnv := ""
	currentModule := ""

	for _, config := range configs {
		// 显示项目分组
		if config.ProjectCode != currentProject {
			currentProject = config.ProjectCode
			currentEnv = ""
			currentModule = ""
			fmt.Printf("\n[Project: %s]\n", currentProject)
		}

		// 显示环境分组
		if config.EnvCode != currentEnv {
			currentEnv = config.EnvCode
			currentModule = ""
			fmt.Printf("  [Environment: %s]\n", currentEnv)
		}

		// 显示模块分组
		if config.ModuleCode != currentModule {
			currentModule = config.ModuleCode
			fmt.Printf("    [Module: %s]\n", currentModule)
		}

		// 显示配置项详情
		fmt.Printf("      Key: %s\n", config.ConfigKey)
		fmt.Printf("      Value: %s\n", config.ConfigValue)
		if config.ConfigAlias != "" {
			fmt.Printf("      Alias: %s\n", config.ConfigAlias)
		}
		if config.AutoAlias != "" {
			fmt.Printf("      AutoAlias: %s\n", config.AutoAlias)
		}
		fmt.Println()
	}
}
