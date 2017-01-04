package main

import (
	"fmt"
	"github.com/vimond/configmap-generator"
	"gopkg.in/urfave/cli.v1"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	//configmap_generator.Vault()

	app := cli.NewApp()
	app.Name = "configmap-generator (cmapgen)"
	app.Version = "1.0.0"
	app.Usage = "Generates config maps for Kubernetes from Ansible"
	app.Flags = []cli.Flag {
		cli.BoolFlag{
			Name: "noheader, H",
			Usage: "Do not show header",
		},
	}
	/*
	app.Action = func(c *cli.Context) error {
		if c.Bool("list") {
			listNames()
		} else {
			fmt.Println("tis false")
		}
		return nil
	}
	*/
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Usage:   "list supported apps",
			Action:  func(c *cli.Context) error {
				if !c.GlobalBool("noheader") {
					fmt.Println("Supported applications\n----------------------")
				}

				listNames()
				return nil
			},
		},
		{
			Name: "generate",
			Usage: "Generate configmap",
			Flags:   []cli.Flag {
				cli.StringFlag{
					Name: "name, n",
					Usage: "`Name` of application",

				},
				cli.StringFlag{
					Name: "environment, e",
					Usage: "`Environment` to use",
					EnvVar: "CG_ENV",
				},
				cli.StringFlag{
					Name: "ansible, a",
					Usage: "`Folder` where vimond-ansible is",
					EnvVar: "CG_VIMOND_ANSIBLE_FOLDER",
				},
			},
			Action:  func(c *cli.Context) error {
				if !c.GlobalBool("noheader") {
					fmt.Println("Generating configMap\n----------------------")
				}
				generateConfigMap(c.String("name"), c.String("environment"), c.String("ansible"))
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func listNames() {
	appConfig := configmap_generator.LoadConfig()
	fmt.Print(strings.Join(appConfig.AppNames(), "\n"))
}

func generateConfigMap(name, env, ansibleFolder string) {
	appConfig := configmap_generator.LoadConfig()
	allVars := configmap_generator.LoadVars(ansibleFolder, env)
	vars := configmap_generator.FilterVariables(appConfig, allVars, name)
	app := configmap_generator.ConfigMapData{
		AppName: name,
		Vars: vars,
	}
	configmap_generator.Generate(app)
}
