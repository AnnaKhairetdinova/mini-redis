package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnnaKhairetdinova/mini-redis/server"
	"github.com/AnnaKhairetdinova/mini-redis/store"
)

func main() {
	port := flag.String("port", "6379", "Port to connect to mini-redis")
	if *port == "" {
		fmt.Println("Ошибка: флаг -port обязательный")
		flag.Usage()
		os.Exit(1)
	}

	s := store.New()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s.StartCleaner(ctx, 60*time.Second)

	err := server.Start(ctx, *port, s)
	if err != nil {
		fmt.Println(err)
	}
}
