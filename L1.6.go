package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func stopByCondition() {
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("шаг: ", i)
			time.Sleep(time.Second)
		}
		fmt.Println("выходим по условию")
	}()
}

// Если горутина читает из канала, она завершится, когда канал закроется
// Можно юзать при работе с потоками данных, когда программа должна завершаться только тогда,
// когда входящие данные заканчиваются
func stopByClosingChannel() {
	ch := make(chan string)

	go func() {
		for msg := range ch {
			fmt.Println(msg)
		}
		fmt.Println("остановка по закрытию канала")
	}()

	ch <- "Hello"
	close(ch) // после этого горутина выйдет
}

// Отдельный канал только для сигнализации о завершении
// Можно юзать, когда несколько горутин должны завершиться по одному сигналу
func stopByDoneChannel() {
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("остановка по done каналу")
				return
			default:
				fmt.Println("работаем")
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	time.Sleep(1 * time.Second)
	close(done)
}

// Универсальный способ
// Можно юзать при работе с запросами, цепочками горутин, при graceful shutdown
func stopByContext() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("остановка по контексту (переделаи Done())", ctx.Err())
				return
			default:
				fmt.Println("в процессе")
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()
}

// Горутина завершается по истечении времени.
func stopByTimer() {
	timeout := time.After(1 * time.Second)

	go func() {
		for {
			select {
			case <-timeout:
				fmt.Println("остановка по таймеру")
				return
			default:
				fmt.Println("тик")
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()
}

// Завершает текущую горутину немедленно, выполняя все defer
func stopByGoexit() {
	go func() {
		defer fmt.Println("это выполнится до Goexit")
		fmt.Println("остановили при помощи Goexit")
		runtime.Goexit()
		fmt.Println("это никогда не напечатается")
	}()
}

// return из горутины
// Если горутина — просто анонимная функция, можно завершить её возвратом
func stopByReturn() {
	go func() {
		fmt.Println("работаем")
		return
	}()
}

// Паника и перехват
func stopByPanic() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("восстановились из паники", r)
			}
		}()
		panic("паника! аварийно завершаемся")
	}()
}
