package main

import (
	"flag"
	"per.zdh.org/influxdb-sentinel/global"
	"per.zdh.org/influxdb-sentinel/logger"
	influxDB "per.zdh.org/influxdb-sentinel/db/influxdb"
	"fmt"
	"os"
	"strconv"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
	"net"
	"net/http"
	"github.com/go-openapi/spec"
	"per.zdh.org/influxdb-sentinel/dg"
	"per.zdh.org/influxdb-sentinel/dg/daemonScript"
	"per.zdh.org/influxdb-sentinel/dg/dataSync"
	storeDataRepo "per.zdh.org/influxdb-sentinel/model/storeData"
	storeDataEndpoint "per.zdh.org/influxdb-sentinel/http/storeData"
	"runtime"
)

var filePath string

func init() {
	if runtime.GOOS != "windows" {
		flag.StringVar(&filePath, "config", "/usr/GoPath/src/per.zdh.org/influxdb-sentinel/cmd/etc/ins.json", "指定配置文件的路径")
	} else {
		flag.StringVar(&filePath, "config", "etc/ins.json", "指定配置文件的路径")
	}
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	options := global.New(filePath)
	// 初始化日志
	err := logger.New(options.Log)
	if err != nil {
		panic(fmt.Sprintf("日志配置失败.%+v", err))
	}

	rootPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("获取当前目录失败.%+v", err))
	}

	// 初始化InfluxDB数据库
	influxCli := influxDB.New(options.DB.InfluxDB)
	// emqtt消息队列
	mqttOps := options.MQ.Emqtt
	//// 初始化RabbitMQ消息队列
	//rabbitMQOps := rabbitmq.New(options.MQ.RabbitMQ)
	//数据同步管道
	syncChan := make(chan int64, 0)

	// 服务器地址
	serviceURL := net.JoinHostPort(options.Server.Host, strconv.Itoa(options.Server.Port))

	// data.json数据接口
	influxOps := influxDB.NewOps(*influxCli)
	storeDataRepository := storeDataRepo.New(influxOps)
	storeDataService := storeDataEndpoint.New(options.Version, storeDataRepository)
	storeDataWS := storeDataService.WebService()
	restful.Add(storeDataWS)

	if options.MQType == global.EMQ {
		go mqttOps.Receive(storeDataRepository)
	}

	//InfluxDB存活状态检测服务
	daemonScriptService := daemonScript.New(options.DB.InfluxDB.Path, syncChan)
	//InfluxDB集群宕机数据同步服务
	dataSyncService := dataSync.New(influxOps, &options.ClusterNodes, &options.SentinelState, syncChan)

	schedule := &dg.Schedule{
		Time:                options.DG.DGTime,
		DaemonScriptService: daemonScriptService,
	}
	// 开启定时任务
	go schedule.Starts()
	//开启数据同步服务
	go dataSyncService.DataSyncService()

	restful.Filter(restful.OPTIONSFilter())
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-Custom-Header"},
		AllowedHeaders: []string{"X-Custom-Header", "X-Additional-Header", restful.HEADER_ContentType, restful.HEADER_AccessControlAllowOrigin},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer}
	restful.Filter(cors.Filter)

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		WebServicesURL:                serviceURL,
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject,
	}

	restful.Add(restfulspec.NewOpenAPIService(config))
	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir(rootPath+"/swagger"))))
	fmt.Printf("服务已启动,访问地址:http://%s/apidocs \n", serviceURL)
	http.ListenAndServe(serviceURL, nil)
}

func enrichSwaggerObject(swagger *spec.Swagger) {
	swagger.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "influxdb-sentinel",
			Description: "influxdb-sentinel接口文档",
			Contact: &spec.ContactInfo{
				Name:  "zhoudenghuang",
				Email: "zhoudenghuang@per.zdh.org",
				URL:   "",
			},
			Version: "1.0.0",
		},
	}

	swagger.Tags = []spec.Tag{
		spec.Tag{
			TagProps: spec.TagProps{
				Name:        "influxdb-sentinel",
				Description: "influxdb-sentinel接口",
			},
		},
	}
}
