package utils

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

const (
	marketUrl = "https://api.tradeking.com/v1/market/clock.xml"
)

type Response struct {
	Next        string `xml:"status>next"`
	Message     string `xml:"message"`
	Error       string `xml:"error"`
	Id          string `xml:"id,attr"`
	Elapsedtime string `xml:"elapsedtime"`
	Date        string `xml:"date"`
	Change_at   string `xml:"status>change_at"`
	Current     string `xml:"status>current"`
	Unixtime    string `xml:"unixtime"`
}

func MarketIsOpen() bool {
	resp, err := http.Get(marketUrl)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	var r Response
	xml.Unmarshal(data, &r)
	return r.Current != "close"
}
