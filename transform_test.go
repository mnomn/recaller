package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTransIFTTT2AIO(t *testing.T) {
	// Transfer from ifttt format {value1:123, value2:321, value3:987}
	// to adafruit.io {value:123 ...} (assuming value2 and value 3 can be left in the body)
	// Also: Convert value 2 to {value2:321}. That way one post to route2cloud
	// can result in two posts to "cloud".

	var route Route
	route_string := `{"in":"regextest","regexpFind":"value1", "regexpReplace":"value"}`
	if err := json.Unmarshal([]byte(route_string), &route); err != nil {
		t.Errorf("Bad input! %v", err)
	}
	body := `{"value1":123, "value2":321, "value3":987}`
	res := TransformBody(body, route)
	if strings.Index(res, `"value":123`) < 0 {
		t.Errorf("Transformation 1 failed %v\n", res)
	}

	route_string = `{"in":"regextest","regexpFind":"value2", "regexpReplace":"value"}`
	if err := json.Unmarshal([]byte(route_string), &route); err != nil {
		t.Errorf("Bad input! %v", err)
	}
	res = TransformBody(body, route)
	if strings.Index(res, `"value":321`) < 0 {
		t.Errorf("Transformation 2 failed %v\n", res)
	}
}

// func TestTransformBody(postString string, url_config map[string]interface{}) {

func TestTransformBody(t *testing.T) {
	var route Route
	body := `{"value1":123, "value2":321, "value3":987}`
	raw := `{"in":"regextest","regexpFind":"value1", "regexpReplace":"banan"}`
	if err := json.Unmarshal([]byte(raw), &route); err != nil {
		t.Errorf("Bad input! %v", err)
	}
	newBody := TransformBody(body, route)

	if newBody != `{"banan":123, "value2":321, "value3":987}` {
		t.Errorf("Bad TransformBody %v", newBody)
	}

}
