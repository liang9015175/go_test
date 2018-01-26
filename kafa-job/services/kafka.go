package services

import (
	"awesomeProject/kafa-job/models"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"time"
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"fmt"
	"github.com/astaxie/beego"
)

const (
	MyTOPIC = "my_topic"
)

var p sarama.AsyncProducer

func asyncProducer(item *models.ElasticItem) {
	data, _ := json.Marshal(item)
	msg := &sarama.ProducerMessage{
		Topic: MyTOPIC,
		Value: sarama.ByteEncoder(string(data)),
	}
	beego.Info(fmt.Sprintf("推送消息到kafka中: %s",string(data)))
	p.Input()<-msg

}

func connectKafka() {
	//init (custom) config ,enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	//init consumer
	brokers := []string{"127.0.0.1:9092"}
	topics := []string{MyTOPIC}
	consumer, err := cluster.NewConsumer(brokers, "my-consumer-group", topics, config)
	if err != nil {
		panic(err)
	}

	//trap sigint to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	//consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Println("error: "+ err.Error())
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			log.Println(ntf)
		}
	}()
	go func() {
		for {
			select {
			case msg, ok := <-consumer.Messages():
				if ok {
					beego.Info( fmt.Sprintf("接收到消息: %s", msg.Value))

					index(msg.Value)
				}
			}
		}
	}()

}
func init() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	//connect to broker
	p, _ = sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	//必须有这个匿名函数内容
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					beego.Error(err)
				}
			case <-success:
			}
		}
	}(p)

	connectKafka()

}
