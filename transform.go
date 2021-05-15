package main

import (
	"regexp"
)

func TransformBody(postBody string, route Route) (ret string) {
	ret = ""

	find := route.RegexpFind
	replace := route.RegexpReplace
	if find == "" {
		return postBody
	}

	re := regexp.MustCompile(find)
	ret = re.ReplaceAllString(postBody, replace)

	return ret
}
