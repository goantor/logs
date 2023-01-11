package logs

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
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

	buf.WriteString(
		fmt.Sprintf("[%v] %v",
			entry.Time.Format(f.timestampFormat), entry.Message,
		),
	)

	buf.WriteByte('\n')
	return buf.Bytes(), nil
}
