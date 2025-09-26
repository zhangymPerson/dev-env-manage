package cmd

import (
	"os"
	"strings"

	"github.com/zhangymPerson/dev-env-manage/src/db"
	"github.com/zhangymPerson/dev-env-manage/src/log"
	"github.com/zhangymPerson/dev-env-manage/src/models"
)

// HandleAddCommand handles the add command
func HandleAddCommand(project, env, module string, key, alias, value string) {
	log.Info("key: %s, value: %s, alias: %s", key, value, alias)
	if len(os.Args) < 4 {
		log.Fatal("Usage: dem add <key> <value> [alias]")
	}
	if len(os.Args) > 4 {
		alias = os.Args[4]
	} else {
		alias = generateDefaultAlias(key)
	}

	config := models.ConfigMaster{
		ProjectCode: project,
		EnvCode:     env,
		ModuleCode:  module,
		ConfigKey:   key,
		ConfigValue: value,
		ConfigAlias: alias,
		AutoAlias:   generateDefaultAlias(key),
		ConfigType:  "string",
		IsEncrypted: 0,
		IsDeleted:   0,
	}

	if err := db.AddConfig(config); err != nil {
		log.Fatal("Failed to add config: %v", err)
	}
	log.Debug("Added config: %+v\n", config)
}

func generateDefaultAlias(key string) string {
	parts := strings.Split(key, ".")
	var aliasParts []string
	for _, part := range parts {
		if len(part) > 0 {
			aliasParts = append(aliasParts, string(part[0]))
		}
	}
	return strings.Join(aliasParts, ".")
}
