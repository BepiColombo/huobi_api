/**
 * @time 2019-08-26 14:09
 * @author jarvis4901
 * @description
 */
package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leizongmin/huobiapi"
	"sync"
	"websocket_go/utils"
)


/**
 * Description: K线socket接口，用于币种详情页面的行情折线图的实现
 * @Author: jarvis
 * @Date: 2019-08-28 14:16
 */
func Kline(c *gin.Context) {
	var (
		isSubscriped    bool
		marketCloseOnce sync.Once
	)
	ws, err := utils.InitWebsocket(c)
	if err != nil {
		fmt.Println("创建socket连接失败")
	}
	market, err := huobiapi.NewMarket()
	if err != nil {
		fmt.Println("创建火币socket失败")
	}
	isSubscriped = false

	defer func() {
		//关闭操作
		//fmt.Println("关闭操作")
		marketCloseOnce.Do(func() {
			market.Close()
		})
	}()

	for !ws.Closed {

		mType, message, err := ws.ReadMsg()
		if err != nil {
			ws.Close()
			break
		}
		if mType == "data" {
			if err != nil {
				fmt.Println("参数解析错误", message)
			}
			if !isSubscriped {
				//fmt.Println("订阅")
				market.Subscribe(message.(string), func(topic string, json *huobiapi.JSON) {
					// 收到数据更新时回调
					isSubscriped = true
					fmt.Println(json)
					jsonByte, _ := json.MarshalJSON()
					err = ws.WriteMsg(1, jsonByte)
					if err != nil {
						ws.Close()
					}
				})
				//market.Close()
			}
			fmt.Println(message)
		}
	}

}


/**
 * Description:市场概要socket接口，提供24小时内最新市场概要。
 * @Author: jarvis
 * @Date: 2019-08-28 14:15
 */
func MarketDetail(c *gin.Context) {
	var (
		marketCloseOnce sync.Once
	)

	topics := [5]string{"market.btcusdt.detail", "market.ethusdt.detail","market.eosusdt.detail", "market.ltcusdt.detail","market.xrpusdt.detail"}
	ws, err := utils.InitWebsocket(c)
	if err != nil {
		fmt.Println("创建socket连接失败")
	}

	market, err := huobiapi.NewMarket()
	if err != nil {
		fmt.Println("创建火币socket失败")
	}

	defer func() {
		//关闭操作
		//fmt.Println("关闭操作")
		marketCloseOnce.Do(func() {
			market.Close()
		})
	}()

	for _, v := range topics {
		// 订阅主题
		//fmt.Println("订阅消息事件...")
		market.Subscribe(v, func(topic string, json *huobiapi.JSON) {
			// 收到数据更新时回调
			//fmt.Println(json)
			jsonByte, _ := json.MarshalJSON()
			err = ws.WriteMsg(1, jsonByte)
			if err != nil {
				ws.Close()
			}
		})
	}
	for !ws.Closed {
		market.Loop()
	}
}