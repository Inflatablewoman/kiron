package main

import (
	"github.com/Inflatablewoman/kiron/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"flag"
	"fmt"

	"github.com/rcrowley/go-tigertonic"
)

var (
	host = flag.String("host", "localhost", "Host Address")
	port = flag.Int("port", 1979, "The Post")
)

func main() {

	flag.Parse()

	// Setup logging
	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf("Starting Kiron Service - Port: %s", fmt.Sprintf(":%d", *port))

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
	server := tigertonic.NewServer(fmt.Sprintf(":%d", *port), aMux)
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
