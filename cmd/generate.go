package cmd

import (
	"fmt"
	"log"
	"errors"
	"strings"

	"io/ioutil"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	"github.com/vimond/configmap-generator/generator"
)

var (
	appName           string
	environment       string
	vaultPasswordFile string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkRequiredArg("name", appName); err != nil {
			log.Fatal(err)
		}
		config := configmap_generator.New(cfgFile)
		if !config.CheckNameExists(appName) && appName != "all" {
			log.Fatal(errors.New("Error: App not found: " + appName))
		}
		if err := checkRequiredArg("environment", environment); err != nil {
			log.Fatal(err)
		}
		if err := checkRequiredArg("group-vars", groupVars); err != nil {
			log.Fatal(err)
		}
		if err := checkRequiredArg("vault-password-file", vaultPasswordFile); err != nil {
			log.Fatal(err)
		}
		vaultPassword, err := getVaultPassword(vaultPasswordFile)
		if err != nil {
			log.Fatal(err)
		}
		if !noHeader {
			fmt.Println("Generating configMap\n----------------------")
		}
		fmt.Println(generateConfigMap(appName, environment, groupVars, vaultPassword, config))
	},
}

func init() {
	generateCmd.Flags().StringVarP(
		&groupVars,
		"group-vars",
		"g",
		"",
		"`FOLDER` where group_vars reside (required)",
	)
	generateCmd.Flags().StringVarP(
		&appName,
		"name",
		"n",
		"",
		"`NAME` of application, or 'all' (required)",
	)
	generateCmd.Flags().StringVarP(
		&environment,
		"environment",
		"e",
		"",
		"`ENVIRONMENT` to use (required)",
	)
	generateCmd.Flags().StringVarP(
		&vaultPasswordFile,
		"vault-password-file",
		"p",
		"",
		"load password from `VAULT_PASSWORD_FILE`",
	)
	RootCmd.AddCommand(generateCmd)
}

func generateConfigMap(name, env, groupVarsFolder, vaultPassword string, appConfig *configmap_generator.AppConfig) (string){

	allVars := configmap_generator.LoadVars(groupVarsFolder, env, vaultPassword)
	var result string
	if name != "all" {
		result = getConfigMap(name, allVars, appConfig)
	} else {
		result = getAllConfigMaps(allVars, appConfig)
	}
	return result
}

func getConfigMap(name string, allVars map[string]interface{}, appConfig *configmap_generator.AppConfig) (string) {
	allVars["service_name"] = name
	allVars = configmap_generator.SubstituteVars(allVars)
	vars := configmap_generator.FilterVariables(appConfig, allVars, name)
	vars2,_ := yaml.Marshal(vars)
	app := configmap_generator.ConfigMapData{
		AppName: name,
		Data: string(vars2[:]),
	}

	return configmap_generator.Generate(app)
}

func getAllConfigMaps(allVars map[string]interface{}, appConfig *configmap_generator.AppConfig) (string) {
	configMaps := make([]string, len(appConfig.Applications))
	for i, v := range appConfig.Applications {
		configMaps[i] = getConfigMap(v.Name, allVars, appConfig)
	}
	return strings.Join(configMaps, "\n")
}


func checkRequiredArg(name, value string) (err error) {
	if value == "" {
		err = errors.New("Error: missing required argument: " + name)
	}
	return
}

func getVaultPassword(vaultPasswordFile string) (pw string, err error) {
	if vaultPasswordFile != "" {
		data, err := ioutil.ReadFile(vaultPasswordFile)
		if err == nil {
			pw = strings.TrimSpace(string(data))
		}
	}
	return
}
