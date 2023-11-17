package logs

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type IOption interface {
	TakeStdout() bool
	TakeHooks() []logrus.Hook
	TakeLevel() string
	TakeFormatter() logrus.Formatter
	TakeReportCaller() bool
}

func NewBuilder(opt IOption) *Builder {
	return &Builder{opt: opt}
}

type Builder struct {
	opt IOption
}

func (b Builder) makeHooks() (hooks logrus.LevelHooks) {
	hooks = make(logrus.LevelHooks)
	for _, hook := range b.opt.TakeHooks() {
		hooks.Add(hook)
	}

	return
}

func (b Builder) makeLevel() (level logrus.Level) {
	level, _ = logrus.ParseLevel(b.opt.TakeLevel())
	return
}

func (b Builder) makeStdout() io.Writer {
	if b.opt.TakeStdout() {
		return os.Stderr
	}

	file, _ := os.OpenFile(os.DevNull, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	return file
}

func (b Builder) Make() (entity *logrus.Logger) {
	entity = &logrus.Logger{
		Out:          b.makeStdout(),
		Hooks:        b.makeHooks(),
		Level:        b.makeLevel(),
		ExitFunc:     os.Exit,
		ReportCaller: b.opt.TakeReportCaller(),
	}

	return
}
