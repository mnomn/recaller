package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTransformBody(t *testing.T) {
	var route Route
	body := `
	{
		"sensor":"S4",
		"values": {
			"T":23.4,
			"unit":"C"
		}
	}`
	raw := `
	{
		"in":"test",
		"bodyTemplate":"sensor_values,sensor_id={{.sensor}} temperature={{.values.T}},client=r2c"
	}`
	if err := json.Unmarshal([]byte(raw), &route); err != nil {
		panic(err)
	}
	transformedBody, err := TransformBody(body, route)

	if err != nil || !strings.Contains(transformedBody, "sensor") || !strings.Contains(transformedBody, "temperature=23.4") {
		t.Errorf("Bad TransformBody %v", transformedBody)
	}
}

func TestTransformWrongInput(t *testing.T) {
	// Verify graceful handeling if incomming json does not fit template
	var route Route
	body := `{"bananas":11}`

	raw := `
	{
		"in":"test2",
		"bodyTemplate":"sensor_values,sensor_id={{.sensor}} temperature={{.values.T}},client=r2c"
	}`
	if err := json.Unmarshal([]byte(raw), &route); err != nil {
		panic(err)
	}

	transformedBody, err := TransformBody(body, route)

	if err == nil || transformedBody != "" {
		t.Error("Expected error, wrong input")
	}
}

func TestNoTemplate(t *testing.T) {
	// Body should be intact if there is no bodyTemplate
	var route Route
	body := `
	{
		"bananas":11,
		"apples": 22
	}`
	raw := `
	{
		"in":"test3"
	}`
	if err := json.Unmarshal([]byte(raw), &route); err != nil {
		panic(err)
	}

	transformedBody, err := TransformBody(body, route)

	if err != nil || body != transformedBody {
		t.Error("Expected error, wrong input")
	}
}
