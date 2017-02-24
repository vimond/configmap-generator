package configmap_generator

import (
	"text/template"
	"bytes"
)

type ConfigMapData struct {
	AppName string
	Vars Variables
}

func Generate(data ConfigMapData) (string)  {
	tmpl, err := template.New("ConfigMap.tmpl").ParseFiles("config/ConfigMap.tmpl")
	if err != nil { panic(err) }
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, data)
	if err != nil { panic(err) }
	return doc.String()
}
