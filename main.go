package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
	"github.com/leekchan/accounting"
	"github.com/shubik22/robinhood-client"
)

func main() {
	loadEnv()
	twitterApi := getTwitterApi()
	rhClient := getRobinhoodClient()

	ticker := time.NewTicker(1 * time.Hour)
	for {
		tweetBalance(twitterApi, rhClient)
		<-ticker.C
	}
}

func tweetBalance(twitterApi *anaconda.TwitterApi, rhClient *robinhood.Client) {
	u := getRobinhoodData(rhClient)
	tweetText := createBalancesText(u)
	tweets := createTweetsFromText(tweetText)

	for _, tweetStr := range tweets {
		tweet, err := twitterApi.PostTweet(tweetStr, nil)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Just tweeted: ", tweet.Text)
		}
	}
}

func getTwitterApi() *anaconda.TwitterApi {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessSecret)

	return api
}

func getRobinhoodClient() *robinhood.Client {
	rhUsername := os.Getenv("ROBINHOOD_USERNAME")
	rhPassword := os.Getenv("ROBINHOOD_PW")

	client := robinhood.NewClient(rhUsername, rhPassword)
	return client
}

func createBalancesText(u *robinhood.User) string {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	var tweetStr string
	tweetStr += "The total value of my account is "
	tweetStr += ac.FormatMoney(u.TotalBalance)
	tweetStr += ".  "

	tweetStr += "Cash: "
	tweetStr += ac.FormatMoney(u.CashBalance)
	tweetStr += ", "

	var pStrings []string
	for _, p := range u.Positions {
		pStrings = append(pStrings, fmt.Sprintf("%v shares of %v", p.Quantity, p.Symbol))
	}
	tweetStr += strings.Join(pStrings, ", ")
	tweetStr += "."

	return tweetStr
}

func createTweetsFromText(text string) []string {
	var tweets []string
	words := strings.Split(text, " ")
	currentTweet := getPhrase()
	for _, word := range words {
		if len(currentTweet) > 132 {
			tweets = append(tweets, currentTweet)
			currentTweet = ""
		}
		currentTweet += word
		currentTweet += " "
	}
	tweets = append(tweets, currentTweet)

	numTweets := len(tweets)
	if numTweets > 1 {
		for idx, tweet := range tweets {
			tweet = fmt.Sprintf("%v(%v/%v)", tweet, idx+1, numTweets)
			tweets[idx] = tweet
		}
	}

	return tweets
}

func getPhrase() string {
	phrases := [...]string{
		"U know I been tradin. ",
		"Takes money 2 make money. ",
		"How efficient is this market lol. ",
		"Can a bot ever have 2 much money?  I'm about to find out... ",
		"Watup @KimKardashian ",
	}
	return phrases[rand.Intn(len(phrases))]
}

func getRobinhoodData(c *robinhood.Client) *robinhood.User {
	username := c.UserName

	u := &robinhood.User{
		Username: username,
	}
	ar, _, err := c.Accounts.ListAccounts()
	if err != nil {
		panic(err)
	}

	account := ar.Results[0]
	pr, _, err := c.Positions.ListPositions()
	if err != nil {
		panic(err)
	}

	cashBalance, _ := strconv.ParseFloat(account.Cash, 64)
	u.CashBalance = cashBalance

	var positionBalance float64
	for _, p := range pr.Results {
		quantity, _ := strconv.ParseFloat(p.Quantity, 64)
		if quantity == 0 {
			continue
		}

		q, _, err := c.Quotes.GetQuote(&p)
		if err != nil {
			panic(err)
		}

		lastPrice, _ := strconv.ParseFloat(q.LastTradePrice, 64)
		avgBuyPrice, _ := strconv.ParseFloat(p.AverageBuyPrice, 64)
		simplePosition := robinhood.SimplePosition{
			PurchaseTime:    p.CreatedAt,
			Quantity:        quantity,
			Symbol:          q.Symbol,
			AverageBuyPrice: avgBuyPrice,
			LastTradePrice:  lastPrice,
		}
		u.Positions = append(u.Positions, simplePosition)
		positionBalance += (quantity * lastPrice)
	}

	totalBalance := cashBalance + positionBalance
	u.PositionBalance = positionBalance
	u.TotalBalance = totalBalance

	return u
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
