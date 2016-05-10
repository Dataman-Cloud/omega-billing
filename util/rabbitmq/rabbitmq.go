package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/Dataman-Cloud/omega-billing/config"
	"github.com/Dataman-Cloud/omega-billing/dao"
	"github.com/Dataman-Cloud/omega-billing/model"
	"github.com/Dataman-Cloud/omega-billing/util"
	"github.com/Jeffail/gabs"
	log "github.com/cihub/seelog"
	"github.com/streadway/amqp"
	"strconv"
	"strings"
	"time"
)

var channel *amqp.Channel
var connection *amqp.Connection

func init() {
	var err error
	connection, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", GetConfig().Mq.User, GetConfig().Mq.PassWord, GetConfig().Mq.Host, GetConfig().Mq.Port))
	if err != nil {
		log.Errorf("connection rabbitmq error: %v", err)
		log.Flush()
		panic(-1)
	}
	channel, err = connection.Channel()
	if err != nil {
		log.Errorf("get rabbitmq channel error: %v", err)
		log.Flush()
		panic(-1)
	}
	queue, err := DeclareQueue(channel, GetConfig().Mq.ConsumeName)
	if err != nil {
		log.Errorf("declare queue %s error: %v", GetConfig().Mq.ConsumeName, err)
		log.Flush()
		panic(-1)
	}
	err = channel.QueueBind(queue.Name, GetConfig().Mq.Routingkey, GetConfig().Mq.Exchange, false, nil)
	if err != nil {
		log.Errorf("queue bind queuename:%s key:%s exchangename:%s error: %v", queue.Name, GetConfig().Mq.Routingkey, GetConfig().Mq.Exchange, err)
		log.Flush()
		panic(-1)
	}
}

func DeclareQueue(channel *amqp.Channel, name string) (amqp.Queue, error) {
	args := amqp.Table{
		"x-message-ttl": GetConfig().Mq.MessageTTL,
		"x-expires":     GetConfig().Mq.QueueTTL,
	}
	return channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		args,
	)
}

func ReciveQueue(queue string) (<-chan amqp.Delivery, error) {
	cs, err := channel.Consume(queue, "", true, false, false, false, nil)
	return cs, err
}

func Close() {
	if channel != nil {
		channel.Close()
	}
	if connection != nil {
		connection.Close()
	}
}

func Run() {
	cs, err := ReciveQueue(GetConfig().Mq.ConsumeName)
	if err != nil {
		log.Errorf("rabbitmq subscribe error: %v", err)
		panic(-1)
	}
	for {
		select {
		case c := <-cs:
			EventProcess(c.Body)
		}
	}
}

func EventProcess(body []byte) error {
	var message model.Message
	err := json.Unmarshal(body, &message)
	message.Method = message.Task["method"].(string)
	message.Meta = message.Task["metadata"].(string)
	message.Path = message.Task["path"].(string)
	if err != nil {
		log.Errorf("unmarshal message error: %v", err)
		return err
	}
	switch message.Method {
	case "POST":
		event, err := newEvent(&message)
		if err == nil {
			dao.AddEvent(event)
		}
	case "DELETE":
		event, err := newDeleteEvent(&message)
		if err == nil {
			event.Active = false
			dao.DeleteApp(event)
		}
	case "PUT":
		event, err := newUpdateEvent(&message)
		if err == nil {
			event.Active = false
			dao.UpdateApp(event)
		}
	}
	return nil
}

func newUpdateEvent(message *model.Message) (*model.Event, error) {
	cid, err := strconv.ParseUint(message.ClusterId, 10, 64)
	if err != nil {
		log.Errorf("string cluserid parse to uint error: %v", err)
		return nil, err
	}
	mjson, err := gabs.ParseJSON([]byte(message.Meta))
	if err != nil {
		log.Errorf("string marathon parse to json error: %v", err)
		return nil, err
	}
	id, err := util.ParseAppAlias(strings.Replace(message.Path, "/v2/apps/", "", 1))
	if err != nil {
		log.Errorf("base32 stdencoding id error1: %v %s", err, strings.Replace(message.Path, "/v2/apps/", "", 1))
		return nil, err
	}
	ids := strings.SplitN(id, ":", 2)
	if len(ids) != 2 {
		log.Errorf("split marathon id is not 2 len: %d", len(ids))
		return nil, errors.New("split marathon id is not 2 len")
	}
	uid, err := strconv.ParseUint(ids[0], 10, 64)
	if err != nil {
		log.Errorf("parse uid string to uint64 error: %v", err)
		return nil, err
	}
	//timen := time.Now()
	timen := time.Unix(message.Timestamp, 0)
	event := &model.Event{
		Cid:        cid,
		CreateTime: timen,
		EndTime:    timen,
		Active:     true,
		Uid:        uid,
		AppName:    ids[1],
	}
	billing, err := dao.GetBilling(event)
	if cpus := mjson.Path("cpus").Data(); cpus != nil {
		event.Cpus = cpus.(float64)
	} else {
		if err != nil {
			log.Errorf("get billing error1: %v", err)
			return nil, err
		}
		event.Cpus = billing.Cpus
	}
	if mem := mjson.Path("mem").Data(); mem != nil {
		event.Mem = mem.(float64)
	} else {
		if err != nil {
			log.Errorf("get billing error2: %v", err)
			return nil, err
		}
		event.Mem = billing.Mem
	}
	if instances := mjson.Path("instances").Data(); instances != nil {
		event.Instances = uint32(instances.(float64))
	} else {
		if err != nil {
			log.Errorf("get billing error3: %v", err)
			return nil, err
		}
		event.Instances = billing.Instances
	}
	event.CreateTime = billing.CreateTime
	return event, nil
}

func newDeleteEvent(message *model.Message) (*model.Event, error) {
	cid, err := strconv.ParseUint(message.ClusterId, 10, 64)
	if err != nil {
		log.Errorf("string cluserid parse to uint error: %v", err)
		return nil, err
	}
	id, err := util.ParseAppAlias(strings.Replace(message.Path, "/v2/apps/", "", 1))
	if err != nil {
		log.Errorf("base32 stdencoding id error1: %v %s", err, strings.Replace(message.Path, "/v2/apps/", "", 1))
		return nil, err
	}
	ids := strings.SplitN(id, ":", 2)
	if len(ids) != 2 {
		log.Errorf("split marathon id is not 2 len: %d", len(ids))
		return nil, errors.New("split marathon id is not 2 len")
	}
	uid, err := strconv.ParseUint(ids[0], 10, 64)
	if err != nil {
		log.Errorf("parse uid string to uint64 error: %v", err)
		return nil, err
	}
	//timen := time.Now()
	timen := time.Unix(message.Timestamp, 0)
	event := &model.Event{
		Cid:        cid,
		CreateTime: timen,
		EndTime:    timen,
		Active:     true,
		Uid:        uid,
		AppName:    ids[1],
	}
	billing, err := dao.GetBilling(event)
	if err != nil {
		log.Errorf("get billing error1: %v", err)
		return nil, err
	}
	event.Cpus = billing.Cpus
	event.Mem = billing.Mem
	event.Instances = billing.Instances
	return event, nil
}

func newEvent(message *model.Message) (*model.Event, error) {
	cid, err := strconv.ParseUint(message.ClusterId, 10, 64)
	if err != nil {
		log.Errorf("string cluserid parse to uint error: %v", err)
		return nil, err
	}
	mjson, err := gabs.ParseJSON([]byte(message.Meta))
	if err != nil {
		log.Errorf("string marathon parse to json error: %v", err)
		return nil, err
	}
	id, err := util.ParseAppAlias(mjson.Path("id").Data().(string))
	if err != nil {
		log.Errorf("base32 stdencoding id error2: %v", err)
		return nil, err
	}
	ids := strings.SplitN(id, ":", 2)
	if len(ids) != 2 {
		log.Errorf("split marathon id is not 2 len: %d", len(ids))
		return nil, errors.New("split marathon id is not 2 len")
	}
	uid, err := strconv.ParseUint(ids[0], 10, 64)
	if err != nil {
		log.Errorf("parse uid string to uint64 error: %v", err)
		return nil, err
	}
	cpus := mjson.Path("cpus").Data().(float64)
	mem := mjson.Path("mem").Data().(float64)
	instances := mjson.Path("instances").Data().(float64)
	//timen := time.Now()
	timen := time.Unix(message.Timestamp, 0)
	event := &model.Event{
		Cid:        cid,
		CreateTime: timen,
		EndTime:    timen,
		Active:     true,
		Uid:        uid,
		AppName:    ids[1],
		Cpus:       cpus,
		Mem:        mem,
		Instances:  uint32(instances),
	}
	return event, nil
}
