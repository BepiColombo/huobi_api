/**
 * @time 2019-08-26 14:01
 * @author jarvis4901
 * @description
 */
package router

import (
	"github.com/gin-gonic/gin"
	"websocket_go/controller/history"
	"websocket_go/controller/websocket"
	"websocket_go/middleware"
)

func Init() {
	router := gin.Default()
	// CrossDomain跨域处理，options请求处理
	router.Use(middleware.CrossDomain())
	// v1群组
	v1 := router.Group("/v1")
	{
		v1.GET("/kline", websocket.Kline)
		v1.GET("/marketDetail", websocket.MarketDetail)
		v1.GET("/coinHistory", history.CoinHistory)
	}
	router.Run(":8000")
}
