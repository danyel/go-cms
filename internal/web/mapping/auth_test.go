package mapping

import "testing"

func TestLoginRequestIdentity(t *testing.T) {
	i := (LoginRequest{Subject: "s", Email: "e"}).Identity()
	if i.Subject != "s" || i.Email != "e" {
		t.Fatal("identity mapping failed")
	}
}
