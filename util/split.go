package util

import "per.zdh.org/influxdb-sentinel/model/storeData"

func Split(data, len int) []int {
	dataList := []int{}

	if data < 0 {
		for i := 0; i < len; i++ {
			dataList = append(dataList, 0)
		}
	} else if data < len {
		for i := 0; i < len; i++ {
			if i < data {
				dataList = append(dataList, 1)
			} else {
				dataList = append(dataList, 0)
			}
		}
	} else {
		d := data % len
		m := data - d
		n := m / len

		for i := 0; i < len; i++ {
			if i < d {
				dataList = append(dataList, n+1)
			} else {
				dataList = append(dataList, n+0)
			}
		}
	}

	return dataList
}

func BaseIntSplit(list *[]storeData.BaseInt, value, size int, unixTime int64) {
	if len(*list) == 0 {
		*list = append(*list, storeData.BaseInt{Name: GetUnixToDayTime(unixTime), Value: value})
	} else {
		lastData := (*list)[len(*list)-1]
		if lastData.Name == GetUnixToDayTime(unixTime) {
			*list = (*list)[0:len(*list)-1]
		}
		*list = append(*list, storeData.BaseInt{Name: GetUnixToDayTime(unixTime), Value: value})

		if len(*list) > size {
			*list = (*list)[1:]
		}
	}
}

func BaseIntSplitStr(list *[]storeData.BaseInt, value, size int, strTime string) {
	if len(*list) == 0 {
		*list = append(*list, storeData.BaseInt{Name: strTime, Value: value})
	} else {
		lastData := (*list)[len(*list)-1]
		if lastData.Name == strTime {
			*list = (*list)[0:len(*list)-1]
		}
		*list = append(*list, storeData.BaseInt{Name: strTime, Value: value})

		if len(*list) > size {
			*list = (*list)[1:]
		}
	}
}
