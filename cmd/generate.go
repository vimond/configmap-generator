package cmd

import (
	"fmt"
	"log"
	"errors"
	"strings"

	"io/ioutil"

	"github.com/spf13/cobra"
	configmap_generator "github.com/vimond/configmap-generator/generator"
)

var (
	config			configmap_generator.GeneratorConfig
	vaultPasswordFile	string
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
		var err error

		if err = checkRequiredArg("name", config.AppName); err != nil {
			log.Fatal(err)
		}
		config.AppConfig, err = configmap_generator.New(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		if !config.AppConfig.CheckNameExists(config.AppName) && config.AppName != "all" {
			log.Fatal(errors.New("Error: App not found: " + config.AppName))
		}
		if err = checkRequiredArg("environment", config.Environment); err != nil {
			log.Fatal(err)
		}
		if err = checkRequiredArg("group-vars", config.GroupVars); err != nil {
			log.Fatal(err)
		}
		if err = checkRequiredArg("vault-password-file", vaultPasswordFile); err != nil {
			log.Fatal(err)
		}
		config.VaultPassword, err = getVaultPassword(vaultPasswordFile)
		if err != nil {
			log.Fatal(err)
		}
		if !noHeader {
			fmt.Println("Generating configMap\n----------------------")
		}
		result, err := config.GenerateConfigMap()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	},
}

func init() {
	generateCmd.Flags().StringVarP(
		&config.GroupVars,
		"group-vars",
		"g",
		"",
		"`FOLDER` where group_vars reside (required)",
	)
	generateCmd.Flags().StringVarP(
		&config.AppName,
		"name",
		"n",
		"",
		"`NAME` of application, or 'all' (required)",
	)
	generateCmd.Flags().StringVarP(
		&config.Environment,
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
