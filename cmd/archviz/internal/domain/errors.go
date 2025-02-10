package domain

import "github.com/go-faster/errors"

var (
	ErrComponentNotFound   = errors.New("component not found")
	ErrComponentUnresolved = errors.New("component has not been resolved yet")
)
