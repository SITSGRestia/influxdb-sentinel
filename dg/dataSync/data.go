package dataSync

import (
	"sync"
	"per.zdh.org/influxdb-sentinel/db/influxdb"
	"per.zdh.org/influxdb-sentinel/global"
	"github.com/influxdata/influxdb/client/v2"
	"per.zdh.org/influxdb-sentinel/util"
	"github.com/sirupsen/logrus"
	"time"
	"fmt"
	"errors"
	"strconv"
)

type Service struct {
	influxOps     *influxdb.Ops
	clusterNodes  *[]influxdb.Nodes
	sentinelState *influxdb.SentinelState
	syncChannel   chan int64
	mutex         sync.RWMutex
}

func New(influxOps *influxdb.Ops, clusterNodes *[]influxdb.Nodes, sentinelState *influxdb.SentinelState, syncChannel chan int64, ) *Service {
	return &Service{influxOps: influxOps, clusterNodes: clusterNodes, sentinelState: sentinelState, syncChannel: syncChannel}
}

//目前仅支持s级精度
func (service *Service) DataSyncService() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("数据同步进程panic错误:%s", err)
		}
	}()
	for {
		<-service.syncChannel
		newTime := int64(0)
		downTime := int64(0)
		var clnt *client.Client
		for _, v := range service.sentinelState.DataInfo {
			for _, measurement := range v.MeasurementInfo {
				for _, policyName := range v.RetentionPolicy {
					for _, node := range *service.clusterNodes {
						flag, err := util.IsLocalIP(node.Host)
						if err != nil {
							logrus.Fatalf("判断本地地址错误:%s", err)
						}
						if flag {
							continue
						}
						for {
							otherClnt := influxdb.NewSyn(node)
							nodeTime, isEmpty, err := service.GetOtherSyncData(otherClnt, policyName, v.Database, measurement)
							if isEmpty == global.HANORESULT {
								break
							}
							if err != nil {
								logrus.Errorf("同步数据错误(查询其他节点最新时间):%s", err)
								time.Sleep(time.Second * 2)
								continue
							}
							if nodeTime > newTime {
								newTime = nodeTime
								clnt = otherClnt
							}
							break
						}
					}
					if clnt == nil {
						logrus.Errorf("同步数据错误(查询其他节点最新时间，所有节点数据均为空).database:%s,measurement:%s,policyName:%s", v.Database, measurement.Name, policyName)
						continue
					}
					for {
						localTime, isEmpty, err := service.GetLocalSyncData(policyName, v.Database, measurement)
						if isEmpty == global.HANORESULT {
							break
						}
						if err != nil {
							logrus.Errorf("同步数据错误(查询本地节点最新时间):%s", err)
							time.Sleep(time.Second * 2)
							continue
						}
						downTime = localTime
						break
					}
					for {
						startTime := strconv.Itoa(int(downTime)) + "s"
						endTime := strconv.Itoa(int(newTime)) + "s"
						isEmpty, err := service.SynData(*clnt, startTime, endTime, policyName, v.Database, measurement)
						if isEmpty == global.HANORESULT {
							break
						}
						if err != nil {
							logrus.Errorf("同步数据错误:%s", err)
							time.Sleep(time.Second * 2)
							continue
						}
						break
					}
				}
			}
		}
	}
}

//同步数据(查询本地节点最新时间)
func (service *Service) GetLocalSyncData(policyName, database string, measurement influxdb.Measurement) (int64, int, error) {
	res, err := service.influxOps.QueryDB("select * from \""+policyName+"\"."+measurement.Name+" order by time desc LIMIT 1", database)
	if err != nil {
		return 0, global.HAHASRESULT, errors.New(fmt.Sprintf("同步数据(查询最新时间)，查询时序数据库失败.%+v", err))
	}
	resMap, err, isEmpty := GetDataSyn(res)
	if err != nil {
		return 0, isEmpty, errors.New(fmt.Sprintf("同步数据(查询最新时间)，格式化数据失败.%+v", err))
	}
	lastTime, err := GetTimeSyn(*resMap)
	if err != nil {
		return 0, global.HAHASRESULT, errors.New(fmt.Sprintf("同步数据(查询最新时间)，格式化数据失败.%+v", err))
	}
	return lastTime, global.HAHASRESULT, nil
}

//同步数据(查询其他节点最新时间)
func (service *Service) GetOtherSyncData(clnt *client.Client, policyName, database string, measurement influxdb.Measurement) (int64, int, error) {
	res, err := influxdb.QueryDBSyn("select * from \""+policyName+"\"."+measurement.Name+" order by time desc LIMIT 1", database, *clnt)
	if err != nil {
		return 0, global.HAHASRESULT, errors.New(fmt.Sprintf("同步数据(查询最新时间)，查询时序数据库失败.%+v", err))
	}
	resMap, err, isEmpty := GetDataSyn(res)
	if err != nil {
		return 0, isEmpty, errors.New(fmt.Sprintf("同步数据(查询最新时间)，格式化数据失败.%+v", err))
	}
	lastTime, err := GetTimeSyn(*resMap)
	if err != nil {
		return 0, global.HAHASRESULT, errors.New(fmt.Sprintf("同步数据(查询最新时间)，格式化数据失败.%+v", err))
	}
	return lastTime, global.HAHASRESULT, nil
}

//同步数据
func (service *Service) SynData(clnt client.Client, startTime, endTime, policyName, database string, measurement influxdb.Measurement) (int, error) {
	res, err := influxdb.QueryDBSyn("select * from \""+policyName+"\"."+measurement.Name+" where time >= "+startTime+" and time <= "+endTime+" order by time asc", database, clnt)
	if err != nil {
		return global.HAHASRESULT, errors.New(fmt.Sprintf("同步数据，查询时序数据库失败.%+v", err))
	}
	resMap, err, isEmpty := GetDataSyn(res)
	if err != nil {
		return isEmpty, errors.New(fmt.Sprintf("同步数据，格式化数据失败.%+v", err))
	}
	err = service.influxOps.SetPointsSyn(*resMap, database, measurement, policyName)
	if err != nil {
		return global.HAHASRESULT, errors.New(fmt.Sprintf("同步数据，存储数据到本地失败.%+v", err))
	}
	return global.HAHASRESULT, nil
}

//格式化Result（集群同步）
func GetDataSyn(res []client.Result) (*[]map[string]interface{}, error, int) {
	if res == nil || len(res) == 0 || len(res[0].Series) == 0 || len(res[0].Series[0].Values) == 0 {
		return nil, errors.New(fmt.Sprint("查询结果为空（集群同步）.")), global.HANORESULT
	}
	resMap := []map[string]interface{}{}
	for _, row := range res[0].Series[0].Values {
		innerMap := map[string]interface{}{}
		for j, innerRow := range res[0].Series[0].Columns {
			innerMap[innerRow] = row[j]
		}
		resMap = append(resMap, innerMap)
	}
	return &resMap, nil, global.HAHASRESULT
}

func GetTimeSyn(data []map[string]interface{}) (int64, error) {
	var timestamp *time.Time
	for _, m := range data {
		var err error
		for k, v := range m {
			if "time" == k {
				timestamp, err = util.ConvertStringToTime(time.RFC3339, v.(string), time.Local)
				if err != nil {
					return 0, errors.New(fmt.Sprint("查询出的时间转time.Time错误（查询最新时间）"))
				}
				break
			}
		}
		break
	}
	return timestamp.Unix(), nil
}
