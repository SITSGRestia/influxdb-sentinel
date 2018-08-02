package influxdb

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/sirupsen/logrus"
	"errors"
	"time"
)

type Ops struct {
	client client.Client
	//options InfluxDB
}

type BaseData struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func NewOps(clnt client.Client) *Ops {
	return &Ops{client: clnt}
}

// queryDB convenience function to query the database
func (ops *Ops) QueryDB(cmd, database string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: database,
	}

	if response, err := ops.client.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

// queryDB convenience function to query the database（集群同步）
func QueryDBSyn(cmd, database string, clnt client.Client) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: database,
	}

	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

//格式化Result
func (ops *Ops) GetData(res []client.Result, timeFormat string) (*[][]BaseData, error) {
	if res == nil || len(res) == 0 || len(res[0].Series) == 0 || len(res[0].Series[0].Values) == 0 {
		return nil, errors.New(fmt.Sprint("查询结果为空."))
	}
	data := make([][]BaseData, len(res[0].Series[0].Values[0])-1)
	for i, row := range res[0].Series[0].Values {
		for j, _ := range res[0].Series[0].Columns {
			if j == 0 {
				t, err := time.Parse(time.RFC3339, row[0].(string))
				if err != nil {
					logrus.Error(err)
				}
				for k := 0; k < len(row)-1; k++ {
					data[k] = append(data[k], BaseData{Name: t.Local().Format(timeFormat)})
				}
			} else {
				data[j-1][i].Value = row[j]
			}
		}
	}
	return &data, nil
}

//存储Result（集群同步）
func (ops *Ops) SetPointsSyn(data []map[string]interface{}, database string, measurement Measurement, policyName string) error {
	tags := []map[string]string{}
	fields := []map[string]interface{}{}
	timestamps := []time.Time{}
	for _, m := range data {
		tag := map[string]string{}
		field := map[string]interface{}{}
		var timestamp *time.Time
		var err error
		for k, v := range m {
			if measurement.PrimaryKey == k {
				tag[k] = v.(string)
				continue
			}
			if "time" == k {
				timestamp, err = ConvertStringToTime(time.RFC3339, v.(string), time.Local)
				if err != nil {
					return errors.New(fmt.Sprint("查询出的时间转time.Time错误（集群同步）"))
				}
				continue
			}
			field[k] = v
		}
		tags = append(tags, tag)
		fields = append(fields, field)
		timestamps = append(timestamps, *timestamp)
	}
	err := ops.WritePointsRetentionBatch("s", policyName, measurement.Name, database, tags, fields, timestamps)
	if err != nil {
		return errors.New(fmt.Sprintf("批量添加操作失败(Retention)(集群同步).%+v", err))
	}
	return nil
}

//批量存储
func (ops *Ops) WritePointsRetentionBatch(precision, policyName, measurement, database string, tags []map[string]string, fields []map[string]interface{}, timestamps []time.Time) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:        database,
		Precision:       precision,
		RetentionPolicy: policyName,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("新建批量添加操作失败(Retention).%+v", err))
	}

	for i := 0; i < len(tags); i++ {
		pt, err := client.NewPoint(
			measurement,
			tags[i],
			fields[i],
			timestamps[i],
		)
		if err != nil {
			return errors.New(fmt.Sprintf("新建Point失败(Retention).%+v", err))
		}
		bp.AddPoint(pt)
	}

	if err := ops.client.Write(bp); err != nil {
		return errors.New(fmt.Sprintf("存储Point失败(Retention).%+v", err))
	}
	return nil
}

//格式字符串转Time
//time.ParseInLocation("20060102150405","20180223100000",time.Local)
func ConvertStringToTime(layout, value string, location *time.Location) (*time.Time, error) {
	t, err := time.ParseInLocation(layout, value, location)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
