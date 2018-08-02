package influxdb

import (
	"github.com/influxdata/influxdb/client/v2"
	"fmt"
	"time"
)

type InfluxDB struct {
	Host     string `json:"host" description:"数据库地址"`
	Port     int    `json:"port" description:"数据库端口"`
	Username string `json:"username" description:"数据库访问用户名"`
	Password string `json:"password" description:"数据库访问密码"`
	Path     string `json:"path" description:"influxd执行文件绝对路径"`
}

type StoreData struct {
	Precision   string              `json:"precision"`
	PolicyName  string              `json:"policyName"`
	Measurement string              `json:"measurement"`
	Database    string              `json:"database"`
	Tags        []map[string]string `json:"tags"`
	Fields      []map[string]interface{} `json:"fields"`
	Timestamps  []time.Time         `json:"timestamps"`
}

type Nodes struct {
	Host     string `json:"host" description:"集群节点访问地址"`
	Port     int    `json:"port" description:"集群节点访问端口"`
	Username string `json:"username" description:"集群节点数据库访问用户名"`
	Password string `json:"password" description:"集群节点数据库访问密码"`
}

type SentinelState struct {
	DataInfo []DataState `json:"dataInfo" description:"数据信息"`
}

type DataState struct {
	Database        string        `json:"database" description:"数据库名"`
	MeasurementInfo []Measurement `json:"measurementInfo" description:"数据库表名"`
	RetentionPolicy []string      `json:"retentionPolicy" description:"数据库保留策略名"`
}

type Measurement struct {
	Name       string `json:"name" description:"数据库表名"`
	PrimaryKey string `json:"primaryKey" description:"数据库主键名"`
}

func New(option InfluxDB) *client.Client {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%d", option.Host, option.Port),
		Username: option.Username,
		Password: option.Password,
	})
	if err != nil {
		panic(fmt.Sprintf("InfluxDB数据库连接失败.%+v", err))
	}
	return &c
}

func NewSyn(option Nodes) *client.Client {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%d", option.Host, option.Port),
		Username: option.Username,
		Password: option.Password,
	})
	if err != nil {
		panic(fmt.Sprintf("InfluxDB数据库连接失败.%+v", err))
	}
	return &c
}
