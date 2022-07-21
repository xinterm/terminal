package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	"github.com/creack/pty"
	"github.com/nsf/termbox-go"
	"github.com/sirupsen/logrus"

	"github.com/xinterm/terminal"
)

const logFilename = "term-example.log"

func setupLogger(logger *logrus.Logger, f *os.File) {
	level := logrus.DebugLevel
	timeFormat := "2006-01-02 15:04:05.000 -0700 MST"

	logger.SetLevel(level)
	logger.SetFormatter(&simpleFormatter{
		TimeFormat: timeFormat,
	})

	logger.SetOutput(f)
}

func main() {
	logger := logrus.New()
	f, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Panicln(err)
	}
	defer f.Close()
	setupLogger(logger, f)

	logger.Infof("%s", strings.Repeat("=", 80))
	defer logger.Infof("%s", strings.Repeat("-", 80))

	err = termbox.Init()
	if err != nil {
		logger.Panicln("Termbox init failed:", err)
	}
	defer termbox.Close()

	c := exec.Command("bash")

	ptmx, err := pty.StartWithSize(c, &pty.Winsize{
		Rows: 80,
		Cols: 24,
	})
	if err != nil {
		logger.Panicln("Start PTY failed:", err)
	}
	defer ptmx.Close()

	t := terminal.New(80, 24)
	t.SetLogger(logger)
	t.SetMinUpdateInterval(10000000)

	t.Start()
	defer t.Stop()

	go func() {
		_, _ = io.Copy(t, ptmx)
	}()

	go func() {
		for {
			t.WaitUpdate()
			logger.Debugf("Should update now")
		}
	}()

	quit := false
	for !quit {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			//logger.Infof("Key ch: %c, 0x%x", event.Ch, event.Ch)
			//logger.Infof("Key: 0x%x", event.Key)
			if event.Key == termbox.KeyCtrlQ {
				quit = true
				break
			}
			if event.Ch != 0 {
				chUTF8 := make([]byte, utf8.UTFMax)
				size := utf8.EncodeRune(chUTF8, event.Ch)
				ptmx.Write(chUTF8[:size])
			} else if event.Key != 0 {
				ptmx.Write([]byte{byte(event.Key)})
			}
		case termbox.EventError:
			logger.Errorf("Termbox event error: %s", err)
		}
	}
}
