package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

const modulePath = "go.pantheon.tech/stonework"

var (
	logstyleFilename = color.Style{color.LightWhite}
	logstyleFilenum  = color.Style{color.White}
	logstyleFuncname = color.Style{color.LightBlue}
	logstyleFuncpkg  = color.Style{color.Blue}
)

func init() {
	formatter := &logFormatter{&logrus.TextFormatter{
		EnvironmentOverrideColors: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			call := strings.TrimPrefix(frame.Function, modulePath)
			parts := strings.SplitN(strings.TrimPrefix(call, "/"), ".", 2)
			function = fmt.Sprintf("%s.%s()", logstyleFuncpkg.Sprint(parts[0]), logstyleFuncname.Sprint(parts[1]))
			_, file = filepath.Split(frame.File)
			file = fmt.Sprintf("%s:%s", logstyleFilename.Sprint(file), logstyleFilenum.Sprint(frame.Line))
			return function, file
		},
	}}
	logrus.SetFormatter(formatter)
}

type logFormatter struct {
	*logrus.TextFormatter
}

const (
	traceLvlPrefix = "\x1b[37mTRAC"
	debugLvlPrefix = "\x1b[37mDEBU"
)

func (l *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b, err := l.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("LOG: %q\n", string(b))
	if bytes.HasPrefix(b, []byte(traceLvlPrefix)) {
		b[2] = '9'
		b[3] = '0'
	} else if bytes.HasPrefix(b, []byte(debugLvlPrefix)) {
		b[2] = '3'
		b[3] = '6'
	}
	return b, nil
}
