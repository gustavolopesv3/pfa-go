package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"gitgub.com/gustavolopesv3/pfa-go/internal/order/infra/database"
	"gitgub.com/gustavolopesv3/pfa-go/internal/order/usecase"
	"gitgub.com/gustavolopesv3/pfa-go/pkg/rabbitmq"
	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

func worker(deliveryMessage <-chan amqp.Delivery, uc *usecase.CalculateFinalPriceUseCase, workerId int) {
	for msg := range deliveryMessage {
		var input usecase.OrderInputDTO
		err := json.Unmarshal(msg.Body, &input)
		if err != nil {
			fmt.Println("Error unmarshalling message", err)
		}
		input.Tax = 10.0
		_, err = uc.Execute(input)
		if err != nil {
			fmt.Println("Error Execute message", err)
		}
		msg.Ack(false)
		fmt.Println("Worker", workerId, "processed order", input.ID)
		// time.Sleep(1 * time.Second)
	}
}

func main() {
	maxWorkers := 100
	waitgroup := sync.WaitGroup{}
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/orders")

	if err != nil {
		log.Fatal("ERRO AO ABRIR CONEXÃO COM MYSQL", err)
		panic(err)
	}
	defer db.Close()

	repository := database.NewOrderRepository(db)
	usecaseCaller := usecase.NewCalculateFinalPriceUseCase(repository)
	channelRabbitMq, err := rabbitmq.OpenChannel()
	channelRabbitMq.Qos(10, 0, false)
	if err != nil {
		panic("Error on open a new channel")
	}
	defer channelRabbitMq.Close()

	outChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(channelRabbitMq, outChannel)

	waitgroup.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		defer waitgroup.Done()
		go worker(outChannel, usecaseCaller, i)
	}

	waitgroup.Wait()
	println("PROCESSADAS")

}
