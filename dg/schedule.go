package dg

import (
	"github.com/robfig/cron"
	"time"
	"per.zdh.org/influxdb-sentinel/global"
	"per.zdh.org/influxdb-sentinel/dg/daemonScript"
)

type Engine interface {
	Start(t int64)
}

type Schedule struct {
	Time               global.DGTime
	DaemonScriptService *daemonScript.Service
}

func (schedule *Schedule) Starts() {
	//10秒
	timeTenSeconds := []Engine{
		schedule.DaemonScriptService,
	}
	schedule.start(schedule.Time.TenSecond, timeTenSeconds)
}

func (schedule *Schedule) start(t string, jobs []Engine) {
	c := cron.New()
	c.AddFunc(t, func() {
		nowTime := time.Now().Unix()
		for _, t := range jobs {
			go t.Start(nowTime)
		}
	})
	c.Start()
}
