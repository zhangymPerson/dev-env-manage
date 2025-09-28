package cmd

import (
	"fmt"

	"github.com/zhangymPerson/dev-env-manage/src/db"
	"github.com/zhangymPerson/dev-env-manage/src/log"
)

type ConfigItem struct {
	Project     *string
	Env         *string
	Module      *string
	ConfigKey   *string
	ConfigValue *string
	ConfigAlias *string
	AutoAlias   *string
}

func HandleListCommand(project, env, module string, verbose bool, show bool) {
	// 根据参数动态构建查询条件
	query := "SELECT project, env, module, config_key, config_value, config_alias, auto_alias FROM config_master WHERE 1=1"

	// 如果指定了项目名，则添加项目条件
	if project != "default" {
		query += fmt.Sprintf(" AND project = '%s'", project)
	}

	// 如果指定了环境名，则添加环境条件
	if env != "default" {
		query += fmt.Sprintf(" AND env = '%s'", env)
	}

	// 如果指定了模块名，则添加模块条件
	if module != "default" {
		query += fmt.Sprintf(" AND module = '%s'", module)
	}

	// 添加排序条件
	query += " ORDER BY project, env, module, config_key"

	// 查询配置项
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Fatalf("Failed to query config items: %v", err)
	}
	defer rows.Close()

	var configs []ConfigItem
	for rows.Next() {
		var config ConfigItem
		err := rows.Scan(
			&config.Project, &config.Env,
			&config.Module, &config.ConfigKey,
			&config.ConfigValue, &config.ConfigAlias,
			&config.AutoAlias,
		)
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
		printSimpleList(configs, show)
	}
}

func printSimpleList(configs []ConfigItem, show bool) {
	fmt.Print(show)
	// 简单输出：只显示config_key
	for _, config := range configs {
		if show {
			project := ""
			env := ""
			module := ""
			key := ""

			if config.Project != nil {
				project = *config.Project
			}
			if config.Env != nil {
				env = *config.Env
			}
			if config.Module != nil {
				module = *config.Module
			}
			if config.ConfigKey != nil {
				key = *config.ConfigKey
			}

			fmt.Println(project, env, module, key)
		} else {
			key := ""
			if config.ConfigKey != nil {
				key = *config.ConfigKey
			}
			fmt.Println(key)
		}
	}
}

func printVerboseList(configs []ConfigItem) {
	// 详细输出：显示所有信息
	var currentProject *string
	var currentEnv *string
	var currentModule *string

	for _, config := range configs {
		// 显示项目分组
		if !stringPtrEqual(config.Project, currentProject) {
			currentProject = config.Project
			currentEnv = nil
			currentModule = nil
			projectStr := ""
			if currentProject != nil {
				projectStr = *currentProject
			}
			fmt.Printf("[Project: %s]\n", projectStr)
		}

		// 显示环境分组
		if !stringPtrEqual(config.Env, currentEnv) {
			currentEnv = config.Env
			currentModule = nil
			envStr := ""
			if currentEnv != nil {
				envStr = *currentEnv
			}
			fmt.Printf("  [Environment: %s]\n", envStr)
		}

		// 显示模块分组
		if !stringPtrEqual(config.Module, currentModule) {
			currentModule = config.Module
			moduleStr := ""
			if currentModule != nil {
				moduleStr = *currentModule
			}
			fmt.Printf("    [Module: %s]\n", moduleStr)
		}

		// 显示配置项详情
		keyStr := ""
		if config.ConfigKey != nil {
			keyStr = *config.ConfigKey
		}
		valueStr := ""
		if config.ConfigValue != nil {
			valueStr = *config.ConfigValue
		}

		fmt.Printf("      Key: %s\n", keyStr)
		fmt.Printf("      Value: %s\n", valueStr)
		if config.ConfigAlias != nil && *config.ConfigAlias != "" {
			fmt.Printf("      Alias: %s\n", *config.ConfigAlias)
		}
		if config.AutoAlias != nil && *config.AutoAlias != "" {
			fmt.Printf("      AutoAlias: %s\n", *config.AutoAlias)
		}
		fmt.Println()
	}
}

// Helper function to compare two string pointers for equality
func stringPtrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func HandleListProjects() {
	// 获取所有项目名
	rows, err := db.DB.Query("SELECT DISTINCT project FROM config_master WHERE project IS NOT NULL")
	if err != nil {
		log.Fatalf("Failed to query projects: %v", err)
	}
	defer rows.Close()

	var projects []*string
	for rows.Next() {
		var project *string
		err := rows.Scan(&project)
		if err != nil {
			log.Fatalf("Failed to scan project name: %v", err)
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating projects: %v", err)
	}

	// fmt.Println("Projects:")
	for _, project := range projects {
		if project != nil {
			fmt.Printf("%s\n", *project)
		}
	}
}

func HandleListEnvs(project string) {
	// 基础查询语句
	query := "SELECT DISTINCT env FROM config_master WHERE env IS NOT NULL"
	args := []interface{}{}

	// 只有当project参数不为空且不为默认值时才添加条件
	if project != "" && project != "default" {
		query += " AND project = ?"
		args = append(args, project)
	}

	// 执行查询
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Fatalf("Failed to query environments: %v", err)
	}
	defer rows.Close()

	var envs []*string
	for rows.Next() {
		var env *string
		err := rows.Scan(&env)
		if err != nil {
			log.Fatalf("Failed to scan environment name: %v", err)
		}
		envs = append(envs, env)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating environments: %v", err)
	}

	// 输出结果
	if len(envs) == 0 {
		fmt.Println("No environments found.")
		return
	}

	// // 根据是否有project参数决定输出格式
	// if project != "" && project != "default" {
	// 	fmt.Printf("Environments for Project '%s':\n", project)
	// } else {
	// 	fmt.Println("All Environments:")
	// }

	for _, env := range envs {
		if env != nil {
			fmt.Printf("%s\n", *env)
		}
	}
}

func HandleListModules(project, env string) {
	// 基础查询语句
	query := "SELECT DISTINCT module FROM config_master WHERE module IS NOT NULL"
	args := []interface{}{}

	// 添加project条件（如果参数有效）
	if project != "" && project != "default" {
		query += " AND project = ?"
		args = append(args, project)
	}

	// 添加env条件（如果参数有效）
	if env != "" && env != "default" {
		query += " AND env = ?"
		args = append(args, env)
	}

	// 执行查询
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Fatalf("Failed to query modules: %v", err)
	}
	defer rows.Close()

	var modules []*string
	for rows.Next() {
		var module *string
		err := rows.Scan(&module)
		if err != nil {
			log.Fatalf("Failed to scan module name: %v", err)
		}
		modules = append(modules, module)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating modules: %v", err)
	}

	// 输出结果
	if len(modules) == 0 {
		fmt.Println("No modules found.")
		return
	}

	// 根据参数情况决定输出标题
	title := "All Modules:"
	if project != "" && project != "default" && env != "" && env != "default" {
		title = fmt.Sprintf("Modules for Project '%s' and Environment '%s':", project, env)
	} else if project != "" && project != "default" {
		title = fmt.Sprintf("Modules for Project '%s':", project)
	} else if env != "" && env != "default" {
		title = fmt.Sprintf("Modules for Environment '%s':", env)
	}

	log.Info(title)
	for _, module := range modules {
		if module != nil {
			fmt.Printf("%s\n", *module)
		}
	}
}
