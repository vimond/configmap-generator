package main

import (
	"fmt"
	"github.com/vimond/configmap-generator"
	"gopkg.in/urfave/cli.v1"
	"os"
	"strings"
	"errors"
	"io/ioutil"
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
					Usage: "`Name` of application (required)",

				},
				cli.StringFlag{
					Name: "environment, e",
					Usage: "`Environment` to use (required)",
					EnvVar: "CG_ENV",
				},
				cli.StringFlag{
					Name: "group-vars, g",
					Usage: "`Folder` where group_vars reside (required)",
					EnvVar: "CG_GROUP_VARS_FOLDER",
				},
				cli.StringFlag{
					Name: "vault-password-file, p",
					Usage: "load password from `VAULT_PASSWORD_FILE`",
					EnvVar: "CG_VAULT_PASSWORD_FILE",
				},
			},
			Action:  func(c *cli.Context) error {
				name, err := checkRequiredArg("name", c.String("name"))
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				env, err := checkRequiredArg("environment", c.String("environment"))
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				groupVarsFolder, err := checkRequiredArg("group-vars", c.String("group-vars"))
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				vaultPassword := getVaultPassword(c.String("vault-password-file"))
				if !c.GlobalBool("noheader") {
					fmt.Println("Generating configMap\n----------------------")
				}
				generateConfigMap(name, env, groupVarsFolder, vaultPassword)
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func getVaultPassword(vaultPasswordFile string) (string) {
	if vaultPasswordFile == "" {
		return ""
	} else {
		data, err := ioutil.ReadFile(vaultPasswordFile)
		check(err)
		return strings.TrimSpace(string(data))
	}
}

func checkRequiredArg(name, value string) (string, error) {
	if value == "" {
		msg := "Error: missing required argument: " + name
		return "", errors.New(msg)
	} else {
		return value, nil
	}
}

func listNames() {
	appConfig := configmap_generator.LoadConfig()
	fmt.Print(strings.Join(appConfig.AppNames(), "\n"))
}

func generateConfigMap(name, env, groupVarsFolder, vaultPassword string) {
	appConfig := configmap_generator.LoadConfig()
	allVars := configmap_generator.LoadVars(groupVarsFolder, env, vaultPassword)
	allVars["service_name"] = name
	allVars = configmap_generator.SubstituteVars(allVars)
	vars := configmap_generator.FilterVariables(appConfig, allVars, name)

	app := configmap_generator.ConfigMapData{
		AppName: name,
		Vars: vars,
	}
	configmap_generator.Generate(app)
}
