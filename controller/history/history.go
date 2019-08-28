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

//币种历史行情
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
	c.JSON(200, res)

}

//调取火币history 接口
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
