package rabbit

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var MQConn *amqp.Connection

func ConnectRabbitMQ(dsn string) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("连接RabbitMQ失败，err：%s", err))
	}
	MQConn = conn
}

func Consumer(ip string, port int64, username string, password string) {
	// 创建连接
	ConnectRabbitMQ(fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, ip, port))
	defer MQConn.Close()

	// 创建通道
	channel, err := MQConn.Channel()
	if err != nil {
		log.Fatal(fmt.Errorf("创建通道失败，err：%s", err.Error()))
	}

	// 获取消息
	msgs, err := channel.Consume("test", "consume1", false, false, false, false, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("获取消息，err：%s", err.Error()))
	}
	for msg := range msgs {
		fmt.Println(msg.DeliveryTag, string(msg.Body))
		msg.Ack(false)
	}
}

func Producer(ip string, port int64, username string, password string) {
	// 创建连接
	ConnectRabbitMQ(fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, ip, port))
	defer MQConn.Close()

	// 创建通道
	channel, err := MQConn.Channel()
	if err != nil {
		log.Fatal(fmt.Errorf("创建通道失败，err：%s", err.Error()))
	}

	// 创建队列
	query, err := channel.QueueDeclare("test", false, false, false, false, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("创建队列失败，err：%s", err.Error()))
	}

	// 发送消息
	if err := channel.Publish("", query.Name, false, false, amqp.Publishing{
		ContentType: "text/plan",
		Body:        []byte("001"),
	}); err != nil {
		log.Fatal(fmt.Errorf("发送消息失败，%s", err.Error()))
	}
	log.Fatalln("发送消息成功！")
}
