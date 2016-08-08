package bot

import (
	"math/rand"
)

func getBalancePhrases() []string {
	return []string{
		"U know I been tradin. ",
		"Takes money 2 make money. ",
		"How efficient is this market lol. ",
		"Can a bot ever have 2 much money?  I'm about to find out... ",
		"Watup @KimKardashian ",
	}
}

func getParseErrorPhrases() []string {
	return []string{
		"lot wut",
		"haha no ",
		"dad ",
	}
}

func getMarketClosedPhrases() []string {
	return []string{
		"wait until the market's open champ lol",
		"yea i'd trade that... if the market was open u goof",
		"dawg the market's a lot like a coffeeshop only open 9:30 to 4 on weekdays... rn it's CLOSED",
	}
}

func getTradeErrorPhrases() []string {
	return []string{
		"This is a very complicated case Maude.  Your order did not go through.",
		"I tried to place your trade but it failed.  I failed.  At placing your trade.",
		"ha no sorry no trade :p",
		"listen bud, your trade didn't go through.  i'm not sure it's not your fault, tbh",
	}
}

func getRandomPhrase(phrases []string) string {
	return phrases[rand.Intn(len(phrases))]
}

func GetMarketClosedPhrase() string {
	return getRandomPhrase(getMarketClosedPhrases())
}

func GetBalancePhrase() string {
	return getRandomPhrase(getBalancePhrases())
}

func GetParseErrorPhrase() string {
	return getRandomPhrase(getParseErrorPhrases())
}

func GetTradeErrorPhrase() string {
	return getRandomPhrase(getTradeErrorPhrases())
}
