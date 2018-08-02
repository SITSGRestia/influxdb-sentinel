package storeData

import (
	"sync"
	"per.zdh.org/influxdb-sentinel/db/influxdb"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	influxOps *influxdb.Ops
	mutex     sync.RWMutex
}

func New(influxOps *influxdb.Ops) *Repository {
	return &Repository{influxOps: influxOps}
}

// 保存InfluxDB的Point数据
func (repository *Repository) Save(b []byte) (err error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	data := influxdb.StoreData{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		logrus.Errorf("保存InfluxDB的Point数据,保存错误,解序列化错误.%+v", err)
		return err
	}
	repository.influxOps.WritePointsRetentionBatch(data.Precision, data.PolicyName, data.Measurement, data.Database, data.Tags, data.Fields, data.Timestamps)
	if err != nil {
		logrus.Errorf("保存InfluxDB的Point数据,保存错误,批量保存错误.%+v", err)
		return err
	}
	return nil
}
