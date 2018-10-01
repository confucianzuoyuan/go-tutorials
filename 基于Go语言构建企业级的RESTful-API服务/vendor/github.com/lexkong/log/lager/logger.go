package lager

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

//StackTraceBufferSize is a constant which defines stack track buffer size
const StackTraceBufferSize = 1024 * 100

//Logger is a interface
type Logger interface {
	RegisterSink(Sink)
	Session(task string, data ...Data) Logger
	SessionName() string
	Debug(action string, data ...Data)
	Info(action string, data ...Data)
	Warn(action string, data ...Data)
	Error(action string, err error, data ...Data)
	Fatal(action string, err error, data ...Data)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(err error, format string, args ...interface{})
	Fatalf(err error, format string, args ...interface{})
	WithData(Data) Logger
}

type logger struct {
	component     string
	task          string
	sinks         []Sink
	sessionID     string
	nextSession   uint64
	data          Data
	logFormatText bool
}

//NewLoggerExt is a function which returns logger struct object
func NewLoggerExt(component string, isFormatText bool) Logger {
	return &logger{
		component:     component,
		task:          component,
		sinks:         []Sink{},
		data:          Data{},
		logFormatText: isFormatText,
	}
}

//NewLogger is a function used to get new logger object
func NewLogger(component string) Logger {
	return NewLoggerExt(component, true)
}

//RegisterSink is a function used to register sink
func (l *logger) RegisterSink(sink Sink) {
	l.sinks = append(l.sinks, sink)
}

//SessionName is used to get the session name
func (l *logger) SessionName() string {
	return l.task
}

//Session is a function which returns logger details for that session
func (l *logger) Session(task string, data ...Data) Logger {
	sid := atomic.AddUint64(&l.nextSession, 1)

	var sessionIDstr string

	if l.sessionID != "" {
		sessionIDstr = fmt.Sprintf("%s.%d", l.sessionID, sid)
	} else {
		sessionIDstr = fmt.Sprintf("%d", sid)
	}

	return &logger{
		component: l.component,
		task:      fmt.Sprintf("%s.%s", l.task, task),
		sinks:     l.sinks,
		sessionID: sessionIDstr,
		data:      l.baseData(data...),
	}
}

//WithData which adds data to the logger object
func (l *logger) WithData(data Data) Logger {
	return &logger{
		component: l.component,
		task:      l.task,
		sinks:     l.sinks,
		sessionID: l.sessionID,
		data:      l.baseData(data),
	}
}

// Find the sink need to log
func (l *logger) activeSinks(loglevel LogLevel) []Sink {
	ss := make([]Sink, len(l.sinks))
	idx := 0
	for _, itf := range l.sinks {
		if s, ok := itf.(*writerSink); ok && loglevel < s.minLogLevel {
			continue
		}
		if s, ok := itf.(*ReconfigurableSink); ok && loglevel < LogLevel(s.minLogLevel) {
			continue
		}
		ss[idx] = itf
		idx++
	}
	return ss[:idx]
}

func (l *logger) log(loglevel LogLevel, action string, err error, data ...Data) {
	ss := l.activeSinks(loglevel)
	if len(ss) == 0 {
		return
	}
	l.logs(ss, loglevel, action, err, data...)
}

func (l *logger) logs(ss []Sink, loglevel LogLevel, action string, err error, data ...Data) {
	logData := l.baseData(data...)

	if err != nil {
		logData["error"] = err.Error()
	}

	if loglevel == FATAL {
		stackTrace := make([]byte, StackTraceBufferSize)
		stackSize := runtime.Stack(stackTrace, false)
		stackTrace = stackTrace[:stackSize]

		logData["trace"] = string(stackTrace)
	}

	log := LogFormat{
		Timestamp: currentTimestamp(),
		Message:   action,
		LogLevel:  loglevel,
		Data:      logData,
	}

	// add file, lineno
	addExtLogInfo(&log)
	var logInfo string
	for _, sink := range l.sinks {
		if l.logFormatText {
			levelstr := FormatLogLevel(log.LogLevel)
			extraData, ok := log.Data["error"].(string)
			if ok && extraData != "" {
				extraData = " error: " + extraData
			}
			logInfo = log.Timestamp + " " + levelstr + " " + log.File + " " + log.Message + extraData
			sink.Log(loglevel, []byte(logInfo))

		} else {
			logInfo, jserr := log.ToJSON()
			if jserr != nil {
				fmt.Printf("[lager] ToJSON() ERROR! action: %s, jserr: %s, log: %+v", action, jserr, log)
				// also output json marshal error event to sink
				log.Data = Data{"Data": fmt.Sprint(logData)}
				jsonerrdata, _ := log.ToJSON()
				sink.Log(ERROR, jsonerrdata)
				continue
			}
			sink.Log(loglevel, logInfo)
		}
	}

	if loglevel == FATAL {
		panic(err)
	}
}

func (l *logger) Debug(action string, data ...Data) {
	l.log(DEBUG, action, nil, data...)
}

func (l *logger) Info(action string, data ...Data) {
	l.log(INFO, action, nil, data...)
}

func (l *logger) Warn(action string, data ...Data) {
	l.log(WARN, action, nil, data...)
}

func (l *logger) Error(action string, err error, data ...Data) {
	l.log(ERROR, action, err, data...)
}

func (l *logger) Fatal(action string, err error, data ...Data) {
	l.log(FATAL, action, err, data...)
}

func (l *logger) logf(loglevel LogLevel, err error, format string, args ...interface{}) {
	ss := l.activeSinks(loglevel)
	if len(ss) == 0 {
		return
	}
	logmsg := fmt.Sprintf(format, args...)
	l.logs(ss, loglevel, logmsg, err)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logf(DEBUG, nil, format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logf(INFO, nil, format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.logf(WARN, nil, format, args...)
}

func (l *logger) Errorf(err error, format string, args ...interface{}) {
	l.logf(ERROR, err, format, args...)
}

func (l *logger) Fatalf(err error, format string, args ...interface{}) {
	l.logf(FATAL, err, format, args...)
}

func (l *logger) baseData(givenData ...Data) Data {
	data := Data{}

	for k, v := range l.data {
		data[k] = v
	}

	if len(givenData) > 0 {
		for _, dataArg := range givenData {
			for key, val := range dataArg {
				data[key] = val
			}
		}
	}

	if l.sessionID != "" {
		data["session"] = l.sessionID
	}

	return data
}

func currentTimestamp() string {
	//return time.Now().Format("2006-01-02 15:04:05.000 -07:00")
	return time.Now().Format("2006-01-02 15:04:05.000")
}

func addExtLogInfo(logf *LogFormat) {

	for i := 4; i <= 5; i++ {
		_, file, line, ok := runtime.Caller(i)

		if strings.Index(file, "st_lager.go") > 0 {
			continue
		}

		if ok {
			idx := strings.LastIndex(file, "src")
			switch {
			case idx >= 0:
				logf.File = file[idx+4:]
			default:
				logf.File = file
			}
			// depth: 2
			indexFunc := func(file string) string {
				backup := "/" + file
				lastSlashIndex := strings.LastIndex(backup, "/")
				if lastSlashIndex < 0 {
					return backup
				}
				secondLastSlashIndex := strings.LastIndex(backup[:lastSlashIndex], "/")
				if secondLastSlashIndex < 0 {
					return backup[lastSlashIndex+1:]
				}
				return backup[secondLastSlashIndex+1:]
			}
			logf.File = indexFunc(logf.File) + ":" + strconv.Itoa(line)
		}
		break
	}
}
