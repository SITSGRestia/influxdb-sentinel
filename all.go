package ins

import (
	_ "per.zdh.org/influxdb-sentinel/dg"
	_ "per.zdh.org/influxdb-sentinel/global"
	_ "per.zdh.org/influxdb-sentinel/http/storeData"
	_ "per.zdh.org/influxdb-sentinel/logger"
	_ "per.zdh.org/influxdb-sentinel/model"
	_ "per.zdh.org/influxdb-sentinel/mq/emqtt"
	_ "per.zdh.org/influxdb-sentinel/mq/rabbitmq"
	_ "per.zdh.org/influxdb-sentinel/util"
)
