package tools

import (
	"encoding/base64"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"quantity/common"
	"strconv"
)

type Bot struct {
	*tgbotapi.BotAPI
	publicKey int64
}

var (
	notifyFormat = `
# 币种名称: %s
# 订单价格: %.6f
# 成交价格: %.6f
# 成交数量: %s
# 投入资金: %f
# 订单操作: %s
# 策略名称: %s
# 订单时间: %s
# 下单时间: %s
`
)

const (
	privateKey = "NzEzOTQ0NTczNTpBQUdHdExXaldESGdCUUluQ0FGQjhmT1p6eGFLMkdZaU0wQQ=="
	publicKey  = "LTEwMDIxMzEzODY2NTM="
	timeLayout = "2006-01-02 15:04:05"
)

func NewBot() (bot *Bot, err error) {
	token, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return
	}
	tokenApi := string(token)
	tgBot, err := tgbotapi.NewBotAPI(tokenApi)
	if err != nil {
		return
	}

	uid, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return
	}
	userId, err := strconv.ParseInt(string(uid), 10, 64)
	if err != nil {
		return
	}
	bot = &Bot{BotAPI: tgBot, publicKey: userId}
	return
}

func (bot *Bot) SendMessage(message string) (botMessage tgbotapi.Message, err error) {
	msg := tgbotapi.NewMessage(bot.publicKey, message)
	botMessage, err = bot.Send(msg)
	return
}

func (bot *Bot) DeleteMessage(messageID int) (err error) {
	deleteMessageConfig := tgbotapi.DeleteMessageConfig{
		ChatID:    bot.publicKey,
		MessageID: messageID,
	}
	_, err = bot.Request(deleteMessageConfig)
	return
}

func (bot *Bot) NotifyToken(order *common.Order) (botMessage tgbotapi.Message, err error) {
	action := ""
	switch order.Action {
	case common.Buy:
		action = "买入"
	case common.Sell:
		action = "卖出"
	case common.Hold:
		action = "持有"
	default:
		action = "异常"
	}
	message := fmt.Sprintf(notifyFormat,
		order.Symbol,
		order.OrderPrice,
		order.SubmitPrice,
		order.Amount,
		order.Money,
		action,
		order.StrategyName,
		order.OrderTime.Format(timeLayout),
		order.CreatedAt.Format(timeLayout),
	)
	botMessage, err = bot.SendMessage(message)
	return
}
