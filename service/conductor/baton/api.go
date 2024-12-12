package baton

import (
	"io"

	"crossjoin.com/gorxestra/data"
)

type Baton interface {
	RegisterMusician(m data.Musician) error
	UnregisterMusician(id data.ID) error
	Play(r io.Reader) error
}
