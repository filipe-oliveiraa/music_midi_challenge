package conf

import "strings"

type Case = int8

const (
	LowerCase Case = iota
	UpperCase
	NoCase
)

func joinParents(parents []string, prefix, name, sep string, c Case) string {
	var elms []string

	if prefix != "" {
		elms = append(elms, prefix)
	}

	elms = append(elms, parents...)
	elms = append(elms, name)

	envRaw := strings.Join(elms, sep)

	env := envRaw
	switch c {
	case LowerCase:
		env = strings.ToLower(envRaw)
	case UpperCase:
		env = strings.ToUpper(envRaw)
	}

	return env
}
