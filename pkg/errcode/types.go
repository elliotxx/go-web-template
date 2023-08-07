package errcode

import "fmt"

const InvalidScope = "999"

var codes = map[string]struct{}{}

func mustSetCodeIfNotPresent(code string) {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("The error code %s already exists, please change one", code))
	}
	codes[code] = struct{}{}
}
