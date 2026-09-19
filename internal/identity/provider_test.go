package identity

import (
	"context"
	"testing"
)

func TestGoogleProviderRequiresConfiguration(t *testing.T) {
	_, e := (GoogleProvider{}).Authenticate(context.Background(), "")
	if e == nil {
		t.Fatal("expected configuration error")
	}
}
