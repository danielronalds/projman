package services

import "strings"

type Sanitiser struct{}

func NewSanitiser() Sanitiser {
	return Sanitiser{}
}

func (s Sanitiser) Sanitise(name string) string {
	if name == "/" {
		return "root"
	}

	name = strings.TrimLeft(name, ".")
	name = strings.ReplaceAll(name, ":", "-")

	if name == "" {
		return "session"
	}

	return name
}
