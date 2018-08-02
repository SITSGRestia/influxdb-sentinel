package storeData

import (
	"github.com/emicklei/go-restful"
	"per.zdh.org/influxdb-sentinel/model/storeData"
)

type Endpoints struct {
	version    string
	repository *storeData.Repository
}

func New(version string, repository *storeData.Repository) *Endpoints {
	return &Endpoints{version: version, repository: repository}
}

func (endPoints *Endpoints) WebService() *restful.WebService {
	//tags := []string{"StoreData"}
	ws := new(restful.WebService)
	ws.ApiVersion(endPoints.version).Path("/storeData").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	//ws.Route(ws.GET("").To(endPoints.findAll).
	//	Doc("data.json数据查询接口").Metadata(restfulspec.KeyOpenAPITags, tags).
	//	Returns(http.StatusOK, "查询数据成功.", storeData.DataJson{}))
	return ws
}

//// 查询data.json数据
//func (endPoints *Endpoints) findAll(request *restful.Request, response *restful.Response) {
//	data, err := endPoints.repository.Finds()
//	if err != nil {
//		response.WriteHeaderAndEntity(global.BAD_REQUEST, global.NewErrorMsg(global.BAD_REQUEST, "查询数据失败.", err.Error()))
//		return
//	}
//	response.WriteAsJson(data)
//}
