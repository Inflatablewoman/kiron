package main

import (
	"fmt"
	"github.com/Inflatablewoman/kiron/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rcrowley/go-tigertonic"
)

func main() {

	host := os.Getenv("KIRON_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("KIRON_PORT")
	if port == "" {
		port = "80"
	}

	// Setup logging
	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf("Starting Kiron Service - Port: %s", fmt.Sprintf("%s:%s", host, port))

	err := server.InitDatabase()
	if err != nil {
		log.Fatalf("Unable to connect to database %v", err)
	}

	// Create handlers
	mux := tigertonic.NewTrieServeMux()
	server.RegisterHTTPHandlers(mux)

	// Log apache style
	aMux := tigertonic.ApacheLogged(mux)

	// Create server and listen to requests
	server := tigertonic.NewServer(fmt.Sprintf("%s:%s", host, port), aMux)
	// server.Close to stop gracefully.
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Kiron Service Fatal: %v", err)
		}
	}()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	log.Println(<-ch)
	server.Close()
}
