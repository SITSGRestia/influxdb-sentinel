package emqtt

import (
	"fmt"
	"net"
	"github.com/jeffallen/mqtt"
	proto "github.com/huin/mqtt"
	"bytes"
	"per.zdh.org/influxdb-sentinel/model/storeData"
	"strconv"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"time"
)

func (ops *Emqtt) Send(msg string) (err error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(ops.Host, strconv.Itoa(ops.Port)))
	if err != nil {
		return err
	}
	cc := mqtt.NewClientConn(conn)
	if err := cc.Connect(ops.Username, ops.Password); err != nil {
		return errors.New(fmt.Sprintf("EMQTT连接错误.%+v", err))
	}
	cc.Publish(&proto.Publish{
		Header: proto.Header{
			Retain: true,
		},
		TopicName: ops.TopicName,
		Payload:   proto.BytesPayload([]byte(msg)),
	})
	cc.Disconnect()
	return nil
}

func (ops *Emqtt) SendKeyValue(key string, val interface{}) (err error) {
	m := map[string]interface{}{key: val}
	mb, err := json.Marshal(m)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", net.JoinHostPort(ops.Host, strconv.Itoa(ops.Port)))
	if err != nil {
		return err
	}
	cc := mqtt.NewClientConn(conn)
	if err := cc.Connect(ops.Username, ops.Password); err != nil {
		return errors.New(fmt.Sprintf("EMQTT连接错误.%+v", err))
	}
	cc.Publish(&proto.Publish{
		Header: proto.Header{
			Retain: true,
		},
		TopicName: ops.TopicName,
		Payload:   proto.BytesPayload(mb),
	})
	cc.Disconnect()
	return nil
}

func (ops *Emqtt) Receive(repository *storeData.Repository) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("InfluxDB数据同步器出现故障，接收器宕机:%s", err)
		}
	}()
	reconnectionTime := time.Duration(10)
	//追加断线重连机制
	for {
		conn, err := net.Dial("tcp", net.JoinHostPort(ops.Host, strconv.Itoa(ops.Port)))
		if err != nil {
			logrus.Errorf("Failed to connect to EMQ, %s: %s", "Reconnection will be executed after several seconds", err)
			time.Sleep(reconnectionTime * time.Second)
			continue
		}
		cc := mqtt.NewClientConn(conn)
		tq := make([]proto.TopicQos, 1)
		tq[0].Topic = ops.TopicName
		tq[0].Qos = proto.QosAtMostOnce

		if err := cc.Connect(ops.Username, ops.Username); err != nil {
			logrus.Errorf("EMQTT消费者创建失败.%+v", err)
			time.Sleep(reconnectionTime * time.Second)
			continue
		}
		cc.Subscribe(tq)

		for m := range cc.Incoming {
			b := make([]byte, m.Payload.Size())
			s := &bytes.Buffer{}
			err := m.Payload.WritePayload(s)
			if err != nil {
				continue
			}
			_, err = s.Read(b)
			if err != nil {
				continue
			}
			for {
				err = repository.Save(b)
				if err != nil {
					logrus.Errorf("存储数据错误(存储数据到当前节点):%s", err)
					time.Sleep(time.Second * 2)
					continue
				}
				break
			}
		}
		logrus.Errorf("MQ接收器，channel关闭，可能已经和EMQ失去连接,%s:%s", "Reconnection will be executed after several seconds", err)
		time.Sleep(reconnectionTime * time.Second)
	}
}
