package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	runMapWriterSafeWithMutex()

	runMapWriterSafeWithSyncMap()
}

// Функция для запуска безопасной записи с использованием мьютекса
func runMapWriterSafeWithMutex() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mapa := make(map[string]int)

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				mapWriterSafeWithMutex(mapa, fmt.Sprintf("apple%d", i), i, &mutex)
			}
		}()
	}

	wg.Wait()
	fmt.Println(mapa)
}

// Функция для запуска безопасной записи с SyncMap
func runMapWriterSafeWithSyncMap() {
	var mapa sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mapWriterSafeWithSyncMap(&mapa, fmt.Sprintf("apple%d", i), i)
		}(i)
	}

	wg.Wait()

	mapa.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})
}

// Базовая запись с блокировкой мьютекса
func mapWriterSafeWithMutex(mapa map[string]int, key string, value int, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	mapa[key] = value
}

// Запись с использованием SyncMap
func mapWriterSafeWithSyncMap(mapa *sync.Map, key string, value int) {
	mapa.Store(key, value)
}
