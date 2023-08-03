package publicid

import (
	nanoid "github.com/matoous/go-nanoid/v2"
)

// Fixed nanoid parameters used in the Rails application.
const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	length   = 12
)

// New generates a unique public ID.
// func New() (string, error) { return nanoid.Generate(alphabet, length) }

// Must is the same as New, but panics on error.
func Must() string { return nanoid.MustGenerate(alphabet, length) }
