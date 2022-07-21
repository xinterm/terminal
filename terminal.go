package terminal

import (
	"sync"
	"time"

	"github.com/xinterm/terminal/internal/sequence"
	"github.com/xinterm/terminal/util"
)

// Terminal describes a full featured virtual terminal
type Terminal struct {
	st     *state
	seqCtl *sequence.Control

	minUpdateInterval time.Duration
	updateEvent       chan struct{}

	closeTerminal chan struct{}

	wg sync.WaitGroup

	log internalLog
}

// New creates a new terminal
func New(rows, cols int) *Terminal {
	t := &Terminal{
		updateEvent:   make(chan struct{}),
		closeTerminal: make(chan struct{}),
	}

	t.st = newState(&t.log)
	t.st.resize(rows, cols)

	t.seqCtl = sequence.NewControl(&t.log)

	return t
}

// SetLogger accepts a Logger interface
func (t *Terminal) SetLogger(log util.Logger) {
	t.log.log = log
}

// SetHistoryLimit sets the scroll back size
func (t *Terminal) SetHistoryLimit(size int) {
	t.st.setScrollBackSize(size)
}

// SetMinUpdateInterval sets mininum update interval
func (t *Terminal) SetMinUpdateInterval(interval time.Duration) {
	if interval < 0 {
		t.minUpdateInterval = 0
		return
	}
	t.minUpdateInterval = interval
}

func (t *Terminal) sendUpdateEvent(sendOver chan struct{}) {
	defer t.wg.Done()

	lastUpdate := time.Now().Add(-t.minUpdateInterval - 1000000)

	for {
		select {
		case <-sendOver:
		case <-t.closeTerminal:
			t.log.Debugf("Receive close terminal signal in sendUpdateEvent send over loop")
			return
		}

		interval := time.Since(lastUpdate)
		if interval < t.minUpdateInterval {
			timer := time.NewTimer(t.minUpdateInterval - interval)
		timerLoop:
			for {
				select {
				case <-timer.C:
					break timerLoop
				case <-sendOver:
				case <-t.closeTerminal:
					if !timer.Stop() {
						<-timer.C
					}
					t.log.Debugf("Receive close terminal signal in sendUpdateEvent timer loop")
					return
				}
			}
		}

	updateLoop:
		for {
			select {
			case t.updateEvent <- struct{}{}:
				lastUpdate = time.Now()
				break updateLoop
			case <-sendOver:
			case <-t.closeTerminal:
				t.log.Debugf("Receive close terminal signal in sendUpdateEvent send update loop")
				return
			}
		}
	}
}

func (t *Terminal) pollAndSendResult(sendOver chan struct{}) {
	defer t.wg.Done()

	for {
		select {
		case r := <-t.seqCtl.ResultEvent():
			t.st.sendResult(r)
		case <-t.closeTerminal:
			t.log.Debugf("Receive close terminal signal in pollAndSendResult when receive result")
			return
		}

		select {
		case sendOver <- struct{}{}:
		case <-t.closeTerminal:
			t.log.Debugf("Receive close terminal signal in pollAndSendResult when send over")
			return
		}
	}
}

// Start the terminal loop
func (t *Terminal) Start() {
	sendOver := make(chan struct{})

	t.wg.Add(2)

	go t.sendUpdateEvent(sendOver)

	go t.pollAndSendResult(sendOver)
}

// Stop the terminal
func (t *Terminal) Stop() {
	close(t.closeTerminal)

	t.log.Debugf("Wait for exiting terminal...")

	t.wg.Wait()
}

// Write implements io.Writer
func (t *Terminal) Write(p []byte) (n int, err error) {
	t.log.Debugf("Write: 0x%x", p)
	for _, c := range p {
		t.seqCtl.Parse(c)
	}
	return len(p), nil
}

// Resize the terminal
func (t *Terminal) Resize(rows, cols int) {
	t.st.resize(rows, cols)
}

// WaitUpdate waits for the update event
func (t *Terminal) WaitUpdate() {
	<-t.updateEvent
}

// GridVisible gets the visible grid
func (t *Terminal) GridVisible(flip bool) *Grid {
	return t.st.gridVisible(flip)
}

// GridFromAbsStart gets the required grid
func (t *Terminal) GridFromAbsStart(startLine LineAbsNo, startLineSubNo, displayLines int, flip bool) *Grid {
	return nil
}

// GridToAbsEnd gets the required grid
func (t *Terminal) GridToAbsEnd(endLine LineAbsNo, endLineSubNo, displayLines int, flip bool) *Grid {
	return nil
}

// GridFromStart gets the required grid
func (t *Terminal) GridFromStart(startLine, displayLines int, flip bool) *Grid {
	return nil
}

// GridToEnd gets the required grid
func (t *Terminal) GridToEnd(endLine, displayLines int, flip bool) *Grid {
	return t.st.gridToEnd(endLine, displayLines, flip)
}
