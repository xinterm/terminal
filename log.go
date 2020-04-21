package terminal

import (
	"github.com/xinterm/terminal/util"
)

type internalLog struct {
	log util.Logger
}

func (l *internalLog) Debugf(f string, v ...interface{}) {
	if l.log != nil {
		l.log.Debugf(f, v...)
	}
}

func (l *internalLog) Infof(f string, v ...interface{}) {
	if l.log != nil {
		l.log.Infof(f, v...)
	}
}

func (l *internalLog) Warnf(f string, v ...interface{}) {
	if l.log != nil {
		l.log.Warnf(f, v...)
	}
}

func (l *internalLog) Errorf(f string, v ...interface{}) {
	if l.log != nil {
		l.log.Errorf(f, v...)
	}
}
