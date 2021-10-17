package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

func TransformBody(postBody string, route Route) (ret string, err error) {
	err = nil

	if route.BodyTemplate == "" {
		return postBody, nil
	}

	// Todo: reuse
	tmpl, err := template.New("bodyTemplate").Option("missingkey=error").Parse(route.BodyTemplate)
	if err != nil {
		fmt.Printf("Failed to parse %v\n", err)
	}
	json_in := map[string]interface{}{}
	if err = json.Unmarshal([]byte(postBody), &json_in); err != nil {
		fmt.Printf("Invalid imput json %v\n", err)
		return
	}

	var resultBuffer bytes.Buffer
	if err = tmpl.Execute(&resultBuffer, json_in); err != nil {
		fmt.Printf("Failed convert input body with template %v\n", err)
		return
	}

	ret = resultBuffer.String()
	return
}
