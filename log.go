package logs

import (
	"fmt"
	rotate "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	entity  *logrus.Logger
	options *Options
	//once      sync.Once
	//singleton *Entity
)

type Options struct {
	Path            string `json:"path" yaml:"path"`
	Level           string `json:"level" yaml:"level"`
	Stdout          bool   `json:"stdout" yaml:"stdout"`
	SaveDay         uint   `json:"save_day" yaml:"save_day"`
	TimestampFormat string `json:"timestamp_format" yaml:"timestamp_format"`
}

type Entity struct {
	opt *Options
}

func NewEntity(opt *Options) *Entity {
	return &Entity{opt: opt}
}

func (e Entity) Initialize() {
	entity = logrus.New()

	if err := e.display(); err != nil {
		panic(fmt.Sprintf("log initialize failed: %v", err))
	}

	level, _ := logrus.ParseLevel(e.opt.Level)
	entity.SetLevel(level)

	formatter := &Formatter{e.opt.TimestampFormat}
	levels := []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel}
	lfHook := lfshook.NewHook(e.outputMap(levels...), formatter)
	entity.AddHook(lfHook)
}

// display 关闭控制输出
func (e Entity) display() (err error) {
	if e.opt.Stdout {
		return
	}

	file, err := os.OpenFile(os.DevNull, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	entity.SetOutput(file)
	return
}

func (e Entity) getWriter(path, name string, rotationCount uint) *rotate.RotateLogs {
	now := time.Now().Format("200601/02")
	writer, err := rotate.New(
		path+"/%Y%m/%d/"+name+".%H.log",
		rotate.WithLinkName(path+"/"+now+"/"+name+".log"),
		rotate.WithMaxAge(time.Duration(rotationCount)*24*time.Hour),
		rotate.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(fmt.Sprintf("InitLog writer error:%s, logPath:%s",
			err.Error(), path))
	}

	return writer
}

// outputMap 切割字典
func (e Entity) outputMap(levels ...logrus.Level) lfshook.WriterMap {
	writers := make(lfshook.WriterMap)
	var b []byte
	for _, level := range levels {
		b, _ = level.MarshalText()
		writers[level] = e.getWriter(
			e.opt.Path, string(b), e.opt.SaveDay,
		)
	}

	return writers
}
