package main

import (
	"log"
	"regexp"
)

func TransformBody(postBody string, route_config map[string]interface{}) (ret string) {
	ret = ""

	find, exist1 := route_config["regexp_find"]
	replace, exist2 := route_config["regexp_replace"]
	if exist1 && !exist2 {
		log.Printf("No regexp_replace in regexp_replace %d", route_config["in"])
	}
	if !exist1 || !exist2 {
		return postBody
	}

	re := regexp.MustCompile(find.(string))
	ret = re.ReplaceAllString(postBody, replace.(string))

	return ret
}
