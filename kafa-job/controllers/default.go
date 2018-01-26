package controllers

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"github.com/bsm/sarama-cluster"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"
	"awesomeProject/kafa-job/services"
	"github.com/robfig/cron"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

const (
	MyTOPIC     = "my_topic"
)

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
	defer consumer.Close() //close the connection

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
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintln(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}
	}

}
func syncProducer() {
	//settings
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	//connect to broker
	p, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)

	//close the connection

	if err != nil {
		panic(err)
	}
	defer p.Close()
	v := "sync: " + strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000))
	fmt.Fprintln(os.Stdout, v)
	msg := sarama.ProducerMessage{
		Topic: MyTOPIC,
		Value: sarama.ByteEncoder(v),
	}
	if _, _, err := p.SendMessage(&msg); err != nil {
		log.Println( "error: "+err.Error())
		return
	}

}
func init()  {
	cron := cron.New()
	spec := "0/5 * * * * ?"
	cron.AddFunc(spec, services.ScheduleImport)
	cron.Start()


}
