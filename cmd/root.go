package cmd

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/vimond/configmap-generator/generator"
)

var (
	cfgFile		string
	noHeader	bool
)

var RootCmd = &cobra.Command{
	Use:   "configmap-generator",
	Short: "Generates config maps for Kubernetes from Ansible",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		"./config/app-config.yml",
		"Application `CONFIG_FILE` to use",
	)
	RootCmd.PersistentFlags().BoolVarP(&noHeader, "noheader", "H", false, "Do not show header")
	RootCmd.PersistentFlags().BoolVarP(&configmap_generator.Debug, "debug", "d", false, "Print debug messages to console")
}
