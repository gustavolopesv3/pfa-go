package main

import (
	"database/sql"
	"log"

	"gitgub.com/gustavolopesv3/pfa-go/internal/order/infra/database"
	"gitgub.com/gustavolopesv3/pfa-go/internal/order/usecase"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "root:root@tcp(mysql-go-intensivo:3306)/orders")

	if err != nil {
		log.Fatal("ERRO AO ABRIR CONEXÃO COM MYSQL", err)
		panic(err)
	}
	defer db.Close()

	repository := database.NewOrderRepository(db)
	usecaseCaller := usecase.NewCalculateFinalPriceUseCase(repository)

	input := usecase.OrderInputDTO{
		ID:    "1234",
		Price: 100,
		Tax:   10,
	}

	output, err := usecaseCaller.Execute(input)
	if err != nil {
		log.Fatal("ERRO AO EXECUTAR USE CASE", err)
		panic(err)
	}

	log.Println("FINAL PRICE: ", output.FinalPrice)

}
