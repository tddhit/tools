package msgqueue

import (
	"time"

	etcd "github.com/coreos/etcd/clientv3"

	"github.com/tddhit/go-nsq"
	"github.com/tddhit/tools/log"
	"github.com/tddhit/tools/msgqueue/option"
	"github.com/tddhit/wox/naming"
)

type NsqConsumer struct {
	consumer map[string]*nsq.Consumer
	msgChan  map[string]chan *Message
}

func NewNsqConsumer(etcdClient *etcd.Client, opt option.NsqConsumer) (c *NsqConsumer, err error) {
	c = &NsqConsumer{
		consumer: make(map[string]*nsq.Consumer),
		msgChan:  make(map[string]chan *Message),
	}
	r := &naming.Resolver{
		Client:  etcdClient,
		Timeout: 2 * time.Second,
	}
	addrs := r.Resolve(opt.Registry)
	cfg := nsq.NewConfig()
	for _, topic := range opt.Topics {
		var consumer *nsq.Consumer
		if consumer, err = nsq.NewConsumer(topic, opt.Channel, cfg); err != nil {
			log.Error(err)
			return
		}
		c.msgChan[topic] = make(chan *Message, 1000)
		consumer.AddHandler(&nsqHandler{c.msgChan[topic]})
		if err = consumer.ConnectToNSQDs(addrs); err != nil {
			log.Error(err)
			return
		}
		c.consumer[topic] = consumer
	}
	return
}

func (c *NsqConsumer) Messages(topic string) <-chan *Message {
	return c.msgChan[topic]
}

type nsqHandler struct {
	msgChan chan<- *Message
}

func (h *nsqHandler) HandleMessage(message *nsq.Message) error {
	msg := &Message{
		Body: message.Body,
	}
	h.msgChan <- msg
	return nil
}
