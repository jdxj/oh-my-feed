package validator

import (
	"context"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mmcdole/gofeed"
)

var (
	varValidator  = validator.New()
	feedValidator = gofeed.NewParser()
)

func ValidateFeed(ctx context.Context, address string) (string, error) {
	address = strings.TrimSuffix(address, "/")
	err := varValidator.Var(address, "url")
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	_, err = feedValidator.ParseURLWithContext(address, ctx)
	return address, err
}

func ValidateFeedWithoutParse(address string) (string, error) {
	address = strings.TrimSuffix(address, "/")
	return address, varValidator.Var(address, "url")
}
