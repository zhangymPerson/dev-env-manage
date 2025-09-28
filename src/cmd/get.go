package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/zhangymPerson/dev-env-manage/src/db"
	"github.com/zhangymPerson/dev-env-manage/src/log"
	"github.com/zhangymPerson/dev-env-manage/src/models"
)

func HandleGetCommand(project, env, module string, verbose bool, key string) {
	if key == "" {
		fmt.Println("Usage: dem get <key>")
	}

	// 构建查询条件和参数
	query, params := buildQueryConditions(project, env, module, key)

	// 输出最终执行的SQL
	if verbose {
		log.Info("=== 执行的SQL查询 ===\n")
		log.Info("第一级查询 (config_key): %s\n", query.configKeyQuery)
		log.Info("参数: %v\n", append([]interface{}{}, params...))
		log.Info("第二级查询 (config_alias): %s\n", query.configAliasQuery)
		log.Info("参数: %v\n", append([]interface{}{}, params...))
		log.Info("第三级查询 (auto_alias): %s\n", query.autoAliasQuery)
		log.Info("参数: %v\n", append([]interface{}{}, params...))
		log.Info("==================\n\n")
	}

	// 三级查询逻辑：config_key -> config_alias -> auto_alias
	// 收集所有匹配的配置项
	var configs []models.ConfigMaster

	// 第一级：查询 config_key = key
	rows, err := db.DB.Query(query.configKeyQuery, params...)
	if err != nil {
		log.Fatalf("Failed to execute config_key query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var config models.ConfigMaster
		err = rows.Scan(
			&config.ProjectCode, &config.EnvCode, &config.ModuleCode,
			&config.ConfigKey, &config.ConfigValue, &config.ConfigAlias, &config.AutoAlias,
		)
		if err != nil {
			log.Fatalf("Failed to scan config: %v", err)
		}
		configs = append(configs, config)
	}

	printInfo(configs, verbose)

	// 第二级：查询 config_alias = key
	rows, err = db.DB.Query(query.configAliasQuery, params...)
	if err != nil {
		log.Fatalf("Failed to execute config_alias query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var config models.ConfigMaster
		err = rows.Scan(
			&config.ProjectCode, &config.EnvCode, &config.ModuleCode,
			&config.ConfigKey, &config.ConfigValue, &config.ConfigAlias, &config.AutoAlias,
		)
		if err != nil {
			log.Fatalf("Failed to scan config: %v", err)
		}
		configs = append(configs, config)
	}
	printInfo(configs, verbose)
	// 第三级：查询 auto_alias = key
	rows, err = db.DB.Query(query.autoAliasQuery, params...)
	if err != nil {
		log.Fatalf("Failed to execute auto_alias query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var config models.ConfigMaster
		err = rows.Scan(
			&config.ProjectCode, &config.EnvCode, &config.ModuleCode,
			&config.ConfigKey, &config.ConfigValue, &config.ConfigAlias, &config.AutoAlias,
		)
		if err != nil {
			log.Fatalf("Failed to scan config: %v", err)
		}
		configs = append(configs, config)
	}
	printInfo(configs, verbose)

}

func printInfo(configs []models.ConfigMaster, verbose bool) {
	// 根据匹配数量决定输出格式
	if len(configs) == 0 {
		// 不存在匹配项，返回空
		return
	} else if len(configs) == 1 {
		// 存在且数量为1，只输出value
		fmt.Printf("%s\n", configs[0].ConfigValue)
		os.Exit(0)
	} else {
		// 存在且数量不为1，执行详细输出逻辑
		for _, config := range configs {
			if verbose {
				fmt.Printf("Config details:\nProject: %s\nEnv: %s\nModule: %s\nKey: %s\nValue: %s\nAlias: %s\nAutoAlias: %s\n\n",
					config.ProjectCode, config.EnvCode, config.ModuleCode, config.ConfigKey, config.ConfigValue, config.ConfigAlias, config.AutoAlias)
			} else {
				fmt.Println(config.ConfigValue)
			}
		}
	}
}

// QueryConditions 定义查询条件和参数
type QueryConditions struct {
	configKeyQuery   string
	configAliasQuery string
	autoAliasQuery   string
}

// buildQueryConditions 根据参数构建查询条件
func buildQueryConditions(project, env, module, key string) (QueryConditions, []interface{}) {
	var conditions []string
	var params []interface{}

	// 添加项目条件（如果不是默认值）
	if project != "default" {
		conditions = append(conditions, "project_code=?")
		params = append(params, project)
	}

	// 添加环境条件（如果不是默认值）
	if env != "default" {
		conditions = append(conditions, "env_code=?")
		params = append(params, env)
	}

	// 添加模块条件（如果不是默认值）
	if module != "default" {
		conditions = append(conditions, "module_code=?")
		params = append(params, module)
	}

	// 构建基础查询模板
	baseQuery := "SELECT project_code, env_code, module_code, config_key, config_value, config_alias, auto_alias FROM config_master WHERE is_deleted=0"

	// 如果有条件，添加AND连接
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// 构建三级查询的完整SQL
	configKeyQuery := baseQuery + " AND config_key=?"
	configAliasQuery := baseQuery + " AND config_alias=?"
	autoAliasQuery := baseQuery + " AND auto_alias=?"

	// 为每个查询添加key参数
	configKeyParams := append([]interface{}{}, params...)
	configKeyParams = append(configKeyParams, key)

	return QueryConditions{
		configKeyQuery:   configKeyQuery,
		configAliasQuery: configAliasQuery,
		autoAliasQuery:   autoAliasQuery,
	}, configKeyParams
}
