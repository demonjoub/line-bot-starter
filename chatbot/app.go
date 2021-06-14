package chatbot

import (
	"github.com/demonjoub/chatbot/logger"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

type ChatbotService struct {
}

func NewChatbotService() *ChatbotService {
	return &ChatbotService{}
}

func (s *ChatbotService) replyMessage(message string, lineBot *linebot.Client, replyToken string) error {
	if _, err := lineBot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		logger.Error("reply message", zap.String("Message", err.Error()))
		return err
	}
	return nil
}
