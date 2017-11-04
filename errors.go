package gzog

import (
	e "errors"
)

var (
	MissingZogOptionsClose = e.New("missing ZogOptions closing brace")
	MissingZogOptionsOpen  = e.New("missing ZogOptions opening brace")
	NotAValidBoolean       = e.New("not a valid boolean option value")
	NotAnOption            = e.New("not a valid option name")
)
