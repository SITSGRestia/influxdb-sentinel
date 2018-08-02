package influxdb

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
)

func newClient() (*Ops, error) {
	influxDbConfig := InfluxDB{
		"127.0.0.1",
		8086,
		"admin",
		"dell@123",
		"datacenter",
	}
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%d", influxDbConfig.Host, influxDbConfig.Port),
		Username: influxDbConfig.Username,
		Password: influxDbConfig.Password,
	})
	if err != nil {
		return nil, err
	}

	influxOps := NewOps(c)
	return influxOps, nil
}

//func TestInfluxDB_CreateDB(t *testing.T) {
//	client, err := newClient()
//	if err != nil {
//		t.Fatalf("InfluxDB连接错误.%+v", err)
//	}
//	err = client.CreateDB()
//	if err != nil {
//		t.Fatalf("创建数据库失败%+v", err)
//	}
//	t.Logf("创建数据库成功")
//}

//func TestInfluxDB_GetData(t *testing.T) {
//	client, err := newClient()
//	if err != nil {
//		t.Fatalf("InfluxDB连接错误.%+v", err)
//	}
//
//	err = client.CreateRetentionPolicy("1month", "31d", 1)
//	if err != nil {
//		t.Fatalf("创建保留策略失败%+v", err)
//	}
//	t.Logf("创建保留策略成功")
//
//	tags := map[string]string{
//		"id": "id1",
//	}
//	fields := map[string]interface{}{
//		"val": 1.0,
//		"gu":  2.0,
//		"gd":  3.0,
//		"su":  4.0,
//		"sd":  5.0,
//	}
//	timestamp := time.Now()
//	err = client.WritePointsRetention("s", "1month", "dc2", tags, fields, timestamp)
//	if err != nil {
//		t.Fatalf("存储Point失败(Retention)%+v", err)
//	}
//	t.Logf("存储Point成功(Retention)")
//
//	tags = map[string]string{
//		"id": "id2",
//	}
//	fields = map[string]interface{}{
//		"val": 6.0,
//		"gu":  7.0,
//		"gd":  8.0,
//		"su":  9.0,
//		"sd":  10.0,
//	}
//	err = client.WritePointsRetention("s", "1month", "dc2", tags, fields, timestamp)
//	if err != nil {
//		t.Fatalf("存储Point失败(Retention)%+v", err)
//	}
//	t.Logf("存储Point成功(Retention)")
//
//	res, err := client.QueryDB("select val,gu,gd,su,sd from \"1month\".dc2 where id = 'id1'")
//	if err != nil {
//		t.Fatalf("查询指定数据失败%+v", err)
//	}
//	t.Logf("查询指定数据成功")
//
//	_, err = client.GetData(res, "2006/01/02 15:04:05")
//	if err != nil {
//		t.Fatalf("数据格式化失败%+v", err)
//	}
//	t.Logf("数据格式化成功")
//}

//func TestInfluxDB_BatchSetData(t *testing.T) {
//	client, err := newClient()
//	if err != nil {
//		t.Fatalf("InfluxDB连接错误.%+v", err)
//	}
//
//	tag1 := map[string]string{
//		"id": "id3",
//	}
//	field1 := map[string]interface{}{
//		"val": 11.0,
//		"gu":  22.0,
//		"gd":  33.0,
//		"su":  44.0,
//		"sd":  55.0,
//	}
//	timestamp1 := time.Now()
//
//	tag2 := map[string]string{
//		"id": "id4",
//	}
//	field2 := map[string]interface{}{
//		"val": 16.0,
//		"gu":  17.0,
//		"gd":  18.0,
//		"su":  19.0,
//		"sd":  20.0,
//	}
//
//	timestamp2 := time.Now()
//
//	tags := []map[string]string{
//		tag1,
//		tag2,
//	}
//	fields := []map[string]interface{}{
//		field1,
//		field2,
//	}
//
//	timestamp := []time.Time{
//		timestamp1,
//		timestamp2,
//	}
//	err = client.WritePointsRetentionBatch("s", "1month", "dc2", tags, fields, timestamp)
//	if err != nil {
//		t.Fatalf("存储Point失败(Retention)%+v", err)
//	}
//	t.Logf("存储Point成功(Retention)")
//
//	res, err := client.QueryDB("select val,gu,gd,su,sd from \"1month\".dc2 where id = 'id1'")
//	if err != nil {
//		t.Fatalf("查询指定数据失败%+v", err)
//	}
//	t.Logf("查询指定数据成功")
//
//	_, err = client.GetData(res, "2006/01/02 15:04:05")
//	if err != nil {
//		t.Fatalf("数据格式化失败%+v", err)
//	}
//	t.Logf("数据格式化成功")
//}
