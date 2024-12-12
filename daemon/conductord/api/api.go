package api

import (
	"crossjoin.com/gorxestra/data"
)

type NodeInterface interface {
	RegisterMusician(m data.Musician) error
	UnregisterMusician(id data.ID) error
	PlayMusic(name string) error
}
