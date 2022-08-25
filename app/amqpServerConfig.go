package app

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func ConnectToAMQPServer() *amqp.Connection {
	amqpUser := os.Getenv("AMQP_USER")
	amqpPassword := os.Getenv("AMQP_PASSWORD")
	amqpHost := os.Getenv("AMQP_HOST")
	amqpPort := os.Getenv("AMQP_PORT")
	amqpVhost := os.Getenv("AMQP_VHOST")

	amqpConnection, err := amqp.Dial(getServerConnectionURL(amqpUser, amqpPassword, amqpHost, amqpPort, amqpVhost))
	if err != nil {
		panic(err)
	}
	return amqpConnection
}

func getServerConnectionURL(amqpUser string, amqpPassword string, amqpHost string, amqpPort string, amqpVhost string) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", amqpUser, amqpPassword, amqpHost, amqpPort, amqpVhost)
}
