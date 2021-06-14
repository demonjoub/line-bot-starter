package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/demonjoub/chatbot/chatbot"
	"github.com/demonjoub/chatbot/config"
	"github.com/demonjoub/chatbot/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

var project *config.Config

func main() {
	// read config
	config.GetSecretValue()
	project = config.NewConfig()
	port := project.App.Port

	e := newRouter()
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "ok",
		})
	})

	chatbotService := chatbot.NewChatbotService()
	chatbotHandler := chatbot.NewChatbotHandler(chatbotService, connectLineBot())
	e.POST("/webhook", chatbotHandler.Webhook)

	go run(e, port)

	shutdown(e)
}

func connectLineBot() *linebot.Client {
	channelSecret := project.Line.ChannelSecret
	channelToken := project.Line.ChannelAccessToken
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

func run(e *echo.Echo, port int) {
	logger.Info("Starting...")
	startPort := fmt.Sprintf(":%d", port)
	logger.Info(fmt.Sprintf("Start Server On %s", startPort))
	e.Logger.Fatal(e.Start(startPort))
}

func newRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Secure())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
		HSTSExcludeSubdomains: true,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))
	return e
}

func shutdown(e *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	if err := e.Shutdown(context.Background()); err != nil {
		logger.Error(err.Error(), zap.String("tag", "shutdown Server"))
	}
}
