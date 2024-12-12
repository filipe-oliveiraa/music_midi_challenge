package conf

import "errors"

// ErrHelp is the error returned if a source requests help
var ErrHelp = errors.New("help requested")
