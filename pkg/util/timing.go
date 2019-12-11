package util

import (
	"fmt"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

// Logger is a simple logging interface
type Logger interface {
	Printf(format string, args ...interface{})
}

// FuncTimer is a global FunctionTimer, non-nil if initialized
var FuncTimer *FunctionTimer

func init() {
	FuncTimer = NewFunctionTimer()
}

type logger struct{}

func (l *logger) Printf(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// FunctionTimer collects infomation about the elapsed time of a function call
type FunctionTimer struct {
	timings chan *functionTiming
	done    chan int
	logger  Logger
}

type functionTiming struct {
	file     string
	lineNo   int
	funcName string
	duration time.Duration
}

func (t *functionTiming) String() string {
	return fmt.Sprintf("%s %s:%d %s", fmtDuration(t.duration), t.file, t.lineNo, t.funcName)
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond
	d -= ms * time.Millisecond
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

// NewFunctionTimer creates a FunctionTimer
func NewFunctionTimer() *FunctionTimer {
	return &FunctionTimer{
		timings: make(chan *functionTiming),
		done:    make(chan int),
		logger:  &logger{},
	}
}

// Receive collects the timing information collected in the program, should be run in a go routine
func (ft *FunctionTimer) Receive() {
	for {
		select {
		case timing := <-ft.timings:
			ft.logger.Printf("TIMING %s", timing.String())
		case <-ft.done:
			return
		}
	}
}

// Terminate stops the Receive method
func (ft *FunctionTimer) Terminate() {
	ft.done <- 0
}

// Instrument returns a stop function, insert code at the beginning of a function to collect timing data
// if util.FuncTimer != nil {
// 	timer := util.FuncTimer.Instrument()
// 	defer timer()
// }
func (ft *FunctionTimer) Instrument() func() {
	funcName := "none"
	file := ""
	lineNo := 0

	fpcs := make([]uintptr, 1)
	n := runtime.Callers(2, fpcs) // start with caller of Instrument()
	if n > 0 {
		frames := runtime.CallersFrames(fpcs[0:1]) // just one frame
		frame, _ := frames.Next()
		// fmt.Println(more) // it's false
		file = frame.File
		funcName = frame.Func.Name()
		lineNo = frame.Line
	}

	start := time.Now()
	return func() {
		stop := time.Now()
		ft.timings <- &functionTiming{file, lineNo, funcName, stop.Sub(start)}
	}
}
