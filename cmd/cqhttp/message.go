package cqhttp

import (
	"QQ-ChatGPT-Bot/cmd/chatgpt"
	"QQ-ChatGPT-Bot/config"
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"time"
)

type RcvMsg struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	Time        int64  `json:"time"`
	SelfId      int64  `json:"self_id"`
	SubType     string `json:"sub_type"`
	UserId      int64  `json:"user_id"`
	TargetId    int64  `json:"target_id"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
	Font        int    `json:"font"`
	Sender      struct {
		Age      int    `json:"age"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		UserId   int64  `json:"user_id"`
	}
	GroupId   int64 `json:"group_id"`
	MessageId int64 `json:"message_id"`
}
type SendMsg struct {
	Action string `json:"action"`
	Params struct {
		UserId  int64  `json:"user_id"`
		GroupId int64  `json:"group_id"`
		Message string `json:"message"`
	} `json:"params"`
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //初始化日志格式
}

// HandleMsg 对CqHttp发送的json进行处理
func (bot *Bot) HandleMsg(rawData []byte) {
	var rcvMsg RcvMsg
	err := json.Unmarshal(rawData, &rcvMsg)
	if err != nil {
		log.Println(err)
	}
	// 准备处理消息
	isAt, err := regexp.MatchString(`CQ:at,qq=`+strconv.FormatInt(rcvMsg.SelfId, 10), rcvMsg.RawMessage)
	if err != nil {
		log.Println(err)
	}
	// 去除消息CQ码
	rcvMsg.Message = regexp.MustCompile(`\[CQ:.*?\]`).ReplaceAllString(rcvMsg.Message, "")
	switch rcvMsg.MessageType {
	case "private":
		go CheckTimeOut("private", rcvMsg.MessageId, rcvMsg.Sender.UserId)
		err := bot.SendPrivateMsg(rcvMsg.Sender.UserId, "[CQ:reply,id="+strconv.FormatInt(rcvMsg.MessageId, 10)+"]"+chatgpt.GenerateText(rcvMsg.Message))
		bot.MQ <- &rcvMsg
		if err != nil {
			log.Println(err)
		}
		//<-bot.MQ
	case "group":
		// 群消息@机器人才处理
		if !isAt && config.Cfg.CqHttp.AtOnly || rcvMsg.Sender.UserId == bot.QQ {
			return
		}
		go CheckTimeOut("group", rcvMsg.MessageId, rcvMsg.Sender.UserId)
		err := bot.SendGroupMsg(rcvMsg.GroupId, "[CQ:reply,id="+strconv.FormatInt(rcvMsg.MessageId, 10)+"]"+chatgpt.GenerateText(rcvMsg.Message))
		bot.MQ <- &rcvMsg
		if err != nil {
			log.Println(err)
		}
	}

}
func CheckTimeOut(msgType string, msgId int64, userId int64) {
	// 检查超时
	select {
	case <-bot.MQ:
	case <-time.After(time.Second * time.Duration(config.Cfg.CqHttp.TimeOut)):
		switch msgType {
		case "private":
			err := bot.SendPrivateMsg(userId, "[CQ:reply,id="+strconv.FormatInt(msgId, 10)+"]"+"思考中，请耐心等待~")
			if err != nil {
				log.Println(err)
			}
		case "group":
			err := bot.SendGroupMsg(userId, "思考中，请耐心等待~")
			if err != nil {
				log.Println(err)
			}
		}
	}
}
func (bot *Bot) SendPrivateMsg(userId int64, text string) error {
	sendMsg := SendMsg{
		Action: "send_private_msg",
		Params: struct {
			UserId  int64  `json:"user_id"`
			GroupId int64  `json:"group_id"`
			Message string `json:"message"`
		}{UserId: userId, Message: text},
	}

	rawMsg, _ := json.Marshal(sendMsg)
	err := bot.Conn.WriteMessage(1, rawMsg)
	if err != nil {
		return err
	}
	return nil
}
func (bot *Bot) SendGroupMsg(groupId int64, text string) error {
	sendMsg := SendMsg{
		Action: "send_group_msg",
		Params: struct {
			UserId  int64  `json:"user_id"`
			GroupId int64  `json:"group_id"`
			Message string `json:"message"`
		}{GroupId: groupId, Message: text},
	}

	rawMsg, _ := json.Marshal(sendMsg)
	err := bot.Conn.WriteMessage(1, rawMsg)
	if err != nil {
		return err
	}
	return nil
}
