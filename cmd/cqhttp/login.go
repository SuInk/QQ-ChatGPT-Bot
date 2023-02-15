package cqhttp

import (
	"QQ-ChatGPT-Bot/config"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type Bot struct {
	// 机器人QQ号
	QQ   int64
	Conn *websocket.Conn
	MQ   chan *RcvMsg //收到的消息队列
	MX   sync.Mutex   //消息队列锁
}

var bot = &Bot{
	MQ: make(chan *RcvMsg, 50),
}

func Run() {

	for i := 1; ; i++ {
		log.Printf("第%d次尝试连接%s中...\n", i, config.Cfg.CqHttp.WebSocket)
		var err error
		bot.Conn, _, err = websocket.DefaultDialer.Dial(config.Cfg.CqHttp.WebSocket, nil)
		if err != nil {
			log.Printf("连接失败, 5秒后重试:%v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("CqHttp连接成功!")
		go bot.Read(bot.Conn)
		break
	}
}
func (bot *Bot) Read(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		var rcvMsg RcvMsg
		err = json.Unmarshal(msg, &rcvMsg)
		if err != nil {
			log.Println(err)
		}
		//处理收到的消息
		if rcvMsg.PostType == "message" {
			// 消息预处理Parser
			isAt, err := regexp.MatchString(`CQ:at,qq=`+strconv.FormatInt(rcvMsg.SelfId, 10), rcvMsg.RawMessage)
			if err != nil {
				log.Println(err)
			}
			// 去除消息CQ码
			rcvMsg.Message = regexp.MustCompile(`\[CQ:.*?]`).ReplaceAllString(rcvMsg.Message, "")
			if rcvMsg.Message != " " && rcvMsg.Message != "" && rcvMsg.Message != "   " {
				go bot.HandleMsg(isAt, rcvMsg)
			}
		}
	}
}
