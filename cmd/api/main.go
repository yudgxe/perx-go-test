package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/yudgxe/perx-go-test/internal/handlers"
	"github.com/yudgxe/perx-go-test/internal/middleware"
	"github.com/yudgxe/perx-go-test/internal/queue"
	"github.com/yudgxe/perx-go-test/internal/store"
	"github.com/yudgxe/perx-go-test/internal/worker"
)

const (
	workerQuantityDefault = 2
	queueSizeDefault      = 10
)

var (
	addr           string
	listenPrefix   string
	queueSize      int
	workerQuantity int
)

func init() {
	listenPrefix = os.Getenv("LISTEN_PREFIX")
	if listenPrefix == "" {
		listenPrefix = "/api"
	}
}

func init() {
	addr = os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	flag.StringVar(&addr, "addr", addr, "Адрес для запускаемого сервера")

	getEnvIntDefault := func(key string, def int) (int, error) {
		if env := os.Getenv(key); env != "" {
			var err error
			value, err := strconv.Atoi(env)
			if err != nil {
				return 0, fmt.Errorf("Ошибка c переменной %s: %w", key, err)
			}
			if value < 0 {
				return 0, fmt.Errorf("Ошибка c переменной %s: не может быть меньше нуля", key)
			}
			return value, nil
		} else {
			return def, nil
		}
	}

	env, err := getEnvIntDefault("QUEUE_SIZE", queueSizeDefault)
	if err != nil {
		log.Fatal(err)
	}
	queueSize = env
	flag.IntVar(&queueSize, "s", queueSize, "Размер очереди")

	env, err = getEnvIntDefault("QUEUE_SIZE", queueSizeDefault)
	if err != nil {
		log.Fatal(err)
	}
	workerQuantity = env
	flag.IntVar(&workerQuantity, "q", workerQuantityDefault, "Количество воркеров")

	flag.Parse()
}

func main() {
	queue := queue.New(queueSize)
	store := store.NewTaskStore()
	worker := worker.New(queue, store)

	// Запуск воркеров
	for i := 0; i < workerQuantity; i++ {
		worker.StartWorker()
	}

	handler := http.NewServeMux()
	handler.Handle("/", middleware.Logging(handlers.NewHandler(queue, store)))
	server := &http.Server{
		Addr:    addr,
		Handler: http.StripPrefix(listenPrefix, handler),
	}
	waitconns := make(chan struct{})
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-sigchan
		log.Println("Завершение работы сервера")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Shutdown(): %v", err)
		}
		close(waitconns)
	}()
	log.Printf("Запуск сервера на %s%s", addr, listenPrefix)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("ListenAndServe(): %v", err)
	}
	<-waitconns
}
