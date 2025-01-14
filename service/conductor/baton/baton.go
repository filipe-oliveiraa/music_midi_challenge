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
	defer b.playing.Store(false)

	tracks := smf.ReadTracksFrom(r)
	channelMap := make(map[int]chan Note)
	var wg sync.WaitGroup

	// Initialize a channel for each musician
	for i := range b.musicians {
		channelMap[i] = make(chan Note)
		wg.Add(1)
		go b.handleMusician(i, channelMap[i], &wg)
	}
	// Initialize a track for each channel(each musician)
	trackouts := make(map[int]drivers.Out)
	for i := range tracks.SMF().Tracks {
		trackouts[i] = &Track{
			ch:    channelMap[i],
			index: i,
			wg:    &wg,
		}
	}
	// Send notes to channels concurrently
	tracks.Do(func(ev smf.TrackEvent) {
		b.log.Infof("track %v @%vms %s\n", ev.TrackNo, ev.AbsMicroSeconds/1000, ev.Message)
	}).MultiPlay(trackouts)

	// Close channels
	for _, ch := range channelMap {
		close(ch)
	}
	// Wait for all musicians to finish playing
	wg.Wait()
}

func (b *baton) handleMusician(index int, ch chan Note, wg *sync.WaitGroup) {
	defer wg.Done()
	// Read notes from channel
	for note := range ch {
		cli, err := client.New(b.musicians[index].Address)
		if err != nil {
			b.log.With("error", err).Error("creating musician client")
			continue
		}
		// Send note to musician
		b.log.With("note", note).Info("sending note")
		if len(note.note) > 0 {
			err = cli.Play(note.note)
		}

		if err != nil {
			b.log.With("error", err).Error("playing note")
		}
	}
}
