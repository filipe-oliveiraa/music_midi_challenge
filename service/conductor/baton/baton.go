package baton

import (
	"io"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
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
	paused    atomic.Bool
	controlCh chan string
}

func New(log logging.Logger) Baton {
	b := &baton{
		log:       log,
		mu:        sync.Mutex{},
		musicians: make([]data.Musician, 0, 100),
		playing:   atomic.Bool{},
		paused:    atomic.Bool{},
		controlCh: make(chan string),
	}

	go b.handleSignals() // Start signal handler

	return b
}

func (b *baton) handleSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGUSR1, syscall.SIGUSR2) // Register signals to be caught by the process

	for {
		sig := <-sigCh // Wait for signal to be caught
		b.log.Infof("Received signal: %v", sig)
		switch sig {
		case syscall.SIGUSR1:
			b.Pause()
		case syscall.SIGUSR2:
			b.Resume()
		}
	}
}

// Pauses the music
func (b *baton) Pause() {
	if b.playing.Load() { //use atomic operation to check if music is playing
		b.log.Info("Pausing music")
		b.paused.Store(true)
	}
}

// Resumes the music
func (b *baton) Resume() {
	if b.playing.Load() && b.paused.Load() {
		b.log.Info("Resuming music")
		b.paused.Store(false)
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
	channelMap := make(map[int]chan Note) // map of channels for each musician
	var wg sync.WaitGroup

	// Initialize channels for each musician
	for i := range b.musicians {
		channelMap[i] = make(chan Note)
		wg.Add(1)
		go b.handleMusician(i, channelMap[i], &wg)
	}

	// Initialize channels for each track
	trackouts := make(map[int]drivers.Out)
	for i := range tracks.SMF().Tracks {
		trackouts[i] = &Track{
			ch:    channelMap[i],
			index: i,
			//wg:    &wg,
		}
	}

	// Send notes to the musician channels
	tracks.Do(func(ev smf.TrackEvent) {
		b.log.Infof("track %v @%vms %s\n", ev.TrackNo, ev.AbsMicroSeconds/1000, ev.Message)
	}).MultiPlay(trackouts)

	// Close all channels
	for _, ch := range channelMap {
		close(ch)
	}

	// Wait for all musicians to finish
	wg.Wait()
}

func (b *baton) handleMusician(index int, ch chan Note, wg *sync.WaitGroup) {
	defer wg.Done()
	for note := range ch {
		// Check if paused before sending the note
		for b.paused.Load() {
			b.log.Debug("Music is paused")
			time.Sleep(100 * time.Millisecond) // Prevent consuming CPU
		}

		cli, err := client.New(b.musicians[index].Address) // Create client for musician

		if err != nil {
			b.log.With("error", err).Error("creating musician client")
			continue
		}

		b.log.With("note", note).Info("sending note")
		if len(note.note) > 0 {
			err = cli.Play(note.note) // Send note to musician
		}

		if err != nil {
			b.log.With("error", err).Error("playing note")
		}
	}
}
