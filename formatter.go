package logs

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type Formatter struct {
	timestampFormat string
}

func (f *Formatter) buf(entry *logrus.Entry) *bytes.Buffer {
	if entry.Buffer == nil {
		return &bytes.Buffer{}
	}

	return entry.Buffer
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	buf := f.buf(entry)
	var message string

	//if entry.Level < 4 {
	//	message = fmt.Sprintf("[%v] %v stack:%s", entry.Time.Format(f.timestampFormat), entry.Message, makeCallStack())
	//} else {
	message = fmt.Sprintf("[%v] %v", entry.Time.Format(f.timestampFormat), entry.Message)
	//}

	buf.WriteString(message)
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

func makeCallStack() string {

	pc := make([]uintptr, 10) // 可根据实际情况调整大小
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	callStack := make([]string, 0)
	for {
		frame, more := frames.Next()
		callStack = append(callStack, fmt.Sprintf("%s:%d", frame.File, frame.Line))
		if !more {
			break
		}
	}

	return strings.Join(callStack, "\n")
}

func shortenFileName(file string) string {
	parts := strings.Split(file, "/")
	if len(parts) > 1 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	return file
}
