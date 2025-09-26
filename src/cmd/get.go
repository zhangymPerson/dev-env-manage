package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/zhangymPerson/dev-env-manage/src/db"
	"github.com/zhangymPerson/dev-env-manage/src/models"
)

func HandleGetCommand(project, env, module string, verbose bool, key string) {
	if len(os.Args) < 3 {
		log.Fatal("Usage: dem get <key>")
	}

	var config models.ConfigMaster
	err := db.DB.QueryRow(`
		SELECT project_code, env_code, module_code, config_key, config_value, config_alias 
		FROM config_master 
		WHERE project_code=? AND env_code=? AND module_code=? AND config_key=? AND is_deleted=0`,
		project, env, module, key,
	).Scan(&config.ProjectCode, &config.EnvCode, &config.ModuleCode, &config.ConfigKey, &config.ConfigValue, &config.ConfigAlias)

	switch {
	case err == sql.ErrNoRows:
		log.Fatalf("Config not found for key: %s", key)
	case err != nil:
		log.Fatalf("Failed to get config: %v", err)
	}

	if verbose {
		fmt.Printf("Config details:\nProject: %s\nEnv: %s\nModule: %s\nKey: %s\nValue: %s\nAlias: %s\n",
			config.ProjectCode, config.EnvCode, config.ModuleCode, config.ConfigKey, config.ConfigValue, config.ConfigAlias)
	} else {
		fmt.Println(config.ConfigValue)
	}
}
