package model

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/jdxj/oh-my-feed/internal/pkg/config"
)

func TestMain(t *testing.M) {
	config.Init("../config/config.yaml")
	Init()

	setDebug()
	os.Exit(t.Run())
}

func TestAddFeed(t *testing.T) {
	_, err := AddFeed(db, "")
	if err != nil {
		t.Fatalf("ii: %s\n", err)
	}
}

func TestGetFeed(t *testing.T) {
	feed, err := GetFeed(db, 16)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%s\n", feed.UpdatedAt.Format(time.StampNano))
}

func TestAddUserFeed(t *testing.T) {
	err := AddUserFeed(context.Background(), 456, "ggg")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}

func TestList(t *testing.T) {
	rsp, err := ListUserFeed(context.Background(), ListUserFeedReq{})
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGoFeed(t *testing.T) {
	p := gofeed.NewParser()
	feed, err := p.ParseURL("")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%+v\n", feed)
}

func TestGetFeeds(t *testing.T) {
	feeds, err := GetFeeds(context.Background())
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	for _, v := range feeds {
		fmt.Println(v.Address)
	}
}

func TestUpdateLatestPost(t *testing.T) {
	err := UpdateLatestPost(context.Background(), 16, "kk")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}
