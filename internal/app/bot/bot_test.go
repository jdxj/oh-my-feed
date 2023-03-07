package bot

import (
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jdxj/oh-my-feed/internal/pkg/config"
	"github.com/jdxj/oh-my-feed/internal/pkg/db"
	"github.com/jdxj/oh-my-feed/internal/pkg/log"
)

func TestChan(t *testing.T) {
	c := make(chan int, 1)
	close(c)
	c <- 1
	fmt.Printf("dd")
}

func TestFlag(t *testing.T) {
	myFlag := flag.NewFlagSet("myFlag", flag.ContinueOnError)
	args := []string{"123", "def", "", "hh"}
	err := myFlag.Parse(args)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%v\n", myFlag.Args())

	s := "a\tb"
	ss := strings.Split(s, " ")
	fmt.Printf("%v\n", ss)
	fmt.Printf("%d\n", len(ss))
}

func TestStartCanal(t *testing.T) {
	config.Init("../config/config.yaml")
	log.Init()
	db.Init()
	Init()

	time.Sleep(time.Hour)

	Stop()
}

func TestAnyEqual(t *testing.T) {
	var i1 any = "abc"
	var i2 any = "abc2"
	fmt.Printf("%t\n", i1 == i2)
}
