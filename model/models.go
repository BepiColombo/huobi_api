/**
 * @time 2019-08-28 09:20
 * @author jarvis4901
 * @description
 */
package model

type Message struct {
	Ping  int64  `json:ping`
	Topic string `json:topic`
}

type HistoryResponse struct {
	Status string `json:"status"`
	Ch string `json:"ch"`
	Ts int64 `json:"ts"`
	Data []*HistoryItem `json:"data"`
}

type HistoryItem struct {
	Id int64 `json:"id"`
	Open float64 `json:"open"`
	Close float64 `json:"close"`
	Low float64 `json:"low"`
	Hign float64 `json:"hign"`
	Amount float64 `json:"amount"`
	Vol float64 `json:"vol"`
	Count int64 `json:"count"`
}


type ExChangeRateStruct struct {
	Code int               `json:"code"`
	Data []*RateDataStruct `json:"data"`
}


type RateDataStruct struct {
	Data_Time uint64  `json:"data_time"`
	Name      string  `json:"name"`
	Rate      float64 `json:"rate"`
	Time      uint64  `json:"time"`
}