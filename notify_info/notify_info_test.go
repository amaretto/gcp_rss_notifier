package p

import (
	"context"
	"testing"
)

// This test doesn't run yet
func TestNotifyInfo(t *testing.T) {
	m := PubSubMessage{}
	ctx := context.Background()

	err := NotifyInfo(ctx, m)
	if err != nil {
		t.Fatalf("failed test : %v", err)
	}
}
