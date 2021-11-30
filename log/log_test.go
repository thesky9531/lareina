package log

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLog(t *testing.T) {
	err := Init(&Config{
		Output:  "test.log",
		MaxSize: 1,
		Level:   "debug",
	})
	if err != nil {
		t.Error(err)
	}
	for {
		logrus.Error("test")
	}
}
