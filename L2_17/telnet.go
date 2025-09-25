package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Аргументы: host, port, timeout
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		fmt.Println("Connection error:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("Connected to %s\n", address)

	done := make(chan struct{})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nSIGINT received, closing connection...")
		conn.Close()
		close(done)
	}()

	// Горутина для чтения из сокета --> stdout
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Read error:", err)
		}
		fmt.Println("Connection closed by remote host")
		close(done)
	}()

	// Чтение stdin --> сокет
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err := fmt.Fprintln(conn, scanner.Text())
			if err != nil {
				fmt.Println("Write error:", err)
				break
			}
		}
		// Ctrl+D (EOF) — выходим
		fmt.Println("EOF received, closing connection...")
		conn.Close()
		close(done)
	}()

	<-done
}
