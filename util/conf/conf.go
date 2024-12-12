package conf

import (
	"errors"
	"fmt"
	"strings"
)

// ParseConfig parses configurations from multiple sources
// and stores in v.
// The default opts reads from env variables and flag arguments.
func ParseConfig(v any, opts ...Option) (string, error) {
	fields, err := getFields(v)
	if err != nil {
		return "", err
	}

	_ = setFieldDefaultValues(fields)

	o := applyOptions(opts)

	for i := range o.sources {
		err := o.sources[i].Apply(v, fields)
		if err != nil {
			if errors.Is(err, ErrHelp) {
				return buildHelp(o, fields), ErrHelp
			}

			return "", err
		}
	}

	return "", nil
}

const (
	HelpOptionHeader  = "OPTIONS\n"
	HelpVersionHeader = "VERSION\n"
	HelpBackspace     = "  "
	HelpSourceSpliter = "/"
)

func buildHelp(o opts, fields []Field) string {
	var str strings.Builder

	str.WriteString(HelpVersionHeader)
	str.WriteString(HelpBackspace)
	str.WriteString(GetCurrentVersion().String())
	str.WriteString(HelpOptionHeader)

	for i := range fields {
		str.WriteString(HelpBackspace)

		for j := range o.sources {
			nextHasHelp := j+1 < len(o.sources) && o.sources[j+1].HasHelp()
			if !o.sources[j].HasHelp() {
				if nextHasHelp {
					str.WriteString(HelpSourceSpliter)
				}
				continue
			}

			str.WriteString(o.sources[j].Help(fields[i]))

			if nextHasHelp {
				str.WriteString(HelpSourceSpliter)
			}
		}

		str.WriteRune(' ')
		str.WriteString(fmt.Sprintf("<%s>", fields[i].Value.Type().String()))

		defaultV := fields[i].Options.DefaultValue
		if defaultV != "" {
			str.WriteRune(' ')
			str.WriteString(fmt.Sprintf("(default: %s)", defaultV))
		}

		str.WriteRune('\n')
	}

	return str.String()
}

// OutputTo outputs the v default configs to Outputer
func OutputTo(v any, os ...Outputter) error {
	if _, err := ParseConfig(v); err != nil {
		return err
	}

	for i := range os {
		err := os[i].Output(v)
		if err != nil {
			return err
		}
	}

	return nil
}
