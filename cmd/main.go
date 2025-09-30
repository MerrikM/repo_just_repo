package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"repo_just_repo/internal/config"
	"repo_just_repo/internal/handler"
	"repo_just_repo/internal/middleware"
	"repo_just_repo/internal/repository"
	"repo_just_repo/internal/service"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig("app/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	db, err := config.SetupDatabase(cfg.DatabaseConfig.DSN)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии БД: %v", err)
		}
	}()

	calendarRepo := repository.NewCalendarRepository(db)
	calendarService := service.NewCalendarService(calendarRepo)
	eventHandler := handler.NewEventHandler(calendarService)

	srv, router := config.SetupServer(cfg.ServerAddr)
	setupEventRouter(eventHandler, router)
	runServer(ctx, srv)
}

func setupEventRouter(h *handler.EventHandler, r *chi.Mux) *chi.Mux {
	r.Use(middleware.Logging)

	r.Post("/create_event", h.CreateEvent)
	r.Put("/update_event", h.UpdateEvent)
	r.Delete("/delete_event", h.DeleteEvent)
	r.Get("/events_for_day", h.EventsForDay)
	r.Get("/events_for_week", h.EventsForWeek)
	r.Get("/events_for_month", h.EventsForMonth)

	return r
}

func runServer(ctx context.Context, server *http.Server) {
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("сервер запущен на " + server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil {
			log.Fatalf("ошибка работы сервера: %v", err)
		}
	case sig := <-signalChannel:
		log.Printf("получен сигнал %v остановки работы сервера ", sig)
	}

	shutDownCtx, shutDownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutDownCancel()

	if err := server.Shutdown(shutDownCtx); err != nil {
		log.Printf("ошибка при остановке сервера: %v", err)
	} else {
		log.Println("Сервер успешно остановлен")
	}
}
