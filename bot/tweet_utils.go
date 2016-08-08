package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/leekchan/accounting"
	"github.com/shubik22/robinhood-client"
)

func parseTweet(text string) (*TradeInputs, error) {
	text = strings.Split(text, mentionString)[1]

	words := strings.Split(text, " ")
	if len(words) != 3 {
		return nil, fmt.Errorf("invalid order format")
	}

	var ti TradeInputs
	orderType := strings.ToLower(words[0])
	if orderType == "buy" || orderType == "sell" {
		ti.OrderType = orderType
	} else {
		return nil, fmt.Errorf("invalid order type (need buy or sell)")
	}

	quantity, err := strconv.Atoi(words[1])
	if err != nil {
		return nil, fmt.Errorf("invalid quantity (need integer)")
	} else {
		ti.Quantity = quantity
	}

	ti.Symbol = words[2]

	return &ti, nil
}

func shuffle(a []robinhood.SimplePosition) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
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
	shuffle(u.Positions)
	for _, p := range u.Positions {
		pStrings = append(pStrings, fmt.Sprintf("%v of %v", pluralize(int(p.Quantity), "share"), p.Symbol))
	}
	tweetStr += strings.Join(pStrings, ", ")
	tweetStr += "."

	return tweetStr
}

func createTweetsFromText(text string) []string {
	var tweets []string
	words := strings.Split(text, " ")
	currentTweet := GetBalancePhrase()
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

func pluralize(num int, word string) string {
	if num == 1 {
		return fmt.Sprintf("%v %v", num, word)
	} else {
		return fmt.Sprintf("%v %vs", num, word)
	}
}
