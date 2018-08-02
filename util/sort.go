package util

import "sort"

import (
	"per.zdh.org/influxdb-sentinel/model/storeData"
	"strconv"
	"github.com/sirupsen/logrus"
	"strings"
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	return p[i].Value > p[j].Value
}

func SortIntByValue(baseData []storeData.BaseInt) (list []storeData.BaseInt) {
	var pairList PairList = make([]Pair, len(baseData))
	for i, data := range baseData {
		pairList[i] = Pair{Key: data.Name, Value: data.Value}
	}
	sort.Sort(pairList)
	list = make([]storeData.BaseInt, len(pairList))
	for i, v := range pairList {
		list[i] = storeData.BaseInt{Name: v.Key, Value: v.Value}
	}
	return list
}

type PairFloat struct {
	Key   string
	Value float64
}

type PairListFloat []PairFloat

func (p PairListFloat) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PairListFloat) Len() int {
	return len(p)
}

func (p PairListFloat) Less(i, j int) bool {
	return p[i].Value > p[j].Value
}

func SortFloat64ByValue(baseData []storeData.BaseFloat64) (list []storeData.BaseFloat64) {
	var pairList PairListFloat = make([]PairFloat, len(baseData))
	for i, data := range baseData {
		pairList[i] = PairFloat{Key: data.Name, Value: data.Value}
	}
	sort.Sort(pairList)
	list = make([]storeData.BaseFloat64, len(pairList))
	for i, v := range pairList {
		list[i] = storeData.BaseFloat64{Name: v.Key, Value: v.Value}
	}
	return list
}

type PairFloat32 struct {
	Key   string
	Value float32
}

type PairListFloat32 []PairFloat32

func (p PairListFloat32) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PairListFloat32) Len() int {
	return len(p)
}

func (p PairListFloat32) Less(i, j int) bool {
	return p[i].Value > p[j].Value
}

func SortFloat32ByValue(baseData []storeData.BaseFloat32) (list []storeData.BaseFloat32) {
	var pairList PairListFloat32 = make([]PairFloat32, len(baseData))
	for i, data := range baseData {
		pairList[i] = PairFloat32{Key: data.Name, Value: data.Value}
	}
	sort.Sort(pairList)
	list = make([]storeData.BaseFloat32, len(pairList))
	for i, v := range pairList {
		list[i] = storeData.BaseFloat32{Name: v.Key, Value: v.Value}
	}
	return list
}

//病毒查杀趋势使用
type BaseInt struct {
	Name  string
	Value int
}
type SortByName []BaseInt

func (p SortByName) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p SortByName) Len() int {
	return len(p)
}

func (p SortByName) Less(i, j int) bool {

	ib := p[i].Name[:2]
	ie := p[i].Name[3:]

	jb := p[j].Name[:2]
	je := p[j].Name[3:]

	ibi, err := strconv.Atoi(ib)
	if err != nil {
		logrus.Error("趋势排序错误")
		return false
	}

	iei, err := strconv.Atoi(ie)
	if err != nil {
		logrus.Error("趋势排序错误")
		return false
	}

	jbi, err := strconv.Atoi(jb)
	if err != nil {
		logrus.Error("趋势排序错误")
		return false
	}

	jei, err := strconv.Atoi(je)
	if err != nil {
		logrus.Error("趋势排序错误")
		return false
	}

	if ibi != jbi {
		return ibi < jbi
	} else {
		return iei < jei
	}

	return p[i].Name > p[j].Name
}

func SortIntByName(baseDateList []storeData.BaseInt) []storeData.BaseInt {

	var pairList SortByName

	for _, data := range baseDateList {
		pairList = append(pairList, BaseInt{Name: data.Name, Value: data.Value})
	}

	sort.Sort(pairList)
	var basedatelist []storeData.BaseInt
	for _, v := range pairList {
		basedatelist = append(basedatelist, storeData.BaseInt{Name: strings.Replace(v.Name, "-", ".", -1), Value: v.Value})
	}
	return basedatelist
}
