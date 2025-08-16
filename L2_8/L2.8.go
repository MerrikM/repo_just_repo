package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	const ntpServer = "pool.ntp.org"

	currentTime, err := ntp.Time(ntpServer)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Printf("Ошибка получения времени с NTP: %v", err)
		os.Exit(1)
	}

	fmt.Println("Текущее время:", currentTime.Format(time.RFC1123))
}
