package storeData

import (
)

//func client() *Repository {
//	var filePath string
//	flag.StringVar(&filePath, "config", "etc/ic_local_lyf.json", "指定配置文件的路径")
//	flag.Parse()
//
//	ops1 := redis.New(redis.Redis{"127.0.0.1", 6379, "dell@123", 0, 100, 1})
//
//	dataJsonRepository := New(false, `D:\GOPATH\src\per.zdh.org\ic\cmd\data\data.json`, ops1)
//	//model := New(fmt.Sprintf("%s/src/per.zdh.org/cmdb/cmd/data", home), fmt.Sprintf("file:///%s/src/per.zdh.org/cmdb/cmd/resource", home))
//	return dataJsonRepository
//}
//
//func TestRepository_FindData(t *testing.T) {
//	data, err := client().FindData()
//	if err != nil {
//		fmt.Println("查询结构体失败", err)
//	}
//	fmt.Println(data)
//}
//
//func TestRepository_Finds(t *testing.T) {
//	data, err := client().Finds()
//	if err != nil {
//		fmt.Println("查询data.json失败", err)
//	}
//	fmt.Println(data)
//}