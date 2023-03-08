package bot

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	tbi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"

	"github.com/jdxj/oh-my-feed/internal/pkg/config"
	"github.com/jdxj/oh-my-feed/internal/pkg/log"
)

var (
	stop = make(chan int)
	wg   = &sync.WaitGroup{}

	client *tbi.BotAPI
	server *http.Server
	gp     *ants.Pool
)

func Init() {
	var err error
	client, err = tbi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Fatalf("new bot api err: %s", err)
	}

	gp, err = ants.NewPool(100, ants.WithNonblocking(false), ants.WithPanicHandler(func(i interface{}) {
		log.Desugar().Error("gp-panic", zap.Any("panic", i))
	}))
	if err != nil {
		log.Fatalf("new pool err: %s", err)
	}

	registerCmd()
	startWebhook()
	startCanal()
}

func startWebhook() {
	server = &http.Server{
		Addr:                         "0.0.0.0:8080",
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

	updates := client.ListenForWebhook(webhookPath)
	handlers(updates)

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and server webhook err: %s", err)
		}
	}()
}

func Stop() {
	close(stop)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Errorf("stop bot err: %s", err)
	}

	err = gp.ReleaseTimeout(time.Second * 10)
	if err != nil {
		log.Errorf("stop gp err: %s", err)
	}

	myCanal.Close()
	wg.Wait()
}
