package cmd

import (
	"log"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vimond/configmap-generator/generator"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list supported apps",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configmap_generator.New(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		if err := listNames(config); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func listNames(appConfig *configmap_generator.AppConfig) error{
	if !noHeader {
		fmt.Println("Supported applications\n----------------------")
		fmt.Println("all (to show all apps combined)")
	}
	fmt.Print(strings.Join(appConfig.AppNames(), "\n"))
	return nil
}