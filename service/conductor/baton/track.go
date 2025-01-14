package baton

import "sync"

type Track struct {
	ch    chan Note
	index int
	wg    *sync.WaitGroup
}

type Note struct {
	index int
	note  []byte
}

func (t *Track) Close() error {
	// Not needed for our implementation
	return nil
}

func (t *Track) IsOpen() bool {
	return true
}

func (t *Track) Number() int {
	// Not needed for our implementation
	return 0
}

func (t *Track) Open() error {
	// Not needed for our implementation
	return nil
}

func (t *Track) Send(bs []byte) error {
	t.wg.Add(1)
	t.ch <- Note{
		index: t.index,
		note:  bs,
	}
	return nil
}

func (t *Track) String() string {
	// Not needed for our implementation
	return "" //t.out.String()
}

func (t *Track) Underlying() interface{} {
	// Not needed for our implementation
	return nil
}
