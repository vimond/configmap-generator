package cmd

import (
	"log"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vimond/configmap-generator/generator"
)


var prefixesCmd = &cobra.Command{
	Use:   "prefixes",
	Short: "list sorted prefixes in the var store",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		groupVars, err := cmd.Flags().GetString("group-vars")
		if err != nil {
			log.Fatal(err)
		}
		levels, err := cmd.Flags().GetInt("levels")
		if err != nil {
			log.Fatal(err)
		}
		if err := listPrefixes(groupVars, levels); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	prefixesCmd.Flags().StringP(
		"group-vars",
		"g",
		"",
		"`Folder` where group_vars reside (required)",
	)
	prefixesCmd.Flags().IntP(
		"levels",
		"l",
		1,
		"How many '_' to nest prefixes",
	)
	RootCmd.AddCommand(prefixesCmd)
}

func listPrefixes(groupVars string, nestedLevels int) (err error) {
	if err := checkRequiredArg("group-vars", groupVars); err != nil {
		return err
	}
	prefixes, err := configmap_generator.SuggestConfig(groupVars, nestedLevels)
	fmt.Println(strings.Join(prefixes, "\n"))
	return
}