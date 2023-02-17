package logs

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/goantor/ex"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type pkgLogger interface {
	Info(message string, data interface{})
	Trace(message string, data interface{})
	Debug(message string, data interface{})
	Warn(message string, data interface{})
	Fatal(message string, data interface{})
	Panic(message string, data interface{})
	Error(message string, data interface{})
}

type Logger interface {
	pkgLogger

	// GenerateId 创建追踪id
	GenerateId()

	// GetId 获取追踪ID
	GetId() string

	BindId(id string)

	// User 注入数据
	User(user interface{})

	// Params 参数
	Params(params interface{})

	// Auto 自动识别log等级
	Auto(no ex.IErrno, message string, data interface{})
}

type logger struct {
	id     string
	method string
	action string
	ip     string
	user   interface{}
	params interface{}
}

func New(method, action, ip string) Logger {
	l := &logger{
		method: method,
		action: action,
		ip:     ip,
	}

	l.GenerateId()
	return l
}

func (l *logger) GenerateId() {
	buf := make([]byte, 32)
	u := uuid.NewV4().Bytes()
	hex.Encode(buf, u)
	l.id = string(buf)
}

func (l *logger) GetId() string {
	return l.id
}

func (l *logger) BindId(id string) {
	l.id = id
}

func (l *logger) User(user interface{}) {
	l.user = user
}

func (l *logger) Params(params interface{}) {
	l.params = params
}

func (l *logger) format(message string, data interface{}) string {
	jsParam := map[string]interface{}{
		"user":   l.user,
		"params": l.params,
		"data":   data,
	}

	js, _ := json.Marshal(jsParam)
	return fmt.Sprintf("%s Trace: %s %s IP: %s Message: %s Data: %s", l.method, l.id, l.action, l.ip, message, js)
}

func (l *logger) Auto(no ex.IErrno, message string, data interface{}) {
	level, _ := logrus.ParseLevel(no.Level())
	go l.log(level, l.format(message, data))
}

func (l *logger) log(level logrus.Level, message string) {
	entity.Log(level, message)
}

func (l *logger) Info(message string, data interface{}) {
	l.log(logrus.InfoLevel, l.format(message, data))
}

func (l *logger) Trace(message string, data interface{}) {
	l.log(logrus.TraceLevel, l.format(message, data))
}

func (l *logger) Debug(message string, data interface{}) {
	l.log(logrus.DebugLevel, l.format(message, data))
}

func (l *logger) Warn(message string, data interface{}) {
	l.log(logrus.WarnLevel, l.format(message, data))
}

func (l *logger) Fatal(message string, data interface{}) {
	l.log(logrus.FatalLevel, l.format(message, data))
}

func (l *logger) Panic(message string, data interface{}) {
	l.log(logrus.PanicLevel, l.format(message, data))
}

func (l *logger) Error(message string, data interface{}) {
	l.log(logrus.ErrorLevel, l.format(message, data))
}
