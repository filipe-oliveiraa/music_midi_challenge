package conf

import "os"

type opts struct {
	sources []Source
}

type Option func(o *opts)

func applyOptions(o []Option) opts {
	defaultOpts := opts{
		sources: []Source{
			NewFlagSource("", os.Args[1:]),
			NewJsonSource("./conf.json", true),
			NewEnvSource("ENV"),
		},
	}

	for i := range o {
		o[i](&defaultOpts)
	}

	return defaultOpts
}

func WithSources(sources ...Source) Option {
	return func(o *opts) {
		o.sources = sources
	}
}
