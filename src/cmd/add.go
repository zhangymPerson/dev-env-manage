package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/zhangymPerson/dev-env-manage/src/constant"
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

	currentTime := time.Now()

	// Create config using the updated ConfigMaster struct
	config := models.ConfigMaster{
		Project:     constant.ToStrPtr(project),
		Env:         constant.ToStrPtr(env),
		Module:      constant.ToStrPtr(module),
		ConfigKey:   constant.ToStrPtr(key),
		ConfigValue: constant.ToStrPtr(value),
		ConfigAlias: constant.ToStrPtr(alias),
		AutoAlias:   constant.ToStrPtr(generateDefaultAlias(key)),
		ConfigType:  constant.ToStrPtr("string"),
		IsEncrypted: constant.ToIntPtr(1),
		Description: nil, // Set to nil or provide a value if needed
		SortOrder:   nil, // Set to nil or provide a value if needed
		CreatedTime: constant.ToTimePtr(currentTime),
		UpdatedTime: constant.ToTimePtr(currentTime),
	}

	if err := db.AddConfig(config); err != nil {
		log.Fatal("Failed to add config: %v", err)
	}
	log.Debug("Added config: %+v\n", config)
}

func generateDefaultAlias(key string) string {
	parts := strings.Split(key, ".")
	if len(parts) == 1 {
		return key
	}
	var aliasParts []string
	for _, part := range parts {
		if len(part) > 0 {
			aliasParts = append(aliasParts, string(part[0]))
		}
	}
	return strings.Join(aliasParts, ".")
}
