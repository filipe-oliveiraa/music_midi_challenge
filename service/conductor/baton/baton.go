package baton

import (
	"io"
	"sync"
	"sync/atomic"

	"crossjoin.com/gorxestra/daemon/musiciand/api/client/v1"
	"crossjoin.com/gorxestra/data"
	"crossjoin.com/gorxestra/logging"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/smf"
)

type baton struct {
	log       logging.Logger
	mu        sync.Mutex
	musicians []data.Musician
	playing   atomic.Bool
}

func New(log logging.Logger) Baton {
	return &baton{
		log:       log,
		mu:        sync.Mutex{},
		musicians: make([]data.Musician, 0, 100),
		playing:   atomic.Bool{},
	}
}

func (b *baton) RegisterMusician(m data.Musician) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.musicians = append(b.musicians, m)
	return nil
}

func (b *baton) UnregisterMusician(id data.ID) error {
	//b.mu.Lock()
	//defer b.mu.Unlock()
	//TODO
	return nil
}

func (b *baton) Play(r io.Reader) error {
	if b.playing.CompareAndSwap(false, true) {
		go b.play(r)
		return nil
	}

	return data.MusicAlreadyBeingPlayed
}

func (b *baton) play(r io.Reader) {
	// read and play it
	tracks := smf.ReadTracksFrom(r)
	// numTracks := len(tracks.SMF().Tracks)

	channel := make(chan Note)
	var wg sync.WaitGroup

	trackouts := make(map[int]drivers.Out)
	for i := range tracks.SMF().Tracks {
		trackouts[i] = &Track{
			ch:    channel,
			index: i,
			wg:    &wg,
		}
	}

	go func() {
		for note := range channel {
			if note.index >= len(b.musicians) {
				b.log.
					With("index", note.index).
					With("musicians", len(b.musicians)).
					Warn("skipping track")
				wg.Done()
				continue
			}

			go func(note Note) {
				defer wg.Done()
				cli, err := client.New(b.musicians[note.index].Address)
				if err != nil {
					b.log.
						With("error", err).
						Error("creating musician client")
					return
				}

				b.log.With("note", note).Info("sending note")
				if len(note.note) > 0 {
					err = cli.Play(note.note)
				}

				if err != nil {
					b.log.
						With("error", err).
						Error("playing note")
				}
			}(note)
		}
	}()

	// play music
	tracks.Do(func(ev smf.TrackEvent) {
		b.log.Infof("track %v @%vms %s\n", ev.TrackNo, ev.AbsMicroSeconds/1000, ev.Message)
	}).MultiPlay(trackouts)

	// close music channels
	close(channel)

	// wait for all notes to be played
	wg.Wait()

	b.playing.Store(false)
}
