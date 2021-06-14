package chatbot

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Response struct {
	Data    interface{} `json:"data",omitempty`
	Message string      `json:"message, omitempty"`
	Code    string      `json:"code, omitempty"`
}

type ChatbotHandler struct {
	Bot     *linebot.Client
	Service *ChatbotService
}

func NewChatbotHandler(service *ChatbotService, bot *linebot.Client) *ChatbotHandler {
	return &ChatbotHandler{Service: service, Bot: bot}
}

func (s *ChatbotHandler) Webhook(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx != nil {
		ctx = context.Background()
	}

	events, err := s.Bot.ParseRequest(c.Request())
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, Response{
				Message: linebot.ErrInvalidSignature.Error(),
				Code:    "1",
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Message: "internal",
				Code:    "99",
			})
		}
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				println(message.Text)
				replyToken := event.ReplyToken
				s.Service.replyMessage(message.Text, s.Bot, replyToken)
			}
		}
	}

	return c.JSON(http.StatusOK, Response{
		Message: "OK",
		Code:    "0",
	})
}
