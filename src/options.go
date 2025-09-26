package src

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/zhangymPerson/dev-env-manage/src/cmd"
	"github.com/zhangymPerson/dev-env-manage/src/constant"
	"github.com/zhangymPerson/dev-env-manage/src/log"
)

type Config struct {
	Key    string
	Value  string
	Alias  string
	Env    string
	Module string
}

var buildTime = time.Now().String() // 默认值
func Options(gitBranch string, gitCommit string) {
	// Define flags
	defaultValue := constant.EnvDefault.String()
	project := flag.String("p", defaultValue, "Specify project name")
	env := flag.String("e", defaultValue, "Specify environment type [dev|test|prod|other|default]")
	module := flag.String("m", defaultValue, "Specify module name")
	verbose := flag.Bool("v", false, "Enable verbose output")
	version := flag.Bool("version", false, "Show version and build information")
	alias := flag.String("alias", "", "Specify custom alias for the config")
	// configPath := flag.String("config", "", "Specify config file path")
	flag.Parse()

	if *verbose {
		log.SetDebug()
	}

	if *version {
		if len(gitCommit) > 8 {
			gitCommit = gitCommit[:8]
		}
		printVersionInfo(gitBranch, gitCommit) // 提取为独立函数
		os.Exit(0)
	}

	// Check for command
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	// Handle commands
	switch os.Args[1] {
	case "add", "create":
		if len(os.Args) < 3 {
			fmt.Println("Usage: dem add <key> <value>")
			os.Exit(1)
		}
		key := os.Args[2]
		value := strings.Join(os.Args[3:], " ")
		cmd.HandleAddCommand(*project, *env, *module, key, *alias, value)
	case "get", "retrieve":
		if len(os.Args) < 3 {
			fmt.Println("Usage: dem get <key>")
			os.Exit(1)
		}
		key := os.Args[2]
		cmd.HandleGetCommand(*project, *env, *module, *verbose, key)
	case "delete", "remove":
		handleDeleteCommand(*project, *env, *module, *verbose)
	case "list", "ls":
		handleListCommand(*project, *env, *module, *verbose)
	case "info":
		handleInfoCommand(*project, *env, *module, *verbose)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

// 美化版本信息输出（新增函数）
func printVersionInfo(branch, commit string) {
	fmt.Print(constant.VersionHeader)
	fmt.Println(constant.AppDesc)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, constant.VersionTable, "Version:", constant.Version)
	fmt.Fprintf(w, constant.VersionTable, "Git Repo:", constant.GitRepo)
	fmt.Fprintf(w, constant.VersionTable, "Branch:", branch)
	fmt.Fprintf(w, constant.VersionTable, "Commit:", commit)
	fmt.Fprintf(w, constant.VersionTable, "Build Time:", buildTime)
	w.Flush()
}
func printHelp() {
	helpText := `
Usage: dem [OPTIONS] COMMAND [ARGS]...

Key-value configuration management tool

Options:
  -h, --help                     Show this help message and exit
  -p, --project TEXT             Specify project name (default: default)
  -e, --env [dev|test|prod|other|default]
                                 Specify environment type (default: default)
  -m, --module TEXT             Specify module name (default: default)
  -v, --verbose                 Enable verbose output
  --config PATH                 Specify config file path
  --version                     Show version and build information

Commands:
  add, create                   Add key-value configuration (Usage: dem add <key> <value>)
  get, retrieve                Get key-value configuration (Usage: dem get <key>)
  delete, remove               Delete key-value configuration
  list, ls                     List all configurations
  info                         Show configuration details`
	fmt.Println(helpText)
}

func handleAddCommand(project, env, module string, verbose bool) {
	// Implementation for add command
}

func handleGetCommand(project, env, module string, verbose bool) {
	// Implementation for get command
}

func handleDeleteCommand(project, env, module string, verbose bool) {
	// Implementation for delete command
}

func handleListCommand(project, env, module string, verbose bool) {
	// Implementation for list command
}

func handleInfoCommand(project, env, module string, verbose bool) {
	// Implementation for info command
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
