package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/shubik22/robinhood-bot/bot"
)

func main() {
	godotenv.Load()
	bot := bot.NewBot()
	go bot.PeriodicallyTweetBalance(1 * time.Hour)
	go bot.ListenForTweets()
	bindToPort()
}

func bindToPort() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	log.Printf("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}
