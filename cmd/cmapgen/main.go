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

var Version string


func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	groupVarsFlag := cli.StringFlag{
		Name: "group-vars, g",
		Usage: "`Folder` where group_vars reside (required)",
		EnvVar: "CG_GROUP_VARS_FOLDER",
	}
	
	app := cli.NewApp()
	app.Name = "configmap-generator (cmapgen)"
	app.Version = Version
	app.Usage = "Generates config maps for Kubernetes from Ansible"
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "config, c",
			Usage: "Application config file to use",
			EnvVar: "CG_APP_CFG",
			Value: "./config/app-config.yml",
			
		},
		cli.BoolFlag{
			Name: "noheader, H",
			Usage: "Do not show header",
		},
	}
	
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Usage:   "list supported apps",
			Action:  func(c *cli.Context) error {
				config := configmap_generator.New(c.GlobalString("config"))
				return listNames(c, config)
			},
		},
		{
			Name:    "prefixes",
			Usage:   "list sorted prefixes in the var store",
			Action:  listPrefixes,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name: "levels, l",
					Usage: "How many `_` to nest prefixes",
					Value: 1,
				},
				groupVarsFlag,
			},
		},
		{
			Name: "generate",
			Usage: "Generate configmap",
			Flags:   []cli.Flag {
				cli.StringFlag{
					Name: "name, n",
					Usage: "`Name` of application, or 'all' (required)",

				},
				cli.StringFlag{
					Name: "environment, e",
					Usage: "`Environment` to use (required)",
					EnvVar: "CG_ENV",
				},
				groupVarsFlag,
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
				config := configmap_generator.New(c.GlobalString("config"))
				if !config.CheckNameExists(name) && name != "all" {
					return cli.NewExitError("Error: App not found: " + name, 2)
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
				generateConfigMap(name, env, groupVarsFolder, vaultPassword, config)
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

func listPrefixes(c *cli.Context) error{
	groupVarsFolder, err := checkRequiredArg("group-vars", c.String("group-vars"))
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	
	prefixes, err := configmap_generator.SuggestConfig(groupVarsFolder, c.Int("levels"))
	fmt.Println(strings.Join(prefixes, "\n"))
	
	return err
}


func listNames(c *cli.Context, appConfig *configmap_generator.AppConfig) error{
	if !c.GlobalBool("noheader") {
		fmt.Println("Supported applications\n----------------------")
		fmt.Println("all (to show all apps combined)")
	}
	
	fmt.Print(strings.Join(appConfig.AppNames(), "\n"))
	return nil
}


func generateConfigMap(name, env, groupVarsFolder, vaultPassword string, appConfig *configmap_generator.AppConfig) {
	
	allVars := configmap_generator.LoadVars(groupVarsFolder, env, vaultPassword)
	var result string
	if name != "all" {
		result = getConfigMap(name, allVars, appConfig)
	} else {
		result = getAllConfigMaps(allVars, appConfig)
	}
	fmt.Println(result)
}

func getConfigMap(name string, allVars configmap_generator.Variables, appConfig *configmap_generator.AppConfig) (string) {
	allVars["service_name"] = configmap_generator.VarVal(name)
	allVars = configmap_generator.SubstituteVars(allVars)
	vars := configmap_generator.FilterVariables(appConfig, allVars, name)
	app := configmap_generator.ConfigMapData{
		AppName: name,
		Vars: vars,
	}
	return configmap_generator.Generate(app)
}

func getAllConfigMaps(allVars configmap_generator.Variables, appConfig *configmap_generator.AppConfig) (string) {
	configMaps := make([]string, len(appConfig.Applications))
	for i, v := range appConfig.Applications {
		configMaps[i] = getConfigMap(v.Name, allVars, appConfig)
	}
	return strings.Join(configMaps, "---\n")
}
