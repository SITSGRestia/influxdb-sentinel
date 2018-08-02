package global

import (
	"per.zdh.org/influxdb-sentinel/logger"
	"per.zdh.org/influxdb-sentinel/db/influxdb"
	"per.zdh.org/influxdb-sentinel/mq/emqtt"
	"per.zdh.org/influxdb-sentinel/mq/rabbitmq"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
)

type App struct {
	Version       string                 `json:"version" description:"api版本"`
	Server        Server                 `json:"server" description:"接口服务地址"`
	MQType        string                 `json:"mqtype" description:"选用的消息队列类型"`
	ClusterNodes  []influxdb.Nodes       `json:"clusterNodes" description:"集群节点访问信息"`
	SentinelState influxdb.SentinelState `json:"sentinelState" description:"集群节点数据状态"`
	Log           logger.Logger          `json:"log" description:"日志配置"`
	DB            DB                     `json:"db" description:"数据库配置"`
	MQ            MQ                     `json:"mq" description:消息队列配置"`
	DG            DG                     `json:"dg" description:数据采集"`
}

type Server struct {
	Host string `json:"host" description:"接口访问地址"`
	Port int    `json:"port" description:"接口访问端口"`
}

type DB struct {
	InfluxDB influxdb.InfluxDB `json:"influxDB" description:"InfluxDB数据库"`
}

type MQ struct {
	Emqtt    emqtt.Emqtt       `json:"emqtt" description:"emqtt配置"`
	RabbitMQ rabbitmq.RabbitMQ `json:"rabbitMQ" description:"emqtt配置"`
}

type DG struct {
	DGTime DGTime `json:"time" description:数据调度时间"`
}

type DGTime struct {
	FiveSecond string `json:"fiveSecond" description:5s"`
	TenSecond  string `json:"tenSecond" description:10s"`
	OneMin     string `json:"oneMin" description:1分钟"`
	ThreeMin   string `json:"threeMin" description:3分钟"`
	FiveMin    string `json:"fiveMin" description:5分钟"`
	TenMin     string `json:"tenMin" description:10分钟"`
	Day        string `json:"day" description:每天"`
	ZeroHour   string `json:"zeroHour" description:每天零点"`
}

func New(filePath string) (options *App) {
	if strings.TrimSpace(filePath) == "" {
		panic("未找到配置文件")
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("配置文件读取错误.%+v", err))
	}
	if data == nil || len(data) <= 0 {
		panic(fmt.Sprintf("配置文件错误.数据不存在"))
	}

	options = &App{}
	err = json.Unmarshal(data, options)
	if err != nil {
		panic(fmt.Sprintf("配置文件解析错误.%+v", err))
	}
	return options
}
