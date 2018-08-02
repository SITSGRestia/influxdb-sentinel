package daemonScript

import (
	"github.com/sirupsen/logrus"
	"os/exec"
	"runtime"
	"bytes"
	"strings"
	"fmt"
	"time"
)

func (service *Service) Start(unixTime int64) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("influxd守护进程panic错误:%s", err)
		}
	}()
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("heart.bat")
	} else {
		cmd = exec.Command("/bin/sh", "-c", `ps -C influxd | wc -l`)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		logrus.Errorf("Error: failed to execute this command: ", fmt.Sprint(err)+": "+stderr.String())
		return
	}

	row := strings.TrimSpace(string(out))
	count := "1"
	if runtime.GOOS == "windows" {
		count = "0"
	}

	if row == count {
		timeDown := time.Now().Unix()
		for {
			if runtime.GOOS == "windows" {
				cmd = exec.Command("influxd", "-config", service.path)
			} else {
				cmd = exec.Command("/bin/sh", "-c", `systemctl start influxd`)
			}
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err := cmd.Start()
			if err != nil {
				logrus.Errorf("Error: failed to start influxd: ", fmt.Sprint(err)+": "+stderr.String())
				time.Sleep(time.Second * 2)
				continue
			}
			break
		}
		service.syncChannel <- timeDown
	}
}
