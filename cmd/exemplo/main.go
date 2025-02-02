package main

import "time"

func worker(workerId int, msg chan int) {
	for res := range msg {
		println("worker", workerId, "recebeu", res)
		time.Sleep(time.Second)
	}
}

// THREAD PRINCIPAL
func main() {
	canal := make(chan int)

	for i := 1; i <= 3; i++ {
		go worker(i, canal)

	}

	for i := 1; i <= 5; i++ {
		canal <- i
	}
}
