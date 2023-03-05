package validator

import (
	"context"
	"fmt"
	"testing"
)

func TestValidateFeed(t *testing.T) {
	address, err := ValidateFeed(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Println(address)
}
