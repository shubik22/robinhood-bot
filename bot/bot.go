package bot

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/shubik22/robinhood-client"
)

const (
	mentionString = "@TWEETMETRADES "
)

type Bot struct {
	twitterApi *anaconda.TwitterApi
	rhClient   *robinhood.Client
}

type TradeInputs struct {
	Symbol    string
	OrderType string
	Quantity  int
}

func NewBot() *Bot {
	twitterApi := getTwitterApi()
	rhClient := getRobinhoodClient()
	bot := &Bot{
		twitterApi: twitterApi,
		rhClient:   rhClient,
	}
	return bot
}

func (b *Bot) handleMention(t *anaconda.Tweet) {
	text := strings.ToUpper(t.Text)
	if !strings.Contains(text, mentionString) {
		return
	}
	params := url.Values{}
	params.Add("in_reply_to_status_id", t.IdStr)

	ti, err := parseTweet(text)
	if err != nil {
		tweetStr := fmt.Sprintf(".@%v %v... %v", t.User.ScreenName, GetParseErrorPhrase(), err.Error())
		b.twitterApi.PostTweet(tweetStr, params)
		return
	}

	if !MarketIsOpen() {
		tweetStr := fmt.Sprintf(".@%v %v", t.User.ScreenName, GetMarketClosedPhrase())
		b.twitterApi.PostTweet(tweetStr, params)
		return
	}

	or, _, err := b.rhClient.Trades.PlaceTrade(ti.Symbol, ti.OrderType, ti.Quantity)
	var tweetStr string
	if err != nil {
		log.Printf("Error placing %v trade for %v %v: %+v", ti.OrderType, ti.Quantity, ti.Symbol, err)
		tweetStr = fmt.Sprintf(".@%v %v", t.User.ScreenName, GetTradeErrorPhrase())
		b.twitterApi.PostTweet(tweetStr, params)
	} else {
		log.Printf("Response for trade: %+v", or)
		tweetStr := fmt.Sprintf(
			".@%v O yeaa placed a %v order for %v of %v",
			t.User.ScreenName,
			ti.OrderType,
			pluralize(ti.Quantity, "share"),
			ti.Symbol,
		)
		b.twitterApi.PostTweet(tweetStr, params)
		b.tweetBalanceWithDelay(30 * time.Second)
	}
}

func (b *Bot) getRobinhoodData() *robinhood.User {
	username := b.rhClient.UserName

	u := &robinhood.User{
		Username: username,
	}
	ar, _, err := b.rhClient.Accounts.ListAccounts()
	if err != nil {
		panic(err)
	}

	account := ar.Results[0]
	pr, _, err := b.rhClient.Positions.ListPositions()
	if err != nil {
		panic(err)
	}

	cashBalance, _ := strconv.ParseFloat(account.BuyingPower, 64)
	u.CashBalance = cashBalance

	var positionBalance float64
	for _, p := range pr.Results {
		quantity, _ := strconv.ParseFloat(p.Quantity, 64)
		if quantity == 0 {
			continue
		}

		q, _, err := b.rhClient.Quotes.GetQuote(&p)
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

func (b *Bot) tweetBalance() {
	u := b.getRobinhoodData()
	tweetText := createBalancesText(u)
	tweets := createTweetsFromText(tweetText)

	for _, tweetStr := range tweets {
		tweet, err := b.twitterApi.PostTweet(tweetStr, nil)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Just tweeted: ", tweet.Text)
		}
	}
}

func (b *Bot) tweetBalanceWithDelay(d time.Duration) {
	ticker := time.NewTicker(d)
	<-ticker.C
	b.tweetBalance()
}

func (b *Bot) ListenForTweets() {
	stream := b.twitterApi.UserStream(nil)
	for {
		o := <-stream.C
		t, ok := o.(anaconda.Tweet)
		if !ok {
			log.Println("Received non-tweet event")
		} else {
			log.Println("Received tweet")
			b.handleMention(&t)
		}
	}
}

func (b *Bot) PeriodicallyTweetBalance(d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		b.tweetBalance()
		<-ticker.C
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
