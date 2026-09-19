package persistence

import (
	"github.com/example/cms/internal/domain"
	"testing"
	"time"
)

func TestUserMappingRoundTrip(t *testing.T) {
	d := domain.User{ID: 2, Email: "a@b", Name: "A", Provider: "google", ProviderSubject: "s", CreatedAt: time.Now()}
	if got := UserToDomain(UserFromDomain(d)); got.Email != d.Email || got.ProviderSubject != d.ProviderSubject {
		t.Fatalf("round trip mismatch: %+v", got)
	}
}
