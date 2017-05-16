package cmd

import (
	"log"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vimond/configmap-generator/generator"
)

var (
	nestedLevels	int
	groupVars	string
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
		if err := listPrefixes(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	prefixesCmd.Flags().StringVarP(
		&groupVars,
		"group-vars",
		"g",
		"",
		"`Folder` where group_vars reside (required)",
	)
	prefixesCmd.Flags().IntVarP(
		&nestedLevels,
		"levels",
		"l",
		1,
		"How many '_' to nest prefixes",
	)
	RootCmd.AddCommand(prefixesCmd)
}

func listPrefixes() error{
	err := checkRequiredArg("group-vars", groupVars)
	if err != nil {
		return err
	}

	prefixes, err := configmap_generator.SuggestConfig(groupVars, nestedLevels)
	fmt.Println(strings.Join(prefixes, "\n"))

	return err
}