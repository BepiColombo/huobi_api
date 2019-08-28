/**
 * @time 2019-08-27 18:00
 * @author jarvis4901
 * @description
 */
package history

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"websocket_go/model"
)

/**
 * Description:币种历史行情
 * @Author: jarvis
 * @Date: 2019-08-28 14:13
 */
func CoinHistory(c *gin.Context) {

	var (
		err error
	)

	symbol := c.Query("symbol")
	period := c.Query("period")
	sizeQuery := c.Query("size")
	size, err := strconv.Atoi(sizeQuery)
	if err != nil {
		c.String(http.StatusOK, "参数错误")
		return
	}
	data, err := getHistory(symbol, period, size)
	fmt.Println(string(data))
	if err != nil {
		c.String(http.StatusOK, "接口出错")
		return
	}
	res := &model.HistoryResponse{}
	err = json.Unmarshal(data, res)
	if err != nil {
		c.String(http.StatusOK, "解析数据出错")
		return
	}
	c.JSON(http.StatusOK, res)

}

/**
 * Description:调取火币history 接口
 * @Param: symbol 币种符号; period 周期 1day,1year,1month; size 获取数据的条目
 * @Return:
 * @Author: jarvis
 * @Date: 2019-08-28 14:14
 */
func getHistory(symbol string, period string, size int) (data []byte, err error) {
	//构造url
	u, err := url.Parse("https://api.huobi.pro/market/history/kline?")
	if err != nil {
		fmt.Println("url parse fail")
		return nil, err
	}
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("period", period)
	q.Set("size", strconv.Itoa(size))
	u.RawQuery = q.Encode()
	//发起get请求
	resp, err1 := http.Get(u.String())
	if err1 != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("get fail:", err1)
		return nil, err
	}
	defer resp.Body.Close()
	//读取响应体
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("read body fail")
		return nil, err
	}
	//fmt.Println(string(body))
	return body, nil
}

/**
 * Description: 获取usdt_cny汇率
 * @Author: jarvis
 * @Date: 2019-08-28 14:45
 */
func Rate(c *gin.Context) {

	rate_url := "https://www.huobi.br.com/-/x/general/exchange_rate/list"
	resp, err := http.Get(rate_url)

	if err != nil {
		//fmt.Printf("请求实时币种汇率出错: %v", err)
		c.String(http.StatusOK, "请求出错")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		c.String(http.StatusOK, "请求出错")
	}

	result := &model.ExChangeRateStruct{}
	rate := &model.RateDataStruct{}

	json.Unmarshal(bodyBytes, &result)

	if result.Code == 200 {
		for _, v := range result.Data {
			//fmt.Printf("%s <===> %g\n",v.Name,v.Rate)
			if v.Name == "usdt_cny" {
				//fmt.Printf("交易对 %s 汇率：%g",symbol,v.Rate)
				//rate = v.Rate
				rate = v
				break
			}
		}
	} else {
		c.String(http.StatusOK, "请求出错")
	}
	c.JSON(http.StatusOK, rate)
}
