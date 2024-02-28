package main

import (
	"fmt"
	"telegra/internal/dispatcher"
	"time"
)

func main() {
	// получаем время начала выполнения программы
	start := time.Now()
	// создаем диспетчера
	dispatcher := dispatcher.NewDispatcher(10)
	// запускаем диспетчера
	dispatcher.Run()
	// ищем слово в базе
	dispatcher.Search("тест")
	// получаем время окончания выполнения программы
	end := time.Now()
	// выводим время выполнения программы в часы, минуты и секунды
	fmt.Printf("Время выполнения программы: %v\n", end.Sub(start))

}
