package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log Статический экземпляр логгера
var Log = logrus.New()

func init() {
	Log.Out = os.Stdout
}
