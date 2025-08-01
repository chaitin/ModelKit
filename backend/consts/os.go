package consts

import "strings"

type OSType string

func (o OSType) Name() string {
	if strings.HasPrefix(string(o), "Windows") {
		return "Windows"
	}
	if strings.HasPrefix(string(o), "Darwin") {
		return "MacOS"
	}
	return string(o)
}

type OSRelease string
