package configmap_generator

import (
	"text/template"
	"os"
)

type ConfigMapData struct {
	AppName string
	Vars map[string]interface{}
}

func Generate(data ConfigMapData)  {
	tmpl, err := template.New("ConfigMap.tmpl").ParseFiles("config/ConfigMap.tmpl")
	if err != nil { panic(err) }
	err = tmpl.Execute(os.Stdout, data)
	if err != nil { panic(err) }
}
