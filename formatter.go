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
	callStack := make([]string, 0)
	for skip := 1; ; skip++ { // Start from skip = 1 to skip the current frame.
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		if -1 != strings.Index(file, "pkg/mod") {
			continue
		}

		callStack = append(callStack, fmt.Sprintf("%s:%d", file, line))
		if pc == 0 {
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
