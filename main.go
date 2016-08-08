package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/shubik22/robinhood-bot/bot"
)

func main() {
	godotenv.Load()
	bot := bot.NewBot()
	go bot.PeriodicallyTweetBalance(1 * time.Hour)
	bot.ListenForTweets()
}
