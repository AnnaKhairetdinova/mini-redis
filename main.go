package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	err := server.Start(*port, s)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("новый клиент")
}
