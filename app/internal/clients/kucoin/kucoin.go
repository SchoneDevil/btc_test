package kucoin

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type IKucoin struct {
}
type Response struct {
	Code string `json:"code"`
	Data struct {
		Symbol string `json:"symbol"`
		Buy    string `json:"buy"`
		Sell   string `json:"sell"`
	} `json:"data"`
}

//Источник для пары BTC-USDT: <https://api.kucoin.com/api/v1/market/stats?symbol=BTC-USDT>
func (ik IKucoin) GetBtcUsdt() Response {
	resp, err := http.Get("https://api.kucoin.com/api/v1/market/stats?symbol=BTC-USDT")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var res Response
	_ = json.Unmarshal(body, &res)
	return res
}
