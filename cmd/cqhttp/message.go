package cqhttp

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"QQ-ChatGPT-Bot/cmd/chatgpt"
	"QQ-ChatGPT-Bot/config"
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
func (bot *Bot) HandleMsg(isAt bool, rcvMsg RcvMsg) {
	// panic处理 暂时无法判断是否生效
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
			return
		}
	}()
	switch rcvMsg.MessageType {
	case "private":
		// 包含关键词才触发，运算符优先级 && > ||
		if config.Cfg.CqHttp.UseKeyword && !strings.Contains(rcvMsg.Message, config.Cfg.CqHttp.Keyword) && (config.Cfg.CqHttp.KeywordType == "all" || config.Cfg.CqHttp.KeywordType == "private") || rcvMsg.Sender.UserId == bot.QQ {
			return
		}
		rcvMsg.Message = strings.ReplaceAll(rcvMsg.Message, config.Cfg.CqHttp.Keyword, "")
		//输入“/clean”，清理缓存的历史记录（应该是吧？这块不是我写的，但是如果输入里有clear就执行这个clear命令就太破坏了，毕竟clear算常用单词）
		if strings.TrimSpace(rcvMsg.Message) == "/clean" {
			chatgpt.Cache.Clear(strconv.FormatInt(rcvMsg.Sender.UserId, 10))
			err := bot.SendPrivateMsg(rcvMsg.Sender.UserId, "历史记录清理完成")
			if err != nil {
				log.Println(err)
			}
			log.Println("历史记录清理完成")
			return
		}
		bot.MQ <- &rcvMsg
		msg, err := chatgpt.ChooseGenerateWay(strconv.FormatInt(rcvMsg.Sender.UserId, 10), rcvMsg.Message, config.Cfg.Context.PrivateContext)
		if msg != "" {
			err = bot.SendPrivateMsg(rcvMsg.Sender.UserId, "[CQ:reply,id="+strconv.FormatInt(rcvMsg.MessageId, 10)+"]"+msg)
		} else {
			err = bot.SendPrivateMsg(rcvMsg.Sender.UserId, "[CQ:reply,id="+strconv.FormatInt(rcvMsg.MessageId, 10)+"]"+"生成错误！错误信息:\n"+err.Error())
		}
		if err != nil {
			log.Println(err)
		}
		<-bot.MQ
	case "group":
		// 群消息@机器人才处理
		if !isAt && config.Cfg.CqHttp.AtOnly || rcvMsg.Sender.UserId == bot.QQ {
			return
		}
		// 检查是否有关键词
		if config.Cfg.CqHttp.UseKeyword && !strings.Contains(rcvMsg.Message, config.Cfg.CqHttp.Keyword) && (config.Cfg.CqHttp.KeywordType == "all" || config.Cfg.CqHttp.KeywordType == "group") {
			return
		}
		rcvMsg.Message = strings.ReplaceAll(rcvMsg.Message, config.Cfg.CqHttp.Keyword, "")
		if rcvMsg.Message == " /clean" {
			chatgpt.Cache.Clear(strconv.FormatInt(rcvMsg.GroupId, 10))
			err := bot.SendGroupMsg(rcvMsg.GroupId, "历史记录清理完成")
			if err != nil {
				println(err)
			}
			log.Println("历史记录清理完成")
			return
		}
		bot.MQ <- &rcvMsg
		msg, err := chatgpt.ChooseGenerateWay(strconv.FormatInt(rcvMsg.GroupId, 10), rcvMsg.Message, config.Cfg.Context.GroupContext)
		//var err error
		if msg != "" {
			err = bot.SendGroupMsg(rcvMsg.GroupId, "[CQ:reply,id="+strconv.FormatInt(rcvMsg.MessageId, 10)+"]"+msg)
		} else {
			err = bot.SendGroupMsg(rcvMsg.GroupId, "[CQ:reply,id="+strconv.FormatInt(rcvMsg.MessageId, 10)+"]"+"生成错误！错误信息：\n"+err.Error())
		}
		if err != nil {
			log.Println(err)
		}
		<-bot.MQ

	}

}

// TimeOutCheck 检查消息队列中的消息是否超时
func TimeOutCheck() {
	mq := bot.MQ
	for msg := range mq {
		// 搞不懂要不要加锁
		bot.MX.Lock()
		sentTime := time.Unix(msg.Time, 0)
		if time.Now().Sub(sentTime) > time.Duration(config.Cfg.CqHttp.TimeOut)*time.Second && time.Now().Sub(sentTime) < time.Duration(config.Cfg.CqHttp.TimeOut+1)*time.Second {
			log.Println("思考中，请耐心等待~")
			switch msg.MessageType {
			case "private":
				err := bot.SendPrivateMsg(msg.Sender.UserId, "[CQ:reply,id="+strconv.FormatInt(msg.MessageId, 10)+"]"+"思考中，请耐心等待~")
				if err != nil {
					log.Println(err)
				}
			case "group":
				err := bot.SendGroupMsg(msg.GroupId, "[CQ:reply,id="+strconv.FormatInt(msg.MessageId, 10)+"]"+"思考中，请耐心等待~")
				if err != nil {
					log.Println(err)
				}
			}
		}
		mq <- msg
		time.Sleep(time.Second)
		bot.MX.Unlock()
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
