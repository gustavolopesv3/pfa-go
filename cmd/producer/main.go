package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Order struct {
	ID    string
	Price float64
}

func GenerateOrder() Order {
	return Order{
		ID:    uuid.New().String(),
		Price: rand.Float64() * 100,
	}
}

func Notify(ch *amqp.Channel, order Order) error {
	body, err := json.Marshal(order) // convert to json
	if err != nil {
		return err
	}

	err = ch.Publish(
		"amq.direct", // exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()        // fechar conexão no final
	ch, err := conn.Channel() // abre novo canal
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	for i := 0; i < 100; i++ {
		order := GenerateOrder()
		err := Notify(ch, order)
		if err != nil {
			panic(err)
		}
		fmt.Println(order)
	}
}
