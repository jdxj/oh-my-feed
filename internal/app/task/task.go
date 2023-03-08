package task

import (
	"context"
	"errors"
	"time"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"

	"github.com/jdxj/oh-my-feed/internal/app/model"
	"github.com/jdxj/oh-my-feed/internal/pkg/log"
)

var (
	myParser = gofeed.NewParser()

	stop    = make(chan int)
	stopped = make(chan int)
)

var (
	// 延迟启动
	startDelay = time.Second * 5
	// 更新间隔
	updateInterval = time.Hour
	// 更新超时
	updateTimeout = time.Hour * 24
	// 解析超时
	parseTimeout = time.Second * 10
)

func Init() {
	start()
}

func start() {
	go func() {
		timer := time.NewTimer(startDelay)
		defer timer.Stop()

		// 延迟启动
		<-timer.C
		for {
			updateFeedTitle()

			timer.Reset(updateInterval)
			select {
			case <-stop:
				close(stopped)
				return
			case <-timer.C:
			}
		}
	}()
}

func Stop() {
	close(stop)
	<-stopped
}

func updateFeedTitle() {
	log.Infof("start update feed title")

	ctx, cancel := context.WithTimeout(context.Background(), updateTimeout)
	defer cancel()

	feeds, err := model.GetFeeds(ctx)
	if err != nil {
		log.Errorf("get feeds err: %s", err)
		return
	}

	for _, feed := range feeds {
		select {
		case <-stop:
			log.Infof("stop update feed title")
			return
		default:
		}

		latestTitle, err := getLatestPost(ctx, feed.Address)
		if err != nil {
			log.Desugar().Warn(
				"feed-latest-title",
				zap.String("address", feed.Address),
				zap.Error(err),
			)
			continue
		}

		err = model.UpdateLatestPost(ctx, feed.ID, latestTitle)
		if err != nil {
			log.Desugar().Warn(
				"update-feed-title",
				zap.Uint("feedID", feed.ID),
				zap.Error(err),
			)
		}
	}
}

var (
	ErrFeedItemNotFound = errors.New("feed item not found")
)

func getLatestPost(ctx context.Context, address string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, parseTimeout)
	defer cancel()

	feed, err := myParser.ParseURLWithContext(address, ctx)
	if err != nil {
		return "", err
	}
	if len(feed.Items) == 0 {
		return "", ErrFeedItemNotFound
	}

	return feed.Items[0].Link, nil
}

var (
	ErrIntervalTooSmall = errors.New("interval too small")
)

func SetInterval(dur string) error {
	d, err := time.ParseDuration(dur)
	if err != nil {
		return err
	}
	if d < time.Minute*10 {
		return ErrIntervalTooSmall
	}
	updateInterval = d
	return nil
}
