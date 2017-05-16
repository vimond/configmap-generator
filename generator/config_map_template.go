package configmap_generator

import (
	"text/template"
	"bytes"
	"github.com/Masterminds/sprig"
)

type ConfigMapData struct {
	AppName string
	Data string
}

func Generate(data ConfigMapData) (string)  {

	tmpl, err := template.New("ConfigMap.tmpl").Funcs(sprig.TxtFuncMap()).ParseFiles("config/ConfigMap.tmpl")
	if err != nil { panic(err) }
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, data)
	if err != nil { panic(err) }
	return doc.String()
}
