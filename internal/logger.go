package internal

import (
	"strconv"

	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
)

func logDebugErrors(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func LogRequestResponse(code int, remote, method, path, err string) {
	cRemote := ansi.Color(remote, "yellow")
	cMethod := ansi.Color(method, "magenta")
	cPath := ansi.Color(path, "cyan")

	switch code {
	case 200:
		cCode := ansi.Color(strconv.Itoa(code), "green")
		logrus.Infof("%s | %s | %s %s", cCode, cRemote, cMethod, cPath)
	default:
		cCode := ansi.Color(strconv.Itoa(code), "red")
		cErr := ansi.Color(err, "red")
		logrus.Errorf("%s | %s | %s %s | %s", cCode, cRemote, cMethod, cPath, cErr)
	}
}
