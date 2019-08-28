/**
 * @time 2019-08-27 09:46
 * @author jarvis4901
 * @description
 */
package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
	"websocket_go/model"
)

//升级http协议
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


type WsCoon struct {
	conn      *websocket.Conn
	lastPing  int64 // 上次接收到的ping时间戳
	closeOnce sync.Once //once 确保只关闭连接一次
	//WsChan chan int
	Closed bool   //标志是否已关闭
}

/**
 * Description:初始化websocket连接
 * @Return: *WsCoon ; err
 * @Author: jarvis
 * @Date: 2019-08-27 13:24
 */
func InitWebsocket(c *gin.Context) (wsCoon *WsCoon, err error) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	now_time := getUinxMillisecond()
	wsCoon = &WsCoon{
		conn:      ws,
		lastPing:  now_time,
		closeOnce: sync.Once{},
		//WsChan: make(chan int, 1),
		Closed: false,
	}
	go func() {
		t := time.NewTicker(1 * time.Second)
		for _ = range t.C {
			active := wsCoon.HeartBeatTest();
			if !active {
				wsCoon.conn.Close()
				break
			}
		}
	}()
	return wsCoon, err
}

/**
 * Description: 关闭连接
 * @Param: *WsConn
 * @Author: jarvis
 * @Date: 2019-08-27 14:20
 */
func (w *WsCoon) Close() {
	//once确保只关闭一次
	w.closeOnce.Do(func() {
		//fmt.Println("主动关闭连接")
		w.Closed = true
		w.conn.Close()
	})
}

/**
 * Description: 读取消息
 * @Param: *WsConn
 * @Return: mType 数据类型； msg 数据内容；err 错误
 * @Author: jarvis
 * @Date: 2019-08-27 14:27
 */
func (w *WsCoon) ReadMsg() (mType interface{}, msg interface{}, err error) {
	_, message, err := w.conn.ReadMessage()
	if err != nil {
		return nil, nil, err
	}
	_message := &model.Message{}
	json_err := json.Unmarshal(message, &_message)
	if json_err != nil {
		return nil, nil, json_err
	}
	if _message.Ping != 0 {
		//ping消息
		//将lastPing更新为新传来的时间戳
		//fmt.Println("ping消息",_message.Ping)
		w.lastPing = _message.Ping
		return "ping", _message.Ping, nil
	} else if _message.Topic != "" {
		//订阅消息
		return "data", _message.Topic, nil
	}
	return nil, nil, nil
}


/**
 * Description:发送消息
 * @Param: *WsConn
 * @Return: err 错误
 * @Author: jarvis
 * @Date: 2019-08-27 14:41
 */
func (w *WsCoon) WriteMsg(messageType int, data []byte) (err error) {
	err = w.conn.WriteMessage(messageType, data)
	return err
}

/**
 * Description:心跳时间检测：检测连接是否长时间无响应
   超过六秒未收到来自客户端的ping消息则主动断开与客户端的连接
 * @Param: *WsConn
 * @Return: pass 是否透过检测
 * @Author: jarvis
 * @Date: 2019-08-27 15:30
 */
func (w *WsCoon) HeartBeatTest() (pass bool) {
	var t = time.Now().Unix()
	// 检查上次ping时间，如果超过6秒无响应，返回false
	tr := t*1000 - w.lastPing
	if (tr > 60000) {
		//fmt.Println("lastPing：------", w.lastPing)
		//fmt.Println("当前时间:-----", t*1000)
		return false
	}
	return true
}
