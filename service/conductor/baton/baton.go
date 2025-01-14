package baton

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

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

	tracks := smf.ReadTracksFrom(r) //read trackes from the io.Reader
	numTracks := len(tracks.SMF().Tracks)

	b.log.With("numTracks", numTracks).Info("Tracks parsed successfully")

	channel := make(chan Note, 1000) //Buffered channel to handle notes

	trackouts := make(map[int]drivers.Out) //Store the output of channels of each track
	for i := 0; i < numTracks; i++ {
		trackouts[i] = &Track{
			ch:    channel,
			index: i,
		}
		b.log.With("trackIndex", i).Info("Initialized track output")
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for note := range channel {
			if note.index >= len(b.musicians) {
				b.log.
					With("index", note.index).
					With("musicians", len(b.musicians)).
					Warn("skipping track")
				continue
			}

			cli, err := client.New(b.musicians[note.index].Address)
			if err != nil {
				b.log.
					With("error", err).
					Error("creating musician client")
				continue
			}

			b.log.With("note", fmt.Sprintf("%v", note.note)).
				With("index", note.index).
				With("musician_address", b.musicians[note.index].Address).
				With("timestamp_sent", time.Now()). // Log the time when the note is sent
				Info("Sending note to musician")

			//b.log.With("note", note).Info("sending note")
			if len(note.note) > 0 {
				err = cli.Play(note.note)
			}

			if err != nil {
				b.log.
					With("error", err).
					Error("playing note")
			}
		}
	}()

	tracks.Do(func(ev smf.TrackEvent) {
		//b.log.Infof("track %v @%vms %s\n", ev.TrackNo, ev.AbsMicroSeconds/1000, ev.Message)
		fmt.Printf("Track: %d | Time: %d | Message: %s\n", ev.TrackNo, ev.AbsMicroSeconds/1000, ev.Message.String())
	}).MultiPlay(trackouts)

	close(channel)
	wg.Wait()
}
