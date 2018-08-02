package logger

import "testing"
import "encoding/json"
import "github.com/sirupsen/logrus"

func TestNew(t *testing.T) {
	config := Logger{}

	err := json.Unmarshal([]byte(DefaultConfig), &config)
	if err != nil {
		t.Fatal("配置信息解析错误", err)
	}
	err = New(config)
	if err != nil {
		t.Fatal("配置信息错误", err)
	}

	logrus.Debug(1)
	logrus.Info(2)
	logrus.Warn(3)
	logrus.Error(4)
}
