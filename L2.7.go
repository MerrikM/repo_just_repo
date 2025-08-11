package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c { // 53 строка
		fmt.Print(v)
	}
}

/*

Функция asChan создаёт канал и запускает горутину, которая отправляет в этот канал переданные числа с задержкой до 1 секунды, затем закрывает канал

Функция merge принимает два канала a и b, и сливает их значения в один канал c.
Внутри функции, в бесконечном цикле с помощью select читаются значения из каналов a и b.
v, ok := <-a или b: пытаемся прочитать из канала
	Если канал еще открыт и из него успешно прочитано значение v, то тогда пишем это значение в канал c
	Иначе отписываемся от канала, делая его равным nil, чтобы select его игнорировал при следующей итерации внутри цикла на 53 строке
Как только оба канал равны nil, мы закрываем канал c и выходим из горутины.

Программа будет выводить числа от 1 до 8 в любом порядке по мере их поступления

*/
