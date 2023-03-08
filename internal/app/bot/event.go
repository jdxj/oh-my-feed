package bot

import (
	"context"
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/replication"
	tbi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/jdxj/oh-my-feed/internal/app/model"
	"github.com/jdxj/oh-my-feed/internal/pkg/config"
	"github.com/jdxj/oh-my-feed/internal/pkg/log"
)

var (
	myCanal *canal.Canal
)

type myEventHandler struct {
	canal.DummyEventHandler
}

func (hdl *myEventHandler) OnRow(e *canal.RowsEvent) error {
	if e.Header.EventType != replication.UPDATE_ROWS_EVENTv2 {
		return nil
	}

	if len(e.Rows) < 2 {
		log.Errorf("invalid on row len: %d", len(e.Rows))
		return nil
	}

	oldRow := e.Rows[0]
	newRow := e.Rows[1]
	if len(newRow) < 7 {
		log.Errorf("table struct changed")
		return nil
	}

	id, ok := newRow[0].(uint64)
	if !ok {
		log.Errorf("column type changed: id")
		return nil
	}

	if oldRow[6] == newRow[6] {
		// 不是latest_post发生变化
		return nil
	}

	latestPost, ok := newRow[6].(string)
	if !ok {
		log.Errorf("column type changed: latest_post")
		return nil
	}

	log.Desugar().Debug("on-row", zap.String("latest-post", latestPost))

	sendLatestPost(id, latestPost)
	return nil
}

// todo: 应该抽出到一个包里
func startCanal() {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", config.DB.Address, config.DB.Port)
	cfg.User = config.DB.User
	cfg.Password = config.DB.Password
	cfg.Dump.TableDB = config.DB.Dbname
	cfg.Dump.Tables = []string{"feeds"}
	cfg.Dump.ExecutionPath = ""
	cfg.DisableRetrySync = true

	var err error
	myCanal, err = canal.NewCanal(cfg)
	if err != nil {
		log.Fatalf("new canal err: %s", err)
	}

	myCanal.SetEventHandler(&myEventHandler{})

	pos, err := myCanal.GetMasterPos()
	if err != nil {
		log.Fatalf("get master pos err: %s", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := myCanal.RunFrom(pos)
		if err != nil {
			log.Errorf("run canal err: %s", err)
		}
	}()
}

func sendLatestPost(id uint64, url string) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		uf, err := model.ListUserFeed(context.TODO(), model.ListUserFeedReq{FeedID: id})
		if err != nil {
			log.Errorf("list user feed err: %s", err)
			return
		}

		for _, v := range uf.UserFeeds {
			select {
			case <-stop:
				log.Infof("stop send latest post")
				return
			default:
			}

			chatID := v.TelegramID
			_ = gp.Submit(func() {
				msg := tbi.NewMessage(chatID, url)
				_, err = client.Send(msg)
				if err != nil {
					log.Errorf("send latest post err: %s", err)
				}
			})
		}
	}()
}
