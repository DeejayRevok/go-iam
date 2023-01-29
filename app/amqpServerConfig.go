package app

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func ConnectToAMQPServer() *amqp.Connection {
	amqpUser := os.Getenv("IAM_AMQP_USER")
	amqpPassword := os.Getenv("IAM_AMQP_PASSWORD")
	amqpHost := os.Getenv("IAM_AMQP_HOST")
	amqpPort := os.Getenv("IAM_AMQP_PORT")
	amqpVhost := os.Getenv("IAM_AMQP_VHOST")

	amqpConnection, err := amqp.Dial(getServerConnectionURL(amqpUser, amqpPassword, amqpHost, amqpPort, amqpVhost))
	if err != nil {
		panic(err)
	}
	return amqpConnection
}

func getServerConnectionURL(amqpUser string, amqpPassword string, amqpHost string, amqpPort string, amqpVhost string) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", amqpUser, amqpPassword, amqpHost, amqpPort, amqpVhost)
}
