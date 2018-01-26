package controller

import (
	"html/template"
	"io/ioutil"

	"github.com/go-playground/log"
)

const (
	errorPolicyReturn = iota
	errorPolicyPanic
)

func loadTemplate(name, path string, errorPolicy int32) (*template.Template, error) {
	log.Debug("loading template [", name, "] from ", path)
	templateStr, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("error loading tempalte", err)
		if errorPolicy == errorPolicyPanic {
			panic(err)
		}
		return nil, err
	}
	t, err := template.New(name).Parse(string(templateStr))
	if err != nil {
		log.Error("error loading tempalte", err)
		if errorPolicy == errorPolicyPanic {
			panic(err)
		}
	}
	return t, err
}
