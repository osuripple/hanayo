package main

import (
	"net/http"
	"time"

	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

var rates = struct {
	BTC float64
	USD float64
}{
	// approx values
	709,
	1.06050,
}

func getRates(c *gin.Context) {
	c.JSON(200, rates)
}

func init() {
	go func() {
		for {
			updateRates()
			time.Sleep(time.Hour)
		}
	}()
}

func updateRates() {
	resp, err := http.Get("https://www.bitstamp.net/api/v2/ticker/btceur/")
	if err == nil {
		exportRate(&rates.BTC, resp)
	}
	resp, err = http.Get("https://www.bitstamp.net/api/v2/ticker/eurusd/")
	if err == nil {
		exportRate(&rates.USD, resp)
	}
}

func exportRate(x *float64, resp *http.Response) {
	data, _ := ioutil.ReadAll(resp.Body)
	var last struct {
		Last float64 `json:"last"`
	}
	json.Unmarshal(data, &last)
	if last.Last != 0 {
		*x = last.Last
	}
}
