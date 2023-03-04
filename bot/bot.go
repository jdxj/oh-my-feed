package bot

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	tbi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/jdxj/oh-my-feed/config"
	"github.com/jdxj/oh-my-feed/log"
)

var (
	client *tbi.BotAPI
	server *http.Server
)

func Init() {
	var err error
	client, err = tbi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Fatalf("new bot api err: %s", err)
	}

	server = &http.Server{
		Addr:                         "localhost:8080",
		Handler:                      http.DefaultServeMux,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  time.Second * 30,
		ReadHeaderTimeout:            time.Second * 10,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}

	registerCmd()
	start()
}

func start() {
	// 注册webhook
	webhookPath := "/" + config.Telegram.Token
	webhook := strings.TrimSuffix(config.Telegram.Webhook, "/") + webhookPath
	webhookReq, err := tbi.NewWebhook(webhook)
	if err != nil {
		log.Fatalf("new webhook err: %s", err)
	}

	webhookRsp, err := client.Request(webhookReq)
	if err != nil {
		log.Fatalf("request webhook err: %s", err)
	}
	if !webhookRsp.Ok {
		log.Desugar().Fatal(
			"request webhook wrong",
			zap.String("description", webhookRsp.Description),
		)
	}

	webhookInfo, err := client.GetWebhookInfo()
	if err != nil {
		log.Warnf("get webhook info err: %s", err)
	} else {
		log.Desugar().Info("webhook-info", zap.String("url", webhookInfo.URL))

		if webhookInfo.LastErrorDate != 0 {
			log.Desugar().Info(
				"webhook-info",
				zap.String("last-error-date", time.Unix(int64(webhookInfo.LastErrorDate), 0).In(time.Local).Format(time.DateTime)),
				zap.String("last-error-message", webhookInfo.LastErrorMessage),
			)
		}
	}

	// todo: 1. 删除命令 2. 注册命令
	updates := client.ListenForWebhook(webhookPath)
	handlers(updates)

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and server webhook err: %s", err)
		}
	}()
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Errorf("stop bot err: %s", err)
	}
}
